//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package check_anchor_drift implements the
// **`ctx system check-anchor-drift`** hook (added by
// specs/single-source-context-anchor.md).
//
// The hook fires on every UserPromptSubmit. Each operating
// hook line under [internal/assets/claude/hooks/hooks.json]
// exports
// `CTX_DIR="${CLAUDE_PROJECT_DIR:?…}/.context"` inline, which
// **overwrites** the parent shell's CTX_DIR for the hook
// subprocess. That is correct for operating hooks (they must
// write to the right .context/ regardless of what the user
// shell exported), but useless for any hook whose job is to
// *compare* the inherited CTX_DIR against the Claude-injected
// anchor: the comparison would always be tautologically equal.
//
// To break the tautology, this hook's command line in
// hooks.json prepends one extra assignment:
//
//	CTX_DIR_INHERITED="${CTX_DIR:-}" \
//	CTX_DIR="${CLAUDE_PROJECT_DIR:?…}/.context" \
//	ctx system check-anchor-drift
//
// Bash evaluates env-var assignments left-to-right *before*
// invoking the command, so CTX_DIR_INHERITED snapshots the
// parent's CTX_DIR (empty if unset) before the standard
// CTX_DIR injection runs. The hook reads both and emits a
// VERBATIM warning banner only when they disagree.
//
// Behavior matrix:
//
//   - CTX_DIR_INHERITED empty: silent. The user has not run
//     `ctx activate`; there is no shell-level declaration to
//     drift from. Operating hooks still work via the standard
//     injection on every other hook line.
//   - CTX_DIR_INHERITED non-empty and equal to CTX_DIR after
//     [filepath.Clean] on both: silent. Correctly anchored.
//   - CTX_DIR_INHERITED non-empty and unequal to CTX_DIR:
//     emit a warning banner naming both values so the user
//     can see which project's .context/ their CLI /
//     `!`-pragma calls are writing to vs. which project
//     Claude Code is in.
//
// # Public Surface
//
//   - **[Cmd]**: cobra command (hidden under
//     `ctx system`).
//   - **[Run]**: reads [env.CtxDirInherited] and
//     [env.CtxDir] directly via [os.Getenv] (NOT through
//     [rc.ContextDir]), compares, emits the box if drifted.
//
// # Why bypass rc.ContextDir
//
// `rc.ContextDir` is the operating shape validator and would
// reject inherited values that fail the basename guard. This
// hook is a diagnostic: it must accept any observed value
// (including unset, including non-canonical) so it can
// describe reality, not impose policy.
package check_anchor_drift
