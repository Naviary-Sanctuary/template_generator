package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	initForce        bool
	initTemplatesDir string
)

func newInitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize tg configuration in current directory",
		Long: `Initialize creates a tg.config.toml file and template directory structure
in the current directory to start managing templates.

This command will:
  1. Create a tg.config.toml configuration file
	2. Create a templates directory (.tg by default)`,
		Example: `# Initialize with default settings
	tg init
	
	# Initialize with custom template directory
	tg init --template-dir templates
	
	# Force initialization (overwrite existing config)
	tg init --force
	tg init -f`,
		RunE: runInit,
	}

	cmd.Flags().BoolVarP(&initForce, "force", "f", false, "Force initialization (overwrite existing config)")
	cmd.Flags().StringVarP(&initTemplatesDir, "templates-dir", "t", ".tg", "Template directory name")

	return cmd
}

func runInit(cmd *cobra.Command, args []string) error {
	InfoColor.Println("Initializing tg configuration...")

	if !initForce {
		if _, err := os.Stat(configPath); err == nil {
			return fmt.Errorf("config file '%s' already exists. Use --force to overwrite.", configPath)
		}
	}

	if err := createConfigFile(); err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	PrintVerbose("Created config file: %s\n", configPath)

	if err := os.MkdirAll(initTemplatesDir, 0755); err != nil {
		return fmt.Errorf("failed to create templates directory: %w", err)
	}
	PrintVerbose("Created templates directory: %s\n", initTemplatesDir)

	SuccessColor.Println("✓ Configuration initialized successfully!")
	fmt.Printf("  Config file:    %s\n", BoldColor.Sprint(configPath))
	fmt.Printf("  Templates dir:  %s\n", BoldColor.Sprint(initTemplatesDir))

	return nil
}

func createConfigFile() error {
	content := fmt.Sprintf(`# Template Generator Configuration

	# Directory containing templates
	templates_dir = "%s"
	
	# Git remote for fetching templates (optional)
	# git_remote = "https://github.com/yourusername/tg-templates.git"
	
	# Default variables for all templates (optional)
	# [defaults]
	# author = "Your Name"
	# license = "MIT"
	`, initTemplatesDir)

	return os.WriteFile(configPath, []byte(content), 0644)
}
