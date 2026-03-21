//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/file"
	ctxerr "github.com/ActiveMemory/ctx/internal/err/hook"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/hooks/messages"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// ValidationError returns an error for an unknown hook/variant combination.
// It distinguishes between an entirely unknown hook and an unknown variant
// within a known hook.
//
// Parameters:
//   - hook: the hook name to validate
//   - variant: the variant name to validate
//
// Returns:
//   - error: descriptive error with guidance to list available options
func ValidationError(hook, variant string) error {
	if messages.Variants(hook) == nil {
		return ctxerr.Unknown(hook)
	}
	return ctxerr.UnknownVariant(variant, hook)
}

// PrintTemplateVars prints available template variables for a hook message.
// If no variables are defined, prints a "(none)" indicator.
//
// Parameters:
//   - cmd: Cobra command for output
//   - info: hook message info containing template variable names
func PrintTemplateVars(cmd *cobra.Command, info *messages.HookMessageInfo) {
	if len(info.TemplateVars) == 0 {
		cmd.Println(desc.Text(text.DescKeyMessageTemplateVarsNone))
		return
	}
	formatted := make([]string, len(info.TemplateVars))
	for i, v := range info.TemplateVars {
		formatted[i] = "{{." + v + "}}"
	}
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyMessageTemplateVarsLabel), strings.Join(formatted, ", ")))
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
