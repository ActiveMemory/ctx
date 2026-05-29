//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package ctxctl is the root of the maintainer-only audit
// channel logic. Everything under internal/ctxctl/ exists to
// serve the ctxctl binary (a separate Go module at
// tools/ctxctl) and must never be imported by the shipped ctx
// binary; a guard test in internal/compliance enforces that
// cmd/ctx's transitive imports exclude this subtree.
//
// The subtree mirrors ctx's package taxonomy so the relocated
// audit logic keeps its familiar shape:
//
//   - cli/audit, cli/checkaudit: cobra command + hook logic
//   - core/audit: relay-body rendering helpers
//   - config/audit: filesystem layout, enums, format strings
//   - err/audit: error constructors and sentinels
//   - write/audit: CLI output helpers
//
// Unlike ctx, this logic holds no hardcoded user copy and
// makes no desc/i18n calls: ctxctl owns its user-facing text
// as plain English Go constants under tools/ctxctl, passed in
// as parameters. There is no French ctxctl.
//
// See specs/ctxctl-bootstrap.md and DECISIONS.md (2026-05-27).
package ctxctl
