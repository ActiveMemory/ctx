//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package message provides the **terminal-output
// helpers** the `ctx hook message` CLI surface uses to
// render its `list`, `show`, `edit`, and `reset`
// subcommands' output.
//
// All exported functions take a `*cobra.Command` so
// they route through cobra's output stream.
//
// # Public Surface
//
//   - **[TemplateVars]**: renders a template's
//     placeholder variable list (the `%[1]s`-style
//     positional parameters).
//   - **[CtxSpecificWarning]**: the warning
//     shown when the user tries to override a
//     ctx-specific (non-customizable) message.
//   - **[OverrideCreated]**: the
//     "wrote override at PATH" line `edit` and
//     `reset` print after a write.
//   - **[EditHint]**: the "run `$EDITOR PATH` to
//     edit" hint surfaced by `show` when no
//     override exists yet.
//   - **[SourceOverride] / [SourceDefault]**:
//     "[override]" / "[default]" badges shown
//     next to each message in `list` so users
//     know which entries they have customized.
//
// # Concurrency
//
// Pure data → io.Writer. Concurrent calls
// serialize through cobra's output stream.
package message
