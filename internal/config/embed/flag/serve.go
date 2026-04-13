//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package flag

// Serve command flags currently have no overridable descriptions —
// ctx serve only takes a positional [directory] argument.
//
// Hub server flags (port, data-dir, daemon, peers) live in hub.go
// because they belong to `ctx hub start` / `ctx hub stop`, not
// `ctx serve`.
