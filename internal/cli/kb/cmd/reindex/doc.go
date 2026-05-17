//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package reindex implements `ctx kb reindex`. The command
// refreshes the CTX:KB:TOPICS managed block inside
// `.context/kb/index.md` so the kb landing page enumerates
// current topic folders.
//
// # Refusal contract
//
//   - Landing page missing the CTX:KB:TOPICS block:
//     [github.com/ActiveMemory/ctx/internal/err/kb/cli.ErrReindexMissingBlock].
//
// # Related packages
//
//   - [github.com/ActiveMemory/ctx/internal/cli/kb/core/reindex]
//     hosts [ListTopics] and [RenderBlock].
//   - [github.com/ActiveMemory/ctx/internal/config/kb/cli] supplies
//     the marker and format strings.
//   - [github.com/ActiveMemory/ctx/internal/config/regex] hosts the
//     greedy [regex.ManagedKBTopics] matcher.
package reindex
