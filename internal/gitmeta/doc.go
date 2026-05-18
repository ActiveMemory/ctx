//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package gitmeta enforces git as an architectural precondition
// for ctx, and resolves git HEAD into commit + branch for
// provenance metadata.
//
// Two surfaces:
//
//   - [RequireGitTree] returns nil when <projectRoot>/.git exists
//     as a directory (regular repo) or a regular file (worktree
//     pointer). Returns [*MissingGitError] otherwise. Wired into
//     the root command PersistentPreRunE so every non-exempt
//     subcommand inherits the precondition.
//
//   - [ResolveHead] reads HEAD into ([HeadRef]) for closeout and
//     handover provenance. Honors environment overrides
//     CTX_TASK_COMMIT and GITHUB_SHA (when GITHUB_ACTIONS=true)
//     to support CI replay scenarios.
//
// Phase RG (per specs/require-git.md) promotes git from a de
// facto invariant to a de jure one. There is no auto-git-init;
// the user runs git init first, then ctx init.
//
// # Related packages
//
//   - [github.com/ActiveMemory/ctx/internal/exec/git] runs git
//     subcommands; ResolveHead shells through it.
//   - [github.com/ActiveMemory/ctx/internal/config/git] supplies
//     constants for the git binary name, subcommands, and the
//     ".git" directory name (DotDir).
//   - [github.com/ActiveMemory/ctx/internal/err/git] supplies
//     typed errors for git-not-installed and not-in-repo
//     conditions, distinct from gitmeta's MissingGitError which
//     is specifically about the .git tree at a known root.
package gitmeta
