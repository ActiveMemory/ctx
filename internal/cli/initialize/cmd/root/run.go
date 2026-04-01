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

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/catalog"
	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/assets/read/template"
	coreClaude "github.com/ActiveMemory/ctx/internal/cli/initialize/core/claude"
	"github.com/ActiveMemory/ctx/internal/cli/initialize/core/entry"
	coreMerge "github.com/ActiveMemory/ctx/internal/cli/initialize/core/merge"
	"github.com/ActiveMemory/ctx/internal/cli/initialize/core/plugin"
	coreProject "github.com/ActiveMemory/ctx/internal/cli/initialize/core/project"
	"github.com/ActiveMemory/ctx/internal/cli/initialize/core/validate"
	"github.com/ActiveMemory/ctx/internal/config/claude"
	"github.com/ActiveMemory/ctx/internal/config/cli"
	"github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/config/pad"
	"github.com/ActiveMemory/ctx/internal/config/project"
	"github.com/ActiveMemory/ctx/internal/config/sync"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/crypto"
	errCrypto "github.com/ActiveMemory/ctx/internal/err/crypto"
	errFs "github.com/ActiveMemory/ctx/internal/err/fs"
	errPrompt "github.com/ActiveMemory/ctx/internal/err/prompt"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/write/initialize"
)

// Run executes the init command logic.
//
// Creates a .context/ directory with template files. Handles existing
// directories, minimal mode, and CLAUDE.md merge operations.
//
// Parameters:
//   - cmd: Cobra command for output and input streams
//   - force: If true, overwrite existing files without prompting
//   - minimal: If true, only create essential files
//   - merge: If true, auto-merge ctx content into existing files
//   - noPluginEnable: If true, skip auto-enabling the plugin globally
//   - caller: Identifies the calling tool (e.g. "vscode") for template overrides
//
// Returns:
//   - error: Non-nil if directory creation or file operations fail
func Run(
	cmd *cobra.Command, force, minimal, merge, noPluginEnable bool, caller string,
) error {
	// Check if ctx is in PATH (required for hooks to work).
	// Skip when a caller is set — the caller manages its own binary path.
	if caller == "" {
		if pathErr := validate.CheckCtxInPath(cmd); pathErr != nil {
			return pathErr
		}
	}

	contextDir := rc.ContextDir()

	// Check if .context/ already exists and is properly initialized.
	// A directory with only logs/ (created by hooks before init) is
	// treated as uninitialized - no overwrite prompt needed.
	if _, statErr := os.Stat(contextDir); statErr == nil {
		if !force && hasEssentialFiles(contextDir) {
			// When called from an editor (--caller), stdin is unavailable.
			// Skip the interactive prompt to prevent hanging.
			if caller != "" {
				initialize.InfoAborted(cmd)
				return nil
			}
			// Prompt for confirmation
			initialize.InfoOverwritePrompt(cmd, contextDir)
			reader := bufio.NewReader(os.Stdin)
			response, readErr := reader.ReadString(token.NewlineLF[0])
			if readErr != nil {
				return errFs.ReadInput(readErr)
			}
			response = strings.TrimSpace(strings.ToLower(response))
			if response != cli.ConfirmShort && response != cli.ConfirmLong {
				initialize.InfoAborted(cmd)
				return nil
			}
		}
	}

	// Create .context/ directory
	if mkdirErr := os.MkdirAll(contextDir, fs.PermExec); mkdirErr != nil {
		return errFs.Mkdir(contextDir, mkdirErr)
	}

	// Get the list of templates to create
	var templatesToCreate []string
	if minimal {
		templatesToCreate = ctx.FilesRequired
	} else {
		var listErr error
		templatesToCreate, listErr = catalog.List()
		if listErr != nil {
			return errPrompt.ListTemplates(listErr)
		}
	}

	// Create template files
	for _, name := range templatesToCreate {
		targetPath := filepath.Join(contextDir, name)

		// Check if the file exists and --force not set
		if _, statErr := os.Stat(targetPath); statErr == nil && !force {
			initialize.InfoExistsSkipped(cmd, name)
			continue
		}

		content, tplErr := template.Template(name)
		if tplErr != nil {
			return errPrompt.ReadTemplate(name, tplErr)
		}

		if writeErr := os.WriteFile(
			targetPath, content, fs.PermFile,
		); writeErr != nil {
			return errFs.FileWrite(targetPath, writeErr)
		}

		initialize.InfoFileCreated(cmd, name)
	}

	initialize.Initialized(cmd, contextDir)

	// Create entry templates in .context/templates/
	if tplErr := entry.CreateTemplates(cmd, contextDir, force); tplErr != nil {
		// Non-fatal: warn but continue
		label := desc.Text(text.DescKeyInitLabelEntryTemplates)
		initialize.InfoWarnNonFatal(cmd, label, tplErr)
	}

	// Set up scratchpad
	if padErr := initScratchpad(cmd, contextDir); padErr != nil {
		// Non-fatal: warn but continue
		label := desc.Text(text.DescKeyInitLabelScratchpad)
		initialize.InfoWarnNonFatal(cmd, label, padErr)
	}

	// Create project root files
	initialize.InfoCreatingRootFiles(cmd)

	// Create specs/ and ideas/ directories with README.md
	if dirsErr := coreProject.CreateDirs(cmd); dirsErr != nil {
		label := desc.Text(text.DescKeyInitLabelProjectDirs)
		initialize.InfoWarnNonFatal(cmd, label, dirsErr)
	}

	// Merge permissions into settings.local.json (no hook scaffolding)
	initialize.InfoSettingUpPermissions(cmd)
	if permsErr := coreMerge.SettingsPermissions(cmd); permsErr != nil {
		// Non-fatal: warn but continue
		label := desc.Text(text.DescKeyInitLabelPermissions)
		initialize.InfoWarnNonFatal(cmd, label, permsErr)
	}

	// Auto-enable plugin globally unless suppressed
	if !noPluginEnable {
		if pluginErr := plugin.EnableGlobally(cmd); pluginErr != nil {
			// Non-fatal: warn but continue
			label := desc.Text(text.DescKeyInitLabelPluginEnable)
			initialize.InfoWarnNonFatal(cmd, label, pluginErr)
		}
	}

	// Handle CLAUDE.md creation/merge
	if claudeErr := coreClaude.HandleMd(cmd, force, merge); claudeErr != nil {
		// Non-fatal: warn but continue
		initialize.InfoWarnNonFatal(cmd, claude.Md, claudeErr)
	}

	// Deploy Makefile.ctx and amend user Makefile
	if makeErr := coreProject.HandleMakefileCtx(cmd); makeErr != nil {
		// Non-fatal: warn but continue
		initialize.InfoWarnNonFatal(cmd, sync.PatternMakefile, makeErr)
	}

	// Update .gitignore with recommended entries
	if ignoreErr := ensureGitignoreEntries(cmd); ignoreErr != nil {
		initialize.InfoWarnNonFatal(cmd, file.FileGitignore, ignoreErr)
	}

	initialize.InfoNextSteps(cmd)
	initialize.InfoWorkflowTips(cmd)

	// Save the quick-start reference to a project-root file.
	writeGettingStarted(cmd)

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
		if _, statErr := os.Stat(mdPath); statErr != nil {
			if writeErr := os.WriteFile(mdPath, nil, fs.PermFile); writeErr != nil {
				return errFs.Mkdir(mdPath, writeErr)
			}
			initialize.InfoScratchpadPlaintext(cmd, mdPath)
		} else {
			initialize.InfoExistsSkipped(cmd, mdPath)
		}
		return nil
	}

	// Encrypted mode
	kPath := rc.KeyPath()
	encPath := filepath.Join(contextDir, pad.Enc)

	// Check if the key already exists (idempotent)
	if _, keyStatErr := os.Stat(kPath); keyStatErr == nil {
		initialize.InfoExistsSkipped(cmd, kPath)
		return nil
	}

	// Warn if the encrypted file exists but no key
	if _, encStatErr := os.Stat(encPath); encStatErr == nil {
		initialize.InfoScratchpadNoKey(cmd, kPath)
		return nil
	}

	// Ensure the key directory exists.
	if mkdirErr := os.MkdirAll(
		filepath.Dir(kPath), fs.PermKeyDir,
	); mkdirErr != nil {
		return errCrypto.MkdirKeyDir(mkdirErr)
	}

	// Generate key
	key, genErr := crypto.GenerateKey()
	if genErr != nil {
		return errCrypto.GenerateKey(genErr)
	}

	if saveErr := crypto.SaveKey(kPath, key); saveErr != nil {
		return errCrypto.SaveKey(saveErr)
	}
	initialize.InfoScratchpadKeyCreated(cmd, kPath)

	return nil
}

// hasEssentialFiles reports whether contextDir contains at least one of the
// essential context files (TASKS.md, CONSTITUTION.md, DECISIONS.md). A
// directory with only logs/ or other non-essential content is considered
// uninitialized.
//
// Parameters:
//   - contextDir: Absolute path to the context directory to inspect
//
// Returns:
//   - bool: True if at least one essential file exists
func hasEssentialFiles(contextDir string) bool {
	for _, f := range ctx.FilesRequired {
		if _, statErr := os.Stat(filepath.Join(contextDir, f)); statErr == nil {
			return true
		}
	}
	return false
}

// ensureGitignoreEntries appends recommended .gitignore entries that are not
// already present. Creates .gitignore if it does not exist.
//
// Parameters:
//   - cmd: Cobra command for status output
//
// Returns:
//   - error: Non-nil on read or write failure
func ensureGitignoreEntries(cmd *cobra.Command) error {
	content, readErr := os.ReadFile(file.FileGitignore)
	if readErr != nil && !os.IsNotExist(readErr) {
		return readErr
	}

	// Build set of existing trimmed lines.
	existing := make(map[string]bool)
	for _, line := range strings.Split(string(content), token.NewlineLF) {
		existing[strings.TrimSpace(line)] = true
	}

	// Collect missing entries.
	var missing []string
	for _, e := range file.Gitignore {
		if !existing[e] {
			missing = append(missing, e)
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
	sb.WriteString(token.NewlineLF + file.GitignoreHeader + token.NewlineLF)
	for _, e := range missing {
		sb.WriteString(e + token.NewlineLF)
	}

	if writeErr := os.WriteFile( //nolint:gosec // FileGitignore is a project-relative constant, not user-controlled
		file.FileGitignore, append(content, []byte(sb.String())...),
		fs.PermFile,
	); writeErr != nil {
		return writeErr
	}

	initialize.InfoGitignoreUpdated(cmd, len(missing))
	initialize.InfoGitignoreReview(cmd)
	return nil
}

// writeGettingStarted saves the next-steps and workflow-tips text to
// GETTING_STARTED.md in the project root. Best-effort: failures are
// non-fatal since the same content was already printed to stdout.
//
// Parameters:
//   - cmd: Cobra command for status output
func writeGettingStarted(cmd *cobra.Command) {
	content := desc.Text(text.DescKeyWriteInitNextStepsBlock) +
		token.NewlineLF +
		desc.Text(text.DescKeyWriteInitWorkflowTips) +
		token.NewlineLF
	if writeErr := os.WriteFile(
		project.GettingStarted, []byte(content), fs.PermFile,
	); writeErr != nil {
		return
	}
	initialize.InfoGettingStartedSaved(cmd, project.GettingStarted)
}
