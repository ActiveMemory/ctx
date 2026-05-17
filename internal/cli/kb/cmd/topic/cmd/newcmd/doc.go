//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package newcmd implements the `new` leaf under
// `ctx kb topic`. Calls into the kb/core/topic helpers to
// scaffold a topic-page folder from the embedded template.
//
// # Refusal contract
//
//   - Empty / placeholder name reducing to an empty slug:
//     [github.com/ActiveMemory/ctx/internal/err/kb/cli.ErrTopicEmptyName].
//   - Topic already exists:
//     [github.com/ActiveMemory/ctx/internal/err/kb/cli.TopicExists].
//
// # Related packages
//
//   - [github.com/ActiveMemory/ctx/internal/cli/kb/core/topic]
//     hosts [topic.Scaffold] and [topic.Substitute].
//   - [github.com/ActiveMemory/ctx/internal/slug] supplies the
//     slug-derivation primitive ([slug.Path]).
package newcmd
