// Package markdown implements a DataSource that reads .md files with YAML front-matter.
package markdown

import (
	"bytes"
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/user/s4g/internal/core"
	"github.com/yuin/goldmark"
	"gopkg.in/yaml.v3"
)

// FrontMatter represents the YAML block at the top of a Markdown file.
type FrontMatter struct {
	Title string         `yaml:"title"`
	Slug  string         `yaml:"slug"`
	Tags  []string       `yaml:"tags"`
	Meta  map[string]any `yaml:",inline"`
}

// Source is a DataSource that reads *.md files recursively from a directory.
type Source struct {
	// Dir is the path to the directory containing Markdown files.
	Dir string
}

// New returns a new Markdown Source rooted at dir.
func New(dir string) *Source {
	return &Source{Dir: dir}
}

// Name implements core.DataSource.
func (s *Source) Name() string { return "markdown" }

// Fetch implements core.DataSource.
// It walks Dir and its subdirectories, parsing each .md file.
// The Section is derived from the subdirectory path relative to Dir.
func (s *Source) Fetch(_ context.Context) ([]core.Content, error) {
	var pages []core.Content

	err := filepath.WalkDir(s.Dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories themselves, and non-markdown files.
		if d.IsDir() || !strings.HasSuffix(d.Name(), ".md") {
			return nil
		}

		// Calculate the section based on relative path.
		// Example: if Dir is "site/content" and file is "site/content/blog/post.md"
		// relDir becomes "blog".
		relPath, err := filepath.Rel(s.Dir, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path for %q: %w", path, err)
		}

		relDir := filepath.Dir(relPath)
		section := ""
		if relDir != "." {
			section = relDir
		}

		raw, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("markdown source: read file %q: %w", path, err)
		}

		fm, body, err := parseFrontMatter(raw)
		if err != nil {
			return fmt.Errorf("markdown source: parse front-matter in %q: %w", path, err)
		}

		// Default slug to filename without extension.
		slug := fm.Slug
		if slug == "" {
			slug = strings.TrimSuffix(d.Name(), ".md")
		}

		// Default tags to empty slice if not provided.
		tags := fm.Tags
		if tags == nil {
			tags = []string{}
		}

		// Render Markdown body to HTML.
		var buf bytes.Buffer
		if err := goldmark.Convert(body, &buf); err != nil {
			return fmt.Errorf("markdown source: render %q: %w", path, err)
		}

		meta := make(map[string]any)
		for k, v := range fm.Meta {
			meta[k] = v
		}

		pages = append(pages, core.Content{
			Slug:    slug,
			Section: section,
			Title:   fm.Title,
			Body:    buf.String(),
			Tags:    fm.Tags,
			Meta:    meta,
		})

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("markdown source: walk dir %q: %w", s.Dir, err)
	}

	return pages, nil
}

// parseFrontMatter splits a raw Markdown file into its YAML front-matter (metadata) and body.
func parseFrontMatter(raw []byte) (FrontMatter, []byte, error) {
	var fm FrontMatter

	if !bytes.HasPrefix(raw, []byte("---\n")) {
		return fm, raw, nil
	}

	rest := raw[4:]
	idx := bytes.Index(rest, []byte("\n---"))
	if idx == -1 {
		return fm, raw, fmt.Errorf("unclosed front-matter block")
	}

	yamlBlock := rest[:idx]
	body := rest[idx+4:]

	if err := yaml.Unmarshal(yamlBlock, &fm); err != nil {
		return fm, body, err
	}

	body = bytes.TrimPrefix(body, []byte("\n"))
	return fm, body, nil
}
