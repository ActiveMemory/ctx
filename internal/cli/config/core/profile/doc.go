//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package profile manages **`.ctxrc` profile detection,
// copying, and switching**, the engine behind
// `ctx config switch` that lets a user toggle between
// (typically) `dev` and `base` configurations without
// hand-editing `.ctxrc`.
//
// The package is the *mutator*; the read-side cache
// lives in [internal/rc].
//
// # The Profile Convention
//
// Profiles are stored as **per-profile files** in the
// project root:
//
//   - `.ctxrc`: the active configuration.
//   - `.ctxrc.dev`: the dev profile (verbose
//     logs, webhook events, ...).
//   - `.ctxrc.base`: the base / production
//     profile (clean defaults).
//
// `prod` is recognized as an alias for `base`. New
// profiles plug in as `.ctxrc.<name>`.
//
// # Public Surface
//
//   - **[Active]**: returns the name of the
//     currently-active profile (read from `.ctxrc`'s
//     `profile:` field).
//   - **[Detect](root)**: lists every available
//     profile (by glob).
//   - **[Switch](root, name)**: copies
//     `.ctxrc.<name>` over `.ctxrc`. Atomic via the
//     standard write-temp-rename pattern. Refuses
//     to switch to an unknown profile.
//   - **[GitRoot]**: resolves the project's git
//     root for path operations (the profile files
//     live there, not in the current working
//     subdirectory).
//
// # Concurrency
//
// Filesystem-bound and stateless. `ctx` is
// single-process; concurrent switches are not a
// design concern.
package profile
