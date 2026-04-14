//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package marker defines the **HTML-comment marker pairs**
// ctx uses to bracket auto-managed sections inside
// otherwise user-edited markdown files. The markers are
// invisible in rendered Markdown and trivially
// `grep`-able.
//
// The package is the single source of truth for these
// strings so [internal/cli/initialize/core/merge] and
// every other consumer (drift detector, doctor, sync)
// reference the same constants.
//
// # The Convention
//
// Every ctx-managed section in a user-editable file is
// bracketed by:
//
//	<!-- ctx:section-name -->
//	... ctx-managed content ...
//	<!-- ctx:section-name:end -->
//
// Edits *outside* the markers survive `ctx init`
// re-runs; edits *inside* are blown away on the next
// re-run. The contract is documented to the user in
// `docs/home/configuration.md` so the destruction is
// not a surprise.
//
// # Marker Pairs Defined Here
//
//   - **ctx:context** — the persistent-context block
//     in `CLAUDE.md`-style files.
//   - **ctx:copilot** — the persistent-context block
//     in `.github/copilot-instructions.md`.
//   - **ctx:agents**  — the equivalent block in
//     `AGENTS.md`.
//   - **ctx:permissions** — the auto-managed allow/
//     deny entries in `settings.local.json`-style
//     comments.
//   - **INDEX:START / INDEX:END** — the
//     auto-generated index table inside
//     DECISIONS.md / LEARNINGS.md.
//
// # Concurrency
//
// All exports are immutable string constants. Safe
// for any access pattern.
//
// # Related Packages
//
//   - [internal/cli/initialize/core/merge] — the
//     marker-aware editor that respects these
//     constants.
//   - [internal/index]                     — uses
//     `INDEX:START`/`INDEX:END` to locate the
//     index table.
//   - [internal/drift]                     — checks
//     marker pairs are intact (header alignment).
package marker
