package terraform

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Executor struct {
	WorkingDir string
	TfvarsFile string
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
	
	// Set environment
	cmd.Env = os.Environ()
	
	// Capture output
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
	
	// Basic state info (extend as needed)
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
			// Remove * prefix for current workspace
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
	// Try to select existing workspace
	result, err := e.runCommand("workspace", "select", name)
	if err != nil {
		return result, err
	}
	
	if result.Success {
		return result, nil
	}
	
	// If selection failed, try to create new workspace
	return e.runCommand("workspace", "new", name)
}
