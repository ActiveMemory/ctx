//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/activate/core/emit"
	"github.com/ActiveMemory/ctx/internal/cli/activate/core/resolve"
	"github.com/ActiveMemory/ctx/internal/config/env"
	cfgShell "github.com/ActiveMemory/ctx/internal/config/shell"
	writeActivate "github.com/ActiveMemory/ctx/internal/write/activate"
)

// Run executes the `ctx activate` command.
//
// Resolves the target .context/ directory via [resolve.Selected]
// (always project-local scan from CWD under the single-source-anchor
// model), then prints the shell-specific export statement for
// CTX_DIR to stdout.
//
// # Output shape
//
// Two channels:
//
//  1. **stdout**: consumed by `eval "$(ctx activate)"`. Every
//     byte must be valid POSIX shell. Composed in order:
//     (a) zero or one `# ctx: replacing stale CTX_DIR=<old>\n`
//     comment line when the parent shell already has [env.CtxDir]
//     set to a different value than the resolved target;
//     (b) the shell-specific `export CTX_DIR=<value>\n` line.
//
//  2. **stderr**: informational advisories for the user. Always
//     carries a `ctx activated at: <path>` line announcing the
//     bound directory (single-candidate case included), and
//     additionally one `ctx: also visible upward: <path>` line
//     per other `.context/` candidate when more than one is
//     visible upward. `eval` does not capture stderr, so these
//     lines pass through to the terminal where the user sees
//     them. Innermost wins (matches git/make nested-project
//     semantics); the additional candidates are reported, not
//     refused. The comment-on-stdout approach considered
//     earlier was invisible to the only documented invocation
//     form (`eval`), so it informed nobody.
//
// Parameters:
//   - cmd:   cobra command providing stdout / stderr. Nil is a
//     no-op via [writeActivate.Emit] / [writeActivate.AlsoVisible].
//   - shell: value of the --shell flag; empty means auto-detect
//     from $SHELL via [emit.DetectShell].
//
// Returns:
//   - error: non-nil on resolution failure (no `.context/` visible
//     from CWD upward); nil on successful emit.
func Run(cmd *cobra.Command, shell string) error {
	selected, others, err := resolve.Selected()
	if err != nil {
		return err
	}
	out := emit.Set(emit.DetectShell(shell), selected)
	if existing := os.Getenv(env.CtxDir); existing != "" && existing != selected {
		out = fmt.Sprintf(cfgShell.FormatStaleReplaceComment,
			env.CtxDir, existing, out)
	}
	writeActivate.ActivatedAt(cmd, selected)
	writeActivate.AlsoVisible(cmd, others)
	writeActivate.Emit(cmd, out)
	return nil
}
