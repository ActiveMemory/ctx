//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package resume implements the "ctx hook resume" command
// that re-enables context hooks after a pause.
//
// When hooks are paused (via ctx hook pause), all system
// hooks stop firing for the current session. The resume
// command clears the pause flag in .context/state/ so
// that hooks begin firing again on subsequent AI
// prompts.
//
// Resume is idempotent: calling it when hooks are
// already active is a no-op that succeeds silently.
//
// # Subpackages
//
//	cmd/root: cobra command definition and pause
//	  state clearing
//
// [Cmd] returns the cobra command that clears the pause
// flag from the session state directory, allowing hooks
// to fire again on subsequent prompts.
package resume
