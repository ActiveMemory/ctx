//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hook

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// EmbeddedTemplateNotFound returns an error when an embedded hook
// message template cannot be located.
//
// Parameters:
//   - hook: hook name
//   - variant: template variant name
//
// Returns:
//   - error: "embedded template not found for <hook>/<variant>"
func EmbeddedTemplateNotFound(hook, variant string) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrHookEmbeddedTemplateNotFound),
		hook, variant,
	)
}

// OverrideExists returns an error when a message override already
// exists and must be reset before editing.
//
// Parameters:
//   - path: existing override file path
//   - hook: hook name
//   - variant: template variant name
//
// Returns:
//   - error: "override already exists at <path>..."
func OverrideExists(path, hook, variant string) error {
	return fmt.Errorf(desc.Text(text.DescKeyErrHookOverrideExists),
		path, hook, variant)
}

// WriteOverride wraps a message override write failure.
//
// Parameters:
//   - path: the override file path
//   - cause: the underlying error
//
// Returns:
//   - error: "failed to write override <path>: <cause>"
func WriteOverride(path string, cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrHookWriteOverride), path, cause,
	)
}

// RemoveOverride wraps a message override removal failure.
//
// Parameters:
//   - path: the override file path
//   - cause: the underlying error
//
// Returns:
//   - error: "failed to remove override <path>: <cause>"
func RemoveOverride(path string, cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrHookRemoveOverride), path, cause,
	)
}

// Unknown returns an error for an unrecognized hook name.
//
// Parameters:
//   - hook: the unknown hook name
//
// Returns:
//   - error: "unknown hook: <hook>..."
func Unknown(hook string) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrHookUnknownHook), hook,
	)
}

// UnknownVariant returns an error for an unrecognized variant within
// a known hook.
//
// Parameters:
//   - variant: the unknown variant name
//   - hook: the parent hook name
//
// Returns:
//   - error: "unknown variant <variant> for hook <hook>..."
func UnknownVariant(variant, hook string) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrHookUnknownVariant), variant, hook,
	)
}
