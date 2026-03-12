//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package write

import (
	"fmt"
	"github.com/ActiveMemory/ctx/internal/write/config"
	"github.com/spf13/cobra"
)

// InitCreated reports a file created during init.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: created file path
func InitCreated(cmd *cobra.Command, path string) {
	cmd.Println(fmt.Sprintf(config.TplInitFileCreated, path))
}

// InitCreatedWith reports a file created with a qualifier (e.g. " (ralph mode)").
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: created file path
//   - qualifier: additional info appended after the path
func InitCreatedWith(cmd *cobra.Command, path, qualifier string) {
	cmd.Println(fmt.Sprintf(config.TplInitCreatedWith, path, qualifier))
}

// InitSkipped reports a file skipped because it already exists.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: skipped file path
func InitSkipped(cmd *cobra.Command, path string) {
	cmd.Println(fmt.Sprintf(config.TplInitExistsSkipped, path))
}

// InitSkippedPlain reports a file skipped without detail.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: skipped file path
func InitSkippedPlain(cmd *cobra.Command, path string) {
	cmd.Println(fmt.Sprintf(config.TplInitSkippedPlain, path))
}

// InitCtxContentExists reports a file skipped because ctx content exists.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: skipped file path
func InitCtxContentExists(cmd *cobra.Command, path string) {
	cmd.Println(fmt.Sprintf(config.TplInitCtxContentExists, path))
}

// InitMerged reports a file merged during init.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: merged file path
func InitMerged(cmd *cobra.Command, path string) {
	cmd.Println(fmt.Sprintf(config.TplInitMerged, path))
}

// InitBackup reports a backup file created.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: backup file path
func InitBackup(cmd *cobra.Command, path string) {
	cmd.Println(fmt.Sprintf(config.TplInitBackup, path))
}

// InitUpdatedCtxSection reports a file whose ctx section was updated.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: updated file path
func InitUpdatedCtxSection(cmd *cobra.Command, path string) {
	cmd.Println(fmt.Sprintf(config.TplInitUpdatedCtxSection, path))
}

// InitUpdatedPlanSection reports a file whose plan section was updated.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: updated file path
func InitUpdatedPlanSection(cmd *cobra.Command, path string) {
	cmd.Println(fmt.Sprintf(config.TplInitUpdatedPlanSection, path))
}

// InitUpdatedPromptSection reports a file whose prompt section was updated.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: updated file path
func InitUpdatedPromptSection(cmd *cobra.Command, path string) {
	cmd.Println(fmt.Sprintf(config.TplInitUpdatedPromptSection, path))
}

// InitFileExistsNoCtx reports a file exists without ctx content.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: file path
func InitFileExistsNoCtx(cmd *cobra.Command, path string) {
	cmd.Println(fmt.Sprintf(config.TplInitFileExistsNoCtx, path))
}

// InitNoChanges reports a settings file with no changes needed.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: settings file path
func InitNoChanges(cmd *cobra.Command, path string) {
	cmd.Println(fmt.Sprintf(config.TplInitNoChanges, path))
}

// InitPermsMergedDeduped reports permissions merged and deduped.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: settings file path
func InitPermsMergedDeduped(cmd *cobra.Command, path string) {
	cmd.Println(fmt.Sprintf(config.TplInitPermsMergedDeduped, path))
}

// InitPermsDeduped reports duplicate permissions removed.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: settings file path
func InitPermsDeduped(cmd *cobra.Command, path string) {
	cmd.Println(fmt.Sprintf(config.TplInitPermsDeduped, path))
}

// InitPermsAllowDeny reports allow+deny permissions added.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: settings file path
func InitPermsAllowDeny(cmd *cobra.Command, path string) {
	cmd.Println(fmt.Sprintf(config.TplInitPermsAllowDeny, path))
}

// InitPermsDeny reports deny permissions added.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: settings file path
func InitPermsDeny(cmd *cobra.Command, path string) {
	cmd.Println(fmt.Sprintf(config.TplInitPermsDeny, path))
}

// InitPermsAllow reports ctx permissions added.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: settings file path
func InitPermsAllow(cmd *cobra.Command, path string) {
	cmd.Println(fmt.Sprintf(config.TplInitPermsAllow, path))
}

// InitMakefileCreated reports a new Makefile created with ctx include.
//
// Parameters:
//   - cmd: Cobra command for output
func InitMakefileCreated(cmd *cobra.Command) {
	cmd.Println(config.TplInitMakefileCreated)
}

// InitMakefileIncludes reports Makefile already includes the directive.
//
// Parameters:
//   - cmd: Cobra command for output
//   - filename: included filename
func InitMakefileIncludes(cmd *cobra.Command, filename string) {
	cmd.Println(fmt.Sprintf(config.TplInitMakefileIncludes, filename))
}

// InitMakefileAppended reports an include appended to Makefile.
//
// Parameters:
//   - cmd: Cobra command for output
//   - filename: included filename
func InitMakefileAppended(cmd *cobra.Command, filename string) {
	cmd.Println(fmt.Sprintf(config.TplInitMakefileAppended, filename))
}

// InitPluginSkipped reports plugin enablement was skipped.
//
// Parameters:
//   - cmd: Cobra command for output
func InitPluginSkipped(cmd *cobra.Command) {
	cmd.Println(config.TplInitPluginSkipped)
}

// InitPluginAlreadyEnabled reports plugin is already enabled globally.
//
// Parameters:
//   - cmd: Cobra command for output
func InitPluginAlreadyEnabled(cmd *cobra.Command) {
	cmd.Println(config.TplInitPluginAlreadyEnabled)
}

// InitPluginEnabled reports plugin enabled globally.
//
// Parameters:
//   - cmd: Cobra command for output
//   - settingsPath: path to the settings file
func InitPluginEnabled(cmd *cobra.Command, settingsPath string) {
	cmd.Println(fmt.Sprintf(config.TplInitPluginEnabled, settingsPath))
}

// InitSkippedDir reports a directory skipped because it exists.
//
// Parameters:
//   - cmd: Cobra command for output
//   - dir: directory name
func InitSkippedDir(cmd *cobra.Command, dir string) {
	cmd.Println(fmt.Sprintf(config.TplInitSkippedDir, dir))
}

// InitCreatedDir reports a directory created during init.
//
// Parameters:
//   - cmd: Cobra command for output
//   - dir: directory name
func InitCreatedDir(cmd *cobra.Command, dir string) {
	cmd.Println(fmt.Sprintf(config.TplInitCreatedDir, dir))
}
