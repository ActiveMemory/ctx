//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package backup provides error constructors for two narrow file
// operations that still live under a historical "backup" label:
//
//   - [Create]: wraps `.bak` file creation during `ctx init --force`.
//   - [CreateArchiveDir], [WriteArchive]: wrap task-archive directory
//     and file write failures under `.context/archive/`.
//   - [ContextDirNotFound]: the bootstrap-path "context dir missing"
//     error.
//
// The former `ctx backup` command (SMB-driven full-project backup)
// was removed; see docs/operations/runbooks/backup-strategy.md for
// the replacement guidance. The package name is kept to avoid
// churning the non-backup callers that still use these generic
// constructors.
package backup
