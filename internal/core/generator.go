package core

import (
	"context"
	"fmt"
)

// BuildResult summarises a completed build, used by the TUI.
type BuildResult struct {
	SourceName   string
	RendererName string
	Pages        []Content
	Errors       []error
}

// Generator orchestrates one or more DataSources feeding one or more Renderers.
// Adding a new source or renderer is purely additive — no core changes required.
type Generator struct {
	Sources   []DataSource
	Renderers []Renderer
}

// Build fetches from every source and renders every page through every renderer.
// It collects results and errors non-fatally so the TUI can display a full report.
func (g *Generator) Build(ctx context.Context) []BuildResult {
	var results []BuildResult

	for _, src := range g.Sources {
		pages, err := src.Fetch(ctx)
		br := BuildResult{
			SourceName: src.Name(),
			Pages:      pages,
		}
		if err != nil {
			br.Errors = append(br.Errors, fmt.Errorf("source %q fetch: %w", src.Name(), err))
			results = append(results, br)
			continue
		}

		for _, r := range g.Renderers {
			br.RendererName = r.Name()
			for _, p := range pages {
				if rerr := r.Render(ctx, p); rerr != nil {
					br.Errors = append(br.Errors, fmt.Errorf("renderer %q / page %q: %w", r.Name(), p.Slug, rerr))
				}
			}
		}

		results = append(results, br)
	}

	return results
}
