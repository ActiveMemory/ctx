//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package project handles **project-root scaffolding** during
// `ctx init` — creating the `.context/` directory tree with
// the right permissions and deploying optional Makefile
// integration when the host project already uses Make.
//
// The package is the *filesystem layer* of init; the foundation
// **content** comes from [internal/cli/initialize/core/merge]
// and [internal/assets].
//
// # Public Surface
//
//   - **[CreateDirs](contextDir)** — creates the
//     `.context/` tree:
//   - `.context/`          (0o755)
//   - `.context/archive/`  for archived tasks/decisions
//   - `.context/state/`    for per-session markers,
//     events, trace history (mode 0o755 — readable by
//     hooks)
//   - `.context/journal/`  for enriched journal entries
//   - `.context/memory/`   for the Claude-Code memory
//     mirror
//   - `.context/steering/` for steering files
//   - `.context/hooks/`    for project-authored
//     lifecycle scripts
//     Idempotent: existing directories are left in place
//     with their existing permissions.
//   - **[HandleMakefileCtx](projectRoot)** — when a
//     `Makefile` already exists at the project root,
//     deploys `Makefile.ctx` from the embedded template
//     so users can run `make ctx-status`, `make
//     ctx-agent`, etc. Skipped when the project has no
//     Makefile (avoids polluting non-Make projects).
//
// # Permissions Rationale
//
// The hooks directory needs `0o755` (not `0o700`) because
// child hook scripts launched by AI tools may inherit
// reduced privileges; making the directory world-readable
// avoids "cannot stat" failures across user/agent
// boundaries. State files are `0o644` for the same
// reason.
//
// # Concurrency
//
// Filesystem-bound and stateless. Concurrent invocations
// against the same root would race on `MkdirAll` writes;
// in practice ctx is single-process.
//
// # Related Packages
//
//   - [internal/cli/initialize]              — top-level
//     orchestrator.
//   - [internal/cli/initialize/core/merge]   — populates
//     the directories with template content.
//   - [internal/config/dir], [internal/config/fs]   —
//     directory-name and permission constants.
package project
