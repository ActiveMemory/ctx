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
	corePad "github.com/ActiveMemory/ctx/internal/cli/initialize/core/pad"
	"github.com/ActiveMemory/ctx/internal/cli/initialize/core/plugin"
	coreProject "github.com/ActiveMemory/ctx/internal/cli/initialize/core/project"
	"github.com/ActiveMemory/ctx/internal/cli/initialize/core/validate"
	"github.com/ActiveMemory/ctx/internal/config/claude"
	"github.com/ActiveMemory/ctx/internal/config/cli"
	"github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/config/sync"
	"github.com/ActiveMemory/ctx/internal/config/token"
	errFs "github.com/ActiveMemory/ctx/internal/err/fs"
	errPrompt "github.com/ActiveMemory/ctx/internal/err/prompt"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
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
		if !force && validate.EssentialFilesPresent(contextDir) {
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
	if mkdirErr := ctxIo.SafeMkdirAll(contextDir, fs.PermExec); mkdirErr != nil {
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

		if writeErr := ctxIo.SafeWriteFile(
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
	if padErr := corePad.Setup(cmd, contextDir); padErr != nil {
		// Non-fatal: warn but continue
		label := desc.Text(text.DescKeyInitLabelScratchpad)
		initialize.InfoWarnNonFatal(cmd, label, padErr)
	}

	// Create project root files
	initialize.InfoCreatingRootFiles(cmd)

	// Create specs/ and ideas/ directories with README.md
	if dirsErr := coreProject.CreateDirs(cmd); dirsErr != nil {
		// Non-fatal: warn but continue
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
	if ignoreErr := coreProject.EnsureGitignoreEntries(cmd); ignoreErr != nil {
		// Non-fatal: warn but continue
		initialize.InfoWarnNonFatal(cmd, file.FileGitignore, ignoreErr)
	}

	initialize.InfoNextSteps(cmd)
	initialize.InfoWorkflowTips(cmd)

	// Save the quick-start reference to a project-root file.
	coreProject.WriteGettingStarted(cmd)

	return nil
}
