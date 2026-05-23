//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package backend holds the constants consumed by
// `internal/backend/`: backend type labels, HTTP path
// segments, header names and value prefixes, response
// body size limits, default per-request and cold-start
// timeouts. They live in `config/` (per the
// magic-strings audit and ctx convention) so the
// implementation files in `internal/backend/` stay
// free of string and number literals.
package backend
