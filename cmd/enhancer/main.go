package main

import (
	"fmt"
	"os"

	"github.com/MariusOptazi/openapi-enhancer/internal/config"
	"github.com/MariusOptazi/openapi-enhancer/internal/service"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("configuration error: %w", err)
	}

	processor := service.NewProcessor(cfg.InputFile)
	if err := processor.Run(); err != nil {
		return fmt.Errorf("Execution error: %v\n", err)
	}

	return nil
}
