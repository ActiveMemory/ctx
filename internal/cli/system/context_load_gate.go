//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package system

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/notify"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// contextLoadGateCmd returns the "ctx system context-load-gate" command.
//
// Emits a one-shot directive on the first tool use of a session, telling the
// agent to read context files before proceeding. Uses PreToolUse hook timing
// so the directive arrives at the moment of action — not before (when it
// competes with the user's question) and not after (when the agent is already
// deep in task work).
//
// Design rationale: the directive lists file paths directly instead of
// delegating to `ctx system bootstrap`. Experiments showed that delegation
// chains ("run X and follow its output") lose authority at each link — the
// agent runs the command but skips the follow-up steps. A direct file list
// is a single-hop instruction with no authority decay.
//
// The message uses imperative framing ("STOP. You must read these files")
// because experiments showed that advisory framing ("Read your context files
// before proceeding") invites the agent to assess relevance and skip files
// it deems unnecessary — defeating the gate's purpose.
//
// Compliance checkpoint: the agent must ALWAYS output a "Context Loaded"
// block with "Read" and "Skipped" fields. This is unconditional — there is
// no separate skip path. The fill-in-the-blank template is lower friction
// than evaluating a conditional, so models are much more likely to complete
// it. Observable evidence: block present (compliant) or absent (single
// failure mode, addressable via CONSTITUTION.md).
//
// GLOSSARY.md is excluded from the gate — it's reference material for
// lookup, not context that must be loaded at session start.
func contextLoadGateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "context-load-gate",
		Short: "Emit context-load directive on first tool use",
		Long: `Emits a one-shot directive telling the agent to read context files.
Fires on the first tool use per session via PreToolUse hook. Subsequent
tool calls in the same session are silent (tracked by session marker file).

Lists file paths directly — no delegation to bootstrap command.
See specs/context-load-gate.md for design rationale.

Hook event: PreToolUse (.*)
Output: JSON HookResponse (additionalContext) on first tool use, silent otherwise
Silent when: marker exists for session_id, or context not initialized`,
		Hidden: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runContextLoadGate(cmd, os.Stdin)
		},
	}
}

func runContextLoadGate(cmd *cobra.Command, stdin *os.File) error {
	if !isInitialized() {
		return nil
	}

	input := readInput(stdin)
	if input.SessionID == "" {
		return nil
	}

	tmpDir := secureTempDir()
	marker := filepath.Join(tmpDir, "ctx-loaded-"+input.SessionID)

	if _, err := os.Stat(marker); err == nil {
		return nil // already fired this session
	}

	// Create marker before emitting — ensures one-shot even if
	// the agent makes multiple parallel tool calls.
	touchFile(marker)

	dir := rc.ContextDir()
	var files []string
	for _, f := range config.FileReadOrder {
		if f == config.FileGlossary {
			continue // reference material, not required at load time
		}
		files = append(files, filepath.Join(dir, f))
	}
	fileList := strings.Join(files, ", ")

	msg := fmt.Sprintf(
		"STOP. You must read these files in order before proceeding:\n\n"+
			"%s\n\n"+
			"Do not assess relevance. Do not skip. Read all of them.\n\n"+
			"AFTER reading, respond to the user with this block:\n\n"+
			"┌─ Context Loaded ─────────────────────────────────\n"+
			"│ Read: [list files you read]\n"+
			"│ Skipped: [list files you skipped, or 'none']\n"+
			"└───────────────────────────────────────────────────\n\n"+
			"This block is MANDATORY in your next response regardless\n"+
			"of whether you read all files or skipped any.",
		fileList,
	)

	printHookContext(cmd, "PreToolUse", msg)
	_ = notify.Send("relay", "context-load-gate: directed agent to read context files (unconditional checkpoint)", input.SessionID, msg)
	return nil
}
