package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/Naviary-Sanctuary/template_generator/internal/config"
	"github.com/spf13/cobra"
)

var (
	listDetails bool
	listFormat  string
	listFilter  string
)

func newListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List available templates",
		Long: `List displays all available templates in the configured template directory.

Templates are loaded from the directory specified in tg.config.toml.
Each template must have a template.toml configuration file to be recognized.`,
		Example: `  # List all templates
tg list

# List with detailed information
tg list --details

# List in table format
tg list --format table

# Filter templates by name
tg list --filter "web"`,
		RunE: runList,
	}

	cmd.Flags().BoolVarP(&listDetails, "details", "d", false, "Show detailed template information")
	cmd.Flags().StringVarP(&listFormat, "format", "F", "list", "Output format: list, table, json")
	cmd.Flags().StringVarP(&listFilter, "filter", "f", "", "Filter templates by name (case-insensitive)")

	return cmd
}

func runList(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	templatesDir := cfg.TemplatesDir

	if _, err := os.Stat(templatesDir); os.IsNotExist(err) {
		WarnColor.Printf("Templates directory %s does not exist\n", templatesDir)
		fmt.Println("Run 'tg init' first to initialize the configuration")
		return nil
	}

	templates, err := findTemplates(templatesDir)
	if err != nil {
		return fmt.Errorf("failed to find templates: %w", err)
	}

	if listFilter != "" {
		templates = filterTemplates(templates, listFilter)
	}

	if len(templates) == 0 {
		if listFilter != "" {
			fmt.Printf("No templates found matching filter: %s\n", listFilter)
		} else {
			fmt.Println("No templates found.")
			fmt.Println("")
			fmt.Println("Create a new template with:")
			fmt.Println("	  tg new <template-name>")
		}
		return nil
	}

	switch listFormat {
	case "table":
		return displayTemplatesTable(templates)
	case "json":
		return displayTemplatesJSON(templates)
	default:
		return displayTemplatesList(templates)
	}
}

type TemplateInfo struct {
	Name        string
	Path        string
	Description string
	Author      string
	Version     string
	Variables   int
}

func findTemplates(templatesDir string) ([]TemplateInfo, error) {
	var templates []TemplateInfo

	entries, err := os.ReadDir(templatesDir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		templatePath := filepath.Join(templatesDir, entry.Name())
		configPath := filepath.Join(templatePath, "template.toml")

		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			PrintVerbose("Skipping %s: no template.toml found\n", entry.Name())
			continue
		}

		template, err := config.LoadTemplate(templatePath)
		if err != nil {
			return nil, fmt.Errorf("failed to load template: %w", err)
		}

		info := TemplateInfo{
			Name:        template.Metadata.Name,
			Path:        templatePath,
			Description: template.Metadata.Description,
			Author:      template.Metadata.Author,
			Version:     template.Version,
			Variables:   len(template.Variables),
		}

		templates = append(templates, info)
	}

	return templates, nil
}

func filterTemplates(templates []TemplateInfo, filter string) []TemplateInfo {
	var filtered []TemplateInfo
	lowerFilter := strings.ToLower(filter)

	for _, tmpl := range templates {
		if strings.Contains(strings.ToLower(tmpl.Name), lowerFilter) ||
			strings.Contains(strings.ToLower(tmpl.Description), lowerFilter) {
			filtered = append(filtered, tmpl)
		}
	}

	return filtered
}

func displayTemplatesList(templates []TemplateInfo) error {
	InfoColor.Printf("Found %d template(s):\n", len(templates))
	fmt.Println()

	for _, tmpl := range templates {
		fmt.Printf("  • %s", BoldColor.Sprint(tmpl.Name))
		if tmpl.Version != "" && tmpl.Version != "1.0.0" {
			fmt.Printf(" (v%s)", tmpl.Version)
		}
		fmt.Println()

		if listDetails {
			if tmpl.Description != "" {
				fmt.Printf("    %s\n", tmpl.Description)
			}
			if tmpl.Author != "" && tmpl.Author != "Unknown" {
				fmt.Printf("    Author: %s\n", tmpl.Author)
			}
			fmt.Printf("    Variables: %d\n", tmpl.Variables)
			fmt.Printf("    Path: %s\n", tmpl.Path)
			fmt.Println()
		}
	}

	if !listDetails {
		fmt.Println()
		fmt.Println("Use 'tg list --details' for more information")
	}

	return nil
}

// TODO: use json library
func displayTemplatesJSON(templates []TemplateInfo) error {
	fmt.Println("{")
	fmt.Printf("  \"templates\": [\n")
	for i, tmpl := range templates {
		fmt.Printf("    {\n")
		fmt.Printf("      \"name\": \"%s\",\n", tmpl.Name)
		fmt.Printf("      \"version\": \"%s\",\n", tmpl.Version)
		fmt.Printf("      \"author\": \"%s\",\n", tmpl.Author)
		fmt.Printf("      \"description\": \"%s\",\n", tmpl.Description)
		fmt.Printf("      \"variables\": %d,\n", tmpl.Variables)
		fmt.Printf("      \"path\": \"%s\"\n", tmpl.Path)
		fmt.Printf("    }")
		if i < len(templates)-1 {
			fmt.Printf(",")
		}
		fmt.Println()
	}
	fmt.Printf("  ]\n")
	fmt.Println("}")
	return nil
}

func displayTemplatesTable(templates []TemplateInfo) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)

	fmt.Fprintln(w, "NAME\tVERSION\tAUTHOR\tVARIABLES\tDESCRIPTION")
	fmt.Fprintln(w, "----\t-------\t------\t---------\t-----------")

	for _, tmpl := range templates {
		description := tmpl.Description
		if len(description) > 40 && !listDetails {
			description = description[:37] + "..."
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%s\n",
			tmpl.Name,
			tmpl.Version,
			tmpl.Author,
			tmpl.Variables,
			description,
		)
	}

	return w.Flush()
}
