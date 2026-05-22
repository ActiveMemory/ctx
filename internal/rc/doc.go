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
// # Context-Directory Resolution (cwd-anchored)
//
// Under the cwd-anchored resolution model
// (spec: specs/cwd-anchored-context.md), rc does NOT walk the
// filesystem and does NOT consult any environment variable.
// `$PWD/.context/` is the answer, full stop. The user (or their
// AI tool) is responsible for being at the project root.
//
// [ContextDir] returns the absolute path to `$PWD/.context/` after
// a single [os.Stat], or a typed errCtx error:
//
//   - [errCtx.ErrNoCtxHere] when the directory is absent;
//   - [errCtx.ErrContextDirNotADirectory] when the path exists but
//     is a regular file;
//   - [errCtx.ErrContextDirStat] when stat fails for permission /
//     I/O reasons.
//
// [RequireContextDir] is a thin wrapper retained as the canonical
// "I need a usable directory" call shape for operating commands.
//
// # Configuration File (.ctxrc)
//
// When `$PWD/.context/` is present, [load] reads `.ctxrc` from
// `$PWD/.ctxrc`: the project root, which by contract is the parent
// of [ContextDir]. When `.context/` is absent, `.ctxrc` is not
// read at all and defaults apply.
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
