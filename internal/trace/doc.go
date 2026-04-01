//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package trace provides commit context tracing — linking git commits
// back to the decisions, tasks, learnings, and sessions that motivated them.
//
// Key exports: [Collect], [FormatTrailer], [Record], [Resolve], [ShortHash],
// [ReadHistory], [WriteHistory], [ReadOverrides], [WriteOverride],
// [CollectRefsForCommit], [ResolveCommitHash], [CommitMessage], [CommitDate].
// See source files for implementation details.
// Part of the internal subsystem.
package trace
