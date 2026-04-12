//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package initialize

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// Created reports a file created during init.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: created file path
func Created(cmd *cobra.Command, path string) {
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWriteInitFileCreated), path))
}

// Skipped reports a file skipped because it already exists.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: skipped file path
func Skipped(cmd *cobra.Command, path string) {
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWriteInitExistsSkipped), path))
}

// SkippedPlain reports a file skipped without detail.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: skipped file path
func SkippedPlain(cmd *cobra.Command, path string) {
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWriteInitSkippedPlain), path))
}

// CtxContentExists reports a file skipped because ctx content exists.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: skipped file path
func CtxContentExists(cmd *cobra.Command, path string) {
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteInitCtxContentExists),
		path))
}

// Merged reports a file merged during init.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: merged file path
func Merged(cmd *cobra.Command, path string) {
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWriteInitMerged), path))
}

// Backup reports a backup file created.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: backup file path
func Backup(cmd *cobra.Command, path string) {
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWriteInitBackup), path))
}

// UpdatedSection reports a file whose marked section was updated.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: updated file path
//   - textKey: text descriptor key for the update message
func UpdatedSection(cmd *cobra.Command, path, textKey string) {
	cmd.Println(fmt.Sprintf(desc.Text(textKey), path))
}

// FileExistsNoCtx reports a file exists without ctx content.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: file path
func FileExistsNoCtx(cmd *cobra.Command, path string) {
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWriteInitFileExistsNoCtx), path))
}

// NoChanges reports a settings file with no changes needed.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: settings file path
func NoChanges(cmd *cobra.Command, path string) {
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWriteInitNoChanges), path))
}

// PermsMergedDeduped reports permissions merged and deduped.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: settings file path
func PermsMergedDeduped(cmd *cobra.Command, path string) {
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteInitPermsMergedDeduped),
		path))
}

// PermsDeduped reports duplicate permissions removed.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: settings file path
func PermsDeduped(cmd *cobra.Command, path string) {
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWriteInitPermsDeduped), path))
}

// PermsAllowDeny reports allow+deny permissions added.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: settings file path
func PermsAllowDeny(cmd *cobra.Command, path string) {
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWriteInitPermsAllowDeny), path))
}

// PermsDeny reports deny permissions added.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: settings file path
func PermsDeny(cmd *cobra.Command, path string) {
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWriteInitPermsDeny), path))
}

// PermsAllow reports ctx permissions added.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: settings file path
func PermsAllow(cmd *cobra.Command, path string) {
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWriteInitPermsAllow), path))
}

// MakefileCreated reports a new Makefile created with ctx include.
//
// Parameters:
//   - cmd: Cobra command for output
func MakefileCreated(cmd *cobra.Command) {
	cmd.Println(desc.Text(text.DescKeyWriteInitMakefileCreated))
}

// MakefileIncludes reports Makefile already includes the directive.
//
// Parameters:
//   - cmd: Cobra command for output
//   - filename: included filename
func MakefileIncludes(cmd *cobra.Command, filename string) {
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteInitMakefileIncludes),
		filename))
}

// MakefileAppended reports an include appended to Makefile.
//
// Parameters:
//   - cmd: Cobra command for output
//   - filename: included filename
func MakefileAppended(cmd *cobra.Command, filename string) {
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteInitMakefileAppended),
		filename))
}

// PluginSkipped reports plugin enablement was skipped.
//
// Parameters:
//   - cmd: Cobra command for output
func PluginSkipped(cmd *cobra.Command) {
	cmd.Println(desc.Text(text.DescKeyWriteInitPluginSkipped))
}

// PluginAlreadyEnabled reports plugin is already enabled globally.
//
// Parameters:
//   - cmd: Cobra command for output
func PluginAlreadyEnabled(cmd *cobra.Command) {
	cmd.Println(desc.Text(text.DescKeyWriteInitPluginAlreadyEnabled))
}

// PluginEnabled reports plugin enabled globally.
//
// Parameters:
//   - cmd: Cobra command for output
//   - settingsPath: path to the settings file
func PluginEnabled(cmd *cobra.Command, settingsPath string) {
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteInitPluginEnabled),
		settingsPath))
}

// PluginLocalAlreadyEnabled reports plugin is already enabled
// in project settings.
//
// Parameters:
//   - cmd: Cobra command for output
func PluginLocalAlreadyEnabled(cmd *cobra.Command) {
	cmd.Println(desc.Text(
		text.DescKeyWriteInitPluginLocalAlreadyEnabled))
}

// PluginLocalEnabled reports plugin enabled in project
// settings.
//
// Parameters:
//   - cmd: Cobra command for output
//   - settingsPath: path to the local settings file
func PluginLocalEnabled(
	cmd *cobra.Command, settingsPath string,
) {
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteInitPluginLocalEnabled),
		settingsPath))
}

// ClaudeAbsent prints the guidance shown when the `claude`
// binary is not on PATH. Walks the user through installing
// Claude Code first, then the ctx plugin.
//
// Parameters:
//   - cmd: Cobra command for output
func ClaudeAbsent(cmd *cobra.Command) {
	cmd.Println(desc.Text(text.DescKeyWriteInitClaudeAbsent))
}

// ClaudePluginMissing prints the guidance shown when
// `claude` is present but the ctx plugin is not installed.
// Walks the user through the dev-symlink install flow with
// user-scope guidance to avoid the local-install enablement
// gotcha.
//
// Parameters:
//   - cmd: Cobra command for output
func ClaudePluginMissing(cmd *cobra.Command) {
	cmd.Println(desc.Text(
		text.DescKeyWriteInitClaudePluginMissing))
}

// ClaudeReady prints the multi-line confirmation for the
// init post-script when Claude Code and the ctx plugin are
// both detected and enabled. Displays scope, version,
// source type, clone path, and enablement flags so the user
// can quickly see which clone feeds their plugin.
//
// Callers should fall back to [ClaudeReadyMinimal] if
// [claude_check.Details] couldn't parse plugin metadata.
//
// Parameters:
//   - cmd: Cobra command for output
//   - scope: installation scope (e.g. "user")
//   - version: plugin version with short git SHA
//   - source: marketplace source type
//     (e.g. "directory (hot-reload)")
//   - clonePath: filesystem path to the source clone
//   - enabled: enablement summary
//     (e.g. "global + project")
func ClaudeReady(
	cmd *cobra.Command,
	scope, version, source, clonePath, enabled string,
) {
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteInitClaudeReady),
		scope, version, source, clonePath, enabled,
	))
}

// ClaudeReadyMinimal prints the fallback one-line
// confirmation used when [claude_check.Details] couldn't
// parse plugin metadata from the Claude Code plugin
// registry files. Never fails.
//
// Parameters:
//   - cmd: Cobra command for output
func ClaudeReadyMinimal(cmd *cobra.Command) {
	cmd.Println(desc.Text(
		text.DescKeyWriteInitClaudeReadyMinimal))
}

// SetupClaudeReady prints the `ctx setup claude-code`
// success block with the detailed plugin metadata. Tells
// the user the plugin is installed and enabled, shows the
// scope/version/source/clone path/enablement, and
// redirects them to run `ctx init` in a project.
//
// Parameters:
//   - cmd: Cobra command for output
//   - scope: installation scope
//   - version: plugin version with short git SHA
//   - source: marketplace source type
//   - clonePath: filesystem path to the source clone
//   - enabled: enablement summary
func SetupClaudeReady(
	cmd *cobra.Command,
	scope, version, source, clonePath, enabled string,
) {
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteSetupClaudeReady),
		scope, version, source, clonePath, enabled,
	))
}

// SetupClaudeReadyMinimal prints the fallback setup
// success block when plugin metadata can't be parsed.
//
// Parameters:
//   - cmd: Cobra command for output
func SetupClaudeReadyMinimal(cmd *cobra.Command) {
	cmd.Println(desc.Text(
		text.DescKeyWriteSetupClaudeReadyMinimal))
}

// SkippedDir reports a directory skipped because it exists.
//
// Parameters:
//   - cmd: Cobra command for output
//   - dir: directory name
func SkippedDir(cmd *cobra.Command, dir string) {
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWriteInitSkippedDir), dir))
}

// CreatedDir reports a directory created during init.
//
// Parameters:
//   - cmd: Cobra command for output
//   - dir: directory name
func CreatedDir(cmd *cobra.Command, dir string) {
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWriteInitCreatedDir), dir))
}

// NotInPath prints a multi-line diagnostic to stderr
// explaining that ctx is not in PATH, with installation
// instructions.
//
// Parameters:
//   - cmd: Cobra command whose stderr stream receives the output.
//     Nil is a no-op.
func NotInPath(cmd *cobra.Command) {
	if cmd == nil {
		return
	}

	cmd.PrintErrln(desc.Text(text.DescKeyErrInitCtxNotInPath))
}

// MergePrompt prints a merge confirmation prompt with [y/N] suffix.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - prompt: The confirmation prompt text.
func MergePrompt(cmd *cobra.Command, prompt string) {
	if cmd == nil {
		return
	}
	cmd.Println(prompt)
	cmd.Print(desc.Text(text.DescKeyConfirmProceed))
}
