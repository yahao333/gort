package terraform

import (
	"fmt"
	"os/exec"
	"strings"
)

type Workspace struct {
    name    string
    workDir string
}

func (p *TerraformProvider) EnsureWorkspace(name string) error {
    // List existing workspaces
    cmd := exec.Command(p.binPath, "workspace", "list")
    cmd.Dir = p.workDir
    output, err := cmd.Output()
    if err != nil {
        return fmt.Errorf("failed to list workspaces: %w", err)
    }

    // Check if workspace exists
    workspaces := strings.Split(string(output), "\n")
    exists := false
    for _, ws := range workspaces {
        if strings.TrimSpace(strings.TrimPrefix(ws, "*")) == name {
            exists = true
            break
        }
    }

    // Create workspace if it doesn't exist
    if !exists {
        cmd = exec.Command(p.binPath, "workspace", "new", name)
        cmd.Dir = p.workDir
        if err := cmd.Run(); err != nil {
            return fmt.Errorf("failed to create workspace: %w", err)
        }
    }

    // Select workspace
    cmd = exec.Command(p.binPath, "workspace", "select", name)
    cmd.Dir = p.workDir
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("failed to select workspace: %w", err)
    }

    return nil
}