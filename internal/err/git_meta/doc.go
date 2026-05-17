//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package gitmeta defines the typed error constructors for the
// git-as-architectural-precondition surface (Phase RG). The
// surface has two pieces:
//
//   - The require check
//     ([github.com/ActiveMemory/ctx/internal/git_meta.RequireGitTree])
//     verifies `<projectRoot>/.git` exists; failure surfaces as
//     [ErrMissingGitTree] (sentinel) or a wrapped stat error.
//   - The head resolver
//     ([github.com/ActiveMemory/ctx/internal/git_meta.ResolveHead])
//     resolves the current commit + branch; failure surfaces as
//     [ErrResolveHeadEmpty] (sentinel) or [ResolveHeadFailed]
//     (wrapping the underlying exec error).
//
// # Related packages
//
//   - [github.com/ActiveMemory/ctx/internal/config/git_meta]
//     supplies the sentinel-message + format-string constants.
//   - [github.com/ActiveMemory/ctx/internal/git_meta] is the
//     primary caller; the root command's PersistentPreRunE
//     also calls into these constructors when wrapping the
//     missing-tree error with the failing subcommand name.
package gitmeta
