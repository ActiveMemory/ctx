//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package desc provides command, flag, text, and example description
// lookups backed by embedded YAML.
//
// All user-facing strings are externalized to YAML files loaded at
// init time via [lookup.Init]. The four accessors — [Command],
// [Flag], [Example], and [Text] — resolve DescKey constants to
// their localized values. Missing keys return the key itself as
// a fallback, making gaps visible without crashing.
package desc
