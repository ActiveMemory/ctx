//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hook

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config"
)

// copilotInstructions is the comprehensive GitHub Copilot integration
// template for .github/copilot-instructions.md.
//
// It instructs Copilot Chat (agent mode) to load .context/ files at
// session start, persist session summaries to .context/sessions/, and
// proactively update context files during work.
const copilotInstructions = `# Project Context

<!-- ctx:copilot -->
<!-- DO NOT REMOVE: This marker indicates ctx-managed content -->

## Context System

This project uses Context (` + "`ctx`" + `) for persistent AI context
management. Your memory is NOT ephemeral — it lives in ` + "`.context/`" + ` files.

## On Session Start

Read these files **in order** before starting any work:

1. ` + "`.context/CONSTITUTION.md`" + ` — Hard rules, NEVER violate
2. ` + "`.context/TASKS.md`" + ` — Current work items
3. ` + "`.context/CONVENTIONS.md`" + ` — Code patterns and standards
4. ` + "`.context/ARCHITECTURE.md`" + ` — System structure
5. ` + "`.context/DECISIONS.md`" + ` — Architectural decisions with rationale
6. ` + "`.context/LEARNINGS.md`" + ` — Gotchas, tips, lessons learned
7. ` + "`.context/GLOSSARY.md`" + ` — Domain terms and abbreviations
8. ` + "`.context/AGENT_PLAYBOOK.md`" + ` — How to use this context system

After reading, confirm: "I have read the required context files and I'm
following project conventions."

## When Asked "Do You Remember?"

**Do this FIRST (silently):**
- Read ` + "`.context/TASKS.md`" + `
- Read ` + "`.context/DECISIONS.md`" + ` and ` + "`.context/LEARNINGS.md`" + `
- Check ` + "`.context/sessions/`" + ` for recent session files

**Then respond with a structured readback:**
1. **Last session**: cite the most recent session topic and date
2. **Active work**: list pending or in-progress tasks
3. **Recent context**: mention 1-2 recent decisions or learnings
4. **Next step**: offer to continue or ask what to focus on

**Never** lead with "I don't have memory." The ` + "`.context/`" + ` files are
your memory. Read them silently, then present what you found as recall.

## Session Persistence

After completing meaningful work, save a session summary to
` + "`.context/sessions/`" + `.

### Session File Format

Create a file named ` + "`YYYY-MM-DD-topic.md`" + `:

` + "```" + `markdown
# Session: YYYY-MM-DD — Brief Topic Description

## What Was Done
- Describe completed work items

## Decisions
- Key decisions made and their rationale

## Learnings
- Gotchas, tips, or insights discovered

## Next Steps
- Follow-up work or remaining items
` + "```" + `

### When to Save

- After completing a task or feature
- After making architectural decisions
- After a debugging session
- Before ending the session
- At natural breakpoints in long sessions

## Context Updates During Work

Proactively update context files as you work:

| Event                       | Action                              |
|-----------------------------|-------------------------------------|
| Made architectural decision | Add to ` + "`.context/DECISIONS.md`" + `  |
| Discovered gotcha/bug       | Add to ` + "`.context/LEARNINGS.md`" + `  |
| Established new pattern     | Add to ` + "`.context/CONVENTIONS.md`" + ` |
| Completed task              | Mark [x] in ` + "`.context/TASKS.md`" + ` |

## Self-Check

Periodically ask yourself:

> "If this session ended right now, would the next session know what happened?"

If no — save a session file or update context files before continuing.

## CLI Commands

If ` + "`ctx`" + ` is installed, use these commands:

` + "```" + `bash
ctx status        # Context summary and health check
ctx agent         # AI-ready context packet
ctx drift         # Check for stale context
ctx recall list   # Recent session history
` + "```" + `

<!-- ctx:copilot:end -->
`

// ToolConfigFiles maps tool names to their configuration file paths.
var ToolConfigFiles = map[string]string{
	"copilot":  filepath.Join(".github", "copilot-instructions.md"),
	"cursor":   ".cursorrules",
	"windsurf": ".windsurfrules",
}

// runHook executes the hook command logic.
//
// Outputs integration instructions and configuration snippets for the
// specified AI tool. With --write, generates the configuration file
// directly.
//
// Parameters:
//   - cmd: Cobra command for output stream
//   - args: Command arguments; args[0] is the tool name
//   - write: If true, write the configuration file instead of printing
//
// Returns:
//   - error: Non-nil if the tool is not supported or file write fails
func runHook(cmd *cobra.Command, args []string, write bool) error {
	tool := strings.ToLower(args[0])

	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()

	switch tool {
	case "claude-code", "claude":
		cmd.Println(cyan("Claude Code Integration"))
		cmd.Println(cyan("======================="))
		cmd.Println()
		cmd.Println("Claude Code integration is now provided via the ctx plugin.")
		cmd.Println()
		cmd.Println("Install the plugin:")
		cmd.Println(green("  /plugin marketplace add ActiveMemory/ctx"))
		cmd.Println(green("  /plugin install ctx@activememory-ctx"))
		cmd.Println()
		cmd.Println("The plugin provides hooks (context monitoring, persistence")
		cmd.Println("nudges, post-commit capture) and 25 skills automatically.")

	case "cursor":
		cmd.Println(cyan("Cursor IDE Integration"))
		cmd.Println(cyan("======================"))
		cmd.Println()
		cmd.Println("Add to your .cursorrules file:")
		cmd.Println()
		cmd.Println(green("```markdown"))
		cmd.Print(`# Project Context

Always read these files before making changes:
- .context/CONSTITUTION.md (NEVER violate these rules)
- .context/TASKS.md (current work)
- .context/CONVENTIONS.md (how we write code)
- .context/ARCHITECTURE.md (system structure)

Run 'ctx agent' for a context summary.
Run 'ctx drift' to check for stale context.
`)
		cmd.Println(green("```"))

	case "aider":
		cmd.Println(cyan("Aider Integration"))
		cmd.Println(cyan("================="))
		cmd.Println()
		cmd.Println("Add to your .aider.conf.yml:")
		cmd.Println()
		cmd.Println(green("```yaml"))
		cmd.Println(`read:
  - .context/CONSTITUTION.md
  - .context/TASKS.md
  - .context/CONVENTIONS.md
  - .context/ARCHITECTURE.md
  - .context/DECISIONS.md`)
		cmd.Println(green("```"))
		cmd.Println()
		cmd.Println("Or pass context via command line:")
		cmd.Println()
		cmd.Println(green("```bash"))
		cmd.Println(`ctx agent | aider --message "$(cat -)"`)
		cmd.Println(green("```"))

	case "copilot":
		if write {
			return writeCopilotInstructions(cmd)
		}
		cmd.Println(cyan("GitHub Copilot Integration"))
		cmd.Println(cyan("=========================="))
		cmd.Println()
		cmd.Println("Add the following to .github/copilot-instructions.md,")
		cmd.Println("or run with --write to generate the file directly:")
		cmd.Println()
		cmd.Println(green("  ctx hook copilot --write"))
		cmd.Println()
		cmd.Print(copilotInstructions)

	case "windsurf":
		cmd.Println(cyan("Windsurf Integration"))
		cmd.Println(cyan("===================="))
		cmd.Println()
		cmd.Println("Add to your .windsurfrules file:")
		cmd.Println()
		cmd.Println(green("```markdown"))
		cmd.Print(`# Context

Read order for context:
1. .context/CONSTITUTION.md
2. .context/TASKS.md
3. .context/CONVENTIONS.md
4. .context/ARCHITECTURE.md
5. .context/DECISIONS.md

Run 'ctx agent' for AI-ready context packet.
`)
		cmd.Println(green("```"))

	default:
		cmd.Printf("Unknown tool: %s\n\n", tool)
		cmd.Println("Supported tools:")
		cmd.Println("  claude-code  - Anthropic's Claude Code CLI (use plugin instead)")
		cmd.Println("  cursor       - Cursor IDE")
		cmd.Println("  aider        - Aider AI coding assistant")
		cmd.Println("  copilot      - GitHub Copilot")
		cmd.Println("  windsurf     - Windsurf IDE")
		return fmt.Errorf("unsupported tool: %s", tool)
	}

	return nil
}

// writeCopilotInstructions generates .github/copilot-instructions.md.
//
// Creates the .github/ directory if needed and writes the comprehensive
// Copilot instructions file. Preserves existing non-ctx content by
// checking for ctx markers.
//
// Parameters:
//   - cmd: Cobra command for output messages
//
// Returns:
//   - error: Non-nil if directory creation or file write fails
func writeCopilotInstructions(cmd *cobra.Command) error {
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	targetDir := ".github"
	targetFile := filepath.Join(targetDir, "copilot-instructions.md")

	// Create .github/ directory if needed
	if err := os.MkdirAll(targetDir, config.PermExec); err != nil {
		return fmt.Errorf("failed to create %s: %w", targetDir, err)
	}

	// Check if file exists
	existingContent, err := os.ReadFile(filepath.Clean(targetFile)) //nolint:gosec // targetFile is constructed from constants, not user input
	fileExists := err == nil

	if fileExists {
		existingStr := string(existingContent)
		if strings.Contains(existingStr, "<!-- ctx:copilot -->") {
			cmd.Println(fmt.Sprintf(
				"  %s %s (ctx content exists, skipped)", yellow("○"), targetFile,
			))
			cmd.Println("  Use --force to overwrite (not yet implemented).")
			return nil
		}

		// File exists without ctx markers: append ctx content
		merged := existingStr + config.NewlineLF + copilotInstructions
		if err := os.WriteFile(targetFile, []byte(merged), config.PermFile); err != nil {
			return fmt.Errorf("failed to write %s: %w", targetFile, err)
		}
		cmd.Println(fmt.Sprintf("  %s %s (merged)", green("✓"), targetFile))
		return nil
	}

	// File doesn't exist: create it
	if err := os.WriteFile(
		targetFile, []byte(copilotInstructions), config.PermFile,
	); err != nil {
		return fmt.Errorf("failed to write %s: %w", targetFile, err)
	}
	cmd.Println(fmt.Sprintf("  %s %s", green("✓"), targetFile))

	// Also create .context/sessions/ if it doesn't exist
	sessionsDir := filepath.Join(config.DirContext, config.DirSessions)
	if err := os.MkdirAll(sessionsDir, config.PermExec); err != nil {
		cmd.Println(fmt.Sprintf(
			"  %s %s: %v", yellow("⚠"), sessionsDir, err,
		))
	} else {
		cmd.Println(fmt.Sprintf("  %s %s/", green("✓"), sessionsDir))
	}

	cmd.Println()
	cmd.Println("Copilot Chat (agent mode) will now:")
	cmd.Println("  1. Read .context/ files at session start")
	cmd.Println("  2. Save session summaries to .context/sessions/")
	cmd.Println("  3. Proactively update context during work")

	return nil
}
