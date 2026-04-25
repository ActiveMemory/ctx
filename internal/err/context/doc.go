//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package context defines the typed error constructors
// for .context/ directory validation. These errors
// fire during bootstrap when the CLI verifies that the
// context directory exists, is not a symlink, and
// resides within the project root.
//
// # Domain
//
// Errors fall into two categories:
//
//   - **Not found**: the .context/ directory does
//     not exist. The [NotFoundError] struct
//     implements the error interface and supports
//     errors.As matching. Constructor: [NotFound].
//   - **Security validation**: the directory or a
//     file inside it is a symlink. Constructors:
//     [DirSymlink], [FileSymlink].
//
// # Typed Error: NotFoundError
//
// [NotFoundError] is the only typed struct in the
// err/ tree that carries a Dir field. This lets
// callers distinguish "directory missing" from
// other errors with errors.As and inspect which
// path was checked.
//
// # Wrapping Strategy
//
// Security validators return plain errors (no cause
// wrapping) because the failure is a policy
// violation, not an IO failure. All user-facing
// text is resolved through [internal/assets/read/desc].
//
// # Concurrency
//
// Pure constructors. Concurrent callers never race.
package context
