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
	destroy := flag.Bool("destroy", false, "Destroy terraform")
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
		return
	}
	
	// Handle destroy flag
	if *destroy {
		fmt.Println("Destroying deployment...")
		result, err = executor.Destroy()
		if err != nil {
			log.Fatal(err)
		}
		if !result.Success {
			log.Fatalf("Terraform destroy failed: %s", result.Error)
		}
		fmt.Println("✓ Environment destroyed")
		return
	}
	
	log.Println("Ready for development")
}
