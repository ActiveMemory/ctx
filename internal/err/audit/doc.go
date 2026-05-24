//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package audit holds error constructors for the audit
// channel: report-file I/O, frontmatter parsing, dismissal
// ledger I/O, and CLI input validation. All messages route
// through the i18n descriptor pipeline; see
// `internal/config/embed/text/err_audit.go` for keys and
// `internal/assets/commands/text/errors.yaml` for the
// English bodies.
package audit
