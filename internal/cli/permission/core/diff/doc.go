//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package diff computes set differences between
// permission slices for golden-vs-local comparison.
//
// The "ctx permission sanitize" command compares a
// golden (reference) settings file against the local
// settings file to detect permission entries that have
// been added or removed. This package provides the
// set-difference logic for that comparison.
//
// # Set Difference
//
// [StringSlices] is the sole exported function. It
// accepts two string slices -- golden and local -- and
// returns two result slices:
//
//   - restored: entries present in golden but absent
//     from local. These are permissions that the
//     golden file grants but the local file lacks.
//   - dropped: entries present in local but absent
//     from golden. These are permissions that the
//     local file grants beyond what the golden file
//     allows.
//
// Both output slices preserve the source ordering of
// their respective inputs. Internally, the function
// builds a set for each input and performs membership
// checks to compute the symmetric difference.
//
// # Data Flow
//
// The cmd/permission layer reads the golden and local
// settings files, extracts permission string slices,
// and calls StringSlices. The restored and dropped
// results are passed to the write/permission package
// for reporting to the user.
package diff
