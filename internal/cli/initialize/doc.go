//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package initialize implements **`ctx init`** — the first
// command a user runs against a project to bootstrap the
// `.context/` directory, scaffold the foundation files, and
// optionally wire the Claude Code plugin and other tool
// integrations.
//
// `ctx init` is the entry point that turns "a directory" into
// "a ctx-managed project". Its idempotency is a hard
// requirement: running it twice in a row must produce no
// destructive changes — only fresh foundation files where
// they were missing, and merge-aware updates to settings
// files that already exist.
//
// # What `ctx init` Creates
//
// On a clean directory the command produces:
//
//   - **`.context/` tree** — the dir itself plus
//     `archive/`, `state/`, `journal/`, `memory/`,
//     `steering/`, `hooks/` subdirectories with sane
//     permissions ([core/project]).
//   - **Foundation files** — `CONSTITUTION.md`,
//     `TASKS.md`, `DECISIONS.md`, `LEARNINGS.md`,
//     `CONVENTIONS.md`, `ARCHITECTURE.md`, `GLOSSARY.md`,
//     each from a template with the project name
//     interpolated.
//   - **Steering scaffold** — four foundation steering
//     files (`product.md`, `tech.md`, `structure.md`,
//     `workflow.md`) under `.context/steering/`.
//   - **`Makefile.ctx`** — optional; deployed when the
//     project has a `Makefile` so users can `make
//     ctx-status` etc.
//   - **Tool wiring** — Claude Code plugin enablement,
//     Copilot instructions, VS Code tasks, MCP config,
//     etc., depending on what the host environment has
//     installed.
//
// # Sub-Packages
//
//   - **[cmd/root]**         — the cobra command +
//     flag wiring.
//   - **[core/project]**     — directory tree and
//     foundation file creation.
//   - **[core/plugin]**      — Claude Code plugin
//     detection and global enablement.
//   - **[core/claude_check]** — stage-aware detection of
//     Claude Code state used to print contextual
//     guidance during init.
//   - **[core/merge]**       — create-or-merge file
//     operations with marker-bracketed sections so
//     re-running init never clobbers user edits.
//   - **[core/vscode]**      — `.vscode/` workspace
//     artifacts (tasks.json, mcp.json, extensions.json).
//
// # Idempotency Contract
//
// Every action performed by init must satisfy:
//
//  1. **Existing files are merged, not overwritten** — the
//     [core/merge] helpers find the marker pair, replace
//     only the bracketed content, and leave everything
//     else alone.
//  2. **Permissions are deduplicated** — Claude Code
//     `allow`/`deny` lists are merged; existing entries
//     are preserved.
//  3. **Templated values are stable** — the project name
//     interpolation uses `git remote` data when
//     available so re-running produces byte-identical
//     output.
//  4. **No destructive operations without an explicit
//     `--force`** — `init` does not delete or move user
//     files.
//
// # Related Packages
//
//   - [internal/bootstrap]                — wires
//     `ctx init` into the root cobra tree.
//   - [internal/assets]                   — the embedded
//     templates init writes to disk.
//   - [internal/steering]                 — supplies
//     [FoundationFiles] for the steering scaffold.
//   - [internal/cli/setup]                — separate CLI
//     for re-running just the tool-integration step
//     after `ctx init`.
package initialize
