//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package dream centralizes external-process execution for the
// ctx-dream pass. It resolves the configured executor binary on PATH
// and builds the bounded command that runs one headless triage pass.
//
// All exec.Command and exec.LookPath calls for the dream live here so
// the nolint:gosec annotation and argument sanitization stay in one
// place, per the internal/exec convention.
package dream
