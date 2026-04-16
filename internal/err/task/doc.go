//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package task defines the typed error constructors
// for TASKS.md file operations. These errors fire
// when reading, writing, querying, or archiving
// tasks in the context directory.
//
// # Domain
//
// Errors fall into three categories:
//
//   - **File IO**: TASKS.md does not exist, or
//     reading/writing it failed. Constructors:
//     [FileNotFound], [FileRead], [FileWrite],
//     [SnapshotWrite].
//   - **Query**: no task matches the search
//     query, multiple tasks match, or no task was
//     specified. Constructors: [NotFound],
//     [MultipleMatches], [NoneSpecified],
//     [NoMatch].
//   - **Archive**: there are no completed tasks
//     to archive. Constructor: [NoneCompleted].
//
// # Wrapping Strategy
//
// IO constructors ([FileRead], [FileWrite],
// [SnapshotWrite]) wrap their cause with
// fmt.Errorf %w. Query and validation constructors
// return plain errors. All user-facing text is
// resolved through [internal/assets/read/desc].
//
// # Concurrency
//
// Pure constructors. Concurrent callers never race.
package task
