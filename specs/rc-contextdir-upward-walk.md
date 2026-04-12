---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: rc.ContextDir upward walk for subdirectory safety
status: accepted
date: 2026-04-11
owner: jose
scope: bug fix — single function
---

# Spec — `rc.ContextDir()` upward walk

## Problem

`internal/rc/rc.go:63-70` returns the configured context directory as a
**relative string** (default `".context"`) with no path resolution. Every
one of 95 call sites does `filepath.Join(rc.ContextDir(), ...)` and thereby
implicitly trusts that the process CWD is the project root.

When any command or hook is invoked from a project subdirectory, the join
resolves against that subdirectory and the caller creates stray state
files inside — e.g. `docs/cli/.context/state/backup-reminded`,
`site/cli/.context/state/backup-reminded`. This has happened repeatedly
and has been silently ignored.

## Root cause

`ContextDir()` has no concept of "the project root". It returns whatever
the `.ctxrc` said (or the env var, or the default literal) without
translating that into an absolute path anchored to a real project.

## Fix

Change `ContextDir()` to resolve the configured name against an **upward
walk from CWD**. At most one stat call per ancestor.

### Resolution order

1. **CLI override** (`rcOverrideDir`) — returned as absolute. No walk.
2. **Configured absolute path** (`.ctxrc` or env var set an absolute
   `context_dir`) — returned as-is. No walk.
3. **Upward walk** — from CWD, walk parent directories looking for an
   existing directory whose basename matches the configured name. First
   hit wins; returned as absolute.
4. **Fallback** — if nothing is found upward, return
   `filepath.Join(cwd, configuredName)` as absolute. This preserves
   `ctx init`'s ability to create a new context directory in the current
   directory when none exists yet.

### Caching

No cache. An earlier iteration of this spec proposed caching the walk
result behind `sync.Once`, but that broke ~20 existing tests that rely on
`os.Chdir` between subtests in the same process — the cached path from
test 1 leaked into test 2. The walk is cheap (a handful of `os.Stat`
calls bounded by directory depth from CWD to `/`), CLI invocations only
resolve `ContextDir()` a few dozen times per process, and test-hostile
caching added complexity with no real performance win. The resolver runs
on every call; callers that care about stability should capture the
result in a local variable.

### What this does NOT change

- The 95 call sites. They all keep `filepath.Join(rc.ContextDir(), ...)`.
  Changing the return value from relative-`.context` to an absolute path
  is behavior-preserving for project-root invocations and bug-fixing for
  everything else.
- The `.ctxrc` file-loading path in `load.go`. That is a separate (and
  related) bug — today `load()` reads `.ctxrc` from CWD, so a project
  subdirectory still loads defaults. Deferring because:
  - Most users don't have a custom `.ctxrc`, so they get the default
    literal `.context`, which my walk already handles.
  - Fixing `.ctxrc` loading requires walking to find the file BEFORE
    the walk for the dir, creating a chicken-and-egg layering problem
    that is out of scope for this bug fix.
  - Tracked separately as a follow-up task.
- `state.SetDirForTest` and the existing test overrides. Orthogonal.

## Test

New test in `internal/rc/rc_test.go`:

- **`TestContextDir_UpwardWalkFromSubdir`**: create a temp dir with
  `project/.context/` inside, chdir into `project/deep/nested/`, call
  `ContextDir()`, assert the returned path equals the absolute path of
  `project/.context` (not `project/deep/nested/.context`).
- **`TestContextDir_FallbackWhenNotFound`**: chdir into a temp dir with
  no `.context/` anywhere upward, call `ContextDir()`, assert it returns
  `filepath.Join(cwd, ".context")` as an absolute path (i.e., preserves
  `ctx init` semantics).

Both tests call `Reset()` to clear caches before running.

## Regression guard

After this lands, a hook or CLI invoked from any subdirectory should
resolve `.context/` to the correct project-root path. A future reference
test may additionally integration-test this by invoking
`ctx system check-backup-age` from a subdirectory and asserting no stray
`.context/` appears below the project root — deferred.

## Out of scope

- `.ctxrc` file-loading location (separate follow-up).
- `state.Dir()` or any other helper — they already call `rc.ContextDir()`
  and therefore inherit the fix for free.
- Multi-context-per-process scenarios (we intentionally cache once per
  process; tests use `Reset()` to work around this).
- Cleanup of the two stray directories currently under `docs/cli/` and
  `site/cli/` — handled outside this spec as a one-shot `rm -rf`.
