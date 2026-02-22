// Package config loads and holds s4g site configuration from config.yaml.
package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config is the top-level configuration structure for s4g.
type Config struct {
	Site     SiteConfig     `yaml:"site"`
	Source   SourceConfig   `yaml:"source"`
	Renderer RendererConfig `yaml:"renderer"`
}

// SiteConfig holds general site metadata.
type SiteConfig struct {
	Title   string `yaml:"title"`
	BaseURL string `yaml:"base_url"`
}

// SourceConfig selects and configures the active DataSource.
type SourceConfig struct {
	// Type selects the DataSource driver: "markdown", "json", "sql", …
	Type       string `yaml:"type"`
	ContentDir string `yaml:"content_dir"`
}

// RendererConfig selects and configures the active Renderer.
type RendererConfig struct {
	// Type selects the Renderer driver: "html", "rss", …
	Type        string `yaml:"type"`
	TemplateDir string `yaml:"template_dir"`
	OutputDir   string `yaml:"output_dir"`
}

// Load reads and parses the YAML config file at path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config: read %q: %w", path, err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("config: parse %q: %w", path, err)
	}
	return &cfg, nil
}
