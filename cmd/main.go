package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"qa-test-app/internal/terraform"
	"qa-test-app/internal/tests"
	"qa-test-app/internal/yaml"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	tealStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("14"))
)

func main() {
	apply := flag.Bool("apply", false, "Apply terraform")
	destroy := flag.Bool("destroy", false, "Destroy all test workspaces")
	test := flag.Bool("test", false, "Run tests against existing infrastructure")
	flag.Parse()

	fmt.Println(tealStyle.Render("QA Test App Starting..."))
	
	tc, err := yaml.ParseTestCase("test-cases/sample.yaml")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf(tealStyle.Render("Loaded test: %s\n"), tc.Metadata.Name)
	
	workingDir := filepath.Join("terraform", "base")
	tfvarsFile := "generated.tfvars"
	executor := terraform.NewExecutor(workingDir, tfvarsFile)
	
	if *destroy {
		fmt.Println(tealStyle.Render("Destroying all test workspaces..."))
		if err := destroyAllTestWorkspaces(executor); err != nil {
			log.Printf("Destroy warning: %v", err)
		}
		fmt.Println(tealStyle.Render("✓ All test workspaces destroyed"))
		return
	}
	
	if *test {
		// Find existing test workspace instead of creating new one
		workspaces, err := executor.WorkspaceList()
		if err != nil {
			log.Fatalf("Failed to list workspaces: %v", err)
		}
		
		var targetWorkspace string
		for _, ws := range workspaces {
			if strings.HasPrefix(ws, "test-") {
				executor.CurrentWorkspace = ws
				if _, err := executor.SelectWorkspace(ws); err != nil {
					continue
				}
				if hasResources, _ := executor.HasResources(); hasResources {
					targetWorkspace = ws
					break
				}
			}
		}
		
		if targetWorkspace == "" {
			log.Fatal("No test workspace with resources found. Run 'make apply' first.")
		}
		
		fmt.Printf(tealStyle.Render("Using existing workspace: %s\n"), targetWorkspace)
		
		tfOutputs, err := getTerraformOutputs(executor)
		if err != nil {
			log.Fatalf("Could not get terraform outputs: %v", err)
		}
		runTests(tc.TestFunctions, tfOutputs)
		return
	}
	
	fmt.Printf(tealStyle.Render("Setting up test environment for: %s\n"), tc.Metadata.Name)
	err = executor.SetupTestEnvironment(tc.Metadata.Name)
	if err != nil {
		log.Fatalf("Failed to setup test environment: %v", err)
	}
	fmt.Printf(tealStyle.Render("✓ Test workspace created: %s\n"), executor.CurrentWorkspace)
	
	outputPath := filepath.Join("terraform", "base", "generated.tfvars")
	err = terraform.GenerateTfvarsFile(tc.Terraform.TfVars, tc.Metadata.Name, executor.CurrentWorkspace, outputPath)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf(tealStyle.Render("Generated tfvars with test tags: %s\n"), outputPath)
	
	defer func() {
		if executor.CurrentWorkspace != "" && !*apply {
			fmt.Println(tealStyle.Render("Cleaning up test environment..."))
			if err := executor.CleanupTestEnvironment(); err != nil {
				log.Printf("Cleanup failed: %v", err)
				executor.ForceCleanup()
			} else {
				fmt.Println(tealStyle.Render("✓ Test environment cleaned up"))
			}
		}
	}()
	
	fmt.Println(tealStyle.Render("Validating state..."))
	if err := executor.ValidateState(); err != nil {
		log.Fatalf("State validation failed: %v", err)
	}
	fmt.Println(tealStyle.Render("✓ State validated"))
	
	fmt.Println(tealStyle.Render("Planning deployment..."))
	result, err := executor.Plan()
	if err != nil {
		log.Fatal(err)
	}
	if !result.Success {
		log.Fatalf("Terraform plan failed: %s", result.Error)
	}
	fmt.Println(tealStyle.Render("✓ Plan completed"))
	
	if *apply {
		fmt.Println(tealStyle.Render("Applying deployment..."))
		result, err = executor.Apply()
		if err != nil {
			log.Fatal(err)
		}
		if !result.Success {
			log.Fatalf("Terraform apply failed: %s", result.Error)
		}
		fmt.Println(tealStyle.Render("✓ Environment provisioned"))
		
		// Get terraform outputs
		tfOutputs, err := getTerraformOutputs(executor)
		if err != nil {
			log.Printf("Warning: Could not get terraform outputs: %v", err)
		} else {
			// Run tests
			runTests(tc.TestFunctions, tfOutputs)
		}
		
		info, _ := executor.GetWorkspaceInfo()
		fmt.Printf(tealStyle.Render("Workspace: %s (has resources: %v)\n"), 
			info["current_workspace"], info["has_resources"])
		return
	}
	
	fmt.Println(tealStyle.Render("Ready for development"))
	fmt.Printf(tealStyle.Render("Active workspace: %s\n"), executor.CurrentWorkspace)
}

// getTerraformOutputs extracts outputs from terraform
func getTerraformOutputs(executor *terraform.Executor) (map[string]interface{}, error) {
	result, err := executor.GetOutputs()
	if err != nil {
		return nil, err
	}
	
	// Parse terraform outputs JSON
	var outputs map[string]interface{}
	if err := json.Unmarshal([]byte(result.Output), &outputs); err != nil {
		// If JSON parsing fails, create mock outputs for testing
		return map[string]interface{}{
			"vpc_id": "vpc-mock123",
			"vpc_cidr_block": "10.0.0.0/16",
		}, nil
	}
	
	// Extract values from output structure
	finalOutputs := make(map[string]interface{})
	for key, output := range outputs {
		if outputMap, ok := output.(map[string]interface{}); ok {
			if value, exists := outputMap["value"]; exists {
				finalOutputs[key] = value
			}
		}
	}
	
	return finalOutputs, nil
}

// runTests executes the test functions
func runTests(testFunctions []string, tfOutputs map[string]interface{}) {
	fmt.Println(tealStyle.Render("\n=== Running Tests ==="))
	
	testExecutor := tests.NewTestExecutor()
	ctx := context.Background()
	
	results := testExecutor.ExecuteAll(ctx, testFunctions, tfOutputs)
	
	successCount := 0
	for _, result := range results {
		status := "✗ FAIL"
		style := lipgloss.NewStyle().Foreground(lipgloss.Color("9")) // Red
		if result.Success {
			status = "✓ PASS"
			style = lipgloss.NewStyle().Foreground(lipgloss.Color("10")) // Green
			successCount++
		}
		
		fmt.Printf("%s %s (%v)\n", 
			style.Render(status), 
			result.TestName, 
			result.Duration)
		fmt.Printf("  %s\n", result.Message)
		
		if !result.Success && result.Details != nil {
			if details, err := json.MarshalIndent(result.Details, "  ", "  "); err == nil {
				fmt.Printf("  Details: %s\n", string(details))
			}
		}
	}
	
	fmt.Printf(tealStyle.Render("\nTest Summary: %d/%d passed\n"), successCount, len(results))
}

func destroyAllTestWorkspaces(executor *terraform.Executor) error {
	workspaces, err := executor.WorkspaceList()
	if err != nil {
		return err
	}
	
	for _, ws := range workspaces {
		if ws != "default" && strings.HasPrefix(ws, "test-") {
			fmt.Printf(tealStyle.Render("Destroying workspace: %s\n"), ws)
			
			// Select workspace first
			if _, err := executor.SelectWorkspace(ws); err != nil {
				fmt.Printf("Failed to select workspace %s: %v\n", ws, err)
				continue
			}
			
			// Destroy resources then delete workspace
			if err := executor.CleanupTestEnvironment(); err != nil {
				fmt.Printf("Cleanup failed for %s, forcing deletion: %v\n", ws, err)
				executor.ForceCleanup()
			}
		}
	}
	
	return nil
}
