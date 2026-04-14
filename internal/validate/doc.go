//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package validate provides the **input-sanitization**
// helpers ctx uses at every boundary where a user-supplied
// string becomes part of a filesystem path, a filename, or
// a URL: stripping shell metacharacters, enforcing
// boundary-safe components, and rejecting overlong inputs
// that would blow stack frames or filesystem limits.
//
// The package is the safety net under every CLI flag that
// takes free-form text. Direct callers include the
// scratchpad tag normalizer, the journal slug normalizer,
// the pad-import path validator, and the trigger
// path-validator's underlying primitive.
//
// # Public Surface
//
//   - **[FilenameSafe](s)** — turns an arbitrary
//     string into a filename-safe variant: replaces
//     `/`, `\`, NUL, and control characters with
//     `-`, collapses runs, trims leading/trailing
//     separators. Idempotent.
//   - **[PathComponent](s)** — like [FilenameSafe]
//     but additionally rejects `.` and `..` entirely
//     so the result cannot escape the parent
//     directory.
//   - **[Bounded](s, max)** — truncates `s` to `max`
//     **runes** (rune-aware, not byte-aware) so a
//     multi-byte character is never split. Returns
//     the original when already short enough.
//   - **[NoControl](s)** — verifies `s` contains no
//     control characters; used by exec helpers
//     before passing arguments to a child process.
//
// # Why a Dedicated Package
//
// Every CLI tool eventually grows a "what counts as
// a safe filename" function in seventeen files with
// seventeen subtle differences. Hoisting them here
// means every fix lands in one place and the audit
// suite catches duplication.
//
// # Concurrency
//
// All functions are pure. Concurrent callers never
// race.
//
// # Related Packages
//
//   - [internal/cli/pad]     — uses [FilenameSafe]
//     for tag → filename conversion.
//   - [internal/cli/journal/core/slug] — sister
//     primitive that targets URL slugs (lowercase
//     hyphenation, not just safety).
//   - [internal/trigger]     — uses [PathComponent]
//     in the boundary check.
//   - [internal/exec]        — uses [NoControl] before
//     `os/exec` invocations.
package validate
