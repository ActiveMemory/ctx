//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package drift detects ways a project's `.context/` files have
// drifted from the codebase or from project conventions, and
// surfaces the findings as a structured [Report] that the CLI
// (`ctx drift`), the doctor (`ctx doctor`), and the steering /
// trigger nudges all consume.
//
// The package is the *evaluator*; it never modifies state. It
// reads the loaded [entity.Context], runs a battery of
// independent checks, and returns a categorized list of issues.
// Whether a given issue stops the user, prints a yellow nudge,
// or is silently archived is the caller's concern.
//
// # The Public Surface
//
// One function does the work — [Detect](ctx) — plus the result
// types it returns:
//
//   - [Report]        — Warnings, Violations, Passed checks.
//   - [Issue]         — File, Line, Type, Message, Path, Rule.
//   - [Report.Status] — rolls the report up to a single
//     [config/drift.StatusType]: Violation > Warning > Ok.
//
// Everything else in the package is an internal `check*` helper
// that appends to the [Report] passed by reference.
//
// # The Checks
//
// [Detect] runs the following checks in order; each is
// independent and contributes to the same [Report]:
//
//   - **Path references** ([checkPathReferences]) — scans
//     ARCHITECTURE.md and CONVENTIONS.md for backtick-enclosed
//     file paths and verifies each exists on disk. Skips URLs,
//     glob patterns, and template placeholders.
//   - **Staleness** ([checkStaleness]) — flags content that
//     contradicts current code (placeholder markers left in
//     CONSTITUTION.md, missing `.context/` markers, etc).
//   - **Constitution heuristics** ([checkConstitution]) —
//     basic rule presence checks against CONSTITUTION.md.
//   - **Required files** ([checkRequiredFiles]) — flags empty
//     files that the schema expects to be populated.
//   - **File age** ([checkFileAge]) — warns when a context
//     file has not been touched in `stale_age_days` (configured
//     in `.ctxrc`; default 30; 0 disables). [staleAgeExclude]
//     skips files that are intentionally static (CONSTITUTION).
//   - **Entry counts** ([checkEntryCount]) — warns when
//     DECISIONS.md / LEARNINGS.md exceed the per-file
//     thresholds (consolidation nudge).
//   - **Missing internal packages** ([checkMissingPackages]) —
//     flags packages mentioned in ARCHITECTURE.md that no
//     longer exist on disk; also normalizes Go internal
//     package paths via [normalizeInternalPkg].
//   - **Template headers** ([checkTemplateHeaders]) — checks
//     each context file's comment-header banner against the
//     ctx-managed template; mismatch suggests `ctx init
//     --force`.
//   - **Steering tools** ([checkSteeringTools]) — every
//     steering file's `tools:` field must reference a
//     supported tool ID ([supportedTools]).
//   - **Hook permissions** ([checkHookPerms]) — flags any
//     trigger script in `.context/hooks/` that lacks the
//     executable bit (matches the trigger-package security
//     contract).
//   - **Sync staleness** ([checkSyncStaleness]) — warns when
//     a tool-native steering file is older than its source
//     `.context/steering/*.md` (the user needs to run
//     `ctx steering sync`).
//   - **Tool field** ([checkRCTool]) — `.ctxrc`'s `tool:`
//     field must be one of the supported AI tool IDs.
//
// New checks are added by appending one more `checkX` call in
// [Detect] and a constant to [config/drift.CheckName].
//
// # Issues vs Warnings vs Violations
//
// Severity is decided per-check, not per-package:
//
//   - **Violations** — things the user has to fix
//     (constitution rule break, dead path in
//     ARCHITECTURE.md). [Report.Status] returns
//     `StatusViolation` if any violation exists.
//   - **Warnings** — things the user *should* look at but can
//     defer (stale file, oversize entry count). Reported
//     individually; do not block a `ctx doctor` exit.
//   - **Passed** — names of checks that ran clean. Used by
//     `ctx doctor --json` to render a positive checklist.
//
// # Stateless and Concurrency-Safe
//
// The package holds no global state. Callers may invoke
// [Detect] concurrently as long as the [entity.Context] they
// pass is not mutated mid-call. Filesystem reads are scoped to
// the resolved context directory.
package drift
