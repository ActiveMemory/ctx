//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package testctx provides helpers for exercising ctx commands in
// tests under the explicit-context-dir resolution model (spec:
// specs/explicit-context-dir.md).
//
// Under that model [rc.ContextDir] returns "" unless the caller has
// declared a context directory via --context-dir or CTX_DIR. Tests
// that chain multiple ctx commands in the same process (e.g.,
// `ctx init` followed by `ctx add`) must therefore declare CTX_DIR
// before any non-exempt command runs, and must reset rc state between
// test cases so process-global overrides do not leak.
//
// [Declare] is the one-stop helper: it points CTX_DIR at
// `<tempDir>/.context`, resets rc, and registers an end-of-test reset
// via `t.Cleanup`. Callers still need to run `ctx init` (or
// materialize .context/ themselves); Declare only wires the
// environment.
package testctx
