package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/s4g/internal/config"
	"github.com/user/s4g/internal/core"
	htmlrenderer "github.com/user/s4g/internal/renderer/html"
	mdsource "github.com/user/s4g/internal/source/markdown"
	"github.com/user/s4g/internal/tui"
)

var cfgFile string

func newBuildCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "build",
		Short: "Build the static site",
		Long:  `Fetch content from all configured sources and render output files.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(cfgFile)
			if err != nil {
				return fmt.Errorf("loading config: %w", err)
			}

			// --- Wire DataSource --------------------------------------------------
			var sources []core.DataSource
			switch cfg.Source.Type {
			case "markdown":
				sources = append(sources, mdsource.New(cfg.Source.ContentDir))
			default:
				return fmt.Errorf("unknown source type %q (available: markdown)", cfg.Source.Type)
			}

			// --- Wire Renderer ---------------------------------------------------
			var renderers []core.Renderer
			switch cfg.Renderer.Type {
			case "html":
				r := htmlrenderer.New(cfg.Renderer.TemplateDir, cfg.Renderer.OutputDir, cfg.Site.Title)
				if err := r.Init(); err != nil {
					return err
				}
				renderers = append(renderers, r)
			default:
				return fmt.Errorf("unknown renderer type %q (available: html)", cfg.Renderer.Type)
			}

			// --- Run generator ---------------------------------------------------
			gen := &core.Generator{
				Sources:   sources,
				Renderers: renderers,
			}
			results := gen.Build(context.Background())

			// Check for errors before launching TUI
			for _, r := range results {
				for _, e := range r.Errors {
					fmt.Fprintf(os.Stderr, "error: %v\n", e)
				}
			}

			// --- Launch Bubbletea TUI summary -------------------------------------
			return tui.Run(results)
		},
	}

	cmd.Flags().StringVarP(&cfgFile, "config", "c", "config.yaml", "path to config.yaml")
	return cmd
}
