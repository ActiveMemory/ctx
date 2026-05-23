//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package ai groups the `ctx ai` parent command and its
// subcommands (ping, extract, ...). The parent dispatches
// through the AI-backend registry built by
// `internal/cli/ai/core/resolve/` from the user's `.ctxrc`
// backends configuration.
//
// The whole `ctx ai *` surface is optional and fails
// closed: when no backend is configured, AI commands
// surface a typed err/backend sentinel rather than
// degrading silently. The deterministic ctx core
// (`ctx agent`, `ctx status`, ceremony hooks) does not
// import this package or any of its dependencies under
// `internal/backend/`; that invariant is enforced by a
// dedicated audit test (Phase BE Task 8).
package ai
