package output

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"gopkg.in/yaml.v3"
)

type OutputFormat string

const (
	OutputFormatJSON  OutputFormat = "json"
	OutputFormatYAML  OutputFormat = "yaml"
	OutputFormatTable OutputFormat = "table"
)

type Formatter struct {
	format OutputFormat
}

func NewFormatter(format OutputFormat) *Formatter {
	return &Formatter{format: format}
}

func (f *Formatter) Format(data interface{}) error {
	switch f.format {
	case OutputFormatJSON:
		return f.formatJSON(data)
	case OutputFormatYAML:
		return f.formatYAML(data)
	case OutputFormatTable:
		return f.formatTable(data)
	default:
		return fmt.Errorf("unsupported output format: %s", f.format)
	}
}

func (f *Formatter) formatJSON(data interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

func (f *Formatter) formatYAML(data interface{}) error {
	encoder := yaml.NewEncoder(os.Stdout)
	return encoder.Encode(data)
}

func (f *Formatter) formatTable(data interface{}) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()

	// Handle different types of data
	switch v := data.(type) {
	case map[string]interface{}:
		for k, val := range v {
			fmt.Fprintf(w, "%s\t%v\n", k, val)
		}
	case []interface{}:
		// Handle slice data
		// Implementation depends on the structure of your data
	default:
		return fmt.Errorf("unsupported data type for table format")
	}

	return nil
}
