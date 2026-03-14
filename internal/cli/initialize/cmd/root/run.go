//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"github.com/ActiveMemory/ctx/internal/config/cli"
	"github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/config/pad"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/cli/initialize/core"
	"github.com/ActiveMemory/ctx/internal/crypto"
	ctxerr "github.com/ActiveMemory/ctx/internal/err"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/write"
)

// gitignoreHeader is the section comment prepended to ctx-managed entries.
const gitignoreHeader = "# ctx managed entries"

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
//
// Returns:
//   - error: Non-nil if directory creation or file operations fail
func Run(cmd *cobra.Command, force, minimal, merge, ralph, noPluginEnable bool) error {
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
			write.InfoInitOverwritePrompt(cmd, contextDir)
			reader := bufio.NewReader(os.Stdin)
			response, err := reader.ReadString('\n')
			if err != nil {
				return ctxerr.ReadInput(err)
			}
			response = strings.TrimSpace(strings.ToLower(response))
			if response != cli.ConfirmShort && response != cli.ConfirmLong {
				write.InfoInitAborted(cmd)
				return nil
			}
		}
	}

	// Create .context/ directory
	if err := os.MkdirAll(contextDir, fs.PermExec); err != nil {
		return ctxerr.Mkdir(contextDir, err)
	}

	// Get the list of templates to create
	var templatesToCreate []string
	if minimal {
		templatesToCreate = ctx.FilesRequired
	} else {
		var listErr error
		templatesToCreate, listErr = assets.List()
		if listErr != nil {
			return ctxerr.ListTemplates(listErr)
		}
	}

	// Create template files
	for _, name := range templatesToCreate {
		targetPath := filepath.Join(contextDir, name)

		// Check if the file exists and --force not set
		if _, err := os.Stat(targetPath); err == nil && !force {
			write.InfoInitExistsSkipped(cmd, name)
			continue
		}

		content, err := assets.Template(name)
		if err != nil {
			return ctxerr.ReadTemplate(name, err)
		}

		if err := os.WriteFile(targetPath, content, fs.PermFile); err != nil {
			return ctxerr.FileWrite(targetPath, err)
		}

		write.InfoInitFileCreated(cmd, name)
	}

	write.InfoInitialized(cmd, contextDir)

	// Create entry templates in .context/templates/
	if err := core.CreateEntryTemplates(cmd, contextDir, force); err != nil {
		// Non-fatal: warn but continue
		write.InfoInitWarnNonFatal(cmd, "Entry templates", err)
	}

	// Create prompt templates in .context/prompts/
	if err := core.CreatePromptTemplates(cmd, contextDir, force); err != nil {
		// Non-fatal: warn but continue
		write.InfoInitWarnNonFatal(cmd, "Prompt templates", err)
	}

	// Migrate legacy key files and promote to global path.
	crypto.MigrateKeyFile(contextDir)

	// Set up scratchpad
	if err := initScratchpad(cmd, contextDir); err != nil {
		// Non-fatal: warn but continue
		write.InfoInitWarnNonFatal(cmd, "Scratchpad", err)
	}

	// Create project root files
	write.InfoInitCreatingRootFiles(cmd)

	// Create specs/ and ideas/ directories with README.md
	if err := core.CreateProjectDirs(cmd); err != nil {
		write.InfoInitWarnNonFatal(cmd, "Project dirs", err)
	}

	// Create PROMPT.md (uses ralph template if --ralph flag set)
	if err := core.HandlePromptMd(cmd, force, merge, ralph); err != nil {
		// Non-fatal: warn but continue
		write.InfoInitWarnNonFatal(cmd, "PROMPT.md", err)
	}

	// Create IMPLEMENTATION_PLAN.md
	if err := core.HandleImplementationPlan(cmd, force, merge); err != nil {
		// Non-fatal: warn but continue
		write.InfoInitWarnNonFatal(cmd, "IMPLEMENTATION_PLAN.md", err)
	}

	// Merge permissions into settings.local.json (no hook scaffolding)
	write.InfoInitSettingUpPermissions(cmd)
	if err := core.MergeSettingsPermissions(cmd); err != nil {
		// Non-fatal: warn but continue
		write.InfoInitWarnNonFatal(cmd, "Permissions", err)
	}

	// Auto-enable plugin globally unless suppressed
	if !noPluginEnable {
		if pluginErr := core.EnablePluginGlobally(cmd); pluginErr != nil {
			// Non-fatal: warn but continue
			write.InfoInitWarnNonFatal(cmd, "Plugin enablement", pluginErr)
		}
	}

	// Handle CLAUDE.md creation/merge
	if err := core.HandleClaudeMd(cmd, force, merge); err != nil {
		// Non-fatal: warn but continue
		write.InfoInitWarnNonFatal(cmd, "CLAUDE.md", err)
	}

	// Deploy Makefile.ctx and amend user Makefile
	if err := core.HandleMakefileCtx(cmd); err != nil {
		// Non-fatal: warn but continue
		write.InfoInitWarnNonFatal(cmd, "Makefile", err)
	}

	// Update .gitignore with recommended entries
	if err := ensureGitignoreEntries(cmd); err != nil {
		write.InfoInitWarnNonFatal(cmd, ".gitignore", err)
	}

	write.InfoInitNextSteps(cmd)

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
		mdPath := filepath.Join(contextDir, pad.Md)
		if _, err := os.Stat(mdPath); err != nil {
			if err := os.WriteFile(mdPath, nil, fs.PermFile); err != nil {
				return ctxerr.Mkdir(mdPath, err)
			}
			write.InfoInitScratchpadPlaintext(cmd, mdPath)
		} else {
			write.InfoInitExistsSkipped(cmd, mdPath)
		}
		return nil
	}

	// Encrypted mode
	kPath := rc.KeyPath()
	encPath := filepath.Join(contextDir, pad.Enc)

	// Check if key already exists (idempotent)
	if _, err := os.Stat(kPath); err == nil {
		write.InfoInitExistsSkipped(cmd, kPath)
		return nil
	}

	// Warn if encrypted file exists but no key
	if _, err := os.Stat(encPath); err == nil {
		write.InfoInitScratchpadNoKey(cmd, kPath)
		return nil
	}

	// Ensure key directory exists.
	if mkdirErr := os.MkdirAll(filepath.Dir(kPath), fs.PermKeyDir); mkdirErr != nil {
		return ctxerr.MkdirKeyDir(mkdirErr)
	}

	// Generate key
	key, err := crypto.GenerateKey()
	if err != nil {
		return ctxerr.GenerateKey(err)
	}

	if err := crypto.SaveKey(kPath, key); err != nil {
		return ctxerr.SaveKey(err)
	}
	write.InfoInitScratchpadKeyCreated(cmd, kPath)

	return nil
}

// hasEssentialFiles reports whether contextDir contains at least one of the
// essential context files (TASKS.md, CONSTITUTION.md, DECISIONS.md). A
// directory with only logs/ or other non-essential content is considered
// uninitialized.
func hasEssentialFiles(contextDir string) bool {
	for _, f := range ctx.FilesRequired {
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
	for _, line := range strings.Split(string(content), token.NewlineLF) {
		existing[strings.TrimSpace(line)] = true
	}

	// Collect missing entries.
	var missing []string
	for _, entry := range file.Gitignore {
		if !existing[entry] {
			missing = append(missing, entry)
		}
	}

	if len(missing) == 0 {
		return nil
	}

	// Build block to append.
	var sb strings.Builder
	if len(content) > 0 && !strings.HasSuffix(string(content), token.NewlineLF) {
		sb.WriteString(token.NewlineLF)
	}
	sb.WriteString(token.NewlineLF + gitignoreHeader + token.NewlineLF)
	for _, entry := range missing {
		sb.WriteString(entry + token.NewlineLF)
	}

	if err := os.WriteFile(gitignorePath, append(content, []byte(sb.String())...), fs.PermFile); err != nil {
		return err
	}

	write.InfoInitGitignoreUpdated(cmd, len(missing))
	write.InfoInitGitignoreReview(cmd)
	return nil
}
