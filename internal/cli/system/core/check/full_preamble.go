//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package check

import (
	"os"

	"github.com/ActiveMemory/ctx/internal/cli/system/core/state"
	"github.com/ActiveMemory/ctx/internal/config/warn"
	"github.com/ActiveMemory/ctx/internal/entity"
	logWarn "github.com/ActiveMemory/ctx/internal/log/warn"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// FullPreamble runs the standard hook prelude: verifies ctx is
// initialized, reads hook input (via [Preamble]), checks the pause
// state, and resolves the context and state directories. Every
// daily-throttled ctx hook opens with this sequence; the helper
// collapses the twenty-plus lines of gate / probe / log boilerplate
// that each hook would otherwise repeat verbatim.
//
// Returning ctxDir alongside stateDir lets hooks that need to
// [filepath.Join] paths under the context root skip a second
// [rc.ContextDir] call: the Initialized gate above already proves
// ContextDir succeeds, so re-checking ErrDirNotDeclared is dead code.
//
// Returns ok=false when the hook should bail silently. The bail
// reasons and how callers see them:
//
//   - Uninitialized ctx: silent bail.
//   - state.Initialized resolver failure: logs [warn.StateInitializedProbe]
//     then bails.
//   - Paused session: silent bail.
//   - state.Dir resolver failure: logs [warn.StateDirProbe] then bails.
//   - rc.ContextDir resolver failure after Initialized returned true:
//     logs [warn.ContextDirResolve] then bails. Reachable only if a
//     future ContextDir error is added beyond ErrDirNotDeclared.
//
// Recommended call shape: alias `!ok` as `bailSilently` so the
// intent reads as a deliberate bail rather than swallowed-error
// suppression at every site.
//
//	input, _, ctxDir, stateDir, ok := check.FullPreamble(stdin)
//	bailSilently := !ok
//	if bailSilently {
//		return nil
//	}
//
// The returned sessionID is the [Preamble]-normalized value and falls
// back to [cfgSession.IDUnknown] when the hook input omits it. Prefer
// it over input.SessionID when touching state files keyed by session.
//
// The regular [Preamble] stays available for hooks that do not need
// the Initialized gate or a state directory (e.g. check_reminder,
// which emits provenance unconditionally and gates Initialized inline).
//
// Parameters:
//   - stdin: Standard input for hook JSON.
//
// Returns:
//   - entity.HookInput: Parsed hook input (zero value when ok=false).
//   - string: Normalized session ID (IDUnknown when missing).
//   - string: Absolute context directory; always usable when ok=true.
//   - string: Absolute state directory; always usable when ok=true.
//   - bool: true when the caller should proceed.
func FullPreamble(
	stdin *os.File,
) (entity.HookInput, string, string, string, bool) {
	initialized, initErr := state.Initialized()
	if initErr != nil {
		logWarn.Warn(warn.StateInitializedProbe, initErr)
		return entity.HookInput{}, "", "", "", false
	}
	if !initialized {
		return entity.HookInput{}, "", "", "", false
	}

	ctxDir, ctxErr := rc.ContextDir()
	if ctxErr != nil {
		logWarn.Warn(warn.ContextDirResolve, ctxErr)
		return entity.HookInput{}, "", "", "", false
	}

	input, sessionID, paused := Preamble(stdin)
	if paused {
		return entity.HookInput{}, "", "", "", false
	}

	stateDir, dirErr := state.Dir()
	if dirErr != nil {
		logWarn.Warn(warn.StateDirProbe, dirErr)
		return entity.HookInput{}, "", "", "", false
	}

	return input, sessionID, ctxDir, stateDir, true
}
