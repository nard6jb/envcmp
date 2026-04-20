// Package schema provides validation of .env files against a schema
// definition that specifies required keys, their types, and constraints.
package schema

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// FieldType represents the expected type of an env value.
type FieldType string

const (
	TypeString  FieldType = "string"
	TypeInt     FieldType = "int"
	TypeBool    FieldType = "bool"
	TypeURL     FieldType = "url"
)

// Field defines constraints for a single env key.
type Field struct {
	Key      string
	Type     FieldType
	Required bool
	Pattern  string // optional regex pattern
}

// Schema holds a collection of field definitions.
type Schema struct {
	Fields []Field
}

// Issue represents a single schema violation.
type Issue struct {
	Key     string
	Message string
}

// String returns a human-readable representation of the issue.
func (i Issue) String() string {
	return fmt.Sprintf("%s: %s", i.Key, i.Message)
}

var urlPattern = regexp.MustCompile(`^https?://[^\s]+$`)

// Validate checks an env map against the schema and returns any issues found.
func (s *Schema) Validate(env map[string]string) []Issue {
	var issues []Issue

	for _, field := range s.Fields {
		val, exists := env[field.Key]

		if !exists {
			if field.Required {
				issues = append(issues, Issue{Key: field.Key, Message: "required key is missing"})
			}
			continue
		}

		if issue := validateType(field.Key, val, field.Type); issue != nil {
			issues = append(issues, *issue)
		}

		if field.Pattern != "" {
			re, err := regexp.Compile(field.Pattern)
			if err == nil && !re.MatchString(val) {
				issues = append(issues, Issue{
					Key:     field.Key,
					Message: fmt.Sprintf("value does not match pattern %q", field.Pattern),
				})
			}
		}
	}

	return issues
}

func validateType(key, val string, t FieldType) *Issue {
	switch t {
	case TypeInt:
		if _, err := strconv.Atoi(val); err != nil {
			return &Issue{Key: key, Message: fmt.Sprintf("expected int, got %q", val)}
		}
	case TypeBool:
		lower := strings.ToLower(val)
		if lower != "true" && lower != "false" && lower != "1" && lower != "0" {
			return &Issue{Key: key, Message: fmt.Sprintf("expected bool, got %q", val)}
		}
	case TypeURL:
		if !urlPattern.MatchString(val) {
			return &Issue{Key: key, Message: fmt.Sprintf("expected URL, got %q", val)}
		}
	}
	return nil
}
