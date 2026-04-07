//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package trigger

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	errTrigger "github.com/ActiveMemory/ctx/internal/err/trigger"
	execTrigger "github.com/ActiveMemory/ctx/internal/exec/trigger"
)

// runOne executes a single hook script, enforcing the given
// timeout. It writes inputJSON to the script's stdin and reads
// HookOutput JSON from stdout. Returns an error for non-zero
// exit, timeout, or invalid JSON output.
//
// Parameters:
//   - h: hook metadata including the script path
//   - inputJSON: JSON payload piped to the script's stdin
//   - timeout: maximum execution duration
//
// Returns:
//   - *HookOutput: parsed JSON output from the script
//   - error: timeout, non-zero exit, or JSON parse failure
func runOne(
	h HookInfo, inputJSON []byte, timeout time.Duration,
) (*HookOutput, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := execTrigger.CommandContext(ctx, h.Path)
	cmd.Stdin = bytes.NewReader(inputJSON)

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	if runErr := cmd.Run(); runErr != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return nil, errTrigger.Timeout(timeout)
		}
		return nil, errTrigger.Exit(runErr)
	}

	raw := strings.TrimSpace(stdout.String())
	if raw == "" {
		return &HookOutput{}, nil
	}

	var out HookOutput
	if jsonErr := json.Unmarshal([]byte(raw), &out); jsonErr != nil {
		return nil, errTrigger.InvalidJSONOutput(jsonErr)
	}

	return &out, nil
}
