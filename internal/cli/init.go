// Package cli implements the CLI commands for amem.
package cli

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/josealekhine/ActiveMemory/internal/templates"
	"github.com/spf13/cobra"
)

const (
	contextDirName = ".context"
)

var (
	initForce   bool
	initMinimal bool
)

// minimalTemplates are the essential files created with --minimal flag
var minimalTemplates = []string{
	"TASKS.md",
	"DECISIONS.md",
	"CONSTITUTION.md",
}

// InitCmd returns the init command.
func InitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a new .context/ directory with template files",
		Long: `Initialize a new .context/ directory with template files for
maintaining persistent context for AI coding assistants.

The following files are created:
  - CONSTITUTION.md  — Hard invariants that must never be violated
  - TASKS.md         — Current and planned work
  - DECISIONS.md     — Architectural decisions with rationale
  - LEARNINGS.md     — Lessons learned, gotchas, tips
  - CONVENTIONS.md   — Project patterns and standards
  - ARCHITECTURE.md  — System overview
  - GLOSSARY.md      — Domain terms and abbreviations
  - DRIFT.md         — Staleness signals and update triggers
  - AGENT_PLAYBOOK.md — How AI agents should use this system

Use --minimal to only create essential files (TASKS.md, DECISIONS.md, CONSTITUTION.md).`,
		RunE: runInit,
	}

	cmd.Flags().BoolVarP(&initForce, "force", "f", false, "Overwrite existing context files")
	cmd.Flags().BoolVarP(&initMinimal, "minimal", "m", false, "Only create essential files (TASKS.md, DECISIONS.md, CONSTITUTION.md)")

	return cmd
}

func runInit(cmd *cobra.Command, args []string) error {
	contextDir := contextDirName

	// Check if .context/ already exists
	if _, err := os.Stat(contextDir); err == nil {
		if !initForce {
			// Prompt for confirmation
			fmt.Printf("%s already exists. Overwrite? [y/N] ", contextDir)
			reader := bufio.NewReader(os.Stdin)
			response, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("failed to read input: %w", err)
			}
			response = strings.TrimSpace(strings.ToLower(response))
			if response != "y" && response != "yes" {
				fmt.Println("Aborted.")
				return nil
			}
		}
	}

	// Create .context/ directory
	if err := os.MkdirAll(contextDir, 0755); err != nil {
		return fmt.Errorf("failed to create %s: %w", contextDir, err)
	}

	// Get list of templates to create
	var templatesToCreate []string
	if initMinimal {
		templatesToCreate = minimalTemplates
	} else {
		var err error
		templatesToCreate, err = templates.ListTemplates()
		if err != nil {
			return fmt.Errorf("failed to list templates: %w", err)
		}
	}

	// Create template files
	green := color.New(color.FgGreen).SprintFunc()
	for _, name := range templatesToCreate {
		targetPath := filepath.Join(contextDir, name)

		// Check if file exists and --force not set
		if _, err := os.Stat(targetPath); err == nil && !initForce {
			fmt.Printf("  %s %s (exists, skipped)\n", color.YellowString("○"), name)
			continue
		}

		content, err := templates.GetTemplate(name)
		if err != nil {
			return fmt.Errorf("failed to read template %s: %w", name, err)
		}

		if err := os.WriteFile(targetPath, content, 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", targetPath, err)
		}

		fmt.Printf("  %s %s\n", green("✓"), name)
	}

	fmt.Printf("\n%s initialized in %s/\n", green("Active Memory"), contextDir)
	fmt.Println("\nNext steps:")
	fmt.Println("  1. Edit .context/TASKS.md to add your current tasks")
	fmt.Println("  2. Run 'amem status' to see context summary")
	fmt.Println("  3. Run 'amem agent' to get AI-ready context packet")

	return nil
}
