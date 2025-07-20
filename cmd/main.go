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
	
	log.Println("Ready for development")
}
