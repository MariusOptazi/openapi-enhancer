package main

import (
	"os"
	"testing"
)

func TestRun(t *testing.T) {
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	os.Args = []string{
		"openapi-processor",
		"-file", "../../testdata/openapi.yaml",
	}

	err := run()
	if err != nil {
		t.Fatalf("run() failed: %v", err)
	}
}
