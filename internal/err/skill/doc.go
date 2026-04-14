//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package skill defines the **typed error constructors**
// returned by [internal/skill] (the install / list /
// load / remove engine) and its CLI consumers.
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
// Constructors fall into three groups:
//
//   - **Install / Remove** — [CreateDest],
//     [Install], [NotFound], [Remove], [List],
//     [ReadDir], [NotValidDir], [NotValidSource].
//   - **Load / Read** — [Load], [SkillLoad],
//     [Read], [InvalidYAML].
//   - **Manifest validation** — [InvalidManifest],
//     [MissingName], [MissingClosingDelimiter],
//     [MissingOpeningDelimiter].
//
// # Concurrency
//
// Pure constructors. Concurrent callers never race.
//
// # Related Packages
//
//   - [internal/skill]            — chief producer.
//   - [internal/cli/skill]        — also producer.
//   - [internal/write/skill]      — the renderer
//     that maps these to user text.
package skill
