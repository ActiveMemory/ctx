//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package resolve

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/warn"
	errCtx "github.com/ActiveMemory/ctx/internal/err/context"
	logWarn "github.com/ActiveMemory/ctx/internal/log/warn"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// JournalDir returns the path to the journal directory within the
// configured context directory.
//
// Returns:
//   - string: Absolute path to the journal directory
//   - error: non-nil when the context directory is not declared
func JournalDir() (string, error) {
	ctxDir, err := rc.ContextDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(ctxDir, dir.Journal), nil
}

// DirLine returns a one-line context directory identifier.
//
// Emits a warn log on any non-ErrDirNotDeclared resolver error. This
// loudness is intentional: the primary caller is an AI agent whose
// incorrect invocations must be visible to the human reading the
// TUI. Do not silence this; do not move the log to a caller that
// might filter it. The error is also returned so non-rendering
// callers can propagate rather than rely solely on the log channel.
//
// Returns:
//   - string: "Context: <path>" line on success; "" on any error
//   - error: propagated from [rc.ContextDir] unchanged
func DirLine() (string, error) {
	d, err := rc.ContextDir()
	if err != nil {
		if !errors.Is(err, errCtx.ErrDirNotDeclared) {
			logWarn.Warn(warn.ContextDirResolve, err)
		}
		return "", err
	}
	return fmt.Sprintf(desc.Text(text.DescKeyWriteContextDirLabel), d), nil
}

// AppendDir appends a bracketed context directory footer to msg.
//
// Emits a warn log on any non-ErrDirNotDeclared resolver error. This
// loudness is intentional: the primary caller is an AI agent whose
// incorrect invocations must be visible to the human reading the
// TUI. Do not silence this; do not move the log to a caller that
// might filter it. The error is also returned so callers can
// propagate instead of rendering an un-annotated message when the
// context directory is unexpectedly unavailable.
//
// Parameters:
//   - msg: Base message to append the directory footer to
//
// Returns:
//   - string: Message with appended "[Context: <path>]" on success;
//     msg unchanged on any error
//   - error: propagated from [rc.ContextDir] unchanged
func AppendDir(msg string) (string, error) {
	d, err := rc.ContextDir()
	if err != nil {
		if !errors.Is(err, errCtx.ErrDirNotDeclared) {
			logWarn.Warn(warn.ContextDirResolve, err)
		}
		return msg, err
	}
	return msg + fmt.Sprintf(
		desc.Text(text.DescKeyWriteContextDirBracket), d,
	), nil
}
