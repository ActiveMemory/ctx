//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package drift detects **version drift** across the three
// places ctx's version can diverge: the source-of-truth
// `VERSION` file, the installed binary's
// `ctx --version`, and the marketplace plugin manifest. The
// `check_version` hook calls into here to nudge users when
// any of the three drift apart.
//
// (This is the *system-hook* drift package and is unrelated
// to [internal/drift], which detects context-file drift.)
//
// # Public Surface
//
//   - **[CheckVersion]**: runs the full three-way
//     comparison and returns a [DriftReport].
//   - **[ReadVersionFile]**: reads the `VERSION` file
//     from the install dir.
//   - **[ReadMarketplaceVersion]**: reads the
//     plugin manifest's pinned version from
//     `~/.claude/marketplaces/...`.
//   - **[FormatStaleEntries]**: formats a [DriftReport]
//     as the user-facing nudge body (delivered via
//     [internal/cli/system/core/message]).
//
// # The Three Sources, Why They Drift
//
//  1. **`VERSION` file**: bumped by maintainers as part
//     of the release runbook. The source of truth.
//  2. **Installed binary**: the result of the user's
//     last `make install` / `brew upgrade`. Drifts
//     downward if the user has not updated.
//  3. **Marketplace plugin manifest**: pinned by the
//     user's most recent `claude plugin install`.
//     Drifts downward if the user has not run
//     `claude plugin update`.
//
// Any pair-wise mismatch is a candidate nudge; the hook
// picks the most actionable phrasing per case.
//
// # Concurrency
//
// All functions are filesystem-bound and stateless.
// Concurrent invocations never race.
package drift
