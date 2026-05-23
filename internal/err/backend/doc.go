//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package backend owns the typed-string sentinels and
// wrapping constructors for the AI-backend registry's
// resolution surface. The user-facing strings live in
// `internal/assets/commands/text/errors.yaml` keyed by
// `err.backend.*`; sentinels resolve them lazily via
// `desc.Text` at `Error()` call time, sidestepping the
// init-ordering trap of
// `var ErrX = errors.New(desc.Text(...))`.
//
// # Domain
//
// Four sentinels cover the registry failure surface:
//
//   - [ErrBackendNotFound]: Resolve called with a name
//     that has no Config or Factory. Wrapped by
//     [NotFound] with the offending name.
//   - [ErrNoBackends]: Resolve or Default called against
//     a Registry with no backends configured at all.
//   - [ErrAmbiguousDefault]: Default called when more
//     than one backend is configured and no explicit
//     default has been set.
//   - [ErrDuplicateRegistration]: Register called twice
//     for the same type name. Wrapped by
//     [DuplicateRegistration] with the offending name.
//
// # Wrapping strategy
//
// Constructors use `fmt.Errorf` with `%w` so callers can
// `errors.Is(err, ErrX)` against the sentinel and
// `errors.Unwrap` to recover the underlying detail when
// surfacing operator-facing diagnostics.
//
// # Related packages
//
//   - [github.com/ActiveMemory/ctx/internal/backend]
//     is the primary caller.
//   - [github.com/ActiveMemory/ctx/internal/config/embed/text]
//     supplies the DescKey constants resolved by these
//     sentinels.
package backend
