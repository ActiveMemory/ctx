//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package importer implements **`ctx journal import`**,
// the subcommand that ingests raw AI session files from
// `~/.claude/projects/<slug>/` (and the equivalent paths
// for other tools) into enriched, git-tracked journal
// entries under `.context/journal/`.
//
// # Public Surface
//
//   - **[Cmd]**: cobra command with `--all`,
//     `--regenerate`, `--dry-run`, and
//     `--keep-frontmatter` flags.
//
//   - **[Run]**: three-phase orchestration;
//
//     1. **Plan**: diff the source set against the
//     journal state file ([internal/journal/state])
//     to produce an [entity.ImportPlan]: which
//     sources to create, regenerate, or skip.
//     2. **Confirm**: print the plan and ask for
//     confirmation (skipped under `--dry-run`).
//     3. **Execute**: for each action, parse via
//     [internal/journal/parser], reduce/collapse
//     /normalize, write the entry, update the
//     state file. Locked entries
//     ([internal/cli/journal/core/lock]) are
//     skipped with a notice.
//
// # `--regenerate` Semantics
//
// Without `--regenerate`, only sources that have not
// been imported produce new entries. With
// `--regenerate`, **every** source is re-imported,
// preserving any frontmatter the user added by
// default (`--keep-frontmatter true`). Pass
// `--keep-frontmatter=false` to discard enrichments
// (destructive; the importer warns explicitly).
//
// # Concurrency
//
// Sequential. Concurrent imports against the same
// journal directory would race on state-file writes;
// ctx is single-process.
package importer
