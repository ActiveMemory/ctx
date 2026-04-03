//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package trigger

import (
	"context"
	"os/exec"
)

// CommandContext returns an exec.Cmd for a hook script path,
// bound to the given context for timeout enforcement.
//
// Parameters:
//   - ctx: context for deadline/cancellation
//   - path: absolute path to the hook script
//
// Returns:
//   - *exec.Cmd: configured command ready for stdin/stdout wiring
func CommandContext(
	ctx context.Context, path string,
) *exec.Cmd {
	//nolint:gosec // path validated by hook.ValidatePath
	return exec.CommandContext(ctx, path)
}
