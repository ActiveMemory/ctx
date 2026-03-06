//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package parse provides shared text-to-typed-value conversion functions.
//
// Functions here convert string inputs (dates, durations, identifiers)
// into Go types. They are thin wrappers that handle empty inputs and
// use canonical format constants from the config package.
package parse
