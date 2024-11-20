package terraform

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/yahao333/gort/internal/provider"
)

type TerraformProvider struct {
	workDir     string
	binPath     string
	environment string
}

func NewTerraformProvider(workDir string) *TerraformProvider {
	return &TerraformProvider{
		workDir: workDir,
		binPath: "terraform", // Assume terraform is in PATH
	}
}

func (p *TerraformProvider) Initialize() error {
	cmd := exec.Command(p.binPath, "init")
	cmd.Dir = p.workDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (p *TerraformProvider) Plan(env string) (*provider.PlanResult, error) {
	cmd := exec.Command(p.binPath, "plan", "-out=plan.tfplan")
	cmd.Dir = p.workDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("terraform plan failed: %w", err)
	}

	// TODO: Parse plan output to create PlanResult
	return &provider.PlanResult{}, nil
}

func (p *TerraformProvider) Apply(plan *provider.PlanResult) error {
	cmd := exec.Command(p.binPath, "apply", "-auto-approve", "plan.tfplan")
	cmd.Dir = p.workDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
