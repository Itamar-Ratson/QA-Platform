package terraform

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// GenerateTfvarsFile creates a .tfvars file from YAML terraform section with test tags
func GenerateTfvarsFile(tfvars map[string]interface{}, testName, workspace string, outputPath string) error {
	var content strings.Builder
	
	// Add original tfvars
	for key, value := range tfvars {
		line := formatTfvar(key, value)
		content.WriteString(line + "\n")
	}
	
	// Add test identification tags
	testTags := generateTestTags(testName, workspace)
	content.WriteString(formatTfvar("common_tags", testTags) + "\n")
	
	// Ensure directory exists
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}
	
	return os.WriteFile(outputPath, []byte(content.String()), 0644)
}

// generateTestTags creates common tags for test identification
func generateTestTags(testName, workspace string) map[string]interface{} {
	timestamp := time.Now().Format("2006-01-02T15:04:05Z")
	
	return map[string]interface{}{
		"TestCase":        testName,
		"TestWorkspace":   workspace,
		"TestTimestamp":   timestamp,
		"CreatedBy":       "qa-test-app",
		"AutoCleanup":     "true",
		"Environment":     "test",
	}
}

// formatTfvar formats a single terraform variable
func formatTfvar(key string, value interface{}) string {
	switch v := value.(type) {
	case string:
		return fmt.Sprintf(`%s = "%s"`, key, v)
	case []interface{}:
		return fmt.Sprintf(`%s = %s`, key, formatList(v))
	case map[string]interface{}:
		return fmt.Sprintf(`%s = %s`, key, formatMap(v))
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

// formatMap formats a map for terraform
func formatMap(m map[string]interface{}) string {
	var pairs []string
	for k, v := range m {
		switch val := v.(type) {
		case string:
			pairs = append(pairs, fmt.Sprintf(`  "%s" = "%s"`, k, val))
		default:
			pairs = append(pairs, fmt.Sprintf(`  "%s" = "%v"`, k, val))
		}
	}
	return fmt.Sprintf("{\n%s\n}", strings.Join(pairs, "\n"))
}
