//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"github.com/spf13/cobra"

	aiRun "github.com/ActiveMemory/ctx/internal/cli/ai/core/run"
	writeAI "github.com/ActiveMemory/ctx/internal/write/ai"
)

// RunPing executes ctx ai ping.
//
// Parameters:
//   - cmd: cobra command for output stream
//   - backendName: optional backend selector
//
// Returns:
//   - error: backend resolution or ping failure
func RunPing(cmd *cobra.Command, backendName string) error {
	result, pingErr := aiRun.Ping(cmd.Context(), backendName)
	if pingErr != nil {
		return pingErr
	}
	return writeAI.Ping(
		cmd.OutOrStdout(),
		result.Backend,
		result.Endpoint,
		result.FirstModel,
	)
}

// RunPropose executes ctx ai propose.
//
// Parameters:
//   - cmd: cobra command for output stream
//   - input: input file path
//   - backendName: optional backend selector
//   - emit: comma-separated emit kinds
//
// Returns:
//   - error: backend, completion, validation, or artifact write failure
func RunPropose(
	cmd *cobra.Command,
	input string,
	backendName string,
	emit string,
) error {
	_, proposeErr := aiRun.Propose(cmd.Context(), input, backendName, emit)
	return proposeErr
}
