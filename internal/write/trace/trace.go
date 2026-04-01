//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package trace

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	internalTrace "github.com/ActiveMemory/ctx/internal/trace"
)

// CommitHeader prints the commit hash, message, and date for a single commit.
//
// Parameters:
//   - cmd: Cobra command for output
//   - shortHash: abbreviated commit hash
//   - message: commit subject line
//   - date: commit date string
func CommitHeader(cmd *cobra.Command, shortHash, message, date string) {
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWriteTraceCommitHeader), shortHash))
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWriteTraceCommitMessage), message))
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWriteTraceCommitDate), date))
}

// CommitNoContext prints the "no context" message.
//
// Parameters:
//   - cmd: Cobra command for output
func CommitNoContext(cmd *cobra.Command) {
	cmd.Println(desc.Text(text.DescKeyWriteTraceCommitNoContext))
}

// CommitContext prints the "Context:" label.
//
// Parameters:
//   - cmd: Cobra command for output
func CommitContext(cmd *cobra.Command) {
	cmd.Println(desc.Text(text.DescKeyWriteTraceCommitContext))
}

// FileEntry prints a single file trace entry line.
//
// Parameters:
//   - cmd: Cobra command for output
//   - shortHash: abbreviated commit hash
//   - date: commit date string
//   - subject: commit subject line
//   - refStr: formatted ref summary
func FileEntry(cmd *cobra.Command, shortHash, date, subject, refStr string) {
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteTraceFileEntry),
		shortHash, date, subject, refStr,
	))
}

// HooksEnabled reports that trace hooks were installed.
//
// Parameters:
//   - cmd: Cobra command for output
func HooksEnabled(cmd *cobra.Command) {
	cmd.Println(desc.Text(text.DescKeyWriteTraceHooksEnabled))
}

// HooksDisabled reports that trace hooks were removed.
//
// Parameters:
//   - cmd: Cobra command for output
func HooksDisabled(cmd *cobra.Command) {
	cmd.Println(desc.Text(text.DescKeyWriteTraceHooksDisabled))
}

// LastEntry prints a single line in the last-N listing.
//
// Parameters:
//   - cmd: Cobra command for output
//   - shortHash: abbreviated commit hash
//   - message: commit subject line
//   - refSummary: formatted ref summary
func LastEntry(cmd *cobra.Command, shortHash, message, refSummary string) {
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteTraceLastEntry),
		shortHash, message, refSummary,
	))
}

// Resolved prints a single resolved context reference.
//
// Parameters:
//   - cmd: Cobra command for output
//   - rr: resolved reference to display
func Resolved(cmd *cobra.Command, rr internalTrace.ResolvedRef) {
	typeLabel := strings.ToUpper(rr.Type[0:1]) + rr.Type[1:]
	if rr.Found && rr.Title != "" {
		if rr.Detail != "" {
			cmd.Println(fmt.Sprintf(
				desc.Text(text.DescKeyWriteTraceResolvedFull),
				typeLabel, rr.Raw, rr.Title, rr.Detail,
			))
		} else {
			cmd.Println(fmt.Sprintf(
				desc.Text(text.DescKeyWriteTraceResolvedTitle),
				typeLabel, rr.Raw, rr.Title,
			))
		}
	} else {
		cmd.Println(fmt.Sprintf(
			desc.Text(text.DescKeyWriteTraceResolvedRaw),
			typeLabel, rr.Raw,
		))
	}
}

// Tagged reports that a commit was tagged with a context note.
//
// Parameters:
//   - cmd: Cobra command for output
//   - shortHash: abbreviated commit hash
//   - note: the context note that was attached
func Tagged(cmd *cobra.Command, shortHash, note string) {
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWriteTraceTagged), shortHash, note))
}
