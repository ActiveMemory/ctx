//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package ask implements `ctx kb ask "<question>"`.
//
// The Q&A pass itself is performed by the /ctx-kb-ask skill;
// this CLI surface validates input and prints the canonical
// invocation so non-Claude users see the workflow as well.
//
// # Refusal contract
//
//   - Empty question (no positional args or whitespace-only):
//     [github.com/ActiveMemory/ctx/internal/err/kb/cli.ErrAskNoQuestion].
//
// # Related packages
//
//   - [github.com/ActiveMemory/ctx/internal/config/kb/cli] supplies
//     the hint and refusal strings.
//   - [github.com/ActiveMemory/ctx/internal/err/kb/cli] supplies the
//     sentinel constructors.
package ask
