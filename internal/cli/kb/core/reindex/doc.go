//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package reindex carries the topic-enumeration and managed-
// block-rendering helpers used by `ctx kb reindex` to refresh
// the CTX:KB:TOPICS section of `.context/kb/index.md`.
//
// # Files
//
//   - topic.go: ListTopics, the recursive scan that returns topic
//     slugs (slash-separated for grouped layouts) whose index.md
//     exists.
//   - scan.go: unexported walk + group-landing-exclusion helpers
//     behind ListTopics.
//   - block.go: rendering of the managed block contents.
//
// # Related packages
//
//   - [github.com/ActiveMemory/ctx/internal/cli/kb/cmd/reindex]
//     is the CLI surface that drives this core.
//   - [github.com/ActiveMemory/ctx/internal/config/kb/cli] supplies
//     the marker strings.
//   - [github.com/ActiveMemory/ctx/internal/config/regex.ManagedKBTopics]
//     is the matcher.
package reindex
