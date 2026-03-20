//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package initialize

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/spf13/cobra"
)

// InfoOverwritePrompt prints the overwrite confirmation prompt.
//
// Parameters:
//   - cmd: Cobra command for output
//   - contextDir: path to the existing .context/ directory
func InfoOverwritePrompt(cmd *cobra.Command, contextDir string) {
	cmd.Print(fmt.Sprintf(desc.TextDesc(text.DescKeyWriteInitOverwritePrompt), contextDir))
}

// InfoAborted reports that the user cancelled the init operation.
//
// Parameters:
//   - cmd: Cobra command for output
func InfoAborted(cmd *cobra.Command) {
	cmd.Println(desc.TextDesc(text.DescKeyWriteInitAborted))
}

// InfoExistsSkipped reports a template file skipped because it exists.
//
// Parameters:
//   - cmd: Cobra command for output
//   - name: the template filename that was skipped
func InfoExistsSkipped(cmd *cobra.Command, name string) {
	cmd.Println(fmt.Sprintf(desc.TextDesc(text.DescKeyWriteInitExistsSkipped), name))
}

// InfoFileCreated reports a template file that was created.
//
// Parameters:
//   - cmd: Cobra command for output
//   - name: the template filename that was created
func InfoFileCreated(cmd *cobra.Command, name string) {
	cmd.Println(fmt.Sprintf(desc.TextDesc(text.DescKeyWriteInitFileCreated), name))
}

// InfoInitialized reports successful context directory initialization.
//
// Parameters:
//   - cmd: Cobra command for output
//   - contextDir: the path to the initialized .context/ directory
func InfoInitialized(cmd *cobra.Command, contextDir string) {
	cmd.Println()
	cmd.Println(fmt.Sprintf(desc.TextDesc(text.DescKeyWriteInitialized), contextDir))
}

// InfoWarnNonFatal reports a non-fatal warning during init.
//
// Parameters:
//   - cmd: Cobra command for output
//   - label: short description of what failed (e.g. "CLAUDE.md")
//   - err: the non-fatal error
func InfoWarnNonFatal(cmd *cobra.Command, label string, err error) {
	cmd.Println(fmt.Sprintf(desc.TextDesc(text.DescKeyWriteInitWarnNonFatal), label, err))
}

// InfoScratchpadPlaintext reports a plaintext scratchpad was created.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: the scratchpad file path
func InfoScratchpadPlaintext(cmd *cobra.Command, path string) {
	cmd.Println(fmt.Sprintf(desc.TextDesc(text.DescKeyWriteInitScratchpadPlaintext), path))
}

// InfoScratchpadNoKey warns about a missing key for an encrypted scratchpad.
//
// Parameters:
//   - cmd: Cobra command for output
//   - keyPath: the expected key path
func InfoScratchpadNoKey(cmd *cobra.Command, keyPath string) {
	cmd.Println(fmt.Sprintf(desc.TextDesc(text.DescKeyWriteInitScratchpadNoKey), keyPath))
}

// InfoScratchpadKeyCreated reports a scratchpad key was generated.
//
// Parameters:
//   - cmd: Cobra command for output
//   - keyPath: the path where the key was saved
func InfoScratchpadKeyCreated(cmd *cobra.Command, keyPath string) {
	cmd.Println(fmt.Sprintf(desc.TextDesc(text.DescKeyWriteInitScratchpadKeyCreated), keyPath))
}

// InfoCreatingRootFiles prints the heading before root file creation.
//
// Parameters:
//   - cmd: Cobra command for output
func InfoCreatingRootFiles(cmd *cobra.Command) {
	cmd.Println()
	cmd.Println(desc.TextDesc(text.DescKeyWriteInitCreatingRootFiles))
}

// InfoSettingUpPermissions prints the heading before permissions setup.
//
// Parameters:
//   - cmd: Cobra command for output
func InfoSettingUpPermissions(cmd *cobra.Command) {
	cmd.Println()
	cmd.Println(desc.TextDesc(text.DescKeyWriteInitSettingUpPermissions))
}

// InfoGitignoreUpdated reports .gitignore entries were added.
//
// Parameters:
//   - cmd: Cobra command for output
//   - count: number of entries added
func InfoGitignoreUpdated(cmd *cobra.Command, count int) {
	cmd.Println(fmt.Sprintf(desc.TextDesc(text.DescKeyWriteInitGitignoreUpdated), count))
}

// InfoGitignoreReview hints how to review changes.
//
// Parameters:
//   - cmd: Cobra command for output
func InfoGitignoreReview(cmd *cobra.Command) {
	cmd.Println(desc.TextDesc(text.DescKeyWriteInitGitignoreReview))
}

// InfoNextSteps prints the post-init guidance block.
//
// Parameters:
//   - cmd: Cobra command for output
func InfoNextSteps(cmd *cobra.Command) {
	cmd.Println(desc.TextDesc(text.DescKeyWriteInitNextStepsBlock))
}
