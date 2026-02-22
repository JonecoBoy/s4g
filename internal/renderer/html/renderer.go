// Package html implements a Renderer that writes HTML files using Go templates.
package html

import (
	"context"
	"fmt"
	"html/template"
	"os"
	"path/filepath"

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

	tmpl *template.Template
}

// New creates a new HTMLRenderer. Call Init before Render.
func New(templateDir, outputDir, siteTitle string) *Renderer {
	return &Renderer{
		TemplateDir: templateDir,
		OutputDir:   outputDir,
		SiteTitle:   siteTitle,
	}
}

// Name implements core.Renderer.
func (r *Renderer) Name() string { return "html" }

// funcMap provides template helper functions.
var funcMap = template.FuncMap{
	// safeHTML marks a string as safe HTML, bypassing auto-escaping.
	// Only use with content that has already been sanitised (e.g. from Goldmark).
	"safeHTML": func(s string) template.HTML { return template.HTML(s) }, //nolint:gosec
}

// Init parses the template files. Must be called once before Render.
func (r *Renderer) Init() error {
	pattern := filepath.Join(r.TemplateDir, "*.html")
	tmpl, err := template.New("").Funcs(funcMap).ParseGlob(pattern)
	if err != nil {
		return fmt.Errorf("html renderer: parse templates %q: %w", pattern, err)
	}
	r.tmpl = tmpl

	if err := os.MkdirAll(r.OutputDir, 0755); err != nil {
		return fmt.Errorf("html renderer: create output dir %q: %w", r.OutputDir, err)
	}
	return nil
}

// Render implements core.Renderer.
// It executes the "page.html" template and writes slug.html to OutputDir.
// When c.Section is set the file is written to OutputDir/<section>/slug.html.
func (r *Renderer) Render(_ context.Context, all []core.Content, c core.Content) error {
	if r.tmpl == nil {
		return fmt.Errorf("html renderer: not initialised — call Init() first")
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

	if err := r.tmpl.ExecuteTemplate(f, "page.html", data); err != nil {
		return fmt.Errorf("html renderer: execute template for %q: %w", c.Slug, err)
	}
	return nil
}
