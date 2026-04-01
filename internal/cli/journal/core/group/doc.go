//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package group aggregates journal entries for index generation.
//
// [ByMonth] groups entries by year-month for the main index.
// [GroupedIndex] builds topic or key-file aggregations sorted by
// frequency, splitting into popular and long-tail sets.
package group
