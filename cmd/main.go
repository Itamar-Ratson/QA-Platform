package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"qa-test-app/internal/terraform"
	"qa-test-app/internal/yaml"
)

func main() {
	// Parse flags
	apply := flag.Bool("apply", false, "Apply terraform")
	destroy := flag.Bool("destroy", false, "Destroy all test workspaces")
	flag.Parse()

	fmt.Println("QA Test App Starting...")
	
	// Parse YAML test case
	tc, err := yaml.ParseTestCase("test-cases/sample.yaml")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Loaded test: %s\n", tc.Metadata.Name)
	
	// Generate tfvars file
	outputPath := filepath.Join("terraform", "base", "generated.tfvars")
	err = terraform.GenerateTfvarsFile(tc.Terraform.TfVars, outputPath)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Generated tfvars: %s\n", outputPath)
	
	// Create terraform executor
	workingDir := filepath.Join("terraform", "base")
	tfvarsFile := "generated.tfvars"
	executor := terraform.NewExecutor(workingDir, tfvarsFile)
	
	// Handle destroy flag
	if *destroy {
		fmt.Println("Destroying all test workspaces...")
		if err := destroyAllTestWorkspaces(executor); err != nil {
			log.Printf("Destroy warning: %v", err)
		}
		fmt.Println("✓ All test workspaces destroyed")
		return
	}
	
	// Setup isolated test environment for apply/plan
	fmt.Printf("Setting up test environment for: %s\n", tc.Metadata.Name)
	err = executor.SetupTestEnvironment(tc.Metadata.Name)
	if err != nil {
		log.Fatalf("Failed to setup test environment: %v", err)
	}
	fmt.Printf("✓ Test workspace created: %s\n", executor.CurrentWorkspace)
	
	// Ensure cleanup on exit
	defer func() {
		if executor.CurrentWorkspace != "" && !*apply {
			fmt.Println("Cleaning up test environment...")
			if err := executor.CleanupTestEnvironment(); err != nil {
				log.Printf("Cleanup failed: %v", err)
				executor.ForceCleanup()
			} else {
				fmt.Println("✓ Test environment cleaned up")
			}
		}
	}()
	
	// Validate state
	fmt.Println("Validating state...")
	if err := executor.ValidateState(); err != nil {
		log.Fatalf("State validation failed: %v", err)
	}
	fmt.Println("✓ State validated")
	
	// Plan deployment
	fmt.Println("Planning deployment...")
	result, err := executor.Plan()
	if err != nil {
		log.Fatal(err)
	}
	if !result.Success {
		log.Fatalf("Terraform plan failed: %s", result.Error)
	}
	fmt.Println("✓ Plan completed")
	
	// Handle apply flag
	if *apply {
		fmt.Println("Applying deployment...")
		result, err = executor.Apply()
		if err != nil {
			log.Fatal(err)
		}
		if !result.Success {
			log.Fatalf("Terraform apply failed: %s", result.Error)
		}
		fmt.Println("✓ Environment provisioned")
		
		// Show workspace info
		info, _ := executor.GetWorkspaceInfo()
		fmt.Printf("Workspace: %s (has resources: %v)\n", 
			info["current_workspace"], info["has_resources"])
		return
	}
	
	fmt.Println("Ready for development")
	fmt.Printf("Active workspace: %s\n", executor.CurrentWorkspace)
}

// destroyAllTestWorkspaces removes all test-* workspaces
func destroyAllTestWorkspaces(executor *terraform.Executor) error {
	workspaces, err := executor.WorkspaceList()
	if err != nil {
		return err
	}
	
	for _, ws := range workspaces {
		if ws != "default" && (ws[:5] == "test-" || ws[:4] == "test") {
			fmt.Printf("Destroying workspace: %s\n", ws)
			executor.CurrentWorkspace = ws
			executor.ForceCleanup()
		}
	}
	
	return nil
}
