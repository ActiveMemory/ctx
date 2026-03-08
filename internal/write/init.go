//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package write

import (
	"github.com/ActiveMemory/ctx/internal/write/config"
	"github.com/ActiveMemory/ctx/internal/write/io"
	"github.com/spf13/cobra"
)

// InitCreated reports a file created during init.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: created file path
func InitCreated(cmd *cobra.Command, path string) {
	io.sprintf(cmd, config.tplInitFileCreated, path)
}

// InitCreatedWith reports a file created with a qualifier (e.g. " (ralph mode)").
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: created file path
//   - qualifier: additional info appended after the path
func InitCreatedWith(cmd *cobra.Command, path, qualifier string) {
	io.sprintf(cmd, config.tplInitCreatedWith, path, qualifier)
}

// InitSkipped reports a file skipped because it already exists.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: skipped file path
func InitSkipped(cmd *cobra.Command, path string) {
	io.sprintf(cmd, config.tplInitExistsSkipped, path)
}

// InitSkippedPlain reports a file skipped without detail.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: skipped file path
func InitSkippedPlain(cmd *cobra.Command, path string) {
	io.sprintf(cmd, config.tplInitSkippedPlain, path)
}

// InitCtxContentExists reports a file skipped because ctx content exists.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: skipped file path
func InitCtxContentExists(cmd *cobra.Command, path string) {
	io.sprintf(cmd, config.tplInitCtxContentExists, path)
}

// InitMerged reports a file merged during init.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: merged file path
func InitMerged(cmd *cobra.Command, path string) {
	io.sprintf(cmd, config.tplInitMerged, path)
}

// InitBackup reports a backup file created.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: backup file path
func InitBackup(cmd *cobra.Command, path string) {
	io.sprintf(cmd, config.tplInitBackup, path)
}

// InitUpdatedCtxSection reports a file whose ctx section was updated.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: updated file path
func InitUpdatedCtxSection(cmd *cobra.Command, path string) {
	io.sprintf(cmd, config.tplInitUpdatedCtxSection, path)
}

// InitUpdatedPlanSection reports a file whose plan section was updated.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: updated file path
func InitUpdatedPlanSection(cmd *cobra.Command, path string) {
	io.sprintf(cmd, config.tplInitUpdatedPlanSection, path)
}

// InitUpdatedPromptSection reports a file whose prompt section was updated.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: updated file path
func InitUpdatedPromptSection(cmd *cobra.Command, path string) {
	io.sprintf(cmd, config.tplInitUpdatedPromptSection, path)
}

// InitFileExistsNoCtx reports a file exists without ctx content.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: file path
func InitFileExistsNoCtx(cmd *cobra.Command, path string) {
	io.sprintf(cmd, config.tplInitFileExistsNoCtx, path)
}

// InitNoChanges reports a settings file with no changes needed.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: settings file path
func InitNoChanges(cmd *cobra.Command, path string) {
	io.sprintf(cmd, config.tplInitNoChanges, path)
}

// InitPermsMergedDeduped reports permissions merged and deduped.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: settings file path
func InitPermsMergedDeduped(cmd *cobra.Command, path string) {
	io.sprintf(cmd, config.tplInitPermsMergedDeduped, path)
}

// InitPermsDeduped reports duplicate permissions removed.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: settings file path
func InitPermsDeduped(cmd *cobra.Command, path string) {
	io.sprintf(cmd, config.tplInitPermsDeduped, path)
}

// InitPermsAllowDeny reports allow+deny permissions added.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: settings file path
func InitPermsAllowDeny(cmd *cobra.Command, path string) {
	io.sprintf(cmd, config.tplInitPermsAllowDeny, path)
}

// InitPermsDeny reports deny permissions added.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: settings file path
func InitPermsDeny(cmd *cobra.Command, path string) {
	io.sprintf(cmd, config.tplInitPermsDeny, path)
}

// InitPermsAllow reports ctx permissions added.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: settings file path
func InitPermsAllow(cmd *cobra.Command, path string) {
	io.sprintf(cmd, config.tplInitPermsAllow, path)
}

// InitMakefileCreated reports a new Makefile created with ctx include.
//
// Parameters:
//   - cmd: Cobra command for output
func InitMakefileCreated(cmd *cobra.Command) {
	cmd.Println(config.tplInitMakefileCreated)
}

// InitMakefileIncludes reports Makefile already includes the directive.
//
// Parameters:
//   - cmd: Cobra command for output
//   - filename: included filename
func InitMakefileIncludes(cmd *cobra.Command, filename string) {
	io.sprintf(cmd, config.tplInitMakefileIncludes, filename)
}

// InitMakefileAppended reports an include appended to Makefile.
//
// Parameters:
//   - cmd: Cobra command for output
//   - filename: included filename
func InitMakefileAppended(cmd *cobra.Command, filename string) {
	io.sprintf(cmd, config.tplInitMakefileAppended, filename)
}

// InitPluginSkipped reports plugin enablement was skipped.
//
// Parameters:
//   - cmd: Cobra command for output
func InitPluginSkipped(cmd *cobra.Command) {
	cmd.Println(config.tplInitPluginSkipped)
}

// InitPluginAlreadyEnabled reports plugin is already enabled globally.
//
// Parameters:
//   - cmd: Cobra command for output
func InitPluginAlreadyEnabled(cmd *cobra.Command) {
	cmd.Println(config.tplInitPluginAlreadyEnabled)
}

// InitPluginEnabled reports plugin enabled globally.
//
// Parameters:
//   - cmd: Cobra command for output
//   - settingsPath: path to the settings file
func InitPluginEnabled(cmd *cobra.Command, settingsPath string) {
	io.sprintf(cmd, config.tplInitPluginEnabled, settingsPath)
}

// InitSkippedDir reports a directory skipped because it exists.
//
// Parameters:
//   - cmd: Cobra command for output
//   - dir: directory name
func InitSkippedDir(cmd *cobra.Command, dir string) {
	io.sprintf(cmd, config.tplInitSkippedDir, dir)
}

// InitCreatedDir reports a directory created during init.
//
// Parameters:
//   - cmd: Cobra command for output
//   - dir: directory name
func InitCreatedDir(cmd *cobra.Command, dir string) {
	io.sprintf(cmd, config.tplInitCreatedDir, dir)
}
