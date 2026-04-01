//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package show

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/trace"
)

// JSONRef is the JSON representation of a resolved context reference.
type JSONRef struct {
	Raw    string `json:"raw"`
	Type   string `json:"type"`
	Number int    `json:"number,omitempty"`
	Title  string `json:"title,omitempty"`
	Detail string `json:"detail,omitempty"`
	Found  bool   `json:"found"`
}

// JSONCommit is the JSON representation of a commit with its context refs.
type JSONCommit struct {
	Commit  string    `json:"commit"`
	Message string    `json:"message"`
	Refs    []JSONRef `json:"refs"`
}

// ShowCommit displays the context refs for a single commit.
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
func ShowCommit(cmd *cobra.Command, hash, contextDir, traceDir string, jsonOutput bool) error {
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
	cmd.Println(fmt.Sprintf("Commit:  %s", trace.ShortHash(fullHash)))
	cmd.Println(fmt.Sprintf("Message: %s", msg))
	cmd.Println(fmt.Sprintf("Date:    %s", date))

	if len(refs) == 0 {
		cmd.Println("Context: (none)")
		return nil
	}

	cmd.Println("Context:")
	for _, r := range refs {
		rr := trace.Resolve(r, contextDir)
		PrintResolved(cmd, rr)
	}

	return nil
}

// ShowLast displays context refs for the last N commits.
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
func ShowLast(cmd *cobra.Command, n int, contextDir, traceDir string, jsonOutput bool) error {
	//nolint:gosec // n is a user-supplied integer flag, not arbitrary input
	out, err := exec.Command("git", "log", fmt.Sprintf("-%d", n), "--format=%H %s").Output()
	if err != nil {
		return fmt.Errorf("git log: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(out)), token.NewlineLF)

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
		refs := trace.CollectRefsForCommit(hash, traceDir, true)
		refSummary := "(none)"
		if len(refs) > 0 {
			refSummary = strings.Join(refs, ", ")
		}
		cmd.Println(fmt.Sprintf("%s  %s  [%s]", trace.ShortHash(hash), msg, refSummary))
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

// PrintResolved formats a single resolved ref for display.
//
// Parameters:
//   - cmd: Cobra command for output stream
//   - rr: resolved reference to display
func PrintResolved(cmd *cobra.Command, rr trace.ResolvedRef) {
	typeLabel := strings.ToUpper(rr.Type[0:1]) + rr.Type[1:]
	if rr.Found && rr.Title != "" {
		if rr.Detail != "" {
			cmd.Println(fmt.Sprintf("  [%s] %s — %s (%s)", typeLabel, rr.Raw, rr.Title, rr.Detail))
		} else {
			cmd.Println(fmt.Sprintf("  [%s] %s — %s", typeLabel, rr.Raw, rr.Title))
		}
	} else {
		cmd.Println(fmt.Sprintf("  [%s] %s", typeLabel, rr.Raw))
	}
}
