//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package token defines the **string-token vocabulary** ctx
// uses everywhere it concatenates output, parses input, or
// scans content: delimiters (newline, comma, slash), markers
// (frontmatter fence, code fence, ellipsis), prefixes (URL
// schemes, file paths), and small fixed phrases (`Spec:`,
// `Status:`).
//
// Centralizing them eliminates two whole classes of bug:
//
//   - **Typo drift** — `","` vs `", "` vs `" ,"` no longer
//     happen across 40 files; everyone uses
//     [token.CommaSpace].
//   - **Magic-string hunts** — searches for a marker that
//     appears in three places now resolve to one constant
//     declaration with backlinks via `go references`.
//
// The audit suite enforces "no string literal duplication
// across packages" so adding a new common token here is
// the only sustainable path.
//
// # Token Families
//
// Each `*.go` file groups one family:
//
//   - **delimiter** — newlines (`\n` / `\r\n`), commas,
//     spaces, ellipsis, separator runs.
//   - **marker** — Markdown fences, frontmatter fences,
//     ctx HTML markers.
//   - **prefix** — URL schemes (`http://`, `https://`,
//     `ftp://`, `file://`, `//`), absolute path
//     marker.
//   - **slash / quote** — single character constants
//     that have a name to avoid raw `'/'` / `'"'` in
//     calling code.
//   - **content** — common content-pattern groups
//     ([SecretPatterns], [TopicSeparators],
//     [TemplateMarkers]) used by drift, classify, and
//     search.
//
// # Concurrency
//
// All exports are immutable. Safe for any access
// pattern.
package token
