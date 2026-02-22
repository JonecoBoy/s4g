// Package markdown implements a DataSource that reads .md files with YAML front-matter.
package markdown

import (
	"bytes"
	"context"
	"fmt"
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
	Extra map[string]any `yaml:",inline"`
}

// Source is a DataSource that reads *.md files from a directory.
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
// It walks Dir, parses each .md file, and returns a slice of core.Content.
func (s *Source) Fetch(_ context.Context) ([]core.Content, error) {
	entries, err := os.ReadDir(s.Dir)
	if err != nil {
		return nil, fmt.Errorf("markdown source: read dir %q: %w", s.Dir, err)
	}

	var pages []core.Content
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}

		path := filepath.Join(s.Dir, e.Name())
		raw, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("markdown source: read file %q: %w", path, err)
		}

		fm, body, err := parseFrontMatter(raw)
		if err != nil {
			return nil, fmt.Errorf("markdown source: parse front-matter in %q: %w", path, err)
		}

		// Default slug to filename without extension.
		slug := fm.Slug
		if slug == "" {
			slug = strings.TrimSuffix(e.Name(), ".md")
		}

		// Render Markdown body to HTML.
		var buf bytes.Buffer
		if err := goldmark.Convert(body, &buf); err != nil {
			return nil, fmt.Errorf("markdown source: render %q: %w", path, err)
		}

		meta := make(map[string]any)
		for k, v := range fm.Extra {
			meta[k] = v
		}

		pages = append(pages, core.Content{
			Slug:  slug,
			Title: fm.Title,
			Body:  buf.String(),
			Meta:  meta,
		})
	}

	return pages, nil
}

// parseFrontMatter splits a raw Markdown file into its YAML front-matter and body.
// Front-matter is delimited by leading and trailing "---" lines.
func parseFrontMatter(raw []byte) (FrontMatter, []byte, error) {
	var fm FrontMatter

	// Files must start with "---\n"
	if !bytes.HasPrefix(raw, []byte("---\n")) {
		// No front-matter; treat entire file as body.
		return fm, raw, nil
	}

	rest := raw[4:] // skip opening ---
	idx := bytes.Index(rest, []byte("\n---"))
	if idx == -1 {
		return fm, raw, fmt.Errorf("unclosed front-matter block")
	}

	yamlBlock := rest[:idx]
	body := rest[idx+4:] // skip closing ---

	if err := yaml.Unmarshal(yamlBlock, &fm); err != nil {
		return fm, body, err
	}

	// Trim leading newline from body
	body = bytes.TrimPrefix(body, []byte("\n"))

	return fm, body, nil
}
