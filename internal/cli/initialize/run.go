//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package initialize

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/crypto"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// runInit executes the init command logic.
//
// Creates a .context/ directory with template files. Handles existing
// directories, minimal mode, and CLAUDE.md/PROMPT.md merge operations.
//
// Parameters:
//   - cmd: Cobra command for output and input streams
//   - force: If true, overwrite existing files without prompting
//   - minimal: If true, only create essential files
//   - merge: If true, auto-merge ctx content into existing files
//   - ralph: If true, use autonomous loop templates (no questions, signals)
//   - noPluginEnable: If true, skip auto-enabling the plugin globally
//
// Returns:
//   - error: Non-nil if directory creation or file operations fail
func runInit(cmd *cobra.Command, force, minimal, merge, ralph, noPluginEnable bool) error {
	// Check if ctx is in PATH (required for hooks to work)
	if err := checkCtxInPath(cmd); err != nil {
		return err
	}

	contextDir := rc.ContextDir()

	// Check if .context/ already exists and is properly initialized.
	// A directory with only logs/ (created by hooks before init) is
	// treated as uninitialized — no overwrite prompt needed.
	if _, err := os.Stat(contextDir); err == nil {
		if !force && hasEssentialFiles(contextDir) {
			// Prompt for confirmation
			cmd.Printf("%s already exists. Overwrite? [y/N] ", contextDir)
			reader := bufio.NewReader(os.Stdin)
			response, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("failed to read input: %w", err)
			}
			response = strings.TrimSpace(strings.ToLower(response))
			if response != "y" && response != "yes" { //nolint:goconst // trivial user input check
				cmd.Println("Aborted.")
				return nil
			}
		}
	}

	// Create .context/ directory
	if err := os.MkdirAll(contextDir, config.PermExec); err != nil {
		return fmt.Errorf("failed to create %s: %w", contextDir, err)
	}

	// Get the list of templates to create
	var templatesToCreate []string
	if minimal {
		templatesToCreate = config.FilesRequired
	} else {
		var listErr error
		templatesToCreate, listErr = assets.List()
		if listErr != nil {
			return fmt.Errorf("failed to list templates: %w", listErr)
		}
	}

	// Create template files
	green := color.New(color.FgGreen).SprintFunc()
	for _, name := range templatesToCreate {
		targetPath := filepath.Join(contextDir, name)

		// Check if the file exists and --force not set
		if _, err := os.Stat(targetPath); err == nil && !force {
			cmd.Printf(
				"  %s %s (exists, skipped)\n", color.YellowString("○"), name,
			)
			continue
		}

		content, err := assets.Template(name)
		if err != nil {
			return fmt.Errorf("failed to read template %s: %w", name, err)
		}

		if err := os.WriteFile(targetPath, content, config.PermFile); err != nil {
			return fmt.Errorf("failed to write %s: %w", targetPath, err)
		}

		cmd.Printf("  %s %s\n", green("✓"), name)
	}

	cmd.Printf("\n%s initialized in %s/\n", green("Context"), contextDir)

	// Create entry templates in .context/templates/
	if err := createEntryTemplates(cmd, contextDir, force); err != nil {
		// Non-fatal: warn but continue
		cmd.Printf("  %s Entry templates: %v\n", color.YellowString("⚠"), err)
	}

	// Migrate legacy key files and promote to user-level path.
	cwd, _ := os.Getwd()
	config.MigrateKeyFile(contextDir, cwd)

	// Set up scratchpad
	if err := initScratchpad(cmd, contextDir); err != nil {
		// Non-fatal: warn but continue
		cmd.Printf("  %s Scratchpad: %v\n", color.YellowString("⚠"), err)
	}

	// Create project root files
	cmd.Println("\nCreating project root files...")

	// Create PROMPT.md (uses ralph template if --ralph flag set)
	if err := handlePromptMd(cmd, force, merge, ralph); err != nil {
		// Non-fatal: warn but continue
		cmd.Printf("  %s PROMPT.md: %v\n", color.YellowString("⚠"), err)
	}

	// Create IMPLEMENTATION_PLAN.md
	if err := handleImplementationPlan(cmd, force, merge); err != nil {
		// Non-fatal: warn but continue
		cmd.Printf(
			"  %s IMPLEMENTATION_PLAN.md: %v\n", color.YellowString("⚠"), err,
		)
	}

	// Merge permissions into settings.local.json (no hook scaffolding)
	cmd.Println("\nSetting up Claude Code permissions...")
	if err := mergeSettingsPermissions(cmd); err != nil {
		// Non-fatal: warn but continue
		cmd.Printf("  %s Permissions: %v\n", color.YellowString("⚠"), err)
	}

	// Auto-enable plugin globally unless suppressed
	if !noPluginEnable {
		if pluginErr := enablePluginGlobally(cmd); pluginErr != nil {
			// Non-fatal: warn but continue
			cmd.Printf("  %s Plugin enablement: %v\n", color.YellowString("⚠"), pluginErr)
		}
	}

	// Handle CLAUDE.md creation/merge
	if err := handleClaudeMd(cmd, force, merge); err != nil {
		// Non-fatal: warn but continue
		cmd.Printf("  %s CLAUDE.md: %v\n", color.YellowString("⚠"), err)
	}

	// Deploy Makefile.ctx and amend user Makefile
	if err := handleMakefileCtx(cmd); err != nil {
		// Non-fatal: warn but continue
		cmd.Printf("  %s Makefile: %v\n", color.YellowString("⚠"), err)
	}

	// Update .gitignore with recommended entries
	if err := ensureGitignoreEntries(cmd); err != nil {
		cmd.Printf("  %s .gitignore: %v\n", color.YellowString("⚠"), err)
	}

	cmd.Println("\nNext steps:")
	cmd.Println("  1. Edit .context/TASKS.md to add your current tasks")
	cmd.Println("  2. Run 'ctx status' to see context summary")
	cmd.Println("  3. Run 'ctx agent' to get AI-ready context packet")
	cmd.Println()
	cmd.Println("Claude Code users: install the ctx plugin for hooks & skills:")
	cmd.Println("  /plugin marketplace add ActiveMemory/ctx")
	cmd.Println("  /plugin install ctx@activememory-ctx")
	cmd.Println()
	cmd.Println("Note: local plugin installs are not auto-enabled globally.")
	cmd.Println("Run 'ctx init' again after installing the plugin to enable it,")
	cmd.Println("or manually add to ~/.claude/settings.json:")
	cmd.Println("  {\"enabledPlugins\": {\"ctx@activememory-ctx\": true}}")

	return nil
}

// initScratchpad sets up the scratchpad key or plaintext file.
//
// When encryption is enabled (default):
//   - Generates a 256-bit key at .context/.ctx.key if not present
//   - Adds the key file to .gitignore
//   - Warns if .enc exists but no key
//
// When encryption is disabled:
//   - Creates empty .context/scratchpad.md if not present
//
// Parameters:
//   - cmd: Cobra command for output
//   - contextDir: The .context/ directory path
//
// Returns:
//   - error: Non-nil if key generation or file operations fail
func initScratchpad(cmd *cobra.Command, contextDir string) error {
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	if !rc.ScratchpadEncrypt() {
		// Plaintext mode: create empty scratchpad.md if not present
		mdPath := filepath.Join(contextDir, config.FileScratchpadMd)
		if _, err := os.Stat(mdPath); err != nil {
			if err := os.WriteFile(mdPath, nil, config.PermFile); err != nil {
				return fmt.Errorf("failed to create %s: %w", mdPath, err)
			}
			cmd.Printf("  %s %s (plaintext scratchpad)\n", green("✓"), mdPath)
		} else {
			cmd.Printf("  %s %s (exists, skipped)\n", yellow("○"), mdPath)
		}
		return nil
	}

	// Encrypted mode
	kPath := rc.KeyPath()
	encPath := filepath.Join(contextDir, config.FileScratchpadEnc)

	// Check if key already exists (idempotent)
	if _, err := os.Stat(kPath); err == nil {
		cmd.Printf("  %s %s (exists, skipped)\n", yellow("○"), kPath)
		return nil
	}

	// Warn if encrypted file exists but no key
	if _, err := os.Stat(encPath); err == nil {
		cmd.Printf("  %s Encrypted scratchpad found but no key at %s\n",
			yellow("⚠"), kPath)
		return nil
	}

	// Ensure key directory exists.
	if mkdirErr := os.MkdirAll(filepath.Dir(kPath), config.PermKeyDir); mkdirErr != nil {
		return fmt.Errorf("failed to create key dir: %w", mkdirErr)
	}

	// Generate key
	key, err := crypto.GenerateKey()
	if err != nil {
		return fmt.Errorf("failed to generate scratchpad key: %w", err)
	}

	if err := crypto.SaveKey(kPath, key); err != nil {
		return fmt.Errorf("failed to save scratchpad key: %w", err)
	}
	cmd.Printf("  %s Scratchpad key created at %s\n", green("✓"), kPath)

	return nil
}

// hasEssentialFiles reports whether contextDir contains at least one of the
// essential context files (TASKS.md, CONSTITUTION.md, DECISIONS.md). A
// directory with only logs/ or other non-essential content is considered
// uninitialized.
func hasEssentialFiles(contextDir string) bool {
	for _, f := range config.FilesRequired {
		if _, err := os.Stat(filepath.Join(contextDir, f)); err == nil {
			return true
		}
	}
	return false
}

// ensureGitignoreEntries appends recommended .gitignore entries that are not
// already present. Creates .gitignore if it does not exist.
func ensureGitignoreEntries(cmd *cobra.Command) error {
	gitignorePath := ".gitignore"

	content, err := os.ReadFile(gitignorePath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	// Build set of existing trimmed lines.
	existing := make(map[string]bool)
	for _, line := range strings.Split(string(content), "\n") {
		existing[strings.TrimSpace(line)] = true
	}

	// Collect missing entries.
	var missing []string
	for _, entry := range config.GitignoreEntries {
		if !existing[entry] {
			missing = append(missing, entry)
		}
	}

	if len(missing) == 0 {
		return nil
	}

	// Build block to append.
	var sb strings.Builder
	if len(content) > 0 && !strings.HasSuffix(string(content), "\n") {
		sb.WriteString("\n")
	}
	sb.WriteString("\n# ctx managed entries\n")
	for _, entry := range missing {
		sb.WriteString(entry + "\n")
	}

	if err := os.WriteFile(gitignorePath, append(content, []byte(sb.String())...), config.PermFile); err != nil {
		return err
	}

	green := color.New(color.FgGreen).SprintFunc()
	cmd.Printf("  %s .gitignore updated (%d entries added)\n", green("✓"), len(missing))
	cmd.Println("  Review with: cat .gitignore")
	return nil
}

// addToGitignore ensures an entry exists in .gitignore.
//
// Creates .gitignore if it doesn't exist. Checks if the entry is already
// present before adding.
//
// Parameters:
//   - contextDir: The .context/ directory (entry is relative to this)
//   - filename: The filename to add (e.g., ".ctx.key")
func addToGitignore(contextDir, filename string) error {
	entry := filepath.Join(contextDir, filename)
	gitignorePath := ".gitignore"

	// Read existing .gitignore
	content, err := os.ReadFile(gitignorePath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	// Check if already present
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == entry {
			return nil // already present
		}
	}

	// Append entry
	var newContent string
	if len(content) > 0 && !strings.HasSuffix(string(content), "\n") {
		newContent = string(content) + "\n" + entry + "\n"
	} else {
		newContent = string(content) + entry + "\n"
	}

	return os.WriteFile(gitignorePath, []byte(newContent), config.PermFile)
}
