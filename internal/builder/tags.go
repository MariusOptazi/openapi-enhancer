package builder

import (
	"fmt"
	"strings"
)

func BuildValidatorTag(field map[string]interface{}, required bool) string {
	var rules []string

	// 1. Nullability prüfen (OAS 3.0 "nullable: true" & OAS 3.1 "type: ['string', 'null']")
	isNullable := false
	if n, ok := field["nullable"].(bool); ok && n {
		isNullable = true
	}
	if types, ok := field["type"].([]interface{}); ok {
		for _, t := range types {
			if t == "null" {
				isNullable = true
				break
			}
		}
	}

	// 2. Präfix-Logik für required / omitempty / omitnil
	var prefix string
	if required {
		if !isNullable {
			prefix = "required"
		} else {
			// Feld ist Pflicht, darf aber explizit null/leer sein
			prefix = "omitempty"
		}
	} else {
		if isNullable {
			// Optional und darf null sein
			prefix = "omitempty"
		} else {
			// Optional, aber WENN übergeben, darf es nicht null/leer sein
			prefix = "omitnil"
		}
	}

	// 3. Typ-spezifische Regeln extrahieren
	fieldType := extractPrimaryType(field)

	switch fieldType {
	case "string":
		if min, ok := field["minLength"]; ok {
			rules = append(rules, fmt.Sprintf("min=%v", min))
		}
		if max, ok := field["maxLength"]; ok {
			rules = append(rules, fmt.Sprintf("max=%v", max))
		}
		if pattern, ok := field["pattern"].(string); ok && pattern != "" {
			rules = append(rules, fmt.Sprintf("regex=%s", pattern))
		}
		if format, ok := field["format"].(string); ok {
			switch format {
			case "email":
				rules = append(rules, "email")
			case "uuid":
				rules = append(rules, "uuid")
			case "uri", "url":
				rules = append(rules, "uri")
			case "date":
				rules = append(rules, "datetime=2006-01-02")
			case "date-time":
				rules = append(rules, "datetime=2006-01-02T15:04:05Z07:00")
			case "ipv4":
				rules = append(rules, "ipv4")
			case "ipv6":
				rules = append(rules, "ipv6")
			}
		}
		if enumRaw, ok := field["enum"].([]interface{}); ok && len(enumRaw) > 0 {
			values := make([]string, 0, len(enumRaw))
			for _, e := range enumRaw {
				values = append(values, fmt.Sprintf("%v", e))
			}
			rules = append(rules, "oneof="+strings.Join(values, " "))
		}

	case "integer", "number":
		if min, ok := field["minimum"]; ok {
			rules = append(rules, fmt.Sprintf("min=%v", min))
		}
		if max, ok := field["maximum"]; ok {
			rules = append(rules, fmt.Sprintf("max=%v", max))
		}

	case "array":
		if min, ok := field["minItems"]; ok {
			rules = append(rules, fmt.Sprintf("min=%v", min))
		}
		if max, ok := field["maxItems"]; ok {
			rules = append(rules, fmt.Sprintf("max=%v", max))
		}
		if unique, ok := field["uniqueItems"].(bool); ok && unique {
			rules = append(rules, "unique")
		}
	}

	// 4. Tags zusammensetzen
	var finalTags []string

	// "required" immer voranstellen
	if prefix == "required" {
		finalTags = append(finalTags, "required")
	}

	// Wenn keine weiteren Regeln existieren und das Feld optional ist,
	// wird kein Tag benötigt (Feld wird schlicht ignoriert).
	if len(rules) > 0 {
		if prefix == "omitempty" || prefix == "omitnil" {
			finalTags = append(finalTags, prefix)
		}
		finalTags = append(finalTags, rules...)
	}

	return strings.Join(finalTags, ",")
}

func extractPrimaryType(field map[string]interface{}) string {
	switch t := field["type"].(type) {
	case string:
		return t
	case []interface{}:
		for _, item := range t {
			if str, ok := item.(string); ok && str != "null" {
				return str
			}
		}
	}
	return ""
}