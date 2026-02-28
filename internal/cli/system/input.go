//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package system

import (
	"encoding/json"
	"io"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// HookInput represents the JSON payload that Claude Code sends to hook
// commands via stdin.
type HookInput struct {
	SessionID string    `json:"session_id"`
	ToolInput ToolInput `json:"tool_input"`
}

// ToolInput contains the tool-specific fields from a Claude Code hook
// invocation. For Bash hooks, Command holds the shell command.
type ToolInput struct {
	Command string `json:"command"`
}

// HookResponse is the JSON output format for Claude Code hooks.
// Using structured JSON ensures the agent processes the output as a directive
// rather than treating it as ignorable plain text.
type HookResponse struct {
	HookSpecificOutput *HookSpecificOutput `json:"hookSpecificOutput,omitempty"`
}

// HookSpecificOutput carries event-specific fields inside a HookResponse.
type HookSpecificOutput struct {
	HookEventName     string `json:"hookEventName"`
	AdditionalContext string `json:"additionalContext,omitempty"`
}

// printHookContext emits a JSON HookResponse with additionalContext for the
// given hook event. This is the standard way for non-blocking hooks to inject
// directives that the agent will actually process (plain text gets ignored).
func printHookContext(cmd *cobra.Command, event, context string) {
	resp := HookResponse{
		HookSpecificOutput: &HookSpecificOutput{
			HookEventName:     event,
			AdditionalContext: context,
		},
	}
	data, _ := json.Marshal(resp)
	cmd.Println(string(data))
}

// readInput reads and parses the JSON hook input from r.
// Returns a zero-value HookInput on any error (graceful degradation).
//
// Guards against blocking forever on stdin:
//   - Terminal (character device): returns immediately
//   - Pipe/file with no EOF within 2s: times out and returns zero value
//
// Both cases are harmless — hooks degrade gracefully with zero input.
func readInput(r io.Reader) HookInput {
	if f, ok := r.(*os.File); ok {
		if fi, err := f.Stat(); err == nil && fi.Mode()&os.ModeCharDevice != 0 {
			return HookInput{}
		}
	}

	type readResult struct {
		data []byte
		err  error
	}
	ch := make(chan readResult, 1)
	go func() {
		data, err := io.ReadAll(r)
		ch <- readResult{data, err}
	}()

	var input HookInput
	select {
	case res := <-ch:
		if res.err == nil {
			_ = json.Unmarshal(res.data, &input)
		}
	case <-time.After(2 * time.Second):
	}
	return input
}
