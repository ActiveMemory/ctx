//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package post_commit

import (
	"fmt"
	"strings"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/message"
	cfgCtx "github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/file"
	cfgGit "github.com/ActiveMemory/ctx/internal/config/git"
	"github.com/ActiveMemory/ctx/internal/config/regex"
	"github.com/ActiveMemory/ctx/internal/config/stats"
	"github.com/ActiveMemory/ctx/internal/config/token"
	execGit "github.com/ActiveMemory/ctx/internal/exec/git"
)

// scoreCommitViolations reads the last commit and scores it for signs that
// the agent bypassed /ctx-commit. Returns a formatted nudge box for the
// human, or empty string if the commit looks clean.
//
// Returns:
//   - string: Formatted nudge box, or empty string if no violations detected
func scoreCommitViolations() string {
	msgBytes, msgErr := execGit.LastCommitMessage()
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

	diffBytes, diffErr := execGit.DiffTreeHead()
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

	title := fmt.Sprintf(
		desc.Text(text.DescKeyPostCommitAuditTitle),
		score, severity,
	)
	content := fmt.Sprintf(
		desc.Text(text.DescKeyPostCommitAuditContent),
		strings.Join(missing, token.CommaSpace),
	)

	return message.NudgeBox(
		desc.Text(text.DescKeyPostCommitRelayPrefix),
		title,
		content,
	)
}
