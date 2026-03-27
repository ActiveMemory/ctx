//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package post_commit

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	coreCheck "github.com/ActiveMemory/ctx/internal/cli/system/core/check"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/drift"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/message"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/nudge"
	coreSession "github.com/ActiveMemory/ctx/internal/cli/system/core/session"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/state"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	cfgCtx "github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/file"
	cfgGit "github.com/ActiveMemory/ctx/internal/config/git"
	"github.com/ActiveMemory/ctx/internal/config/hook"
	"github.com/ActiveMemory/ctx/internal/config/regex"
	"github.com/ActiveMemory/ctx/internal/config/stats"
	"github.com/ActiveMemory/ctx/internal/config/token"
	ctxContext "github.com/ActiveMemory/ctx/internal/context/resolve"
	"github.com/ActiveMemory/ctx/internal/notify"
	writeHook "github.com/ActiveMemory/ctx/internal/write/hook"
)

// Run executes the post-commit hook logic.
//
// After a successful git commit (non-amend), nudges the agent to offer
// context capture (decision or learning) and to run lints/tests before
// pushing. Also checks for version drift and spec enforcement.
//
// Parameters:
//   - cmd: Cobra command for output
//   - stdin: standard input for hook JSON
//
// Returns:
//   - error: Always nil (hook errors are non-fatal)
func Run(cmd *cobra.Command, stdin *os.File) error {
	if !state.Initialized() {
		return nil
	}
	input, sessionID, paused := coreCheck.Preamble(stdin)
	if paused {
		return nil
	}

	command := input.ToolInput.Command

	if !regex.GitCommit.MatchString(command) {
		return nil
	}

	if regex.GitAmend.MatchString(command) {
		return nil
	}

	hookName, variant := hook.PostCommit, hook.VariantNudge

	fallback := desc.Text(text.DescKeyPostCommitFallback)
	msg := message.LoadMessage(hookName, variant, nil, fallback)
	if msg == "" {
		return nil
	}
	msg = ctxContext.AppendDir(msg)
	writeHook.HookContext(cmd, coreSession.FormatContext(hook.EventPostToolUse, msg))

	ref := notify.NewTemplateRef(hookName, variant, nil)
	nudge.Relay(fmt.Sprintf(desc.Text(text.DescKeyRelayPrefixFormat), hookName, desc.Text(text.DescKeyPostCommitRelayMessage)), input.SessionID, ref)

	if driftResponse := drift.CheckVersionDrift(sessionID); driftResponse != "" {
		writeHook.HookContext(cmd, driftResponse)
	}

	if violations := scoreCommitViolations(); violations != "" {
		writeHook.NudgeBlock(cmd, violations)
	}

	return nil
}

// scoreCommitViolations reads the last commit and scores it for signs that
// the agent bypassed /ctx-commit. Returns a formatted nudge box for the
// human, or empty string if the commit looks clean.
func scoreCommitViolations() string {
	msgBytes, msgErr := exec.Command("git", "log", "-1", "--format=%B").Output() //nolint:gosec // G204: all args are string literals
	if msgErr != nil {
		return ""
	}
	commitMsg := string(msgBytes)

	score := 0
	var missing []string

	if !strings.Contains(commitMsg, cfgGit.TrailerSpec) {
		score += stats.ViolationSpecMissing
		missing = append(missing, desc.Text(text.DescKeyPostCommitMissingSpec))
	}

	if !strings.Contains(commitMsg, cfgGit.TrailerSignedOffBy) {
		score += stats.ViolationSignoffMissing
		missing = append(missing, desc.Text(text.DescKeyPostCommitMissingSignoff))
	}

	lines := strings.Split(strings.TrimSpace(commitMsg), token.NewlineLF)
	if len(lines) <= 1 {
		score += stats.ViolationSingleLine
		missing = append(missing, desc.Text(text.DescKeyPostCommitMissingBody))
	}

	if !regex.TaskRef.MatchString(commitMsg) {
		score += stats.ViolationTaskRefMissing
		missing = append(missing, desc.Text(text.DescKeyPostCommitMissingTaskRef))
	}

	diffBytes, diffErr := exec.Command("git", "diff-tree", "--no-commit-id", "--name-only", "-r", "HEAD").Output() //nolint:gosec // G204: all args are string literals
	if diffErr == nil {
		diffFiles := string(diffBytes)
		hasSource := strings.Contains(diffFiles, file.ExtGo)
		hasTasks := strings.Contains(diffFiles, cfgCtx.Task)
		if hasSource && !hasTasks {
			score += stats.ViolationNoTasksChanged
			missing = append(missing, desc.Text(text.DescKeyPostCommitMissingTaskUpdate))
		}
	}

	if score < stats.ViolationThresholdNudge {
		return ""
	}

	severity := desc.Text(text.DescKeyPostCommitSeverityInformal)
	if score >= stats.ViolationThresholdWarn {
		severity = desc.Text(text.DescKeyPostCommitSeveritySkipped)
	}

	title := fmt.Sprintf(desc.Text(text.DescKeyPostCommitAuditTitle), score, severity)
	content := fmt.Sprintf(desc.Text(text.DescKeyPostCommitAuditContent), strings.Join(missing, ", "))

	return message.NudgeBox(
		desc.Text(text.DescKeyPostCommitRelayPrefix),
		title,
		content,
	)
}
