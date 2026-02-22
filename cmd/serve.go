package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newServeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "Serve the generated site locally (coming soon)",
		Long:  `Start a local HTTP server to preview the generated dist/ directory.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("🚧  serve is not yet implemented — run `s4g build` first, then open dist/ in a browser.")
			return nil
		},
	}
}
