package main

import (
	"fmt"
	"log"
	"path/filepath"
	"qa-test-app/internal/terraform"
	"qa-test-app/internal/yaml"
)

func main() {
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
	
	// Initialize terraform
	fmt.Println("Initializing Terraform...")
	result, err := executor.Init()
	if err != nil {
		log.Fatal(err)
	}
	if !result.Success {
		log.Fatalf("Terraform init failed: %s", result.Error)
	}
	fmt.Println("✓ Terraform initialized")
	
	// Validate configuration
	fmt.Println("Validating configuration...")
	result, err = executor.Validate()
	if err != nil {
		log.Fatal(err)
	}
	if !result.Success {
		log.Fatalf("Configuration invalid: %s", result.Error)
	}
	fmt.Println("✓ Configuration valid")
	
	// Plan deployment
	fmt.Println("Planning deployment...")
	result, err = executor.Plan()
	if err != nil {
		log.Fatal(err)
	}
	if !result.Success {
		log.Fatalf("Terraform plan failed: %s", result.Error)
	}
	fmt.Println("✓ Plan completed")
	
	log.Println("Ready for development")
	log.Println("Next: Run 'Apply()' to provision infrastructure")
}

// Example functions for TUI integration
func provisionEnvironment(executor *terraform.Executor) error {
	fmt.Println("Provisioning environment...")
	result, err := executor.Apply()
	if err != nil {
		return err
	}
	if !result.Success {
		return fmt.Errorf("apply failed: %s", result.Error)
	}
	fmt.Println("✓ Environment provisioned")
	return nil
}

func destroyEnvironment(executor *terraform.Executor) error {
	fmt.Println("Destroying environment...")
	result, err := executor.Destroy()
	if err != nil {
		return err
	}
	if !result.Success {
		return fmt.Errorf("destroy failed: %s", result.Error)
	}
	fmt.Println("✓ Environment destroyed")
	return nil
}
