//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package rc loads, caches, and exposes the runtime configuration
// every other ctx package depends on. It is the single source of
// truth for context directory location, token budget, encryption
// settings, and the dozens of other knobs that shape ctx behavior.
//
// # Context-Directory Resolution (explicit-only)
//
// Under the explicit-context-dir model
// (spec: specs/explicit-context-dir.md), rc does NOT walk the
// filesystem looking for a .context/ directory. Every non-exempt
// command must declare the target explicitly.
//
// [ContextDir] returns the declared path or the empty string:
//
//  1. CLI override set via [OverrideContextDir] (--context-dir
//     flag) wins if present.
//  2. CTX_DIR environment variable is consulted next.
//  3. Otherwise the empty string is returned. Exempt callers
//     (ctx init, activate, deactivate, system bootstrap) handle
//     empty themselves; every other command should call
//     [RequireContextDir] instead, which returns a tailored error
//     whose message depends on how many .context/ candidates are
//     visible from CWD.
//
// [ScanCandidates] is a read-only upward scan used by the
// `ctx activate` subcommand and by [RequireContextDir]'s error
// formatter. It does not resolve, bind, or select a directory.
//
// # Configuration File (.ctxrc)
//
// Once [ContextDir] is declared, [load] reads `.ctxrc` from
// `filepath.Dir(ContextDir())`: the project root, which by contract
// is the parent of [ContextDir]. CWD has no say. When no context
// directory is declared, `.ctxrc` is not read at all and defaults
// apply.
//
// Environment overrides (CTX_TOKEN_BUDGET) are applied after the
// YAML merge so users can tune per-session without editing the
// file.
//
// The singleton [CtxRC] returned by [RC] is memoized via
// sync.Once so YAML is parsed at most once per process.
//
// # Concurrency
//
// [RC] serializes initialization through rcOnce. Read accessors
// hold an RLock; the only writer is the test-only [Reset]. CLI
// override mutation goes through a brief Lock().
package rc
