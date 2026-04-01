//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package show

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	cfgGit "github.com/ActiveMemory/ctx/internal/config/git"
	"github.com/ActiveMemory/ctx/internal/config/token"
	errTrace "github.com/ActiveMemory/ctx/internal/err/trace"
	"github.com/ActiveMemory/ctx/internal/exec/git"
	"github.com/ActiveMemory/ctx/internal/trace"
	writeTrace "github.com/ActiveMemory/ctx/internal/write/trace"
)

// Commit displays the context refs for a single commit.
//
// Parameters:
//   - cmd: Cobra command for output stream
//   - hash: full or abbreviated commit hash
//   - contextDir: absolute path to the context directory
//   - traceDir: absolute path to the trace directory
//   - jsonOutput: whether to format output as JSON
//
// Returns:
//   - error: non-nil on execution failure
func Commit(
	cmd *cobra.Command, hash, contextDir, traceDir string, jsonOutput bool,
) error {
	fullHash, err := trace.ResolveCommitHash(hash)
	if err != nil {
		fullHash = hash
	}

	refs := trace.CollectRefsForCommit(fullHash, traceDir, true)

	if jsonOutput {
		msg, _ := trace.CommitMessage(fullHash)
		out := JSONCommit{
			Commit:  trace.ShortHash(fullHash),
			Message: msg,
			Refs:    ResolveToJSON(refs, contextDir),
		}
		enc := json.NewEncoder(cmd.OutOrStdout())
		enc.SetIndent("", "  ")
		return enc.Encode(out)
	}

	msg, _ := trace.CommitMessage(fullHash)
	date := trace.CommitDate(fullHash)
	writeTrace.CommitHeader(cmd, trace.ShortHash(fullHash), msg, date)

	if len(refs) == 0 {
		writeTrace.CommitNoContext(cmd)
		return nil
	}

	writeTrace.CommitContext(cmd)
	for _, r := range refs {
		rr := trace.Resolve(r, contextDir)
		writeTrace.Resolved(cmd, rr)
	}

	return nil
}

// Last displays context refs for the last N commits.
//
// Parameters:
//   - cmd: Cobra command for output stream
//   - n: number of commits to show
//   - contextDir: absolute path to the context directory
//   - traceDir: absolute path to the trace directory
//   - jsonOutput: whether to format output as JSON
//
// Returns:
//   - error: non-nil on execution failure
func Last(
	cmd *cobra.Command, n int, contextDir, traceDir string, jsonOutput bool,
) error {
	out, err := git.Run(
		cfgGit.Log, fmt.Sprintf("-%d", n), cfgGit.FormatHashSubj,
	)
	if err != nil {
		return errTrace.GitLog(err)
	}

	lines := strings.Split(strings.TrimSpace(string(out)), token.NewlineLF)

	// Bulk listing uses includeTrailers=false: history.jsonl already
	// contains the same refs the post-commit hook read from the trailer,
	// so re-reading trailers would spawn N extra git processes for no gain.
	if jsonOutput {
		commits := make([]JSONCommit, 0, len(lines))
		for _, line := range lines {
			if line == "" {
				continue
			}
			parts := strings.SplitN(line, " ", 2)
			hash := parts[0]
			msg := ""
			if len(parts) > 1 {
				msg = parts[1]
			}
			refs := trace.CollectRefsForCommit(hash, traceDir, false)
			commits = append(commits, JSONCommit{
				Commit:  trace.ShortHash(hash),
				Message: msg,
				Refs:    ResolveToJSON(refs, contextDir),
			})
		}
		enc := json.NewEncoder(cmd.OutOrStdout())
		enc.SetIndent("", "  ")
		return enc.Encode(commits)
	}

	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, " ", 2)
		hash := parts[0]
		msg := ""
		if len(parts) > 1 {
			msg = parts[1]
		}
		refs := trace.CollectRefsForCommit(hash, traceDir, false)
		refSummary := desc.Text(text.DescKeyWriteTraceNoRefs)
		if len(refs) > 0 {
			refSummary = strings.Join(refs, token.CommaSpace)
		}
		writeTrace.LastEntry(cmd, trace.ShortHash(hash), msg, refSummary)
	}

	return nil
}

// ResolveToJSON converts a slice of raw refs to their JSON representations.
//
// Parameters:
//   - refs: raw context reference strings
//   - contextDir: absolute path to the context directory
//
// Returns:
//   - []JSONRef: resolved references ready for JSON encoding
func ResolveToJSON(refs []string, contextDir string) []JSONRef {
	resolved := make([]JSONRef, 0, len(refs))
	for _, r := range refs {
		rr := trace.Resolve(r, contextDir)
		resolved = append(resolved, JSONRef{
			Raw:    rr.Raw,
			Type:   rr.Type,
			Number: rr.Number,
			Title:  rr.Title,
			Detail: rr.Detail,
			Found:  rr.Found,
		})
	}
	return resolved
}
