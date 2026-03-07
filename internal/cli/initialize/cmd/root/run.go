//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/cli/initialize/core"
	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/crypto"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// Run executes the init command logic.
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
//   - caller: Identifies the calling tool (e.g. "vscode") for template overrides
//
// Returns:
//   - error: Non-nil if directory creation or file operations fail
func Run(cmd *cobra.Command, force, minimal, merge, ralph, noPluginEnable bool, caller string) error {
	// Check if ctx is in PATH (required for hooks to work)
	if err := core.CheckCtxInPath(cmd); err != nil {
		return err
	}

	contextDir := rc.ContextDir()

	// Check if .context/ already exists and is properly initialized.
	// A directory with only logs/ (created by hooks before init) is
	// treated as uninitialized — no overwrite prompt needed.
	if _, err := os.Stat(contextDir); err == nil {
		if !force && hasEssentialFiles(contextDir) {
			// Prompt for confirmation
			cmd.Print(fmt.Sprintf("%s already exists. Overwrite? [y/N] ", contextDir))
			reader := bufio.NewReader(os.Stdin)
			response, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("failed to read input: %w", err)
			}
			response = strings.TrimSpace(strings.ToLower(response))
			if response != config.ConfirmShort && response != config.ConfirmLong {
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
	for _, name := range templatesToCreate {
		targetPath := filepath.Join(contextDir, name)

		// Check if the file exists and --force not set
		if _, err := os.Stat(targetPath); err == nil && !force {
			cmd.Println(fmt.Sprintf(
				"  ○ %s (exists, skipped)\n", name,
			))
			continue
		}

		content, err := assets.TemplateForCaller(name, caller)
		if err != nil {
			return fmt.Errorf("failed to read template %s: %w", name, err)
		}

		if err := os.WriteFile(targetPath, content, config.PermFile); err != nil {
			return fmt.Errorf("failed to write %s: %w", targetPath, err)
		}

		cmd.Println(fmt.Sprintf("  ✓ %s", name))
	}

	cmd.Println(fmt.Sprintf("\nContext initialized in %s/", contextDir))

	// Create entry templates in .context/templates/
	if err := core.CreateEntryTemplates(cmd, contextDir, force); err != nil {
		// Non-fatal: warn but continue
		cmd.Println(fmt.Sprintf("  ⚠ Entry templates: %v", err))
	}

	// Create prompt templates in .context/prompts/
	if err := core.CreatePromptTemplates(cmd, contextDir, force); err != nil {
		// Non-fatal: warn but continue
		cmd.Println(fmt.Sprintf("  ⚠ Prompt templates: %v", err))
	}

	// Migrate legacy key files and promote to global path.
	config.MigrateKeyFile(contextDir)

	// Set up scratchpad
	if err := initScratchpad(cmd, contextDir); err != nil {
		// Non-fatal: warn but continue
		cmd.Println(fmt.Sprintf("  ⚠ Scratchpad: %v", err))
	}

	// Create project root files
	cmd.Println("\nCreating project root files...")

	// Create specs/ and ideas/ directories with README.md
	if err := core.CreateProjectDirs(cmd); err != nil {
		cmd.Println(fmt.Sprintf("  ⚠ Project dirs: %v", err))
	}

	// Create PROMPT.md (uses ralph template if --ralph flag set)
	if err := core.HandlePromptMd(cmd, force, merge, ralph); err != nil {
		// Non-fatal: warn but continue
		cmd.Println(fmt.Sprintf("  ⚠ PROMPT.md: %v", err))
	}

	// Create IMPLEMENTATION_PLAN.md
	if err := core.HandleImplementationPlan(cmd, force, merge); err != nil {
		// Non-fatal: warn but continue
		cmd.Println(fmt.Sprintf(
			"  ⚠ IMPLEMENTATION_PLAN.md: %v\n", err,
		))
	}

	// Merge permissions into settings.local.json (no hook scaffolding)
	cmd.Println("\nSetting up Claude Code permissions...")
	if err := core.MergeSettingsPermissions(cmd); err != nil {
		// Non-fatal: warn but continue
		cmd.Println(fmt.Sprintf("  ⚠ Permissions: %v", err))
	}

	// Auto-enable plugin globally unless suppressed
	if !noPluginEnable {
		if pluginErr := core.EnablePluginGlobally(cmd); pluginErr != nil {
			// Non-fatal: warn but continue
			cmd.Println(fmt.Sprintf("  ⚠ Plugin enablement: %v", pluginErr))
		}
	}

	// Handle CLAUDE.md creation/merge
	if err := core.HandleClaudeMd(cmd, force, merge); err != nil {
		// Non-fatal: warn but continue
		cmd.Println(fmt.Sprintf("  ⚠ CLAUDE.md: %v", err))
	}

	// Deploy Makefile.ctx and amend user Makefile
	if err := core.HandleMakefileCtx(cmd); err != nil {
		// Non-fatal: warn but continue
		cmd.Println(fmt.Sprintf("  ⚠ Makefile: %v", err))
	}

	// Update .gitignore with recommended entries
	if err := ensureGitignoreEntries(cmd); err != nil {
		cmd.Println(fmt.Sprintf("  ⚠ .gitignore: %v", err))
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
//   - Generates a 256-bit key at ~/.ctx/ if not present
//   - Adds legacy key path to .gitignore for migration safety
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
	if !rc.ScratchpadEncrypt() {
		// Plaintext mode: create empty scratchpad.md if not present
		mdPath := filepath.Join(contextDir, config.FileScratchpadMd)
		if _, err := os.Stat(mdPath); err != nil {
			if err := os.WriteFile(mdPath, nil, config.PermFile); err != nil {
				return fmt.Errorf("failed to create %s: %w", mdPath, err)
			}
			cmd.Println(fmt.Sprintf("  ✓ %s (plaintext scratchpad)", mdPath))
		} else {
			cmd.Println(fmt.Sprintf("  ○ %s (exists, skipped)", mdPath))
		}
		return nil
	}

	// Encrypted mode
	kPath := rc.KeyPath()
	encPath := filepath.Join(contextDir, config.FileScratchpadEnc)

	// Check if key already exists (idempotent)
	if _, err := os.Stat(kPath); err == nil {
		cmd.Println(fmt.Sprintf("  ○ %s (exists, skipped)", kPath))
		return nil
	}

	// Warn if encrypted file exists but no key
	if _, err := os.Stat(encPath); err == nil {
		cmd.Println(fmt.Sprintf("  ⚠ Encrypted scratchpad found but no key at %s",
			kPath))
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
	cmd.Println(fmt.Sprintf("  ✓ Scratchpad key created at %s", kPath))

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
	for _, line := range strings.Split(string(content), config.NewlineLF) {
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
	if len(content) > 0 && !strings.HasSuffix(string(content), config.NewlineLF) {
		sb.WriteString(config.NewlineLF)
	}
	sb.WriteString("\n# ctx managed entries\n")
	for _, entry := range missing {
		sb.WriteString(entry + config.NewlineLF)
	}

	if err := os.WriteFile(gitignorePath, append(content, []byte(sb.String())...), config.PermFile); err != nil {
		return err
	}

	cmd.Println(fmt.Sprintf("  ✓ .gitignore updated (%d entries added)", len(missing)))
	cmd.Println("  Review with: cat .gitignore")
	return nil
}
