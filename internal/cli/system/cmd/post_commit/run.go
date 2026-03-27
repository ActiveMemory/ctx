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
	"regexp"
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
	"github.com/ActiveMemory/ctx/internal/config/token"
	ctxContext "github.com/ActiveMemory/ctx/internal/context/resolve"
	"github.com/ActiveMemory/ctx/internal/notify"
	writeHook "github.com/ActiveMemory/ctx/internal/write/hook"
)

var (
	reGitCommit = regexp.MustCompile(`git\s+commit`)
	reAmend     = regexp.MustCompile(`--amend`)
	// reTaskRef matches Phase-style task references like HA.1, P-2.5, PD.3, CT.1.
	reTaskRef = regexp.MustCompile(`\b[A-Z]+-?\d+\.?\d*\b`)
)

// Run executes the post-commit hook logic.
//
// After a successful git commit (non-amend), nudges the agent to offer
// context capture (decision or learning) and to run lints/tests before
// pushing. Also checks for version drift.
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

	// Only trigger on git commit commands
	if !reGitCommit.MatchString(command) {
		return nil
	}

	// Skip amend commits
	if reAmend.MatchString(command) {
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

	// Spec enforcement: score the commit for bypass indicators.
	if violations := scoreCommitViolations(); violations != "" {
		writeHook.NudgeBlock(cmd, violations)
	}

	return nil
}

// Violation point values for bypass detection.
const (
	violationSpecMissing    = 3
	violationSignoffMissing = 1
	violationTaskRefMissing = 1
	violationSingleLine     = 1
	violationNoTasksChanged = 1

	violationThresholdNudge = 2
	violationThresholdWarn  = 4
)

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

	// Missing Spec: trailer (3 points).
	if !strings.Contains(commitMsg, cfgGit.TrailerSpec) {
		score += violationSpecMissing
		missing = append(missing, "Spec: trailer")
	}

	// Missing Signed-off-by: trailer (1 point).
	if !strings.Contains(commitMsg, cfgGit.TrailerSignedOffBy) {
		score += violationSignoffMissing
		missing = append(missing, "Signed-off-by: trailer")
	}

	// Single-line message — no body (1 point).
	lines := strings.Split(strings.TrimSpace(commitMsg), token.NewlineLF)
	if len(lines) <= 1 {
		score += violationSingleLine
		missing = append(missing, "commit body")
	}

	// No task reference in message (1 point).
	if !reTaskRef.MatchString(commitMsg) {
		score += violationTaskRefMissing
		missing = append(missing, "task reference")
	}

	// Source files changed but no TASKS.md in diff (1 point).
	diffBytes, diffErr := exec.Command("git", "diff-tree", "--no-commit-id", "--name-only", "-r", "HEAD").Output() //nolint:gosec // G204: all args are string literals
	if diffErr == nil {
		diffFiles := string(diffBytes)
		hasSource := strings.Contains(diffFiles, file.ExtGo)
		hasTasks := strings.Contains(diffFiles, cfgCtx.Task)
		if hasSource && !hasTasks {
			score += violationNoTasksChanged
			missing = append(missing, "TASKS.md update")
		}
	}

	if score < violationThresholdNudge {
		return ""
	}

	severity := "informal"
	if score >= violationThresholdWarn {
		severity = "bypassed /ctx-commit"
	}

	title := fmt.Sprintf("Commit Audit (score: %d — %s)", score, severity)
	content := fmt.Sprintf("Missing: %s", strings.Join(missing, ", "))

	return message.NudgeBox(
		desc.Text(text.DescKeyPostCommitRelayPrefix),
		title,
		content,
	)
}
