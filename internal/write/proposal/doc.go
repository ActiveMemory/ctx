//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package proposal writes per-session AI proposal
// artifacts into `.context/proposals/<TS>-<slug>.md`. AI
// commands (`ctx ai extract`, future B/C verbs) emit
// machine-generated patches here for human ratification;
// `.context/*.md` canonical files are never written
// directly.
//
// The on-disk shape mirrors `.context/handovers/`:
// timestamped filename so concurrent runs don't collide;
// gitignored by default per the
// `2026-05-22-220000` DECISIONS entry.
package proposal
