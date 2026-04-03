//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package regex

import (
	"regexp"

	"github.com/ActiveMemory/ctx/internal/config/project"
	"github.com/ActiveMemory/ctx/internal/config/token"
)

// InternalPkg matches backtick-quoted paths starting with
// "internal/" in documentation. Used by drift checks to
// verify that referenced paths still exist on disk.
var InternalPkg = regexp.MustCompile(
	token.Backtick +
		"(" + project.DirInternalSlash +
		"[^" + token.Backtick + "]+)" +
		token.Backtick,
)
