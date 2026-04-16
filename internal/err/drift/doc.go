//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package drift defines the typed error constructors
// for the drift detection subsystem. Drift detection
// scans context files for stale paths, broken
// cross-references, and constitution violations.
//
// # Domain
//
// A single constructor covers the entire surface:
//
//   - [Violations]: drift detection completed and
//     found one or more violations. The CLI uses
//     this as a non-zero exit signal after printing
//     the violation report.
//
// This package is intentionally minimal. The drift
// scanner itself reports individual violations
// through the writer layer; this sentinel error
// only signals the aggregate outcome.
//
// # Wrapping Strategy
//
// [Violations] returns a plain errors.New value
// with no cause wrapping because the error
// represents a summary, not a single IO failure.
// All user-facing text is resolved through
// [internal/assets/read/desc].
//
// # Concurrency
//
// Pure constructor. Concurrent callers never race.
package drift
