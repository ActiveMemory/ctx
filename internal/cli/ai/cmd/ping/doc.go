//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package ping implements the `ctx ai ping` subcommand,
// a reachability check against the configured AI
// backend's `/v1/models` endpoint. Used by operators to
// confirm `.ctxrc` backend configuration is correct
// after `ctx setup --backend <name>`.
package ping
