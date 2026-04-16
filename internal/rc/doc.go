//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package rc loads, caches, and exposes the runtime
// configuration every other ctx package depends on.
// It is the single source of truth for context
// directory location, token budget, encryption
// settings, and the dozens of other knobs that shape
// ctx behavior.
//
// # Configuration Sources (Resolution Order)
//
//  1. CLI overrides -- set via ctx --context-dir
//     (highest priority, stored in rcOverrideDir).
//  2. Environment variables -- CTX_DIR,
//     CTX_TOKEN_BUDGET override .ctxrc fields.
//  3. .ctxrc (YAML) -- read once at process start
//     by [load]. Parse errors are logged via
//     [internal/write/rc.ParseWarning] and defaults
//     are kept; a malformed .ctxrc never aborts ctx.
//  4. Defaults -- every field has a hardcoded default
//     in [Default] (8000 token budget, 7-day archive,
//     200k context window, etc.).
//
// The result is the singleton [CtxRC] returned by
// [RC], memoized via sync.Once so YAML is parsed at
// most once per process.
//
// # Context-Directory Resolution
//
// [ContextDir] resolves the .context/ path:
//
//  1. CLI override (rcOverrideDir) -- return absolute.
//  2. Configured absolute path -- return as-is.
//  3. Upward walk from CWD ([walkForContextDir]) --
//     find the first ancestor containing a matching
//     directory, bounded by the git root.
//  4. Fallback -- filepath.Join(cwd, name) so that
//     ctx init can create a fresh .context/.
//
// # Key Accessors
//
//   - [TokenBudget], [ContextWindow] -- budgets
//   - [AutoArchive], [ArchiveAfterDays] -- lifecycle
//   - [ScratchpadEncrypt], [KeyPath],
//     [KeyRotationDays] -- encryption
//   - [ClassifyRules], [SpecSignalWords] -- memory
//   - [HooksEnabled], [HooksDir], [HookTimeout] --
//     hook system
//   - [SteeringDir] -- steering layer
//   - [Tool], [ActiveProfile] -- tool and profile
//   - [Validate] -- strict YAML validation with
//     unknown-field warnings
//
// # Concurrency
//
// [RC] serializes initialization through rcOnce.
// Read accessors hold an RLock; the only writer is
// the test-only [Reset]. CLI override mutation goes
// through a brief Lock().
package rc
