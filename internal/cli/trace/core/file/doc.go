//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package file implements the core logic for file-level context tracing.
//
// [ParsePathArg] strips optional :line-range suffixes from path arguments.
// [Trace] runs git log for a file and prints context refs attached to
// each commit that touched it. Results combine history entries and
// override annotations.
//
// Key exports: [ParsePathArg], [Trace].
// Called by the trace file CLI subcommand.
package file
