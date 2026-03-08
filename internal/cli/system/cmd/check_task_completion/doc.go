//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package check_task_completion implements the ctx system
// check-task-completion subcommand.
//
// It counts Edit/Write tool calls and periodically nudges the agent to
// check whether any tasks should be marked done in TASKS.md.
package check_task_completion
