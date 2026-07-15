//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package core groups the ctx dream command's domain logic into
// focused subpackages: paths (root/notebook resolution), pass (one
// executor-agnostic run), and dispose (loading a proposal by id and
// applying accept/reject/amend through the engine). The command
// itself stays a thin Cobra wrapper; all triage behavior lives here
// so it can be exercised without the CLI layer.
package core
