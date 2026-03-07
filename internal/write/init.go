//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package write

import (
	"github.com/spf13/cobra"
)

// InitCreated reports a file created during init.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: created file path
func InitCreated(cmd *cobra.Command, path string) {
	sprintf(cmd, tplInitFileCreated, path)
}

// InitCreatedWith reports a file created with a qualifier (e.g. " (ralph mode)").
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: created file path
//   - qualifier: additional info appended after the path
func InitCreatedWith(cmd *cobra.Command, path, qualifier string) {
	sprintf(cmd, tplInitCreatedWith, path, qualifier)
}

// InitSkipped reports a file skipped because it already exists.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: skipped file path
func InitSkipped(cmd *cobra.Command, path string) {
	sprintf(cmd, tplInitExistsSkipped, path)
}

// InitSkippedPlain reports a file skipped without detail.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: skipped file path
func InitSkippedPlain(cmd *cobra.Command, path string) {
	sprintf(cmd, tplInitSkippedPlain, path)
}

// InitCtxContentExists reports a file skipped because ctx content exists.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: skipped file path
func InitCtxContentExists(cmd *cobra.Command, path string) {
	sprintf(cmd, tplInitCtxContentExists, path)
}

// InitMerged reports a file merged during init.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: merged file path
func InitMerged(cmd *cobra.Command, path string) {
	sprintf(cmd, tplInitMerged, path)
}

// InitBackup reports a backup file created.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: backup file path
func InitBackup(cmd *cobra.Command, path string) {
	sprintf(cmd, tplInitBackup, path)
}

// InitUpdatedCtxSection reports a file whose ctx section was updated.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: updated file path
func InitUpdatedCtxSection(cmd *cobra.Command, path string) {
	sprintf(cmd, tplInitUpdatedCtxSection, path)
}

// InitUpdatedPlanSection reports a file whose plan section was updated.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: updated file path
func InitUpdatedPlanSection(cmd *cobra.Command, path string) {
	sprintf(cmd, tplInitUpdatedPlanSection, path)
}

// InitUpdatedPromptSection reports a file whose prompt section was updated.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: updated file path
func InitUpdatedPromptSection(cmd *cobra.Command, path string) {
	sprintf(cmd, tplInitUpdatedPromptSection, path)
}

// InitFileExistsNoCtx reports a file exists without ctx content.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: file path
func InitFileExistsNoCtx(cmd *cobra.Command, path string) {
	sprintf(cmd, tplInitFileExistsNoCtx, path)
}

// InitNoChanges reports a settings file with no changes needed.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: settings file path
func InitNoChanges(cmd *cobra.Command, path string) {
	sprintf(cmd, tplInitNoChanges, path)
}

// InitPermsMergedDeduped reports permissions merged and deduped.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: settings file path
func InitPermsMergedDeduped(cmd *cobra.Command, path string) {
	sprintf(cmd, tplInitPermsMergedDeduped, path)
}

// InitPermsDeduped reports duplicate permissions removed.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: settings file path
func InitPermsDeduped(cmd *cobra.Command, path string) {
	sprintf(cmd, tplInitPermsDeduped, path)
}

// InitPermsAllowDeny reports allow+deny permissions added.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: settings file path
func InitPermsAllowDeny(cmd *cobra.Command, path string) {
	sprintf(cmd, tplInitPermsAllowDeny, path)
}

// InitPermsDeny reports deny permissions added.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: settings file path
func InitPermsDeny(cmd *cobra.Command, path string) {
	sprintf(cmd, tplInitPermsDeny, path)
}

// InitPermsAllow reports ctx permissions added.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: settings file path
func InitPermsAllow(cmd *cobra.Command, path string) {
	sprintf(cmd, tplInitPermsAllow, path)
}

// InitMakefileCreated reports a new Makefile created with ctx include.
//
// Parameters:
//   - cmd: Cobra command for output
func InitMakefileCreated(cmd *cobra.Command) {
	cmd.Println(tplInitMakefileCreated)
}

// InitMakefileIncludes reports Makefile already includes the directive.
//
// Parameters:
//   - cmd: Cobra command for output
//   - filename: included filename
func InitMakefileIncludes(cmd *cobra.Command, filename string) {
	sprintf(cmd, tplInitMakefileIncludes, filename)
}

// InitMakefileAppended reports an include appended to Makefile.
//
// Parameters:
//   - cmd: Cobra command for output
//   - filename: included filename
func InitMakefileAppended(cmd *cobra.Command, filename string) {
	sprintf(cmd, tplInitMakefileAppended, filename)
}

// InitPluginSkipped reports plugin enablement was skipped.
//
// Parameters:
//   - cmd: Cobra command for output
func InitPluginSkipped(cmd *cobra.Command) {
	cmd.Println(tplInitPluginSkipped)
}

// InitPluginAlreadyEnabled reports plugin is already enabled globally.
//
// Parameters:
//   - cmd: Cobra command for output
func InitPluginAlreadyEnabled(cmd *cobra.Command) {
	cmd.Println(tplInitPluginAlreadyEnabled)
}

// InitPluginEnabled reports plugin enabled globally.
//
// Parameters:
//   - cmd: Cobra command for output
//   - settingsPath: path to the settings file
func InitPluginEnabled(cmd *cobra.Command, settingsPath string) {
	sprintf(cmd, tplInitPluginEnabled, settingsPath)
}

// InitSkippedDir reports a directory skipped because it exists.
//
// Parameters:
//   - cmd: Cobra command for output
//   - dir: directory name
func InitSkippedDir(cmd *cobra.Command, dir string) {
	sprintf(cmd, tplInitSkippedDir, dir)
}

// InitCreatedDir reports a directory created during init.
//
// Parameters:
//   - cmd: Cobra command for output
//   - dir: directory name
func InitCreatedDir(cmd *cobra.Command, dir string) {
	sprintf(cmd, tplInitCreatedDir, dir)
}
