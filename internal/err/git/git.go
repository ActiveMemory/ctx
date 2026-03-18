//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package git

import (
	"errors"
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config/embed"
)

// NotFound returns an error when git is not installed.
// The message is loaded from assets and includes guidance for the user.
//
// Returns:
//   - error: message from the assets key parser.git-not-found
func NotFound() error {
	return errors.New(assets.TextDesc(embed.TextDescKeyParserGitNotFound))
}

// NotInRepo wraps a failure from git rev-parse.
//
// Parameters:
//   - cause: the underlying exec error
//
// Returns:
//   - error: "not in a git repository: <cause>"
func NotInRepo(cause error) error {
	return fmt.Errorf(
		assets.TextDesc(embed.TextDescKeyErrGitNotInGitRepo), cause,
	)
}
