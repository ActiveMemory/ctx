//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package check_map_staleness implements the ctx system check-map-staleness
// subcommand.
//
// It detects when context map files have not been updated within the
// configured staleness window and nudges regeneration.
package check_map_staleness
