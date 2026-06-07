//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package dream

import (
	"context"
	"os/exec"
)

// LookPath resolves the executor binary on PATH, returning the
// resolved absolute path or an error when it is not found. A
// not-found result is the fail-loud signal the caller turns into a
// failmark.
//
// Parameters:
//   - name: the executor binary name (e.g. "claude")
//
// Returns:
//   - string: the resolved path to the binary
//   - error: non-nil when the binary is not on PATH
func LookPath(name string) (string, error) {
	return exec.LookPath(name)
}

// CommandContext returns an exec.Cmd for the resolved executor path
// and its arguments, bound to ctx for timeout/cancellation. The
// caller wires stdout/stderr and working directory.
//
// Parameters:
//   - ctx: context for deadline/cancellation
//   - path: resolved absolute path to the executor binary
//   - args: executor arguments (prompt flag, prompt, budget bound)
//
// Returns:
//   - *exec.Cmd: configured command ready for stream wiring
func CommandContext(
	ctx context.Context, path string, args ...string,
) *exec.Cmd {
	//nolint:gosec // path resolved via LookPath; args from internal config
	return exec.CommandContext(ctx, path, args...)
}
