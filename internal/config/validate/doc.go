//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package validate holds constants used by the
// internal/cli/add/core/validate layer that enforces
// noun-specific body-flag contracts on add subcommands.
//
// The closed set of placeholder values rejected from
// required body flags lives here so call sites stay free
// of magic strings and so the policy can be reviewed in
// one place. Comparison against this set is exact and
// case-insensitive after whitespace trimming; substring
// matches are explicitly allowed.
package validate
