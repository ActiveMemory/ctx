//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package check_anchor_drift

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/anchor"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/nudge"
	coreSession "github.com/ActiveMemory/ctx/internal/cli/system/core/session"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/env"
	"github.com/ActiveMemory/ctx/internal/config/hook"
)

// Run executes the check-anchor-drift hook logic.
//
// Reads the parent-shell CTX_DIR (snapshotted into
// [env.CtxDirInherited] before the standard hook injection)
// and the Claude-injected [env.CtxDir]. Emits a VERBATIM
// warning banner only when both are non-empty and refer to
// genuinely different directories. The banner goes through the
// standard nudge+relay path so the event is recorded in the
// local audit log.
//
// Bypasses [rc.ContextDir]: this is a diagnostic, not an
// operating command. It must accept any observed value
// (including unset, including non-canonical) so it can
// describe reality rather than impose policy.
//
// Symlink-equivalent paths are treated as the same directory
// via [anchor.Equal]. See its package doc for the rationale
// (the canonical case is macOS's `/tmp` → `/private/tmp`).
//
// Parameters:
//   - cmd: cobra command for output. Nil is a no-op.
//   - stdin: standard input for the hook JSON envelope.
//
// Returns:
//   - error: always nil. Diagnostics never fail the hook.
func Run(cmd *cobra.Command, stdin *os.File) error {
	inherited := os.Getenv(env.CtxDirInherited)
	if inherited == "" {
		// No shell-level declaration to drift from.
		return nil
	}
	injected := os.Getenv(env.CtxDir)
	if anchor.Equal(inherited, injected) {
		// Correctly anchored (possibly via symlink-equivalent paths).
		return nil
	}

	content := fmt.Sprintf(
		desc.Text(text.DescKeyCheckAnchorDriftContent),
		inherited, injected,
	)
	input := coreSession.ReadInput(stdin)
	return nudge.Emit(cmd, content,
		desc.Text(text.DescKeyCheckAnchorDriftRelayPrefix),
		desc.Text(text.DescKeyCheckAnchorDriftBoxTitle),
		hook.CheckAnchorDrift, hook.VariantNudge,
		desc.Text(text.DescKeyCheckAnchorDriftRelayMessage),
		input.SessionID, nil, "",
	)
}
