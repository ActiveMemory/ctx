//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package diff computes set differences between permission
// slices for golden-vs-local comparison. Returns restored
// entries (in golden but missing locally) and dropped entries
// (in local but absent from golden) while preserving source
// ordering.
package diff
