package template

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

type TemplateData struct {
	Environment string
	Variables   map[string]interface{}
	Provider    map[string]interface{}
	Tags        map[string]string
}

func RenderTemplate(templatePath string, data *TemplateData) (string, error) {
	content, err := os.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read template file: %w", err)
	}

	tmpl, err := template.New(filepath.Base(templatePath)).Parse(string(content))
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to render template: %w", err)
	}

	return buf.String(), nil
}
