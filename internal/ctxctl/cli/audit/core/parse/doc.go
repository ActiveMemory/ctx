//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package parse splits an audit report file into a typed
// [Header] (parsed YAML frontmatter) and the verbatim body
// that follows. The body is returned untouched so the
// `ctxctl audit-relay` hook can drop it inside the
// verbatim-relay envelope without post-processing.
package parse
