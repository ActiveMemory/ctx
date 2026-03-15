//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package err

import "fmt"

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
	return fmt.Errorf("embedded template not found for %s/%s", hook, variant)
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
	return fmt.Errorf("override already exists at %s\nEdit it directly or use `ctx system message reset %s %s` first",
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
	return fmt.Errorf("failed to write override %s: %w", path, cause)
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
	return fmt.Errorf("failed to remove override %s: %w", path, cause)
}

// UnknownHook returns an error for an unrecognized hook name.
//
// Parameters:
//   - hook: the unknown hook name
//
// Returns:
//   - error: "unknown hook: <hook>..."
func UnknownHook(hook string) error {
	return fmt.Errorf("unknown hook: %s\nRun `ctx system message list` to see available hooks", hook)
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
	return fmt.Errorf("unknown variant %q for hook %q\nRun `ctx system message list` to see available variants", variant, hook)
}
