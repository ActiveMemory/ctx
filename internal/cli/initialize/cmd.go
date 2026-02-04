//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package initialize

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/claude"
)

// createClaudeSkills creates .claude/skills/ with Agent Skills directories.
//
// Creates skill directories following the Agent Skills specification
// (https://agentskills.io), with each skill as a directory containing
// a SKILL.md file with frontmatter and autonomy-focused instructions.
//
// Parameters:
//   - cmd: Cobra command for output messages
//   - force: If true, overwrite existing skill files
//
// Returns:
//   - error: Non-nil if directory creation or file operations fail
func createClaudeSkills(cmd *cobra.Command, force bool) error {
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	skillsDir := ".claude/skills"
	if err := os.MkdirAll(skillsDir, 0755); err != nil {
		return fmt.Errorf("failed to create %s: %w", skillsDir, err)
	}

	// Get the list of embedded skills
	skills, err := claude.Skills()
	if err != nil {
		return fmt.Errorf("failed to list skills: %w", err)
	}

	for _, skillName := range skills {
		// Create skill directory
		skillDir := filepath.Join(skillsDir, skillName)
		if err := os.MkdirAll(skillDir, 0755); err != nil {
			return fmt.Errorf("failed to create %s: %w", skillDir, err)
		}

		// Create SKILL.md file
		skillPath := filepath.Join(skillDir, "SKILL.md")
		if _, err := os.Stat(skillPath); err == nil && !force {
			cmd.Printf("  %s %s (exists, skipped)\n", yellow("○"), skillPath)
			continue
		}

		content, err := claude.SkillContent(skillName)
		if err != nil {
			return fmt.Errorf("failed to get skill %s: %w", skillName, err)
		}

		if err := os.WriteFile(skillPath, content, 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", skillPath, err)
		}
		cmd.Printf("  %s %s\n", green("✓"), skillPath)
	}

	return nil
}
