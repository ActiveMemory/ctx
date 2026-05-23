//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package resolve builds a [backend.Registry] populated
// with the six built-in factories and the per-project
// backend configurations loaded from `.ctxrc`. Used by
// every `ctx ai *` subcommand to obtain a configured
// backend before issuing a call.
//
// The Registry is built lazily per command invocation so
// rc reloads (e.g., after `ctx setup --backend`) are
// visible without restarting ctx.
package resolve
