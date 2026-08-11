package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	InputFile string
}

func Load() (*Config, error) {
	cfg := &Config{}

	flag.StringVar(&cfg.InputFile, "file", "", "Path to OpenAPI YAML file (required)")
	flag.StringVar(&cfg.InputFile, "f", "", "Shorthand for -file")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: openapi-enhancer [options]\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExample:\n")
		fmt.Fprintf(os.Stderr, "  openapi-enhancer -file api/openapi.bundled.yaml\n")
	}

	flag.Parse()

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.InputFile == "" {
		return fmt.Errorf("-file flag is required")
	}

	info, err := os.Stat(c.InputFile)
	if err != nil {
		return fmt.Errorf("file '%s' not found: %w", c.InputFile, err)
	}

	if info.IsDir() {
		return fmt.Errorf("'%s' is a directory, expected a file", c.InputFile)
	}

	ext := filepath.Ext(info.Name())
	if ext != ".yaml" && ext != ".yml" {
		return fmt.Errorf("file '%s' must be .yaml or .yml", c.InputFile)
	}

	return nil
}
