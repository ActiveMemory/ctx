//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package initialize provides the **terminal-output
// helpers** the `ctx init` command uses to narrate every
// step of the initialization workflow: directory
// creation, foundation-file deployment, plugin
// detection, settings merge, hook installation, summary.
//
// The package owns ~40 named output functions, one per
// distinct user-visible event. Centralizing them keeps
// the init flow's terminal text consistent and makes
// localization a single-package change when it
// arrives.
//
// All exported functions take a `*cobra.Command` so
// they route through cobra's output stream (which
// tests can wire to a buffer for assertion).
//
// # Function Families
//
//   - **Prompts**: [InfoOverwritePrompt],
//     [InfoAborted] for the "should I overwrite?"
//     dialog.
//   - **Per-file results**: [InfoFileCreated],
//     [InfoExistsSkipped], [InfoMerged] etc., one
//     line per artifact written.
//   - **Plugin / tool detection**:
//     [InfoPluginInstalled],
//     [InfoPluginEnabled], etc.
//   - **Warnings & non-fatal errors**:
//     [InfoWarnNonFatal] for issues the user
//     should know about but that do not abort
//     init.
//   - **Summary**: [Initialized] (the final
//     "ctx is ready, here's what to do next"
//     banner).
//
// # Concurrency
//
// Pure data → io.Writer. cobra serializes
// concurrent writes through its output stream.
package initialize
