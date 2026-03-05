//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package memory bridges Claude Code's auto memory (MEMORY.md) into
// the .context/ directory with discovery, mirroring, and drift detection.
//
// Claude Code maintains per-project auto memory at
// ~/.claude/projects/<slug>/memory/MEMORY.md. This package locates that
// file from the project root, mirrors it into .context/memory/mirror.md
// (git-tracked), and archives previous versions before each sync.
//
// Discovery encodes the project root path into the Claude Code slug
// format: absolute path with "/" replaced by "-", prefixed with "-".
//
// Sync state is tracked in .context/state/memory-import.json to support
// drift detection and future import/publish phases.
package memory
