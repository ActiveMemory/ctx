//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package paths resolves the ctx dream working locations from the
// cwd-anchored project root: the gitignored dreams/ notebook and the
// ideas/ source directory. The project root is the parent of the
// .context directory, by contract. Resolution is pure path
// computation — it does not create directories or read their
// contents, leaving those side effects to the caller.
package paths
