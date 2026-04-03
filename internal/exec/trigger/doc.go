//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package trigger centralizes process execution for lifecycle trigger
// scripts. All exec.Command calls for trigger runners live here.
//
// [CommandContext] wraps exec.CommandContext to create a hook
// process with the given context and script path, providing
// a single point for testing and security auditing.
package trigger
