//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package initialize hosts compile-time constants consumed by
// the ctx init command — sentinel error messages used by
// the err/initialize package, backup directory naming for
// ctx init --reset, and the canonical reset flag name.
//
// # Why a Separate Config Package
//
// Sentinel error messages live here rather than in the
// embedded YAML loaded via desc.Text because the
// var ErrContextPopulated and var ErrResetRequiresInteractive
// declarations in err/initialize are evaluated at package-load
// time, before the YAML lookup table is populated. Keeping
// the strings as plain Go const values lets the sentinels
// initialize cleanly while the formatted wrapper messages
// (Populated, ResetRequiresInteractive) still flow through
// desc.Text for localization.
package initialize
