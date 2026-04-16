//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package regex centralizes every **compiled regular
// expression** ctx uses anywhere in its codebase, so they are
// compiled exactly once at process start (`regexp.MustCompile`
// at package-init time) and so the patterns themselves live in
// a single, auditable place.
//
// Two motivating problems:
//
//   - **Cost** — `regexp.MustCompile` is non-trivial; calling
//     it inside a hot loop is wasteful. Hoisting every pattern
//     into a package-level `var` guarantees init-time
//     compilation.
//   - **Audit-ability** — patterns scattered through a
//     codebase drift silently. Co-locating them lets
//     reviewers eyeball the full surface in one place and
//     lets the test suite (`regexp_location_test`,
//     `regexp_test`) verify that no other package compiles
//     its own regex.
//
// # File Layout — One Concern per File
//
// Each file groups patterns for a single concern: hook safety
// scanners (`mid_sudo.go`, `git_push.go`,
// `cp_mv_to_bin.go`), context tracing (`task_ref.go`),
// non-PATH-ctx detection (`ctx_absolute_start.go`,
// `ctx_relative_start.go`), source-tree linters
// (`oversize_tokens.go`, `line_number.go`), and so on.
//
// # Naming Convention
//
// Each exported `var` is a `*regexp.Regexp` named for what
// it matches, not how. `MidSudo`, `GitPush`,
// `InstallToLocalBin`, `OversizeTokens`, `TaskRef`. The doc
// comment above each variable documents the pattern, the
// captured groups, and the call sites that consume it.
//
// # Concurrency
//
// `*regexp.Regexp` is safe for concurrent use after
// compilation — the standard library guarantees it. No
// caller needs to lock when invoking `Match`, `Find`,
// `Replace`, or `Submatch`.
//
// # Audit
//
// Two AST tests defend the contract:
//
//   - `regexp_location_test` — fails if any source file
//     outside this package calls `regexp.Compile` or
//     `regexp.MustCompile`.
//   - `regexp_test` — fails if a pattern in this package is
//     dead (unused) or if a referenced variable is missing.
//
// New patterns therefore must land here; new call sites
// reference the variable by name.
package regex
