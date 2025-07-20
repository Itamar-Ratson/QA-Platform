// internal/tests/interface.go
package tests

import (
	"context"
	"time"
)

// TestFunction defines the interface for all test functions
type TestFunction interface {
	Name() string
	Execute(ctx context.Context, tfOutputs map[string]interface{}) TestResult
	Description() string
}

// TestResult represents the result of a test execution
type TestResult struct {
	Success   bool                   `json:"success"`
	Message   string                 `json:"message"`
	Details   map[string]interface{} `json:"details"`
	Duration  time.Duration          `json:"duration"`
	TestName  string                 `json:"test_name"`
	Timestamp time.Time              `json:"timestamp"`
}

// TestExecutor manages and runs test functions
type TestExecutor struct {
	functions map[string]TestFunction
}

// NewTestExecutor creates a new test executor with default functions
func NewTestExecutor() *TestExecutor {
	executor := &TestExecutor{
		functions: make(map[string]TestFunction),
	}
	
	// Register default test functions
	executor.Register(&CIDRValidationTest{})
	executor.Register(&SubnetConnectivityTest{})
	executor.Register(&RouteTableTest{})
	
	return executor
}

// Register adds a test function to the executor
func (te *TestExecutor) Register(fn TestFunction) {
	te.functions[fn.Name()] = fn
}

// Execute runs a test function by name
func (te *TestExecutor) Execute(ctx context.Context, testName string, tfOutputs map[string]interface{}) (TestResult, error) {
	fn, exists := te.functions[testName]
	if !exists {
		return TestResult{
			Success:   false,
			Message:   "Test function not found",
			TestName:  testName,
			Timestamp: time.Now(),
		}, nil
	}
	
	start := time.Now()
	result := fn.Execute(ctx, tfOutputs)
	result.Duration = time.Since(start)
	result.TestName = testName
	result.Timestamp = time.Now()
	
	return result, nil
}

// ExecuteAll runs all specified test functions
func (te *TestExecutor) ExecuteAll(ctx context.Context, testNames []string, tfOutputs map[string]interface{}) []TestResult {
	results := make([]TestResult, 0, len(testNames))
	
	for _, testName := range testNames {
		result, _ := te.Execute(ctx, testName, tfOutputs)
		results = append(results, result)
	}
	
	return results
}

// ListAvailable returns names of all registered test functions
func (te *TestExecutor) ListAvailable() []string {
	names := make([]string, 0, len(te.functions))
	for name := range te.functions {
		names = append(names, name)
	}
	return names
}
