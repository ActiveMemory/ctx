//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package proposal holds the constants for the AI proposal
// queue: the per-session directory under `.context/`
// where AI-produced patches land for human ratification,
// and the default slug used when the caller does not
// supply one.
//
// The queue location was decided in DECISIONS.md
// (2026-05-22-220000): `.context/proposals/<TS>-<slug>.md`,
// gitignored by default, mirroring the
// `.context/handovers/` shape.
package proposal
