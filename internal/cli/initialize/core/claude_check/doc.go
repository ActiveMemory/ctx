//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package claude_check detects the install state of Claude
// Code and the ctx plugin so `ctx init` and `ctx setup
// claude-code` can print **stage-aware** guidance instead of
// dumping every possible setup step at once.
//
// The detector answers four questions in order, with each
// negative answer short-circuiting the cascade:
//
//  1. **Is the `claude` binary on PATH?** If not, suggest
//     installing Claude Code.
//  2. **Is the ctx plugin registered** in
//     `~/.claude/plugins/installed_plugins.json`? If not,
//     suggest `claude plugin install ...`.
//  3. **Is the plugin enabled** globally or in the project's
//     `.claude/settings.local.json`? If not, suggest the
//     enable command.
//  4. **Are MCP, hooks, and slash commands ready?** If not,
//     suggest the missing pieces.
//
// # Public Surface
//
//   - **[State]** — the four-bool detection result plus a
//     [PluginDetails] struct with version, install path,
//     and registration scope.
//   - **[Detect]** — runs the cascade, returns a [State].
//     Pure detection: no installation, no mutation.
//   - **[Details]** — loads rich metadata about the
//     installed plugin (version, marketplace pin,
//     install timestamp). Returns a zero value with
//     `ok == false` when the plugin is not registered.
//
// # Concurrency
//
// All functions are read-only against the user's home
// directory; concurrent calls never race. Results are
// not cached because users frequently install /
// uninstall mid-session and stale-cache bugs are worse
// than the trivial re-read cost.
//
// # Related Packages
//
//   - [internal/cli/initialize]            — top-level
//     `ctx init` orchestrator that consumes [Detect].
//   - [internal/cli/initialize/core/plugin] — the
//     mutator counterpart that flips the enabled bit.
//   - [internal/cli/setup]                 — the
//     `ctx setup claude-code` CLI surface.
package claude_check
