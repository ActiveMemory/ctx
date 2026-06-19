//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package root defines the ctx ai command tree.
//
// It contains only cobra command construction and Run wrappers. Backend
// resolution and proposal writing live in core subpackages.
// Subcommands are constructed inline to keep cmd/ free of helper APIs.
package root
