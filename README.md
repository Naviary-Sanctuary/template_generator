# Template Generator (tg)

A powerful and flexible CLI tool for project scaffolding using template-based file generation with variable substitution.

[![Go Version](https://img.shields.io/badge/Go-1.25%2B-00ADD8?style=flat&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Commands](#commands)
- [Configuration](#configuration)
- [Variable Types](#variable-types)
- [Template Syntax](#template-syntax)
- [Project Structure](#project-structure)
- [Dependencies](#dependencies)
- [Roadmap](#roadmap)
- [Contributing](#contributing)
- [License](#license)

## Features

- **Quick Project Scaffolding**: Generate project structures from predefined templates in seconds
- **Variable Substitution**: Use Go templates for dynamic content generation
- **TOML Configuration**: Simple and readable template definitions
- **Colored Output**: Enhanced CLI experience with color-coded messages
- **Multiple Output Formats**: List templates in list, table, or JSON format
- **Flexible Filtering**: Filter and search templates by name or description
- **Type-Safe Variables**: Support for string, number, boolean, and array types
- **Smart File Handling**: Automatic directory creation and file processing

## Installation

### Using Homebrew

```bash
brew tap naviary-sanctuary/template_generator
```

### From Source

```bash
git clone https://github.com/Naviary-Sanctuary/template_generator.git
cd template_generator
go build -o tg cmd/main.go
```

### Using Go Install

```bash
go install github.com/Naviary-Sanctuary/template_generator/cmd@latest
```

## Quick Start

### 1. Initialize Configuration

```bash
# Initialize with default settings
tg init

# Initialize with custom template directory
tg init --templates-dir my-templates

# Force initialization (overwrite existing config)
tg init --force
```

This creates:

- `tg.config.toml`: Main configuration file
- `.tg/`: Default templates directory

### 2. Create a Template

Create a template directory structure:

```
.tg/
└── my-template/
    ├── template.toml
    ├── README.md
    └── src/
        └── main.go
```

Define your template in `template.toml`:

```toml
version = "1.0.0"

[metadata]
name = "my-template"
description = "A sample Go project template"
author = "Your Name"

[variables]
project_name={ default="my-project", description="Name of the project" }
author={ default="Naviary", description="Project author" }
port={ default=8080, type="number", description="Server port" }

[rules]
ignores = ["*.tmp", ".git"]
includes = ["**/*.go", "**/*.md"]
```

### 3. Use Template Variables

In your template files, use Go template syntax:

**README.md**:

```markdown
# {{.project_name}}

Author: {{.author}}

Server runs on port {{.port}}`
```

**src/main.go**:

```go
package main

import "fmt"

func main() {
    fmt.Println("Welcome to {{.project_name}}!")
}
```

### 4. Apply Template

```bash
# Apply to current directory
tg apply my-template

# Apply to specific directory
tg apply my-template ./new-project

# Override variables
tg apply my-template -v project_name=awesome-app -v author="John Doe" -v port=3000

# With verbose output
tg apply my-template --verbose
```

## Commands

### `tg init`

Initialize tg configuration in the current directory.

**Usage:**

```bash
tg init [flags]
```

**Flags:**

- `-f, --force`: Force initialization (overwrite existing config)
- `-t, --templates-dir string`: Template directory name (default ".tg")

**Examples:**

```bash
tg init
tg init --templates-dir templates
tg init -f
```

### `tg list` (alias: `ls`)

List all available templates.

**Usage:**

```bash
tg list [flags]
tg ls [flags]
```

**Flags:**

- `-d, --details`: Show detailed template information
- `-F, --format string`: Output format: list, table, json (default "list")
- `-f, --filter string`: Filter templates by name (case-insensitive)

**Examples:**

```bash
# Basic list
tg list

# Detailed information
tg list --details

# Table format
tg list --format table

# JSON format
tg list --format json

# Filter by name
tg list --filter "web"
```

### `tg apply`

Apply a template to generate files.

**Usage:**

```bash
tg apply <template-name> [output-dir] [flags]
```

**Flags:**

- `-o, --output string`: Output directory (default ".")
- `-v, --var stringToString`: Set variable values (e.g., -v name=John -v age=30)

**Examples:**

```bash
# Apply to current directory
tg apply hello-world

# Apply to specific directory
tg apply hello-world ./my-project

# Override variables
tg apply web-app -v project_name=MyApp -v port=8080
```

### Global Flags

Available for all commands:

- `-V, --verbose`: Enable verbose output
- `-c, --config string`: Path to config file (default "tg.config.toml")
- `--version`: Display version information

## Configuration

### Main Configuration (`tg.config.toml`)

```toml
# Directory containing templates
templates_dir = ".tg"

# Git remote for fetching templates (optional, coming soon)
# git_remote = "https://github.com/yourusername/tg-templates.git"

# Default variables for all templates (optional)
[defaults]
author = "Your Name"
license = "MIT"
```

### Template Configuration (`template.toml`)

```toml
version = "1.0.0"

[metadata]
name = "template-name"
description = "Template description"
author = "Author Name"

# Define variables with types and defaults
[variables]
var_name={default="default_value" description="Variable description"}

# File processing rules
[rules]
ignores = ["*.tmp", ".git", "node_modules"]
includes = ["**/*.go", "**/*.md", "**/*.json"]
renames = {"README.template.md"="README.md"}
```

## Variable Types

The template system supports the following variable types:

- **string**: Text values
- **number**: Integer or floating-point numbers
- **boolean**: true/false values
- **array**: List of values

Type validation is performed automatically when loading templates.

## Template Syntax

Templates use Go's `text/template` syntax:

```go
// Simple variable substitution
{{.variable_name}}

// Conditional
{{if .enable_feature}}
Feature is enabled
{{end}}

// Range over array
{{range .items}}
- {{.}}
{{end}}

// Pipeline
{{.project_name | upper}}
```

## Project Structure

```
template_generator/
├── cmd/
│   └── main.go                 # Application entry point
├── internal/
│   ├── cli/
│   │   ├── root.go            # Root command and CLI setup
│   │   ├── init.go            # Init command implementation
│   │   ├── list.go            # List command implementation
│   │   └── apply.go           # Apply command implementation
│   ├── config/
│   │   └── config.go          # Configuration and template loading
│   └── template/
│       └── processor.go       # Template processing logic
├── go.mod
├── go.sum
└── README.md
```

## Dependencies

- [spf13/cobra](https://github.com/spf13/cobra): CLI framework
- [pelletier/go-toml](https://github.com/pelletier/go-toml): TOML parser
- [fatih/color](https://github.com/fatih/color): Colored terminal output

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Author

Naviary Sanctuary

## Acknowledgments

- Inspired by modern scaffolding tools like Yeoman, Cookiecutter, and Plop
- Built with Go for performance and cross-platform compatibility
