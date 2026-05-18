//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package topic implements the `ctx kb topic` parent
// command. `new` is the sole writer of folder-shaped
// topic-page scaffolds under `.context/kb/topics/<slug>/`.
//
// # Subcommands
//
//   - newcmd: scaffolds a new topic-page folder.
//
// # Refusal contract (from newcmd)
//
// An empty placeholder name (one that slugifies to the empty
// string) returns
// [github.com/ActiveMemory/ctx/internal/err/kb/cli.ErrTopicEmptyName].
// A name colliding with an existing topic folder returns
// [github.com/ActiveMemory/ctx/internal/err/kb/cli.TopicExists]
// wrapping the slug and target index path.
//
// # Related packages
//
//   - [github.com/ActiveMemory/ctx/internal/cli/kb/core/topic]
//     hosts [Scaffold] and [Substitute].
//   - [github.com/ActiveMemory/ctx/internal/slug] supplies the
//     slug-derivation primitive ([slug.Path]).
//   - [github.com/ActiveMemory/ctx/internal/config/kb/cli] supplies
//     the substitution tokens and format strings.
package topic
