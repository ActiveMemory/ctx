//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package session defines the typed error constructors
// for session lookup and selection. These errors fire
// when scanning for AI tool sessions, resolving a
// session by ID or slug, or validating flag
// combinations on session subcommands.
//
// # Domain
//
// Errors fall into three categories:
//
//   - **Scan failures**: the session scanner could
//     not enumerate available sessions.
//     Constructor: [Find].
//   - **Lookup**: no session matches the query,
//     no sessions exist at all, or the query is
//     ambiguous. Constructors: [NotFound],
//     [NoneFound], [AmbiguousQuery], [IDRequired].
//   - **Flag validation**: mutually exclusive
//     flags were combined (--all with a session ID
//     or pattern, or an invalid --type value).
//     Constructors: [AllWithID], [AllWithPattern],
//     [EventInvalidType].
//
// # Wrapping Strategy
//
// [Find] wraps its cause with fmt.Errorf %w so
// callers can inspect the underlying parser error.
// Lookup and validation constructors return plain
// errors. All user-facing text is resolved through
// [internal/assets/read/desc].
//
// # Concurrency
//
// Pure constructors. Concurrent callers never race.
package session
