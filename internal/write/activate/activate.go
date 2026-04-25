//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package activate

import (
	"fmt"

	"github.com/spf13/cobra"

	cfgShell "github.com/ActiveMemory/ctx/internal/config/shell"
)

// Emit writes pre-formatted shell-eval content to cmd's stdout
// without adding a trailing newline. The emit-layer functions in
// [internal/cli/activate/core/emit] already include the newline
// they need, so this helper must not add another (a stray blank
// line in `eval` output is harmless but ugly in `set -x` traces).
//
// Parameters:
//   - cmd:     cobra command providing the stdout sink. Nil is a
//     no-op so test setups that omit the command don't crash.
//   - content: shell-eval line(s); may be empty (no-op).
func Emit(cmd *cobra.Command, content string) {
	if cmd == nil || content == "" {
		return
	}
	_, _ = fmt.Fprint(cmd.OutOrStdout(), content)
}

// ActivatedAt writes a single informational line to stderr
// announcing the bound `.context/` path. Always called by
// `ctx activate` on success (single-candidate too) so the user
// always sees what just happened, not just an empty terminal.
//
// Stderr (not stdout) because the line is for the user, not the
// shell. `eval` lets stderr pass through to the terminal while
// stripping the eval-captured stdout stream.
//
// Parameters:
//   - cmd:  cobra command providing the stderr sink. Nil is a
//     no-op.
//   - path: absolute path of the bound `.context/` directory.
//     Empty is a no-op (defensive; Run never calls with empty).
func ActivatedAt(cmd *cobra.Command, path string) {
	if cmd == nil || path == "" {
		return
	}
	// ErrOrStderr (not OutOrStderr): cobra's OutOrStderr returns
	// the SetOut writer with stderr fallback (confusingly named).
	// Wrong helper would land the advisory inside the
	// eval-captured stream and make it invisible to anyone
	// running `eval "$(ctx activate)"`.
	_, _ = fmt.Fprintf(cmd.ErrOrStderr(),
		cfgShell.FormatActivatedAtAdvisory, path)
}

// AlsoVisible writes one informational line per additional
// `.context/` candidate to stderr. Used by `ctx activate` when
// more than one candidate is visible upward from CWD: innermost
// wins (the bind goes to stdout via [Emit] and is announced via
// [ActivatedAt]), and the others get surfaced here so the user
// can see what's around but isn't being bound.
//
// Each line follows the shape:
//
//	ctx: also visible upward: <path>
//
// Multiple paths produce multiple lines (one per path) so the
// output stays parseable when anyone scripts around it.
//
// Parameters:
//   - cmd:   cobra command providing the stderr sink. Nil is a
//     no-op.
//   - paths: additional candidates to surface, in the order they
//     came back from the upward scan. Empty / nil is a no-op.
func AlsoVisible(cmd *cobra.Command, paths []string) {
	if cmd == nil || len(paths) == 0 {
		return
	}
	for _, p := range paths {
		_, _ = fmt.Fprintf(cmd.ErrOrStderr(),
			cfgShell.FormatAlsoVisibleAdvisory, p)
	}
}
