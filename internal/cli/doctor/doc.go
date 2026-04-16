//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package doctor implements **`ctx doctor`**, the
// one-stop structural-health command users (and
// onboarding scripts) run when something feels off:
// hooks not firing, drift accumulating, plugin not
// enabled, settings file half-merged, etc.
//
// The doctor is a *shell*: it asks
// [internal/cli/doctor/core/check] for the full battery
// of probes, then renders the results in either
// human-readable checklist form (default) or structured
// JSON form (`--json`).
//
// # Default Output
//
// The checklist groups probes by category (Setup,
// Context, Plugin, State, Resources) and renders each
// with a status glyph (`✓`, `⚠`, `✗`) plus a one-line
// message. The roll-up banner at the end summarizes:
// "all good", "warnings present", or "violations
// present", matching the same severity ladder
// [internal/drift] uses.
//
// # JSON Output
//
// `ctx doctor --json` emits one record per probe with
// `name`, `status`, `message`, and any structured
// detail. Used by CI and by the `_ctx-doctor` skill
// when the AI is the consumer.
//
// # Exit Codes
//
//   - **0**: all checks passed.
//   - **1**: warnings present.
//   - **3**: violations present (so CI scripts can
//     gate on `>= 3`).
//
// # Sub-Packages
//
//   - **[core/check]**: the actual probe battery
//     (no UI, no CLI parsing).
package doctor
