//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package add wires the cobra "ctx task add" subcommand.
//
// The cobra wiring delegates entirely to the shared add core
// at internal/cli/add/core/build, which itself calls into
// internal/cli/add/core/run for the validation, extraction,
// formatting, and insertion pipeline. The package exists only
// to keep the noun-first command tree symmetric with siblings
// like archive, complete, and snapshot.
package add
