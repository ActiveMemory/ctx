//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package err

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets"
)

// GitNotFound returns an error when git is not installed.
// The message is loaded from assets and includes guidance for the user.
//
// Returns:
//   - error: message from assets key parser.git-not-found
func GitNotFound() error {
	return fmt.Errorf("%s", assets.TextDesc(assets.TextDescKeyParserGitNotFound))
}

// NotInGitRepo wraps a failure from git rev-parse.
//
// Parameters:
//   - cause: the underlying exec error
//
// Returns:
//   - error: "not in a git repository: <cause>"
func NotInGitRepo(cause error) error {
	return fmt.Errorf(assets.TextDesc(assets.TextDescKeyErrGitNotInGitRepo), cause)
}
