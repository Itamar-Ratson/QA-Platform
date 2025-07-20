package terraform

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// GenerateTfvarsFile creates a .tfvars file from YAML terraform section
func GenerateTfvarsFile(tfvars map[string]interface{}, outputPath string) error {
	var content strings.Builder
	
	for key, value := range tfvars {
		line := formatTfvar(key, value)
		content.WriteString(line + "\n")
	}
	
	// Ensure directory exists
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}
	
	return os.WriteFile(outputPath, []byte(content.String()), 0644)
}

// formatTfvar formats a single terraform variable
func formatTfvar(key string, value interface{}) string {
	switch v := value.(type) {
	case string:
		return fmt.Sprintf(`%s = "%s"`, key, v)
	case []interface{}:
		return fmt.Sprintf(`%s = %s`, key, formatList(v))
	case bool:
		return fmt.Sprintf(`%s = %t`, key, v)
	case int, int64, float64:
		return fmt.Sprintf(`%s = %v`, key, v)
	default:
		return fmt.Sprintf(`%s = "%v"`, key, v)
	}
}

// formatList formats a list for terraform
func formatList(items []interface{}) string {
	var elements []string
	for _, item := range items {
		switch v := item.(type) {
		case string:
			elements = append(elements, fmt.Sprintf(`"%s"`, v))
		default:
			elements = append(elements, fmt.Sprintf(`%v`, v))
		}
	}
	return fmt.Sprintf(`[%s]`, strings.Join(elements, ", "))
}
