//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package check_resource implements the ctx system check-resource
// subcommand.
//
// It collects system resource metrics (memory, swap, disk, load) and
// emits a warning when any resource hits danger severity.
package check_resource
