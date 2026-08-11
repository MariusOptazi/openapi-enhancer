package service

import (
	"fmt"
	"os"

	"github.com/MariusOptazi/openapi-enhancer/internal/builder"
	"github.com/oasdiff/yaml"
)

type Processor struct {
	FilePath string
}

func NewProcessor(filePath string) *Processor {
	return &Processor{FilePath: filePath}
}

func (p *Processor) Run() error {
	fmt.Printf("Reading: %s\n", p.FilePath)

	data, err := os.ReadFile(p.FilePath)
	if err != nil {
		return fmt.Errorf("reading file: %w", err)
	}

	var openAPI map[string]interface{}
	if _, err := yaml.Unmarshal(data, &openAPI, yaml.DecodeOpts{}); err != nil {
		return fmt.Errorf("parsing YAML: %w", err)
	}

	if err := p.processDTOs(openAPI); err != nil {
		return fmt.Errorf("processing DTOs: %w", err)
	}

	outputData, err := yaml.Marshal(openAPI)
	if err != nil {
		return fmt.Errorf("marshaling YAML: %w", err)
	}

	if err := os.WriteFile(p.FilePath, outputData, 0644); err != nil {
		return fmt.Errorf("writing file: %w", err)
	}

	fmt.Printf("Successfully updated: %s\n", p.FilePath)
	return nil
}

// Helper //
func (p *Processor) processDTOs(openAPI map[string]interface{}) error {
	components, ok := openAPI["components"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("no components section found")
	}

	schemas, ok := components["schemas"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("no schemas section found")
	}

	fmt.Printf("Found %d schemas\n", len(schemas))

	for dtoName, dto := range schemas {
		dtoMap, ok := dto.(map[string]interface{})
		if !ok {
			continue
		}

		properties, ok := dtoMap["properties"].(map[string]interface{})
		if !ok {
			continue
		}

		requiredFields := map[string]bool{}
		if reqList, ok := dtoMap["required"].([]interface{}); ok {
			for _, r := range reqList {
				if name, ok := r.(string); ok {
					requiredFields[name] = true
				}
			}
		}

		for fieldName, field := range properties {
			fieldMap, ok := field.(map[string]interface{})
			if !ok {
				continue
			}

			tag := builder.BuildValidatorTag(fieldMap, requiredFields[fieldName])
			if tag != "" {
				fieldMap["x-oapi-codegen-extra-tags"] = map[string]interface{}{
					"validate": tag,
				}
				fmt.Printf("   ✓ %s.%s → %s\n", dtoName, fieldName, tag)
			}
		}
	}

	return nil
}
