//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package add defines the typed error constructors
// returned by the `ctx add` command family. Every
// failure that can happen when appending entries to
// context files (DECISIONS.md, LEARNINGS.md,
// TASKS.md, CONVENTIONS.md) flows through one of
// these constructors.
//
// # Domain
//
// The add subsystem validates user input, resolves
// context file paths, and writes structured entries.
// Errors fall into three categories:
//
//   - **Missing content** -- the user invoked add
//     without providing a value via flag, argument,
//     or stdin. Constructors: [NoContent],
//     [NoContentProvided].
//   - **Validation** -- the entry type is unknown,
//     a required field is absent, or the --section
//     flag was omitted for tasks. Constructors:
//     [UnknownType], [MissingFields],
//     [SectionRequired].
//   - **File operations** -- the target context
//     file does not exist or its index could not
//     be updated. Constructors: [FileNotFound],
//     [IndexUpdate].
//
// # Wrapping Strategy
//
// Constructors that accept a cause parameter wrap
// it with fmt.Errorf %w so callers can use
// errors.Is / errors.As against the underlying
// system error. Constructors for pure validation
// failures return plain errors.New values.
//
// All user-facing text is resolved through
// [internal/assets/read/desc] at construction
// time, keeping error wording centralized.
//
// # Concurrency
//
// Pure constructors. Concurrent callers never race.
package add
