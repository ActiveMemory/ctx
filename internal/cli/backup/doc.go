//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package backup implements the ctx backup top-level command.
//
// Creates timestamped tar.gz archives of project context and/or
// global Claude Code data. Optionally copies archives to an SMB
// share configured via CTX_BACKUP_SMB_URL.
//
// Key exports: [Cmd], [Run].
package backup
