//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package message

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ActiveMemory/ctx/internal/assets/hooks/messages"
	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// FormatTemplateVars formats available template variables for a hook message.
// If no variables are defined, returns a "(none)" indicator.
//
// Parameters:
//   - info: hook message info containing template variable names
//
// Returns:
//   - string: formatted template variables line
func FormatTemplateVars(info *messages.HookMessageInfo) string {
	if len(info.TemplateVars) == 0 {
		return desc.Text(text.DescKeyMessageTemplateVarsNone)
	}
	formatted := make([]string, len(info.TemplateVars))
	for i, v := range info.TemplateVars {
		formatted[i] = token.GoTplFieldOpen + v + token.GoTplClose
	}
	return fmt.Sprintf(
		desc.Text(text.DescKeyMessageTemplateVarsLabel),
		strings.Join(formatted, token.CommaSpace),
	)
}

// OverridePath returns the user override file path for a hook/variant.
//
// Any resolver error (including [errCtx.ErrDirNotDeclared]) is
// propagated. The previous empty-string return silently produced a
// CWD-relative path when joined by callers, which was exactly the
// "silent write to wrong location" class of bug this branch aims to
// eliminate.
//
// Parameters:
//   - hook: hook name
//   - variant: template variant name
//
// Returns:
//   - string: full filesystem path to the override file
//   - error: non-nil when the context directory cannot be resolved
func OverridePath(hook, variant string) (string, error) {
	ctxDir, err := rc.ContextDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(
		ctxDir, dir.HooksMessages,
		hook, variant+file.ExtTxt,
	), nil
}

// HasOverride checks whether a user override file exists.
//
// Parameters:
//   - hook: hook name
//   - variant: template variant name
//
// Returns:
//   - bool: true if an override file exists
//   - error: non-nil when the context directory cannot be resolved
//     or when the override file cannot be stat'd for a reason other
//     than not-exist (permission, I/O)
func HasOverride(hook, variant string) (bool, error) {
	path, err := OverridePath(hook, variant)
	if err != nil {
		return false, err
	}
	if _, statErr := os.Stat(path); statErr != nil {
		if errors.Is(statErr, os.ErrNotExist) {
			return false, nil
		}
		return false, statErr
	}
	return true, nil
}
