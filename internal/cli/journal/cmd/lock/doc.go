//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package lock implements the ctx journal lock subcommand.
//
// It protects journal entries from being overwritten by export
// --regenerate, marking them as locked in the state file.
package lock
