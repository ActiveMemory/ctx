//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package message

import (
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
		formatted[i] = "{{." + v + "}}"
	}
	return fmt.Sprintf(desc.Text(text.DescKeyMessageTemplateVarsLabel), strings.Join(formatted, token.CommaSpace))
}

// OverridePath returns the user override file path for a hook/variant.
//
// Parameters:
//   - hook: hook name
//   - variant: template variant name
//
// Returns:
//   - string: full filesystem path to the override file
func OverridePath(hook, variant string) string {
	return filepath.Join(rc.ContextDir(), dir.HooksMessages, hook, variant+file.ExtTxt)
}

// HasOverride checks whether a user override file exists.
//
// Parameters:
//   - hook: hook name
//   - variant: template variant name
//
// Returns:
//   - bool: true if an override file exists
func HasOverride(hook, variant string) bool {
	_, statErr := os.Stat(OverridePath(hook, variant))
	return statErr == nil
}
