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
	"github.com/ActiveMemory/ctx/internal/cli/initialize/core/backup"
	coreClaude "github.com/ActiveMemory/ctx/internal/cli/initialize/core/claude"
	coreCC "github.com/ActiveMemory/ctx/internal/cli/initialize/core/claudecheck"
	"github.com/ActiveMemory/ctx/internal/cli/initialize/core/entry"
	coreKB "github.com/ActiveMemory/ctx/internal/cli/initialize/core/kb"
	coreMerge "github.com/ActiveMemory/ctx/internal/cli/initialize/core/merge"
	corePad "github.com/ActiveMemory/ctx/internal/cli/initialize/core/pad"
	"github.com/ActiveMemory/ctx/internal/cli/initialize/core/plugin"
	coreProject "github.com/ActiveMemory/ctx/internal/cli/initialize/core/project"
	"github.com/ActiveMemory/ctx/internal/cli/initialize/core/validate"
	steeringInit "github.com/ActiveMemory/ctx/internal/cli/steering/cmd/initcmd"
	"github.com/ActiveMemory/ctx/internal/config/claude"
	"github.com/ActiveMemory/ctx/internal/config/cli"
	"github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/config/sync"
	"github.com/ActiveMemory/ctx/internal/config/token"
	errFs "github.com/ActiveMemory/ctx/internal/err/fs"
	errInit "github.com/ActiveMemory/ctx/internal/err/initialize"
	errPrompt "github.com/ActiveMemory/ctx/internal/err/prompt"
	"github.com/ActiveMemory/ctx/internal/i18n"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/write/initialize"
)

// Run executes the init command logic.
//
// Creates a .context/ directory with template files. Handles existing
// directories, minimal mode, and CLAUDE.md merge operations.
//
// Under the cwd-anchored resolution model
// (spec: specs/cwd-anchored-context.md), init is the command that
// *creates* `$PWD/.context/`. It does not consult any env var and
// it does not walk: the target is always `$PWD/.context/`. Users
// who want to init a different project `cd` there first.
//
// # Existing-context handling
//
// When `$PWD/.context/` already contains a populated context
// (any file in ctx.FilesRequired exists), behavior depends on
// --reset:
//
//   - Without --reset: refuse with errInit.Populated, listing the
//     populated files and pointing at --reset for recovery. The
//     directory is left untouched.
//   - With --reset and a non-interactive caller (--caller set):
//     refuse with errInit.ResetRequiresInteractive. Reset is
//     destructive and must come from a real terminal session.
//   - With --reset interactively: enumerate the populated files,
//     prompt y/N, and on confirmation copy them to
//     .context/.backup-init-<UTC-ISO>/ before scaffolding.
//
// When the directory is missing or only contains state/ / hook
// scratch (no essential files), init scaffolds the missing
// templates as before.
//
// Spec: specs/ctx-init-overwrite-safety.md.
//
// Parameters:
//   - cmd: Cobra command for output and input streams
//   - reset: If true, attempt destructive reset of an existing
//     populated context (interactive only; backs up first)
//   - minimal: If true, only create essential files
//   - merge: If true, auto-merge ctx content into existing files
//   - noPluginEnable: If true, skip auto-enabling the plugin globally
//   - noSteeringInit: If true, skip scaffolding foundation steering
//     files in .context/steering/
//   - caller: Identifies the calling tool
//     (e.g. "vscode") for template overrides
//
// Returns:
//   - error: Non-nil if refusal triggers, directory creation fails,
//     or file operations fail
func Run(
	cmd *cobra.Command,
	reset, minimal, merge, noPluginEnable, noSteeringInit bool,
	caller string,
) error {
	// Check if ctx is in PATH (required for hooks to work).
	// Skip when a caller is set: the caller manages its own binary path.
	if caller == "" {
		if pathErr := validate.CheckCtxInPath(cmd); pathErr != nil {
			return pathErr
		}
	}

	// Under the cwd-anchored resolution model
	// (spec: specs/cwd-anchored-context.md), init always targets
	// `$PWD/.context/`. No env-var lookup, no fallback dance.
	cwd, cwdErr := os.Getwd()
	if cwdErr != nil {
		return errFs.ReadInput(cwdErr)
	}
	contextDir := filepath.Join(cwd, dir.Context)

	// Existing-context handling: refuse by default; --reset takes a
	// backup and only proceeds on interactive y/N confirmation.
	if _, statErr := os.Stat(contextDir); statErr == nil {
		populated := validate.PopulatedFiles(contextDir)
		if len(populated) > 0 {
			if !reset {
				cmd.SilenceUsage = true
				return errInit.Populated(contextDir, populated)
			}
			if caller != "" {
				cmd.SilenceUsage = true
				return errInit.ResetRequiresInteractive()
			}
			initialize.InfoResetPrompt(cmd, contextDir, populated)
			reader := bufio.NewReader(cmd.InOrStdin())
			response, readErr := reader.ReadString(token.NewlineLF[0])
			if readErr != nil {
				return errFs.ReadInput(readErr)
			}
			response = strings.TrimSpace(i18n.Fold(response))
			if response != cli.ConfirmShort && response != cli.ConfirmLong {
				initialize.InfoAborted(cmd)
				return nil
			}
			backupDir, backupErr := backup.WriteSnapshot(contextDir, populated)
			if backupErr != nil {
				return backupErr
			}
			initialize.InfoBackupWritten(cmd, backupDir)
		}
	}

	// Create .context/ directory
	if mkdirErr := ctxIo.SafeMkdirAll(contextDir, fs.PermExec); mkdirErr != nil {
		return errFs.Mkdir(contextDir, mkdirErr)
	}

	// Create .context/ subdirectories for steering, hooks, and skills.
	// Uses SafeMkdirAll which is a no-op when the directory already exists.
	for _, sub := range []string{dir.Steering, dir.Hooks, dir.Skills} {
		subDir := filepath.Join(contextDir, sub)
		if mkdirErr := ctxIo.SafeMkdirAll(subDir, fs.PermExec); mkdirErr != nil {
			return errFs.Mkdir(subDir, mkdirErr)
		}
	}

	// Scaffold foundation steering files so the directory
	// isn't a confusing empty dir. The initcmd package
	// skips files that already exist, so re-running `ctx
	// init` after the user has edited a foundation file
	// won't clobber their work. Honors --no-steering-init
	// for users who want a bare init with no starter
	// templates.
	if !noSteeringInit {
		steeringErr := steeringInit.RunWithDir(cmd, contextDir)
		if steeringErr != nil {
			// Non-fatal: the rest of init is more
			// important than the steering templates.
			label := desc.Text(text.DescKeyInitLabelSteering)
			initialize.InfoWarnNonFatal(cmd, label, steeringErr)
		}
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

		// Check if the file exists and --reset not set
		if _, statErr := os.Stat(targetPath); statErr == nil && !reset {
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
	if tplErr := entry.CreateTemplates(cmd, contextDir, reset); tplErr != nil {
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

	// Phase KB: scaffold the editorial-pipeline directories
	// and copy embedded templates (KB-RULES, mode prompts,
	// schemas, kb landing). Per-file existence is preserved.
	if kbErr := coreKB.Scaffold(contextDir); kbErr != nil {
		// Non-fatal: warn but continue. The KB surfaces are
		// opt-in for users who run /ctx-kb-ingest et al; a
		// scaffold failure should not block init.
		kbLabel := desc.Text(text.DescKeyInitLabelKB)
		initialize.InfoWarnNonFatal(cmd, kbLabel, kbErr)
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

	// Merge the ctx status line into settings.local.json, backing up
	// any pre-existing statusLine entry (specs/statusline.md).
	if statuslineErr := coreMerge.SettingsStatusLine(cmd); statuslineErr != nil {
		// Non-fatal: warn but continue
		label := desc.Text(text.DescKeyInitLabelStatusline)
		initialize.InfoWarnNonFatal(cmd, label, statuslineErr)
	}

	// Auto-enable plugin globally and locally unless suppressed.
	if !noPluginEnable {
		if pluginErr := plugin.EnableGlobally(cmd); pluginErr != nil {
			label := desc.Text(text.DescKeyInitLabelPluginEnable)
			initialize.InfoWarnNonFatal(cmd, label, pluginErr)
		}
		if localErr := plugin.EnableLocally(cmd); localErr != nil {
			label := desc.Text(text.DescKeyInitLabelPluginEnable)
			initialize.InfoWarnNonFatal(cmd, label, localErr)
		}
	}

	// Handle CLAUDE.md creation/merge
	if claudeErr := coreClaude.HandleMd(cmd, reset, merge); claudeErr != nil {
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
	coreProject.WriteGettingStarted(cmd, contextDir)

	// Post-script: stage-aware Claude Code setup guidance.
	// Never fatal, never an error; a friendly nudge
	// pointing the user at whichever step they're missing.
	// Honors --no-plugin-enable: if plugin detection was
	// suppressed, skip the hint too.
	if !noPluginEnable {
		coreCC.InitHint(cmd)
	}

	return nil
}
