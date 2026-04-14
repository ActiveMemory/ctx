//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package check is the **brain** of `ctx doctor`: a battery of
// independent health probes that together produce a single
// diagnostic report covering everything that can plausibly be
// wrong with a ctx installation or with the project's
// `.context/` state.
//
// The package is the only thing the doctor CLI calls on the
// "produce findings" side. The doctor command itself
// orchestrates output; this package decides what to look at
// and how to grade it.
//
// # The Probe Surface
//
// Each probe answers exactly one yes/no question and emits a
// [CheckResult] with a name, status (Ok / Warning / Error),
// and a one-line message. The full battery, run by [Run]:
//
//   - **Context initialization** — `.context/` exists and
//     is populated.
//   - **Required files** — TASKS, DECISIONS, LEARNINGS,
//     CONVENTIONS, ARCHITECTURE, CONSTITUTION present.
//   - **`.ctxrc` validation** — file parses, all values
//     within range.
//   - **Drift** — wraps [internal/drift.Detect] and
//     surfaces the report's status.
//   - **Plugin enablement** — Claude Code plugin
//     installed AND enabled in `~/.claude/settings.json`.
//   - **Event logging** — if `event_log: true`, the log
//     file exists and is writable.
//   - **Reminders** — pending reminder count and freshness.
//   - **Task completion** — open task count, oldest open
//     task age (consolidation nudge threshold).
//   - **Token budgets** — currently injected size against
//     the configured `injection_token_warn` and
//     `context_window`.
//   - **System resource metrics** — wraps
//     [internal/sysinfo] to surface load/memory/disk
//     pressure.
//
// New probes plug in by adding one more entry to the
// dispatch table in [check.go] and one more constant to
// [config/check.Name] (audited to keep CLI output stable).
//
// # Severity Roll-Up
//
// Each probe produces its own status. The doctor CLI rolls
// the slice up to a single banner per the same rule the
// drift package uses: any **Error** beats any **Warning**
// beats **Ok**. JSON output preserves the per-probe detail
// for tooling.
//
// # Stateless and Concurrency-Safe
//
// Probes hold no state and do not coordinate. They could
// be parallelized; they currently run sequentially because
// the slowest probe (`sysinfo` shelling out on macOS) is
// still under 100ms and the simpler ordering keeps output
// deterministic.
//
// # Related Packages
//
//   - [internal/cli/doctor]       — the `ctx doctor` CLI
//     surface that consumes [Run] and renders the report.
//   - [internal/drift]            — the drift detector
//     this package wraps.
//   - [internal/sysinfo]          — the resource probes
//     this package wraps.
//   - [internal/rc]               — supplies thresholds
//     and feature flags.
//   - [internal/config/check]     — probe-name
//     constants used as keys in the report.
package check
