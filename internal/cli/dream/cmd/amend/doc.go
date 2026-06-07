//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package amend wires the "ctx dream amend <id> --action <action>"
// subcommand.
//
// It builds the cobra command and delegates to the dispose core
// logic, which applies the chosen action in place of the proposal's
// recommendation and records the decision as amended.
package amend
