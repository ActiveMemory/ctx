//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package system

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/notify"
)

// blockNonPathCtxCmd returns the "ctx system block-non-path-ctx" command.
//
// Blocks non-PATH ctx invocations (./ctx, go run ./cmd/ctx, absolute paths)
// to enforce the CONSTITUTION.md rule: "ALWAYS use ctx from PATH".
func blockNonPathCtxCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "block-non-path-ctx",
		Short: "Block non-PATH ctx invocations",
		Long: `Blocks ./ctx, go run ./cmd/ctx, and absolute-path ctx invocations.
Enforces the CONSTITUTION.md rule: always use ctx from PATH.
Outputs a JSON block decision that prevents the tool call.

Hook event: PreToolUse (Bash)
Output: {"decision":"block","reason":"..."} or silent
Silent when: command doesn't invoke ctx via a non-PATH route`,
		Hidden: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runBlockNonPathCtx(cmd, os.Stdin)
		},
	}
}

// blockResponse is the JSON output for blocked commands.
type blockResponse struct {
	Decision string `json:"decision"`
	Reason   string `json:"reason"`
}

// Compiled regex patterns for command-position matching.
var (
	// Pattern 1: ./ctx or ./dist/ctx at start of command
	reRelativeStart = regexp.MustCompile(`^\s*(\./ctx(\s|$)|\./dist/ctx)`)
	// Pattern 1b: ./ctx or ./dist/ctx after command separator
	reRelativeSep = regexp.MustCompile(`(&&|;|\|\||\|)\s*(\./ctx(\s|$)|\./dist/ctx)`)
	// Pattern 2: go run ./cmd/ctx
	reGoRun = regexp.MustCompile(`go run \./cmd/ctx`)
	// Pattern 3: Absolute paths at start of command
	reAbsoluteStart = regexp.MustCompile(`^\s*(/home/|/tmp/|/var/)\S*/ctx(\s|$)`)
	// Pattern 3b: Absolute paths after command separator
	reAbsoluteSep = regexp.MustCompile(`(&&|;|\|\||\|)\s*(/home/|/tmp/|/var/)\S*/ctx(\s|$)`)
	// Exception: /tmp/ctx-test for integration tests
	reTestException = regexp.MustCompile(`/tmp/ctx-test`)
)

func runBlockNonPathCtx(cmd *cobra.Command, stdin *os.File) error {
	input := readInput(stdin)
	command := input.ToolInput.Command

	if command == "" {
		return nil
	}

	var reason string

	// Pattern 1: ./ctx or ./dist/ctx at command position
	if reRelativeStart.MatchString(command) || reRelativeSep.MatchString(command) {
		reason = "Use 'ctx' from PATH, not './ctx' or './dist/ctx'. Ask the user to run: make build && sudo make install"
	}

	// Pattern 2: go run ./cmd/ctx
	if reGoRun.MatchString(command) {
		reason = "Use 'ctx' from PATH, not 'go run ./cmd/ctx'. Ask the user to run: make build && sudo make install"
	}

	// Pattern 3: Absolute paths to ctx binary at command position
	if reason == "" && (reAbsoluteStart.MatchString(command) || reAbsoluteSep.MatchString(command)) {
		if !reTestException.MatchString(command) {
			reason = "Use 'ctx' from PATH, not absolute paths. Ask the user to run: make build && sudo make install"
		}
	}

	if reason != "" {
		resp := blockResponse{
			Decision: "block",
			Reason:   fmt.Sprintf("%s\n\nSee CONSTITUTION.md: ctx Invocation Invariants", reason),
		}
		data, _ := json.Marshal(resp)
		cmd.Println(string(data))
		_ = notify.Send("relay", "block-non-path-ctx: Blocked non-PATH ctx invocation", "")
	}

	return nil
}
