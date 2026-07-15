//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package gitmeta

// HeadRef pairs a short commit SHA with the current branch
// name. Branch is the literal "detached" (see
// [github.com/ActiveMemory/ctx/internal/config/gitmeta.BranchDetached])
// when HEAD points at a commit instead of a symbolic ref.
//
// Fields:
//   - SHA: Short commit SHA of HEAD
//   - Branch: Current branch name, or "detached" when HEAD is detached
type HeadRef struct {
	SHA    string
	Branch string
}
