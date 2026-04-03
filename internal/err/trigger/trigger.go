//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package trigger

import (
	"fmt"
	"time"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// Chmod wraps a failure to change hook file permissions.
//
// Parameters:
//   - cause: the underlying OS error
//
// Returns:
//   - error: "chmod hook: <cause>"
func Chmod(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrHookChmod), cause,
	)
}

// CreateDir wraps a hook directory creation failure.
//
// Parameters:
//   - cause: the underlying OS error
//
// Returns:
//   - error: "create hook directory: <cause>"
func CreateDir(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrHookCreateDir), cause,
	)
}

// DiscoverFailed wraps a hook discovery failure.
//
// Parameters:
//   - cause: the underlying error
//
// Returns:
//   - error: "discover hooks: <cause>"
func DiscoverFailed(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrHookDiscover), cause,
	)
}

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

// Exit wraps a hook script non-zero exit.
//
// Parameters:
//   - cause: the underlying exec error
//
// Returns:
//   - error: "exit: <cause>"
func Exit(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrHookExit), cause,
	)
}

// InvalidJSONOutput wraps an invalid JSON output from a hook.
//
// Parameters:
//   - cause: the underlying JSON parse error
//
// Returns:
//   - error: "invalid JSON output: <cause>"
func InvalidJSONOutput(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrHookInvalidJSONOutput), cause,
	)
}

// InvalidType returns an error for an unrecognized hook type.
//
// Parameters:
//   - hookType: the invalid hook type string
//   - valid: comma-separated list of valid types
//
// Returns:
//   - error: "invalid hook type <hookType>; valid types: <valid>"
func InvalidType(hookType, valid string) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrHookInvalidType), hookType, valid,
	)
}

// MarshalInput wraps a hook input marshal failure.
//
// Parameters:
//   - cause: the underlying marshal error
//
// Returns:
//   - error: "marshal hook input: <cause>"
func MarshalInput(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrHookMarshalInput), cause,
	)
}

// NotFound returns an error when a hook cannot be found by name.
//
// Parameters:
//   - name: the hook name that was not found
//
// Returns:
//   - error: "hook not found: <name>"
func NotFound(name string) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrHookNotFound), name,
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

// ResolveHooksDir wraps a failure to resolve the hooks directory path.
//
// Parameters:
//   - dir: the hooks directory path
//   - cause: the underlying error
//
// Returns:
//   - error: "resolve hooks directory <dir>: <cause>"
func ResolveHooksDir(dir string, cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrHookResolveHooksDir), dir, cause,
	)
}

// ResolvePath wraps a failure to resolve a hook script path.
//
// Parameters:
//   - path: the hook path
//   - cause: the underlying error
//
// Returns:
//   - error: "resolve hook path <path>: <cause>"
func ResolvePath(path string, cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrHookResolvePath), path, cause,
	)
}

// ScriptExists returns an error when a hook script already exists.
//
// Parameters:
//   - path: the existing script path
//
// Returns:
//   - error: "hook script already exists: <path>"
func ScriptExists(path string) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrHookScriptExists), path,
	)
}

// Stat wraps a hook stat failure.
//
// Parameters:
//   - cause: the underlying OS error
//
// Returns:
//   - error: "stat hook: <cause>"
func Stat(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrHookStat), cause,
	)
}

// StatPath wraps a hook path stat failure.
//
// Parameters:
//   - path: the hook path
//   - cause: the underlying OS error
//
// Returns:
//   - error: "stat hook path <path>: <cause>"
func StatPath(path string, cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrHookStatPath), path, cause,
	)
}

// Timeout returns an error when a hook exceeds its execution timeout.
//
// Parameters:
//   - d: the timeout duration
//
// Returns:
//   - error: "timeout after <duration>"
func Timeout(d time.Duration) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrHookTimeout), d,
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

// Validate returns an error for an unknown hook/variant combination.
// It distinguishes between an entirely unknown hook and an unknown
// variant within a known hook.
//
// Parameters:
//   - hookExists: whether the hook name is recognized
//   - hook: the hook name
//   - variant: the variant name
//
// Returns:
//   - error: descriptive error with guidance to list available options
func Validate(hookExists bool, hook, variant string) error {
	if !hookExists {
		return Unknown(hook)
	}
	return UnknownVariant(variant, hook)
}

// WriteScript wraps a hook script write failure.
//
// Parameters:
//   - cause: the underlying OS error
//
// Returns:
//   - error: "write hook script: <cause>"
func WriteScript(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrHookWriteScript), cause,
	)
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

// Symlink returns an error when a hook path is a symlink.
//
// Parameters:
//   - hookPath: the symlink path
//
// Returns:
//   - error: "hook is a symlink: <hookPath>"
func Symlink(hookPath string) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrLifecycleHookSymlink),
		hookPath,
	)
}

// Boundary returns an error when a hook path escapes the
// hooks directory boundary.
//
// Parameters:
//   - hookPath: the escaping path
//   - hooksDir: the hooks root directory
//
// Returns:
//   - error: "hook escapes boundary: <hookPath> not in <hooksDir>"
func Boundary(hookPath, hooksDir string) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrLifecycleHookBoundary),
		hookPath, hooksDir,
	)
}

// NotExecutable returns an error when a hook script lacks
// the executable permission bit.
//
// Parameters:
//   - hookPath: the non-executable path
//
// Returns:
//   - error: "hook not executable: <hookPath>"
func NotExecutable(hookPath string) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrLifecycleHookNotExecutable),
		hookPath,
	)
}
