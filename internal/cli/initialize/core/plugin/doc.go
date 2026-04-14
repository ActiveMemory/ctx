//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package plugin handles **Claude Code plugin enablement**
// during `ctx init` — the read/write side of the same
// settings files that
// [internal/cli/initialize/core/claude_check] only reads.
//
// Claude Code keeps two layers of plugin state:
//
//   - **Global** — `~/.claude/settings.json`'s
//     `enabledPlugins` map. Affects every project on the
//     machine.
//   - **Local** — `<project>/.claude/settings.local.json`'s
//     `enabledPlugins` map. Affects only this project.
//
// Both can independently mark a plugin as enabled. ctx
// prefers global enablement so users do not have to
// re-flip the bit per project, but supports local-only
// enablement for users who segment configs.
//
// # Public Surface
//
//   - **[Installed](pluginID)** — true when the plugin
//     binary is registered in
//     `~/.claude/plugins/installed_plugins.json`.
//   - **[EnabledGlobally](pluginID)** — true when the
//     plugin is enabled in the global settings file.
//   - **[EnabledLocally](projectRoot, pluginID)** —
//     true when the plugin is enabled in the project's
//     local settings file.
//   - **[EnableGlobally](pluginID)** — atomically merges
//     the plugin into the global `enabledPlugins` map.
//     Idempotent. Creates the settings file if missing.
//
// # Settings-File Editing Contract
//
// All writes are **JSON-merge-aware**: existing keys are
// preserved, only `enabledPlugins.<pluginID>` is touched.
// A pre-write backup (`.bak`) is created so a manual
// rollback is one `mv` away.
//
// # Concurrency
//
// Filesystem-bound and stateless; serialized through
// process-level execution.
//
// # Related Packages
//
//   - [internal/cli/initialize]                — top-level
//     orchestrator.
//   - [internal/cli/initialize/core/claude_check] —
//     read-only detection.
//   - [internal/config/claude]                 — settings
//     file paths and JSON key constants.
package plugin
