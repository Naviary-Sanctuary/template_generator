package template

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"text/template"

	"github.com/Naviary-Sanctuary/template_generator/internal/config"
)

type Processor struct {
	template  *config.Template
	variables map[string]any
}

type ProcessResult struct {
	FilesCreated int
	DirsCreated  int
	CreatedFiles []string
}

func NewProcessor(template *config.Template, variables map[string]any) *Processor {
	return &Processor{
		template:  template,
		variables: variables,
	}
}
func (processor *Processor) Process(templateDir, outputDir string) (*ProcessResult, error) {
	result := &ProcessResult{
		CreatedFiles: make([]string, 0),
	}

	err := filepath.WalkDir(templateDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.Name() == config.TemplateConfigFile {
			return nil
		}

		relativePath, err := filepath.Rel(templateDir, path)
		if err != nil {
			return err
		}

		outputPath, err := processor.processString(filepath.Join(outputDir, relativePath))
		if err != nil {
			return fmt.Errorf("failed to process output path %s: %w", d.Name(), err)
		}

		if d.IsDir() {
			if err := os.MkdirAll(outputPath, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", outputPath, err)
			}
			result.DirsCreated++
			return nil
		}

		if err := processor.processFile(path, outputPath); err != nil {
			return fmt.Errorf("failed to process file %s: %w", relativePath, err)
		}

		result.FilesCreated++
		result.CreatedFiles = append(result.CreatedFiles, relativePath)

		return nil
	})

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (processor *Processor) processFile(path, outputPath string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", path, err)
	}

	processed, err := processor.processString(string(content))
	if err != nil {
		return fmt.Errorf("failed to process template %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("failed to create parent directory: %w", err)
	}

	if err := os.WriteFile(outputPath, []byte(processed), 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", outputPath, err)
	}

	return nil
}

func (processor *Processor) processString(content string) (string, error) {
	tmpl, err := template.New("template").Parse(content)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buffer bytes.Buffer
	if err := tmpl.Execute(&buffer, processor.variables); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buffer.String(), nil
}
