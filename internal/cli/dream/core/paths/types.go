//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package paths

// Resolved holds the dream working locations for one invocation.
type Resolved struct {
	// Root is the absolute project root (parent of .context).
	Root string
	// Dreams is the absolute path to the gitignored dreams/ notebook.
	Dreams string
	// Ideas is the absolute path to the ideas/ source directory.
	Ideas string
}
