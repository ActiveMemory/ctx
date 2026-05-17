//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package topic carries the topic-page scaffolding helpers
// used by `ctx kb topic new`: slug normalisation, template
// substitution, and the disk-writing scaffolder.
//
// # Files
//
//   - scaffold.go: [Scaffold] orchestrates the topic-folder
//     create, embedded-template copy, and substitution pass.
//     Slug derivation delegates to
//     [github.com/ActiveMemory/ctx/internal/slug.Path].
//   - template.go: [Substitute] applies <NAME> / <SLUG> token
//     replacement on the rendered topic-page body.
//
// # Related packages
//
//   - [github.com/ActiveMemory/ctx/internal/cli/kb/cmd/topic]
//     is the CLI surface that drives this core.
//   - [github.com/ActiveMemory/ctx/internal/config/kb/cli] supplies
//     the substitution tokens and format strings.
package topic
