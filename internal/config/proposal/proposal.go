//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package proposal

// Subdir is the proposal-queue subdirectory under
// `.context/`. Files inside are named `<TS>-<slug>.md`
// so concurrent agent runs never overwrite each other,
// matching the handovers/ layout.
const Subdir = "proposals"

// DefaultSlug is the fallback slug for proposals when
// the caller does not supply one (e.g., `ctx ai extract`
// without a --slug flag). Kept short and stable so the
// filename's discriminating component is the timestamp.
const DefaultSlug = "extract"
