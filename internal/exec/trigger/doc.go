//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package trigger centralizes process execution for
// lifecycle trigger scripts.
//
// # Command Creation
//
// CommandContext wraps exec.CommandContext to create
// a hook process with the given context and script
// path. The context enables timeout enforcement so
// that runaway hook scripts can be cancelled.
//
//	cmd := trigger.CommandContext(ctx, "/path/hook.sh")
//	cmd.Stdin = input
//	cmd.Stdout = output
//	err := cmd.Run()
//
// # Security
//
// The script path is validated by hook.ValidatePath
// before reaching this package. The exec.Command
// call carries a gosec nolint annotation since the
// path is caller-controlled and pre-validated.
//
// # Centralization
//
// All exec.Command calls for trigger runners live
// here, providing a single point for testing and
// security auditing. Callers wire stdin, stdout,
// and stderr on the returned exec.Cmd as needed.
package trigger
