// Package core defines the fundamental types and interfaces of s4g.
// All data sources and renderers depend only on these types, never on each other.
package core

import "context"

// Content is the canonical, source-agnostic representation of a page.
// A DataSource converts its native format (Markdown, SQL row, JSON object, …)
// into []Content. A Renderer converts []Content into output files.
type Content struct {
	// Slug is the URL-safe identifier used as the output filename (e.g. "about").
	Slug string
	// Title is the human-readable page title.
	Title string
	// Body is the pre-rendered HTML body of the page.
	Body string
	// Meta holds any additional key/value pairs from the source
	// (front-matter fields, DB columns, JSON keys, etc.).
	Meta map[string]any
}

// DataSource is the plugin interface for content providers.
// Implement this to add a new data source (Markdown, SQL, Mongo, REST, …).
type DataSource interface {
	// Name returns a human-readable identifier used in logs and the TUI.
	Name() string
	// Fetch loads and returns all Content items from this source.
	Fetch(ctx context.Context) ([]Content, error)
}

// Renderer is the plugin interface for output generators.
// Implement this to add a new output format (HTML, RSS, JSON feed, …).
type Renderer interface {
	// Name returns a human-readable identifier used in logs and the TUI.
	Name() string
	// Render writes output for a single Content item.
	Render(ctx context.Context, c Content) error
}
