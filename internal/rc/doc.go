//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package rc loads, caches, and exposes the **runtime
// configuration** every other ctx package and command depends
// on. It is the single source of truth for "where does the
// context directory live", "what is the token budget", "is
// scratchpad encryption on", and the dozens of other knobs that
// shape ctx behavior.
//
// The package is foundational: most binaries and hooks call
// [RC] (or one of its accessor wrappers like [ContextDir],
// [TokenBudget], [PriorityOrder]) within the first few
// instructions, and many do so on every prompt. Performance and
// determinism here matter to the rest of the system.
//
// # Configuration Sources, in Resolution Order
//
//  1. **CLI overrides** — set via `ctx --context-dir <path>`
//     (and equivalent flags). Highest priority. Stored in
//     `rcOverrideDir` under [rcMu].
//  2. **Environment variables** — `CTX_DIR`, `CTX_TOKEN_BUDGET`
//     (see [internal/config/env]) override the corresponding
//     `.ctxrc` fields when set.
//  3. **`.ctxrc` (YAML)** — read once at process start by
//     [load] from the current working directory. Parse errors
//     are logged via [internal/write/rc.ParseWarning] and the
//     defaults are kept; **a malformed `.ctxrc` never aborts
//     ctx**.
//  4. **Defaults** — every field has a hardcoded default in
//     [Default]; see the constants in `default.go`
//     (`DefaultTokenBudget`, `DefaultArchiveAfterDays`,
//     `DefaultContextWindow`, etc.).
//
// The result is the singleton `*CtxRC` returned by [RC]. It is
// memoized via [sync.Once] so the YAML is parsed at most once
// per process (tests can call [Reset] to invalidate the cache).
//
// # Context-Directory Resolution
//
// [ContextDir] returns the absolute path of the project's
// `.context/` directory and is the most-called function in the
// package. Its resolution order is:
//
//  1. CLI override (`rcOverrideDir`) → return absolute, no
//     walk.
//  2. Configured **absolute** path (`.ctxrc` or env var) →
//     return as-is.
//  3. **Upward walk from CWD** ([walkForContextDir]) — find
//     the first ancestor directory that contains a folder
//     whose basename matches the configured name. The walk is
//     bounded by the **git root** when one is present: a
//     candidate that falls outside the git root is discarded
//     so commands run in submodules or sibling projects do not
//     leak into the wrong project's context.
//  4. **Fallback** — `filepath.Join(cwd, name)` returned
//     absolute. This preserves `ctx init`'s ability to create
//     a fresh `.context/` at the current location.
//
// The walk result is cached for the life of the process. Hook
// scripts and subcommand binaries invoked from a project
// subdirectory therefore consistently see the same context dir
// the user's terminal does.
//
// # The Configuration Schema
//
// [CtxRC] is the YAML-tagged struct mirrored by `.ctxrc`. A
// non-exhaustive tour of the field families:
//
//   - **Layout** — `ContextDir`, `KeyPathOverride`,
//     `Steering.Dir`, `Hooks.Dir`. Where things live.
//   - **Budgets** — `TokenBudget`, `InjectionTokenWarn`,
//     `BillingTokenWarn`, `ContextWindow`. How big context
//     packets and injections are allowed to be.
//   - **Lifecycle** — `AutoArchive`, `ArchiveAfterDays`,
//     `KeyRotationDays`, `StaleAgeDays`. Time-based nudges and
//     auto-cleanup.
//   - **Per-tool** — `Tool` (claude / cursor / cline / kiro /
//     codex), `Steering.DefaultTools`. Which AI assistant is
//     active and which tools steering files default to.
//   - **Provenance** — [ProvenanceConfig] toggles which of
//     `--session-id`, `--branch`, `--commit` are required when
//     adding entries; default is "all required".
//   - **Notifications** — [NotifyConfig] holds the event
//     filter for `ctx hook notify` (loop / nudge / relay /
//     heartbeat).
//   - **Memory bridge** — `ClassifyRules` overrides the
//     keyword classifier in [internal/memory] when set.
//   - **Spec nudge** — `SpecSignalWords`, `SpecNudgeMinLen`
//     control when `ctx add task` suggests writing a spec.
//   - **Freshness tracking** — [FreshnessFile] entries make
//     the staleness check warn when technology-dependent
//     constants in source files have not been reviewed in N
//     months.
//
// Pointer-typed `*bool` fields ([CtxRC.ScratchpadEncrypt],
// [HooksRC.Enabled], [ProvenanceConfig.SessionID/Branch/Commit])
// distinguish "user explicitly set false" from "unset, use the
// default true" — assignment by value would lose that
// distinction.
//
// # Concurrency
//
// [RC] is safe to call from any goroutine; it serializes
// initialization through `rcOnce`. Read accessors hold an
// `RLock` on `rcMu`; the only writer is the test-only [Reset].
// CLI override mutation goes through a brief `Lock()`.
//
// # Profiles
//
// `.ctxrc` supports an optional `profile:` field plus profile
// overlays (`.ctxrc.dev`, `.ctxrc.base`) wired by `ctx config
// switch`. The active profile name is exposed via
// [ActiveProfile]; see [internal/cli/config] for the
// switching logic and [docs/recipes/configuration-profiles.md]
// for the user-facing story.
//
// # Related Packages
//
//   - [github.com/ActiveMemory/ctx/internal/config/env] — the
//     env-var name constants this package consumes.
//   - [github.com/ActiveMemory/ctx/internal/config/file] —
//     `.ctxrc` file name constant.
//   - [github.com/ActiveMemory/ctx/internal/crypto] —
//     resolves the encryption key path; consumed by
//     [internal/pad] and [internal/notify].
//   - [github.com/ActiveMemory/ctx/internal/cli/config] — the
//     `ctx config` CLI surface (status, switch, schema).
//   - [github.com/ActiveMemory/ctx/internal/write/rc] —
//     terminal output helpers, including the parse-warning
//     formatter used when `.ctxrc` YAML is malformed.
package rc
