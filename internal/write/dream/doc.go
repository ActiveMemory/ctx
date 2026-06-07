//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package dream provides terminal output for the ctx-dream commands
// (ctx dream, dream review, dream accept/reject/amend).
//
// All user-facing strings route through this package so the audit's
// no-cmd.Print-outside-write rule holds. Output is substance-forward:
// the review renders each proposal's id, targets, status, action,
// evidence, confidence, and rationale; the run pass prints a short
// counts digest; the dispositions confirm the applied action and, for
// generative promote/merge, point the user at /ctx-serendipity.
package dream
