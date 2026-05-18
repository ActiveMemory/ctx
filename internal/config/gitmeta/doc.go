//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package gitmeta supplies the environment-variable names,
// special branch-name literals, and short-SHA length used by
// [github.com/ActiveMemory/ctx/internal/gitmeta]. The constants
// live here (not in the gitmeta package itself) to honor the
// project-wide rule that magic strings and values belong in
// internal/config/.
//
// # Related packages
//
//   - [github.com/ActiveMemory/ctx/internal/gitmeta] is the
//     primary consumer.
//   - [github.com/ActiveMemory/ctx/internal/config/git] supplies
//     git binary + subcommand constants; gitmeta complements it
//     with provenance-resolution constants.
package gitmeta
