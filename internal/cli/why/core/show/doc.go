//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package show loads and displays embedded philosophy
// documents for the why command. It bridges the data
// layer (alias resolution) and the strip layer (MkDocs
// cleanup) into a single display call.
//
// # Document Display
//
// [Doc] performs the following steps:
//
//  1. Resolve the user-provided alias to an embedded
//     asset name via data.DocAliases
//  2. Load the embedded document from the philosophy
//     asset reader
//  3. Strip MkDocs-specific syntax (frontmatter,
//     admonitions, tabs, images, relative links) via
//     strip.MkDocs
//  4. Print the cleaned content to the command output
//
// # Error Handling
//
// Doc returns an error when the alias is unknown (not
// in DocAliases) or when the embedded document cannot
// be loaded. Both cases are surfaced to the cmd/ layer
// for user-facing error formatting.
package show
