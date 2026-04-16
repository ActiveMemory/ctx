//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package trigger defines the **typed error
// constructors** returned by [internal/trigger] — every
// validation, discovery, and execution failure the
// trigger lifecycle can produce.
//
// # Why Typed Errors
//
//   - **Stability** — error categories are part of
//     the public API.
//   - **Routing** — write-side packages map error
//     types to localized text via
//     [internal/assets/read/desc].
//   - **Wrapping** — constructors wrap the
//     underlying cause via `%w` so callers can
//     `errors.Is` against system errors when
//     needed.
//
// # Public Surface
//
// Constructors fall into four groups:
//
//   - **Validation** — [Validate], [InvalidType],
//     [Symlink] (boundary check),
//     [ResolveHooksDir], [ResolvePath], [Boundary],
//     [Stat], [StatPath], [NotFound],
//     [ScriptExists].
//   - **Discovery / Lifecycle** — [DiscoverFailed],
//     [Chmod], [CreateDir], [Unknown],
//     [UnknownVariant].
//   - **Override Management** — [OverrideExists],
//     [WriteOverride], [RemoveOverride],
//     [EmbeddedTemplateNotFound], [WriteScript].
//   - **Execution** — [Exit] (non-zero hook
//     exit), [Timeout] (hook ran past the
//     configured timeout), [InvalidJSONOutput]
//     (hook stdout failed to parse),
//     [MarshalInput] (input encoding failed).
//
// # Why So Many Constructors
//
// Triggers run **untrusted code** at a security-
// sensitive boundary. Every distinct failure mode
// gets its own typed error so the user-facing
// message is precise about *which* invariant the
// script violated and *what* to do about it.
//
// # Concurrency
//
// Pure constructors. Concurrent callers never race.
package trigger
