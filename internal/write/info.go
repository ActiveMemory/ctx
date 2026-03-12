//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package write

import (
	"fmt"
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/write/config"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
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
	cmd.Println(fmt.Sprintf(config.TplPathExists, oldPath, filepath.Join(rootDir, newPath)))
}

// InfoAddedTo confirms an entry was added to a context file.
//
// Parameters:
//   - cmd: Cobra command for output
//   - filename: Name of the file the entry was added to
func InfoAddedTo(cmd *cobra.Command, filename string) {
	cmd.Println(fmt.Sprintf(config.TplAddedTo, filename))
}

// InfoMovingTask reports a completed task being moved.
//
// Parameters:
//   - cmd: Cobra command for output
//   - taskText: Truncated task description
func InfoMovingTask(cmd *cobra.Command, taskText string) {
	cmd.Println(fmt.Sprintf(config.TplMovingTask, taskText))
}

// InfoSkippingTask reports a task skipped due to incomplete children.
//
// Parameters:
//   - cmd: Cobra command for output
//   - taskText: Truncated task description
func InfoSkippingTask(cmd *cobra.Command, taskText string) {
	cmd.Println(fmt.Sprintf(assets.TextDesc(assets.TextDescKeyTaskArchiveSkipping), taskText))
}

// InfoArchivedTasks reports the number of tasks archived.
//
// Parameters:
//   - cmd: Cobra command for output
//   - count: Number of tasks archived
//   - archiveFile: Path to the archive file
//   - days: Age threshold in days
func InfoArchivedTasks(cmd *cobra.Command, count int, archiveFile string, days int) {
	cmd.Println(fmt.Sprintf(assets.TextDesc(assets.TextDescKeyTaskArchiveSuccessWithAge), count, archiveFile, days))
}

// InfoCompletedTask reports a task marked complete.
//
// Parameters:
//   - cmd: Cobra command for output
//   - taskText: The completed task description
func InfoCompletedTask(cmd *cobra.Command, taskText string) {
	cmd.Println(fmt.Sprintf(config.TplCompletedTask, taskText))
}

// InfoConfigProfileDev reports that the dev profile is active.
//
// Parameters:
//   - cmd: Cobra command for output
func InfoConfigProfileDev(cmd *cobra.Command) {
	cmd.Println(config.TplConfigProfileDev)
}

// InfoConfigProfileBase reports that the base profile is active.
//
// Parameters:
//   - cmd: Cobra command for output
func InfoConfigProfileBase(cmd *cobra.Command) {
	cmd.Println(config.TplConfigProfileBase)
}

// InfoConfigProfileNone reports that no profile exists.
//
// Parameters:
//   - cmd: Cobra command for output
//   - filename: The .ctxrc filename
func InfoConfigProfileNone(cmd *cobra.Command, filename string) {
	cmd.Println(fmt.Sprintf(config.TplConfigProfileNone, filename))
}

// InfoDepsNoProject reports that no supported project was detected.
//
// Parameters:
//   - cmd: Cobra command for output
//   - builderNames: Comma-separated list of supported project types
func InfoDepsNoProject(cmd *cobra.Command, builderNames string) {
	cmd.Println(config.TplDepsNoProject)
	cmd.Println(config.TplDepsLookingFor)
	cmd.Println(fmt.Sprintf(config.TplDepsUseType, builderNames))
}

// InfoDepsNoDeps reports that no dependencies were found.
//
// Parameters:
//   - cmd: Cobra command for output
func InfoDepsNoDeps(cmd *cobra.Command) {
	cmd.Println(config.TplDepsNoDeps)
}

// InfoSkillsHeader prints the skills list heading.
//
// Parameters:
//   - cmd: Cobra command for output
func InfoSkillsHeader(cmd *cobra.Command) {
	cmd.Println(config.TplSkillsHeader)
	cmd.Println()
}

// InfoSkillLine prints a single skill entry.
//
// Parameters:
//   - cmd: Cobra command for output
//   - name: Skill name
//   - description: Truncated skill description
func InfoSkillLine(cmd *cobra.Command, name, description string) {
	cmd.Println(fmt.Sprintf(config.TplSkillLine, name, description))
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
	cmd.Println(fmt.Sprintf(config.TplExistsWritingAsAlternative, path, alternative))
}

// InfoInitOverwritePrompt prints the overwrite confirmation prompt.
//
// Parameters:
//   - cmd: Cobra command for output
//   - contextDir: path to the existing .context/ directory
func InfoInitOverwritePrompt(cmd *cobra.Command, contextDir string) {
	cmd.Print(fmt.Sprintf(config.TplInitOverwritePrompt, contextDir))
}

// InfoInitAborted reports that the user cancelled the init operation.
//
// Parameters:
//   - cmd: Cobra command for output
func InfoInitAborted(cmd *cobra.Command) {
	cmd.Println(config.TplInitAborted)
}

// InfoInitExistsSkipped reports a template file skipped because it exists.
//
// Parameters:
//   - cmd: Cobra command for output
//   - name: the template filename that was skipped
func InfoInitExistsSkipped(cmd *cobra.Command, name string) {
	cmd.Println(fmt.Sprintf(config.TplInitExistsSkipped, name))
}

// InfoInitFileCreated reports a template file that was created.
//
// Parameters:
//   - cmd: Cobra command for output
//   - name: the template filename that was created
func InfoInitFileCreated(cmd *cobra.Command, name string) {
	cmd.Println(fmt.Sprintf(config.TplInitFileCreated, name))
}

// InfoInitialized reports successful context directory initialization.
//
// Parameters:
//   - cmd: Cobra command for output
//   - contextDir: the path to the initialized .context/ directory
func InfoInitialized(cmd *cobra.Command, contextDir string) {
	cmd.Println()
	cmd.Println(fmt.Sprintf(config.TplInitialized, contextDir))
}

// InfoInitWarnNonFatal reports a non-fatal warning during init.
//
// Parameters:
//   - cmd: Cobra command for output
//   - label: short description of what failed (e.g. "CLAUDE.md")
//   - err: the non-fatal error
func InfoInitWarnNonFatal(cmd *cobra.Command, label string, err error) {
	cmd.Println(fmt.Sprintf(config.TplInitWarnNonFatal, label, err))
}

// InfoInitScratchpadPlaintext reports a plaintext scratchpad was created.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: the scratchpad file path
func InfoInitScratchpadPlaintext(cmd *cobra.Command, path string) {
	cmd.Println(fmt.Sprintf(config.TplInitScratchpadPlaintext, path))
}

// InfoInitScratchpadNoKey warns about a missing key for an encrypted scratchpad.
//
// Parameters:
//   - cmd: Cobra command for output
//   - keyPath: the expected key path
func InfoInitScratchpadNoKey(cmd *cobra.Command, keyPath string) {
	cmd.Println(fmt.Sprintf(config.TplInitScratchpadNoKey, keyPath))
}

// InfoInitScratchpadKeyCreated reports a scratchpad key was generated.
//
// Parameters:
//   - cmd: Cobra command for output
//   - keyPath: the path where the key was saved
func InfoInitScratchpadKeyCreated(cmd *cobra.Command, keyPath string) {
	cmd.Println(fmt.Sprintf(config.TplInitScratchpadKeyCreated, keyPath))
}

// InfoInitCreatingRootFiles prints the heading before root file creation.
//
// Parameters:
//   - cmd: Cobra command for output
func InfoInitCreatingRootFiles(cmd *cobra.Command) {
	cmd.Println()
	cmd.Println(config.TplInitCreatingRootFiles)
}

// InfoInitSettingUpPermissions prints the heading before permissions setup.
//
// Parameters:
//   - cmd: Cobra command for output
func InfoInitSettingUpPermissions(cmd *cobra.Command) {
	cmd.Println()
	cmd.Println(config.TplInitSettingUpPermissions)
}

// InfoInitGitignoreUpdated reports .gitignore entries were added.
//
// Parameters:
//   - cmd: Cobra command for output
//   - count: number of entries added
func InfoInitGitignoreUpdated(cmd *cobra.Command, count int) {
	cmd.Println(fmt.Sprintf(config.TplInitGitignoreUpdated, count))
}

// InfoInitGitignoreReview hints how to review changes.
//
// Parameters:
//   - cmd: Cobra command for output
func InfoInitGitignoreReview(cmd *cobra.Command) {
	cmd.Println(config.TplInitGitignoreReview)
}

// InfoInitNextSteps prints the post-init guidance block.
//
// Parameters:
//   - cmd: Cobra command for output
func InfoInitNextSteps(cmd *cobra.Command) {
	cmd.Println()
	cmd.Println(config.TplInitNextSteps)
	cmd.Println()
	cmd.Println(config.TplInitPluginInfo)
	cmd.Println()
	cmd.Println(config.TplInitPluginNote)
}

// InfoObsidianGenerated reports successful Obsidian vault generation.
//
// Parameters:
//   - cmd: Cobra command for output
//   - count: Number of entries generated
//   - output: Output directory path
func InfoObsidianGenerated(cmd *cobra.Command, count int, output string) {
	cmd.Println(fmt.Sprintf(config.TplObsidianGenerated, count, output))
	cmd.Println()
	cmd.Println("Next steps:")
	cmd.Println(fmt.Sprintf(config.TplObsidianNextSteps, output))
}

// InfoJournalOrphanRemoved reports a removed orphan file.
//
// Parameters:
//   - cmd: Cobra command for output
//   - name: Filename that was removed
func InfoJournalOrphanRemoved(cmd *cobra.Command, name string) {
	cmd.Println(fmt.Sprintf(config.TplJournalOrphanRemoved, name))
}

// InfoJournalSiteGenerated reports successful site generation with next steps.
//
// Parameters:
//   - cmd: Cobra command for output
//   - count: Number of entries generated
//   - output: Output directory path
//   - zensicalBin: Zensical binary name
func InfoJournalSiteGenerated(cmd *cobra.Command, count int, output, zensicalBin string) {
	cmd.Println(fmt.Sprintf(config.TplJournalSiteGenerated, count, output))
	cmd.Println()
	cmd.Println("Next steps:")
	cmd.Println(fmt.Sprintf(config.TplJournalSiteNextSteps, output, zensicalBin))
	cmd.Println("  or")
	cmd.Println(config.TplJournalSiteAlt)
}

// InfoJournalSiteStarting reports the server is starting.
//
// Parameters:
//   - cmd: Cobra command for output
func InfoJournalSiteStarting(cmd *cobra.Command) {
	cmd.Println()
	cmd.Println(config.TplJournalSiteStarting)
}

// InfoJournalSiteBuilding reports a build is in progress.
//
// Parameters:
//   - cmd: Cobra command for output
func InfoJournalSiteBuilding(cmd *cobra.Command) {
	cmd.Println()
	cmd.Println(config.TplJournalSiteBuilding)
}

// InfoLoopGenerated reports successful loop script generation with details.
//
// Parameters:
//   - cmd: Cobra command for output
//   - outputFile: Generated script path
//   - heading: Start heading text
//   - tool: Selected AI tool
//   - promptFile: Prompt file path
//   - maxIterations: Max iterations (0 = unlimited)
//   - completionMsg: Completion signal string
func InfoLoopGenerated(
	cmd *cobra.Command,
	outputFile, heading, tool, promptFile string,
	maxIterations int,
	completionMsg string,
) {
	cmd.Println(fmt.Sprintf(config.TplLoopGenerated, outputFile))
	cmd.Println()
	cmd.Println(heading)
	cmd.Println(fmt.Sprintf(config.TplLoopRunCmd, outputFile))
	cmd.Println()
	cmd.Println(fmt.Sprintf(config.TplLoopTool, tool))
	cmd.Println(fmt.Sprintf(config.TplLoopPrompt, promptFile))
	if maxIterations > 0 {
		cmd.Println(fmt.Sprintf(config.TplLoopMaxIterations, maxIterations))
	} else {
		cmd.Println(config.TplLoopUnlimited)
	}
	cmd.Println(fmt.Sprintf(config.TplLoopCompletion, completionMsg))
}
