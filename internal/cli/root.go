package cli

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	// Version information (set during build with ldflags)
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"

	// Global flags
	verbose    bool
	configPath string

	// Color outputs
	SuccessColor = color.New(color.FgGreen)
	ErrorColor   = color.New(color.FgRed)
	InfoColor    = color.New(color.FgCyan)
	WarnColor    = color.New(color.FgYellow)
	BoldColor    = color.New(color.Bold)

	// Root command
	rootCmd = &cobra.Command{
		Use:   "tg",
		Short: "Template Generator - A CLI tool for project scaffolding",
		Long: `Template Generator (tg) is a CLI tool that helps you quickly scaffold
projects using predefined templates with variable substitution and file filtering.

Templates are defined using TOML configuration files and can include:
  - Variable substitution using Go templates
  - File and directory filtering rules
  - Custom rename patterns
  - Git repository integration (coming soon)`,
		Version: fmt.Sprintf("%s (commit: %s, built: %s)", Version, Commit, Date),
	}
)

// Execute runs the root command
func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		ErrorColor.Fprintf(os.Stderr, "Error: %v\n", err)
		return err
	}
	return nil
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "tg.config.toml", "Path to config file")

	// Add commands
	rootCmd.AddCommand(
		newInitCommand(),
	// newListCommand(),
	// newApplyCommand(),
	// newNewCommand(),
	// newFetchCommand(), // for git integration
	// newValidateCommand(), // for template validation
	)

	// Custom version template
	rootCmd.SetVersionTemplate(`{{with .Name}}{{printf "%s " .}}{{end}}{{printf "version %s" .Version}}
`)
}

// IsVerbose returns true if verbose mode is enabled
func IsVerbose() bool {
	return verbose
}

// GetConfigPath returns the config file path
func GetConfigPath() string {
	return configPath
}

// PrintVerbose prints message only in verbose mode
func PrintVerbose(format string, args ...interface{}) {
	if verbose {
		fmt.Printf(format, args...)
	}
}
