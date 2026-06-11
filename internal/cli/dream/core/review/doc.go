//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package review lists the pending proposals from the latest dream
// run and renders them for a serendipity round.
//
// Pending means not yet recorded in the ledger (dedup-against-seen):
// a proposal already accepted, rejected, amended, or skipped is
// filtered out. Each surviving proposal is rendered substance-forward
// — id, targets, status, action, evidence, confidence, rationale.
package review
