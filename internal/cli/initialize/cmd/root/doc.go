//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package root implements the "ctx init" command that
// creates and populates a .context/ directory for
// persistent AI context.
//
// # What It Does
//
// Orchestrates the full initialization workflow for
// a new or existing project:
//
//  1. Validates that ctx is in PATH.
//  2. Creates .context/ and its subdirectories
//     (steering/, hooks/, skills/).
//  3. Scaffolds foundation steering files.
//  4. Deploys context file templates (TASKS.md,
//     DECISIONS.md, LEARNINGS.md, CONVENTIONS.md,
//     CONSTITUTION.md, etc.).
//  5. Creates entry templates in .context/templates/.
//  6. Sets up the encrypted scratchpad.
//  7. Creates project root directories (specs/,
//     ideas/).
//  8. Merges permissions into settings.local.json.
//  9. Auto-enables the ctx plugin globally and
//     locally.
//  10. Creates or merges CLAUDE.md.
//  11. Deploys Makefile.ctx and amends the user
//     Makefile.
//  12. Updates .gitignore with recommended entries.
//
// # Flags
//
//   - --force, -f: Overwrite existing context files
//     without prompting for confirmation.
//   - --minimal, -m: Only create essential files
//     (TASKS, DECISIONS, CONSTITUTION).
//   - --merge: Auto-merge ctx content into existing
//     CLAUDE.md without prompting.
//   - --no-plugin-enable: Skip auto-enabling the ctx
//     plugin in Claude settings.
//   - --no-steering-init: Skip scaffolding foundation
//     steering files in .context/steering/.
//
// # Output
//
// Prints progress lines for each created file and
// directory, warnings for non-fatal errors, and a
// final "next steps" guide with workflow tips.
package root
