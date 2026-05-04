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
	"github.com/ActiveMemory/ctx/internal/config/token"
)

// InfoResetPrompt prints the reset confirmation prompt with an
// enumeration of the populated essential files that will be
// overwritten. The prompt is intentionally explicit about the
// blast radius — the silent "Overwrite existing context? [y/N]"
// prompt destroyed thousands of lines of curated content in
// the 2026-04-25 incident.
//
// Parameters:
//   - cmd: Cobra command for output
//   - contextDir: path to the existing .context/ directory
//   - files: basenames of populated essential files that will be
//     backed up and overwritten
func InfoResetPrompt(cmd *cobra.Command, contextDir string, files []string) {
	cmd.Println(desc.Text(text.DescKeyWriteInitResetPromptHeader))
	cmd.Println()
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteInitResetPromptDir), contextDir,
	))
	for _, f := range files {
		cmd.Println(fmt.Sprintf(
			desc.Text(text.DescKeyWriteInitResetPromptFile), f,
		))
	}
	cmd.Println()
	cmd.Print(desc.Text(text.DescKeyWriteInitResetPromptFooter))
	cmd.Print(token.Space)
}

// InfoBackupWritten reports the path of the timestamped backup
// directory that holds the pre-reset snapshot of populated files.
//
// Parameters:
//   - cmd: Cobra command for output
//   - backupDir: absolute path of the backup directory
func InfoBackupWritten(cmd *cobra.Command, backupDir string) {
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteInitBackupWritten), backupDir,
	))
}

// InfoAborted reports that the user cancelled the init operation.
//
// Parameters:
//   - cmd: Cobra command for output
func InfoAborted(cmd *cobra.Command) {
	cmd.Println(desc.Text(text.DescKeyWriteInitAborted))
}

// InfoExistsSkipped reports a template file skipped because it exists.
//
// Parameters:
//   - cmd: Cobra command for output
//   - name: the template filename that was skipped
func InfoExistsSkipped(cmd *cobra.Command, name string) {
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWriteInitExistsSkipped), name))
}

// InfoFileCreated reports a template file that was created.
//
// Parameters:
//   - cmd: Cobra command for output
//   - name: the template filename that was created
func InfoFileCreated(cmd *cobra.Command, name string) {
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWriteInitFileCreated), name))
}

// Initialized reports successful context directory initialization.
//
// Parameters:
//   - cmd: Cobra command for output
//   - contextDir: the path to the initialized .context/ directory
func Initialized(cmd *cobra.Command, contextDir string) {
	cmd.Println()
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWriteInitialized), contextDir))
}

// InfoWarnNonFatal reports a non-fatal warning during init.
//
// Parameters:
//   - cmd: Cobra command for output
//   - label: short description of what failed (e.g. "CLAUDE.md")
//   - err: the non-fatal error
func InfoWarnNonFatal(cmd *cobra.Command, label string, err error) {
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteInitWarnNonFatal),
		label, err))
}

// InfoScratchpadPlaintext reports a plaintext scratchpad was created.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: the scratchpad file path
func InfoScratchpadPlaintext(cmd *cobra.Command, path string) {
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteInitScratchpadPlaintext),
		path))
}

// InfoScratchpadNoKey warns about a missing key for an encrypted scratchpad.
//
// Parameters:
//   - cmd: Cobra command for output
//   - keyPath: the expected key path
func InfoScratchpadNoKey(cmd *cobra.Command, keyPath string) {
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteInitScratchpadNoKey),
		keyPath))
}

// InfoScratchpadKeyCreated reports a scratchpad key was generated.
//
// Parameters:
//   - cmd: Cobra command for output
//   - keyPath: the path where the key was saved
func InfoScratchpadKeyCreated(cmd *cobra.Command, keyPath string) {
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteInitScratchpadKeyCreated),
		keyPath))
}

// InfoCreatingRootFiles prints the heading before root file creation.
//
// Parameters:
//   - cmd: Cobra command for output
func InfoCreatingRootFiles(cmd *cobra.Command) {
	cmd.Println()
	cmd.Println(desc.Text(text.DescKeyWriteInitCreatingRootFiles))
}

// InfoSettingUpPermissions prints the heading before permissions setup.
//
// Parameters:
//   - cmd: Cobra command for output
func InfoSettingUpPermissions(cmd *cobra.Command) {
	cmd.Println()
	cmd.Println(desc.Text(text.DescKeyWriteInitSettingUpPermissions))
}

// InfoGitignoreUpdated reports .gitignore entries were added.
//
// Parameters:
//   - cmd: Cobra command for output
//   - count: number of entries added
func InfoGitignoreUpdated(cmd *cobra.Command, count int) {
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteInitGitignoreUpdated),
		count))
}

// InfoGitignoreReview hints how to review changes.
//
// Parameters:
//   - cmd: Cobra command for output
func InfoGitignoreReview(cmd *cobra.Command) {
	cmd.Println(desc.Text(text.DescKeyWriteInitGitignoreReview))
}

// InfoNextSteps prints the post-init guidance block.
//
// Parameters:
//   - cmd: Cobra command for output
func InfoNextSteps(cmd *cobra.Command) {
	cmd.Println(desc.Text(text.DescKeyWriteInitNextStepsBlock))
}

// InfoActivateHint prints the shell-activation block shown right
// after `ctx init` finishes. The block tells the user how to bind
// CTX_DIR for their shell so subsequent ctx commands resolve to the
// freshly-created context directory.
//
// Under the single-source-anchor resolution model
// (specs/single-source-context-anchor.md) this step is not
// optional: every non-exempt ctx command refuses to run without a
// declared CTX_DIR. The hint closes the loop for new users so
// `ctx init` → next command actually works.
//
// Parameters:
//   - cmd: cobra command for output.
//   - contextDir: absolute path to the just-created .context/
//     directory; used in the `export CTX_DIR=...` variant of the
//     hint. The `eval "$(ctx activate)"` variant takes no arg
//     under the single-source-anchor model and discovers the
//     path itself.
func InfoActivateHint(cmd *cobra.Command, contextDir string) {
	tpl := desc.Text(text.DescKeyWriteInitActivateHint)
	cmd.Println(fmt.Sprintf(tpl, contextDir))
}

// InfoWorkflowTips prints the workflow tips block showing key skills
// and the ceremony loop.
//
// Parameters:
//   - cmd: Cobra command for output
func InfoWorkflowTips(cmd *cobra.Command) {
	cmd.Println(desc.Text(text.DescKeyWriteInitWorkflowTips))
}

// InfoGettingStartedSaved reports that the quick-start reference
// file was written.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: File path that was written
func InfoGettingStartedSaved(cmd *cobra.Command, path string) {
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteInitGettingStartedSaved),
		path))
}
