// Package cmd wires all Cobra subcommands for the s4g CLI.
package cmd

import (
	"github.com/spf13/cobra"
)

// NewRootCmd returns the root Cobra command.
func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "s4g",
		Short: "s4g — Super Simple Static Site Generator",
		Long: `s4g is a pluggable static site generator written in Go.
Data sources and renderers are swappable interfaces, so you can feed
content from Markdown, JSON, SQL, REST or any future driver.`,
	}

	root.AddCommand(newBuildCmd())
	root.AddCommand(newServeCmd())

	return root
}
