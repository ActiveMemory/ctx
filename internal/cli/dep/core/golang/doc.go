//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package golang implements GraphBuilder for Go projects.
// Uses go list -json to parse the module dependency graph
// and produces internal-only or full adjacency lists with
// module-prefix-stripped package names.
package golang
