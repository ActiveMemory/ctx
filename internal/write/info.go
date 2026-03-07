//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package write

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
)

// InfoPathConversionExists reports that a path conversion target already
// exists at the destination. Used during init to show which template files
// were skipped.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - rootDir: project root directory for path resolution.
//   - oldPath: original template-relative path.
//   - newPath: destination-relative path joined with rootDir.
func InfoPathConversionExists(
	cmd *cobra.Command, rootDir, oldPath, newPath string,
) {
	if cmd == nil {
		return
	}
	sprintf(cmd, tplPathExists, oldPath, filepath.Join(rootDir, newPath))
}

// InfoAddedTo confirms an entry was added to a context file.
//
// Parameters:
//   - cmd: Cobra command for output
//   - filename: Name of the file the entry was added to
func InfoAddedTo(cmd *cobra.Command, filename string) {
	sprintf(cmd, tplAddedTo, filename)
}

// InfoMovingTask reports a completed task being moved.
//
// Parameters:
//   - cmd: Cobra command for output
//   - taskText: Truncated task description
func InfoMovingTask(cmd *cobra.Command, taskText string) {
	sprintf(cmd, tplMovingTask, taskText)
}

// InfoSkippingTask reports a task skipped due to incomplete children.
//
// Parameters:
//   - cmd: Cobra command for output
//   - taskText: Truncated task description
func InfoSkippingTask(cmd *cobra.Command, taskText string) {
	sprintf(cmd, tplSkippingTask, taskText)
}

// InfoArchivedTasks reports the number of tasks archived.
//
// Parameters:
//   - cmd: Cobra command for output
//   - count: Number of tasks archived
//   - archiveFile: Path to the archive file
//   - days: Age threshold in days
func InfoArchivedTasks(cmd *cobra.Command, count int, archiveFile string, days int) {
	sprintf(cmd, tplArchivedTasks, count, archiveFile, days)
}

// InfoCompletedTask reports a task marked complete.
//
// Parameters:
//   - cmd: Cobra command for output
//   - taskText: The completed task description
func InfoCompletedTask(cmd *cobra.Command, taskText string) {
	sprintf(cmd, tplCompletedTask, taskText)
}

// InfoConfigProfileDev reports that the dev profile is active.
//
// Parameters:
//   - cmd: Cobra command for output
func InfoConfigProfileDev(cmd *cobra.Command) {
	cmd.Println(tplConfigProfileDev)
}

// InfoConfigProfileBase reports that the base profile is active.
//
// Parameters:
//   - cmd: Cobra command for output
func InfoConfigProfileBase(cmd *cobra.Command) {
	cmd.Println(tplConfigProfileBase)
}

// InfoConfigProfileNone reports that no profile exists.
//
// Parameters:
//   - cmd: Cobra command for output
//   - filename: The .ctxrc filename
func InfoConfigProfileNone(cmd *cobra.Command, filename string) {
	sprintf(cmd, tplConfigProfileNone, filename)
}

// InfoDepsNoProject reports that no supported project was detected.
//
// Parameters:
//   - cmd: Cobra command for output
//   - builderNames: Comma-separated list of supported project types
func InfoDepsNoProject(cmd *cobra.Command, builderNames string) {
	cmd.Println(tplDepsNoProject)
	cmd.Println(tplDepsLookingFor)
	sprintf(cmd, tplDepsUseType, builderNames)
}

// InfoDepsNoDeps reports that no dependencies were found.
//
// Parameters:
//   - cmd: Cobra command for output
func InfoDepsNoDeps(cmd *cobra.Command) {
	cmd.Println(tplDepsNoDeps)
}

// InfoSkillsHeader prints the skills list heading.
//
// Parameters:
//   - cmd: Cobra command for output
func InfoSkillsHeader(cmd *cobra.Command) {
	cmd.Println(tplSkillsHeader)
	cmd.Println()
}

// InfoSkillLine prints a single skill entry.
//
// Parameters:
//   - cmd: Cobra command for output
//   - name: Skill name
//   - description: Truncated skill description
func InfoSkillLine(cmd *cobra.Command, name, description string) {
	sprintf(cmd, tplSkillLine, name, description)
}

// InfoExistsWritingAsAlternative reports that a file already exists and the
// content is being written to an alternative filename instead.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - path: the original target path that already exists.
//   - alternative: the fallback path where content was written.
func InfoExistsWritingAsAlternative(
	cmd *cobra.Command, path, alternative string,
) {
	if cmd == nil {
		return
	}
	sprintf(cmd, tplExistsWritingAsAlternative, path, alternative)
}

// InfoInitOverwritePrompt prints the overwrite confirmation prompt.
//
// Parameters:
//   - cmd: Cobra command for output
//   - contextDir: path to the existing .context/ directory
func InfoInitOverwritePrompt(cmd *cobra.Command, contextDir string) {
	cmd.Print(fmt.Sprintf(tplInitOverwritePrompt, contextDir))
}

// InfoInitAborted reports that the user cancelled the init operation.
//
// Parameters:
//   - cmd: Cobra command for output
func InfoInitAborted(cmd *cobra.Command) {
	cmd.Println(tplInitAborted)
}

// InfoInitExistsSkipped reports a template file skipped because it exists.
//
// Parameters:
//   - cmd: Cobra command for output
//   - name: the template filename that was skipped
func InfoInitExistsSkipped(cmd *cobra.Command, name string) {
	sprintf(cmd, tplInitExistsSkipped, name)
}

// InfoInitFileCreated reports a template file that was created.
//
// Parameters:
//   - cmd: Cobra command for output
//   - name: the template filename that was created
func InfoInitFileCreated(cmd *cobra.Command, name string) {
	sprintf(cmd, tplInitFileCreated, name)
}

// InfoInitialized reports successful context directory initialization.
//
// Parameters:
//   - cmd: Cobra command for output
//   - contextDir: the path to the initialized .context/ directory
func InfoInitialized(cmd *cobra.Command, contextDir string) {
	cmd.Println()
	sprintf(cmd, tplInitialized, contextDir)
}

// InfoInitWarnNonFatal reports a non-fatal warning during init.
//
// Parameters:
//   - cmd: Cobra command for output
//   - label: short description of what failed (e.g. "CLAUDE.md")
//   - err: the non-fatal error
func InfoInitWarnNonFatal(cmd *cobra.Command, label string, err error) {
	sprintf(cmd, tplInitWarnNonFatal, label, err)
}

// InfoInitScratchpadPlaintext reports a plaintext scratchpad was created.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: the scratchpad file path
func InfoInitScratchpadPlaintext(cmd *cobra.Command, path string) {
	sprintf(cmd, tplInitScratchpadPlaintext, path)
}

// InfoInitScratchpadNoKey warns about a missing key for an encrypted scratchpad.
//
// Parameters:
//   - cmd: Cobra command for output
//   - keyPath: the expected key path
func InfoInitScratchpadNoKey(cmd *cobra.Command, keyPath string) {
	sprintf(cmd, tplInitScratchpadNoKey, keyPath)
}

// InfoInitScratchpadKeyCreated reports a scratchpad key was generated.
//
// Parameters:
//   - cmd: Cobra command for output
//   - keyPath: the path where the key was saved
func InfoInitScratchpadKeyCreated(cmd *cobra.Command, keyPath string) {
	sprintf(cmd, tplInitScratchpadKeyCreated, keyPath)
}

// InfoInitCreatingRootFiles prints the heading before root file creation.
//
// Parameters:
//   - cmd: Cobra command for output
func InfoInitCreatingRootFiles(cmd *cobra.Command) {
	cmd.Println()
	cmd.Println(tplInitCreatingRootFiles)
}

// InfoInitSettingUpPermissions prints the heading before permissions setup.
//
// Parameters:
//   - cmd: Cobra command for output
func InfoInitSettingUpPermissions(cmd *cobra.Command) {
	cmd.Println()
	cmd.Println(tplInitSettingUpPermissions)
}

// InfoInitGitignoreUpdated reports .gitignore entries were added.
//
// Parameters:
//   - cmd: Cobra command for output
//   - count: number of entries added
func InfoInitGitignoreUpdated(cmd *cobra.Command, count int) {
	sprintf(cmd, tplInitGitignoreUpdated, count)
}

// InfoInitGitignoreReview hints how to review changes.
//
// Parameters:
//   - cmd: Cobra command for output
func InfoInitGitignoreReview(cmd *cobra.Command) {
	cmd.Println(tplInitGitignoreReview)
}

// InfoInitNextSteps prints the post-init guidance block.
//
// Parameters:
//   - cmd: Cobra command for output
func InfoInitNextSteps(cmd *cobra.Command) {
	cmd.Println()
	cmd.Println(tplInitNextSteps)
	cmd.Println()
	cmd.Println(tplInitPluginInfo)
	cmd.Println()
	cmd.Println(tplInitPluginNote)
}

// InfoObsidianGenerated reports successful Obsidian vault generation.
//
// Parameters:
//   - cmd: Cobra command for output
//   - count: Number of entries generated
//   - output: Output directory path
func InfoObsidianGenerated(cmd *cobra.Command, count int, output string) {
	sprintf(cmd, tplObsidianGenerated, count, output)
	cmd.Println()
	cmd.Println("Next steps:")
	sprintf(cmd, tplObsidianNextSteps, output)
}
