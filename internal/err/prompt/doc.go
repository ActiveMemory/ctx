//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package prompt defines the typed error constructors
// for prompt template operations. These errors fire
// when the CLI loads, lists, or validates embedded
// prompt templates used to generate CLAUDE.md and
// other tool-specific instruction files.
//
// # Domain
//
// Errors fall into two categories:
//
//   - **Template IO** -- an embedded template could
//     not be found, listed, or read.
//     Constructors: [NoTemplate], [ListTemplates],
//     [ReadTemplate].
//   - **Template validation** -- a template is
//     missing required section markers (e.g. ctx
//     or prompt markers). Constructors:
//     [TemplateMissingMarkers], [MarkerNotFound].
//
// # Wrapping Strategy
//
// IO constructors wrap their cause with
// fmt.Errorf %w so callers can inspect the
// underlying error. Validation constructors
// return plain formatted errors. All user-facing
// text is resolved through
// [internal/assets/read/desc].
//
// # Concurrency
//
// Pure constructors. Concurrent callers never race.
package prompt
