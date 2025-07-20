package main

import (
    "fmt"
    "log"
    "qa-test-app/internal/yaml"
)

func main() {
    fmt.Println("QA Test App Starting...")
    
    // Test YAML parsing
    tc, err := yaml.ParseTestCase("test-cases/sample.yaml")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Loaded test: %s\n", tc.Metadata.Name)
    
    log.Println("Ready for development")
}
