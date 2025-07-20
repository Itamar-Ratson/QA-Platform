package terraform

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Executor struct {
	WorkingDir      string
	TfvarsFile      string
	CurrentWorkspace string
	TestName        string
}

type ExecutionResult struct {
	Success bool
	Output  string
	Error   string
}

// NewExecutor creates a new terraform executor
func NewExecutor(workingDir, tfvarsFile string) *Executor {
	return &Executor{
		WorkingDir: workingDir,
		TfvarsFile: tfvarsFile,
	}
}

// SetupTestEnvironment creates isolated workspace for test
func (e *Executor) SetupTestEnvironment(testName string) error {
	e.TestName = testName
	
	// Create unique workspace name
	timestamp := time.Now().Unix()
	workspaceName := fmt.Sprintf("test-%s-%d", 
		strings.ToLower(strings.ReplaceAll(testName, " ", "-")), 
		timestamp)
	
	// Initialize if needed
	if _, err := e.Init(); err != nil {
		return fmt.Errorf("init failed: %w", err)
	}
	
	// Create and select workspace
	result, err := e.runCommand("workspace", "new", workspaceName)
	if err != nil {
		return fmt.Errorf("workspace creation failed: %w", err)
	}
	
	if !result.Success {
		return fmt.Errorf("workspace creation failed: %s", result.Error)
	}
	
	e.CurrentWorkspace = workspaceName
	return nil
}

// CleanupTestEnvironment destroys resources and removes workspace
func (e *Executor) CleanupTestEnvironment() error {
	if e.CurrentWorkspace == "" {
		return fmt.Errorf("no active workspace to cleanup")
	}
	
	// Destroy resources first
	destroyResult, err := e.Destroy()
	if err != nil || !destroyResult.Success {
		return fmt.Errorf("destroy failed: %v, %s", err, destroyResult.Error)
	}
	
	// Switch to default workspace
	if _, err := e.runCommand("workspace", "select", "default"); err != nil {
		return fmt.Errorf("failed to switch to default workspace: %w", err)
	}
	
	// Delete test workspace
	if _, err := e.runCommand("workspace", "delete", e.CurrentWorkspace); err != nil {
		return fmt.Errorf("failed to delete workspace: %w", err)
	}
	
	e.CurrentWorkspace = ""
	return nil
}

// ValidateState checks if state is consistent
func (e *Executor) ValidateState() error {
	result, err := e.runCommand("validate")
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}
	
	if !result.Success {
		return fmt.Errorf("state validation failed: %s", result.Error)
	}
	
	return nil
}

// HasResources checks if workspace has any resources
func (e *Executor) HasResources() (bool, error) {
	result, err := e.runCommand("state", "list")
	if err != nil {
		return false, err
	}
	
	return strings.TrimSpace(result.Output) != "", nil
}

// ForceCleanup removes workspace even with resources (emergency cleanup)
func (e *Executor) ForceCleanup() error {
	if e.CurrentWorkspace == "" {
		return nil
	}
	
	// Try normal cleanup first
	if err := e.CleanupTestEnvironment(); err == nil {
		return nil
	}
	
	// Force cleanup if normal cleanup fails
	e.runCommand("workspace", "select", "default")
	e.runCommand("workspace", "delete", "-force", e.CurrentWorkspace)
	e.CurrentWorkspace = ""
	
	return nil
}

// GetWorkspaceInfo returns current workspace details
func (e *Executor) GetWorkspaceInfo() (map[string]interface{}, error) {
	info := map[string]interface{}{
		"current_workspace": e.CurrentWorkspace,
		"test_name":        e.TestName,
		"working_dir":      e.WorkingDir,
	}
	
	// Check if workspace has resources
	hasResources, err := e.HasResources()
	if err != nil {
		info["has_resources"] = "unknown"
		info["error"] = err.Error()
	} else {
		info["has_resources"] = hasResources
	}
	
	return info, nil
}

// Init runs terraform init
func (e *Executor) Init() (*ExecutionResult, error) {
	return e.runCommand("init")
}

// Plan runs terraform plan
func (e *Executor) Plan() (*ExecutionResult, error) {
	return e.runCommand("plan", "-var-file="+e.TfvarsFile)
}

// Apply runs terraform apply
func (e *Executor) Apply() (*ExecutionResult, error) {
	return e.runCommand("apply", "-auto-approve", "-var-file="+e.TfvarsFile)
}

// Destroy runs terraform destroy
func (e *Executor) Destroy() (*ExecutionResult, error) {
	return e.runCommand("destroy", "-auto-approve", "-var-file="+e.TfvarsFile)
}

// runCommand executes terraform with given arguments
func (e *Executor) runCommand(args ...string) (*ExecutionResult, error) {
	cmd := exec.Command("terraform", args...)
	cmd.Dir = e.WorkingDir
	cmd.Env = os.Environ()
	
	output, err := cmd.CombinedOutput()
	
	result := &ExecutionResult{
		Success: err == nil,
		Output:  string(output),
	}
	
	if err != nil {
		result.Error = err.Error()
	}
	
	return result, nil
}

// GetState returns current terraform state info
func (e *Executor) GetState() (map[string]interface{}, error) {
	result, err := e.runCommand("show", "-json")
	if err != nil {
		return nil, err
	}
	
	if !result.Success {
		return nil, fmt.Errorf("terraform show failed: %s", result.Error)
	}
	
	state := map[string]interface{}{
		"has_resources": strings.Contains(result.Output, `"resources"`),
		"output":       result.Output,
	}
	
	return state, nil
}

// Validate checks if terraform configuration is valid
func (e *Executor) Validate() (*ExecutionResult, error) {
	return e.runCommand("validate")
}

// WorkspaceList returns available workspaces
func (e *Executor) WorkspaceList() ([]string, error) {
	result, err := e.runCommand("workspace", "list")
	if err != nil {
		return nil, err
	}
	
	if !result.Success {
		return nil, fmt.Errorf("workspace list failed: %s", result.Error)
	}
	
	var workspaces []string
	lines := strings.Split(result.Output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.Contains(line, "No workspaces") {
			workspace := strings.TrimPrefix(line, "* ")
			workspace = strings.TrimSpace(workspace)
			if workspace != "" {
				workspaces = append(workspaces, workspace)
			}
		}
	}
	
	return workspaces, nil
}

// SelectWorkspace selects or creates a workspace
func (e *Executor) SelectWorkspace(name string) (*ExecutionResult, error) {
	result, err := e.runCommand("workspace", "select", name)
	if err != nil {
		return result, err
	}
	
	if result.Success {
		e.CurrentWorkspace = name
		return result, nil
	}
	
	return e.runCommand("workspace", "new", name)
}
