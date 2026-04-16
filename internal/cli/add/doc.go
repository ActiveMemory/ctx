//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package add provides the "ctx add" command for appending
// entries to context files.
//
// The add command is the primary write interface for
// populating .context/ files. It accepts content via
// positional argument, --file flag, or stdin pipe and
// routes entries to the appropriate file based on the
// entry type argument.
//
// # Supported Entry Types
//
// Entry types map to [config.FileType] values:
//
//   - decision / decisions: appends to DECISIONS.md
//   - task / tasks: inserts into TASKS.md before the
//     first unchecked item, or under a named section
//     when --section is provided
//   - learning / learnings: appends to LEARNINGS.md
//   - convention / conventions: appends to CONVENTIONS.md
//
// # Example Usage
//
//	ctx add decision "Use PostgreSQL for primary DB"
//	ctx add task "Implement auth" --section "Phase 1"
//	ctx add learning --file notes.md
//	echo "Use camelCase" | ctx add convention
//
// # Subpackages
//
//   - cmd/root: cobra command definition and flag binding
//   - core: file-type routing and content insertion logic
package add
