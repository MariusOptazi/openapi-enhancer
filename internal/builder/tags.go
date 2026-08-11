package builder

import (
	"fmt"
	"strings"
)

func BuildValidatorTag(field map[string]interface{}, required bool) string {
	var tags []string

	if required {
		tags = append(tags, "required")
	}

	fieldType, _ := field["type"].(string)

	switch fieldType {
	case "string":
		if min, ok := field["minLength"]; ok {
			tags = append(tags, fmt.Sprintf("min=%v", min))
		}
		if max, ok := field["maxLength"]; ok {
			tags = append(tags, fmt.Sprintf("max=%v", max))
		}
		if format, ok := field["format"].(string); ok {
			switch format {
			case "email":
				tags = append(tags, "email")
			case "uuid":
				tags = append(tags, "uuid")
			case "uri", "url":
				tags = append(tags, "uri")
			case "date":
				tags = append(tags, "datetime=2006-01-02")
			case "date-time":
				tags = append(tags, "datetime=2006-01-02T15:04:05Z07:00")
			}
		}
		if enumRaw, ok := field["enum"].([]interface{}); ok && len(enumRaw) > 0 {
			values := make([]string, 0, len(enumRaw))
			for _, e := range enumRaw {
				values = append(values, fmt.Sprintf("%v", e))
			}
			tags = append(tags, "oneof="+strings.Join(values, " "))
		}

	case "integer", "number":
		if min, ok := field["minimum"]; ok {
			tags = append(tags, fmt.Sprintf("min=%v", min))
		}
		if max, ok := field["maximum"]; ok {
			tags = append(tags, fmt.Sprintf("max=%v", max))
		}

	case "array":
		if min, ok := field["minItems"]; ok {
			tags = append(tags, fmt.Sprintf("min=%v", min))
		}
		if max, ok := field["maxItems"]; ok {
			tags = append(tags, fmt.Sprintf("max=%v", max))
		}
	}

	return strings.Join(tags, ",")
}
