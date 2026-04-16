//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package block_non_path_ctx implements the
// **`ctx system block-non-path-ctx`** hidden hook,
// which blocks invocations of ctx that bypass the
// PATH-installed binary.
//
// # What It Does
//
// The hook reads a JSON envelope from stdin and checks
// the command string against patterns that invoke ctx
// through non-standard paths:
//
//   - **./ctx or ../ctx** -- relative path invocations
//   - **go run ./cmd/ctx** -- source-level execution
//   - **/absolute/path/to/ctx** -- absolute path calls
//     (with an exception for test binaries)
//
// When a match is found, a JSON block response is
// emitted to prevent execution, and a relay
// notification is sent explaining why the invocation
// was rejected. The block reason includes a
// constitution-suffix reminder about using the
// PATH-installed binary.
//
// # Input
//
// A JSON hook envelope on stdin with a ToolInput.Command
// field containing the shell command string.
//
// # Output
//
// On match: a JSON [entity.BlockResponse] with decision
// "block" and a human-readable reason appended with a
// constitution suffix. On no match: no output.
//
// # Delegation
//
// [Cmd] builds the hidden cobra command. [Run] reads
// stdin via [core/session.ReadInput], tests each regex
// pattern, loads the variant message via
// [core/message.Load], and marshals the block response.
// Relay notifications go through [core/nudge.Relay].
package block_non_path_ctx
