// Package html implements a Renderer that writes HTML files using Go templates.
package html

import (
	"context"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/user/s4g/internal/core"
)

// TemplateData is the data passed into each template execution.
type TemplateData struct {
	SiteTitle string
	Pages     []core.Content
	Page      core.Content
}

// Renderer writes one HTML file per Content item using Go html/template.
type Renderer struct {
	// TemplateDir is the directory containing layout.html and page.html.
	TemplateDir string
	// OutputDir is the directory where HTML files are written.
	OutputDir string
	// SiteTitle is injected into every page's template data.
	SiteTitle string

	// base is the universal template set from TemplateDir root.
	base *template.Template
	// cache stores cloned templates with section-specific overrides.
	cache map[string]*template.Template
	mu    sync.RWMutex
}

// New creates a new HTMLRenderer. Call Init before Render.
func New(templateDir, outputDir, siteTitle string) *Renderer {
	return &Renderer{
		TemplateDir: templateDir,
		OutputDir:   outputDir,
		SiteTitle:   siteTitle,
		cache:       make(map[string]*template.Template),
	}
}

// Name implements core.Renderer.
func (r *Renderer) Name() string { return "html" }

// funcMap provides template helper functions.
var funcMap = template.FuncMap{
	"safeHTML": func(s string) template.HTML { return template.HTML(s) }, //nolint:gosec
	"truncate": func(limit int, s string) string {
		if len(s) <= limit {
			return s
		}
		// Try to truncate at the nearest space to avoid cutting words in half
		truncated := s[:limit]
		if lastSpace := strings.LastIndexAny(truncated, " \t\n\r"); lastSpace > 0 {
			truncated = truncated[:lastSpace]
		}
		return truncated + "..."
	},
}

// Init parses the universal template files. Must be called once before Render.
func (r *Renderer) Init() error {
	pattern := filepath.Join(r.TemplateDir, "*.html")
	tmpl, err := template.New("base").Funcs(funcMap).ParseGlob(pattern)
	if err != nil {
		return fmt.Errorf("html renderer: parse universal templates %q: %w", pattern, err)
	}
	r.base = tmpl

	if err := os.MkdirAll(r.OutputDir, 0755); err != nil {
		return fmt.Errorf("html renderer: create output dir %q: %w", r.OutputDir, err)
	}
	return nil
}

// getTemplate returns the template set for a section, cloning and applying overrides if needed.
func (r *Renderer) getTemplate(section string) (*template.Template, error) {
	r.mu.RLock()
	tmpl, ok := r.cache[section]
	r.mu.RUnlock()
	if ok {
		return tmpl, nil
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// Double-check cache inside lock
	if tmpl, ok := r.cache[section]; ok {
		return tmpl, nil
	}

	// Start with a clone of the base template.
	// Cloning is thread-safe and allows us to override definitions like "page.html"
	// without affecting the base or other sections.
	derived, err := r.base.Clone()
	if err != nil {
		return nil, fmt.Errorf("failed to clone base template: %w", err)
	}

	if section != "" {
		// Look for overrides in TemplateDir/<section>/*.html.
		// We use filepath.Join(r.TemplateDir, section) to match the content folder hierarchy.
		sectionDir := filepath.Join(r.TemplateDir, section)
		if info, err := os.Stat(sectionDir); err == nil && info.IsDir() {
			pattern := filepath.Join(sectionDir, "*.html")
			// ParseGlob on the clone will overwrite template definitions with the same name.
			if _, err := derived.ParseGlob(pattern); err != nil {
				return nil, fmt.Errorf("failed to parse overrides in %q: %w", sectionDir, err)
			}
		}
	}

	r.cache[section] = derived
	return derived, nil
}

// Render implements core.Renderer.
// It executes the "page.html" template and writes slug.html to OutputDir.
func (r *Renderer) Render(_ context.Context, all []core.Content, c core.Content) error {
	if r.base == nil {
		return fmt.Errorf("html renderer: not initialised — call Init() first")
	}

	tmpl, err := r.getTemplate(c.Section)
	if err != nil {
		return fmt.Errorf("html renderer: get template for section %q: %w", c.Section, err)
	}

	// Determine the output subdirectory.
	outDir := r.OutputDir
	if c.Section != "" {
		outDir = filepath.Join(r.OutputDir, c.Section)
		if err := os.MkdirAll(outDir, 0755); err != nil {
			return fmt.Errorf("html renderer: create section dir %q: %w", outDir, err)
		}
	}

	outPath := filepath.Join(outDir, c.Slug+".html")
	f, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("html renderer: create %q: %w", outPath, err)
	}
	defer f.Close()

	data := TemplateData{
		SiteTitle: r.SiteTitle,
		Pages:     all,
		Page:      c,
	}

	if err := tmpl.ExecuteTemplate(f, "page.html", data); err != nil {
		return fmt.Errorf("html renderer: execute template for %q: %w", c.Slug, err)
	}
	return nil
}
