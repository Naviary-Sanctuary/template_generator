package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Naviary-Sanctuary/template_generator/internal/config"
	"github.com/Naviary-Sanctuary/template_generator/internal/template"
	"github.com/spf13/cobra"
)

var (
	applyOutputPath string
	applyVariables  map[string]string
)

func newApplyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "apply <template-name> [output-dir]",
		Short: "Apply a template to generate files",
		Long: `Apply reads a template and generates files by substituting variables.

Variables use their default values defined in template.toml.
The output directory defaults to the current directory if not specified.`,
		Example: `  # Apply template to current directory
  tg apply hello-world

  # Apply template to specific directory
  tg apply hello-world ./my-project`,
		Args: cobra.MinimumNArgs(1),
		RunE: runApply,
	}

	cmd.Flags().StringVarP(&applyOutputPath, "output", "o", ".", "Output directory")
	cmd.Flags().StringToStringVarP(&applyVariables, "var", "v", nil, "Set variable values (e.g. -v name=John -v age=30)")

	return cmd
}

func runApply(cmd *cobra.Command, args []string) error {
	templateName := args[0]

	if len(args) > 1 {
		applyOutputPath = args[1]
	}

	InfoColor.Printf("Applying template: %s\n", BoldColor.Sprint(templateName))

	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	templateDir, tmpl, err := resolveTemplateDir(cfg, templateName)
	if err != nil {
		return err
	}

	PrintVerbose("Template loaded: %s\n", tmpl.Metadata.Name)
	PrintVerbose("Description: %s\n", tmpl.Metadata.Description)

	variables := make(map[string]any)
	for name, variable := range tmpl.Variables {
		variables[name] = variable.Default
	}

	for key, value := range applyVariables {
		variables[key] = value
	}

	for name, value := range variables {
		PrintVerbose("Variable %s = %v\n", name, value)
	}

	if err := os.MkdirAll(applyOutputPath, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	processor := template.NewProcessor(tmpl, variables)
	result, err := processor.Process(templateDir, applyOutputPath)
	if err != nil {
		return fmt.Errorf("failed to process template: %w", err)
	}

	SuccessColor.Println("✓ Template applied successfully!")
	PrintVerbose("  Output directory: %s\n", BoldColor.Sprint(applyOutputPath))
	PrintVerbose("  Processed files: %d\n", result.FilesCreated)
	PrintVerbose("  Created directories: %d\n", result.DirsCreated)
	PrintVerbose("  Created files: %d\n", len(result.CreatedFiles))

	for _, file := range result.CreatedFiles {
		PrintVerbose("    %s\n", file)
	}

	return nil
}

func resolveTemplateDir(cfg *config.Config, requestedName string) (string, *config.Template, error) {
	candidateDir := filepath.Join(cfg.TemplatesDir, requestedName)
	if info, err := os.Stat(candidateDir); err == nil && info.IsDir() {
		if _, err := os.Stat(filepath.Join(candidateDir, config.TemplateConfigFile)); err == nil {
			tmpl, err := config.LoadTemplate(candidateDir)
			if err != nil {
				return "", nil, fmt.Errorf("failed to load template: %w", err)
			}
			return candidateDir, tmpl, nil
		}
	}

	entries, err := os.ReadDir(cfg.TemplatesDir)
	if err != nil {
		return "", nil, fmt.Errorf("failed to read templates directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		dir := filepath.Join(cfg.TemplatesDir, entry.Name())
		if _, err := os.Stat(filepath.Join(dir, config.TemplateConfigFile)); err != nil {
			continue
		}
		tmpl, err := config.LoadTemplate(dir)
		if err != nil {
			return "", nil, fmt.Errorf("failed to load template: %w", err)
		}
		if tmpl.Metadata.Name == requestedName {
			return dir, tmpl, nil
		}
	}

	return "", nil, fmt.Errorf("template '%s' not found in '%s'", requestedName, cfg.TemplatesDir)
}
