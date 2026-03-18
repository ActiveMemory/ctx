//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package initialize

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/spf13/cobra"
)

// Created reports a file created during init.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: created file path
func Created(cmd *cobra.Command, path string) {
	cmd.Println(fmt.Sprintf(assets.TextDesc(assets.TextDescKeyWriteInitFileCreated), path))
}

// CreatedWith reports a file created with a qualifier (e.g. " (ralph mode)").
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: created file path
//   - qualifier: additional info appended after the path
func CreatedWith(cmd *cobra.Command, path, qualifier string) {
	cmd.Println(fmt.Sprintf(assets.TextDesc(assets.TextDescKeyWriteInitCreatedWith), path, qualifier))
}

// Skipped reports a file skipped because it already exists.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: skipped file path
func Skipped(cmd *cobra.Command, path string) {
	cmd.Println(fmt.Sprintf(assets.TextDesc(assets.TextDescKeyWriteInitExistsSkipped), path))
}

// SkippedPlain reports a file skipped without detail.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: skipped file path
func SkippedPlain(cmd *cobra.Command, path string) {
	cmd.Println(fmt.Sprintf(assets.TextDesc(assets.TextDescKeyWriteInitSkippedPlain), path))
}

// CtxContentExists reports a file skipped because ctx content exists.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: skipped file path
func CtxContentExists(cmd *cobra.Command, path string) {
	cmd.Println(fmt.Sprintf(assets.TextDesc(assets.TextDescKeyWriteInitCtxContentExists), path))
}

// Merged reports a file merged during init.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: merged file path
func Merged(cmd *cobra.Command, path string) {
	cmd.Println(fmt.Sprintf(assets.TextDesc(assets.TextDescKeyWriteInitMerged), path))
}

// Backup reports a backup file created.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: backup file path
func Backup(cmd *cobra.Command, path string) {
	cmd.Println(fmt.Sprintf(assets.TextDesc(assets.TextDescKeyWriteInitBackup), path))
}

// UpdatedCtxSection reports a file whose ctx section was updated.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: updated file path
func UpdatedCtxSection(cmd *cobra.Command, path string) {
	cmd.Println(fmt.Sprintf(assets.TextDesc(assets.TextDescKeyWriteInitUpdatedCtxSection), path))
}

// UpdatedPlanSection reports a file whose plan section was updated.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: updated file path
func UpdatedPlanSection(cmd *cobra.Command, path string) {
	cmd.Println(fmt.Sprintf(assets.TextDesc(assets.TextDescKeyWriteInitUpdatedPlanSection), path))
}

// UpdatedPromptSection reports a file whose prompt section was updated.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: updated file path
func UpdatedPromptSection(cmd *cobra.Command, path string) {
	cmd.Println(fmt.Sprintf(assets.TextDesc(assets.TextDescKeyWriteInitUpdatedPromptSection), path))
}

// FileExistsNoCtx reports a file exists without ctx content.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: file path
func FileExistsNoCtx(cmd *cobra.Command, path string) {
	cmd.Println(fmt.Sprintf(assets.TextDesc(assets.TextDescKeyWriteInitFileExistsNoCtx), path))
}

// NoChanges reports a settings file with no changes needed.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: settings file path
func NoChanges(cmd *cobra.Command, path string) {
	cmd.Println(fmt.Sprintf(assets.TextDesc(assets.TextDescKeyWriteInitNoChanges), path))
}

// PermsMergedDeduped reports permissions merged and deduped.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: settings file path
func PermsMergedDeduped(cmd *cobra.Command, path string) {
	cmd.Println(fmt.Sprintf(assets.TextDesc(assets.TextDescKeyWriteInitPermsMergedDeduped), path))
}

// PermsDeduped reports duplicate permissions removed.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: settings file path
func PermsDeduped(cmd *cobra.Command, path string) {
	cmd.Println(fmt.Sprintf(assets.TextDesc(assets.TextDescKeyWriteInitPermsDeduped), path))
}

// PermsAllowDeny reports allow+deny permissions added.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: settings file path
func PermsAllowDeny(cmd *cobra.Command, path string) {
	cmd.Println(fmt.Sprintf(assets.TextDesc(assets.TextDescKeyWriteInitPermsAllowDeny), path))
}

// PermsDeny reports deny permissions added.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: settings file path
func PermsDeny(cmd *cobra.Command, path string) {
	cmd.Println(fmt.Sprintf(assets.TextDesc(assets.TextDescKeyWriteInitPermsDeny), path))
}

// PermsAllow reports ctx permissions added.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: settings file path
func PermsAllow(cmd *cobra.Command, path string) {
	cmd.Println(fmt.Sprintf(assets.TextDesc(assets.TextDescKeyWriteInitPermsAllow), path))
}

// MakefileCreated reports a new Makefile created with ctx include.
//
// Parameters:
//   - cmd: Cobra command for output
func MakefileCreated(cmd *cobra.Command) {
	cmd.Println(assets.TextDesc(assets.TextDescKeyWriteInitMakefileCreated))
}

// MakefileIncludes reports Makefile already includes the directive.
//
// Parameters:
//   - cmd: Cobra command for output
//   - filename: included filename
func MakefileIncludes(cmd *cobra.Command, filename string) {
	cmd.Println(fmt.Sprintf(assets.TextDesc(assets.TextDescKeyWriteInitMakefileIncludes), filename))
}

// MakefileAppended reports an include appended to Makefile.
//
// Parameters:
//   - cmd: Cobra command for output
//   - filename: included filename
func MakefileAppended(cmd *cobra.Command, filename string) {
	cmd.Println(fmt.Sprintf(assets.TextDesc(assets.TextDescKeyWriteInitMakefileAppended), filename))
}

// PluginSkipped reports plugin enablement was skipped.
//
// Parameters:
//   - cmd: Cobra command for output
func PluginSkipped(cmd *cobra.Command) {
	cmd.Println(assets.TextDesc(assets.TextDescKeyWriteInitPluginSkipped))
}

// PluginAlreadyEnabled reports plugin is already enabled globally.
//
// Parameters:
//   - cmd: Cobra command for output
func PluginAlreadyEnabled(cmd *cobra.Command) {
	cmd.Println(assets.TextDesc(assets.TextDescKeyWriteInitPluginAlreadyEnabled))
}

// PluginEnabled reports plugin enabled globally.
//
// Parameters:
//   - cmd: Cobra command for output
//   - settingsPath: path to the settings file
func PluginEnabled(cmd *cobra.Command, settingsPath string) {
	cmd.Println(fmt.Sprintf(assets.TextDesc(assets.TextDescKeyWriteInitPluginEnabled), settingsPath))
}

// SkippedDir reports a directory skipped because it exists.
//
// Parameters:
//   - cmd: Cobra command for output
//   - dir: directory name
func SkippedDir(cmd *cobra.Command, dir string) {
	cmd.Println(fmt.Sprintf(assets.TextDesc(assets.TextDescKeyWriteInitSkippedDir), dir))
}

// CreatedDir reports a directory created during init.
//
// Parameters:
//   - cmd: Cobra command for output
//   - dir: directory name
func CreatedDir(cmd *cobra.Command, dir string) {
	cmd.Println(fmt.Sprintf(assets.TextDesc(assets.TextDescKeyWriteInitCreatedDir), dir))
}
