//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package trigger implements the "ctx trigger" command
// group for managing lifecycle triggers that fire at
// specific events during AI sessions.
//
// Triggers are executable scripts in .context/triggers/
// that run when named events occur (e.g. session-start,
// session-end, post-commit). They enable custom
// automation without modifying ctx's core hook system.
//
// # Subcommands
//
//   - add: create a new trigger script from a template,
//     targeting a specific event
//   - list: display all installed triggers with their
//     event bindings and enabled state
//   - test: execute a trigger with mock event input
//     for validation
//   - enable: add the executable bit to a trigger
//     script, activating it
//   - disable: remove the executable bit, deactivating
//     the trigger without deleting it
//
// # Subpackages
//
//	cmd/add: trigger creation and templating
//	cmd/list: trigger enumeration
//	cmd/test: mock execution
//	cmd/enable, cmd/disable: activation control
package trigger
