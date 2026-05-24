//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package show implements `ctx audit show ID`.
//
// The command prints the body of an audit report verbatim,
// with no frontmatter and no decoration, suitable for unix
// pipelines (e.g. `ctx audit show surface | less`).
//
// On an unknown id it surfaces [errAudit.UnknownID] so a
// scripted caller can detect the gap and skip.
package show
