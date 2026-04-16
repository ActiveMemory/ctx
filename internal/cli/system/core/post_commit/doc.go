//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package post_commit scores the last git commit for
// signs that the agent bypassed the /ctx-commit skill.
// It is called by the post-commit hook to nudge the
// human when commit hygiene is poor.
//
// # Violation Scoring
//
// [ScoreCommitViolations] reads the last commit message
// and diff tree, then checks for five violation types:
//
//   - Missing Spec trailer: the commit lacks the
//     structured spec reference
//   - Missing Signed-off-by: no sign-off line
//   - Single-line message: no body after the subject
//   - Missing task reference: no T-NNN or similar ref
//   - Source without TASKS.md: Go files changed but
//     TASKS.md was not updated alongside them
//
// Each violation adds a weighted score. When the total
// exceeds the nudge threshold, the function returns a
// formatted nudge box with a severity label (informal
// or skipped) and a list of missing elements.
//
// # Output
//
// The nudge box is built by message.NudgeBox and
// includes a relay prefix so the hook can forward it
// to the agent session. Returns an empty string when
// the commit looks clean or the score is below the
// nudge threshold.
package post_commit
