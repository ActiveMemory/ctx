//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package ai writes ctx ai command output.
//
// It owns stdout formatting for the ai command family. Command packages
// delegate here so they do not print directly.
// The package has no business logic and only formats terminal output.
package ai
