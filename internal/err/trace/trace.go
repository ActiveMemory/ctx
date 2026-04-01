//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package trace

import (
	"errors"
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// GitDir wraps a git rev-parse --git-dir failure.
//
// Parameters:
//   - cause: the underlying error
//
// Returns:
//   - error: "git rev-parse --git-dir: <cause>"
func GitDir(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrTraceGitDir), cause,
	)
}

// GitLog wraps a git log failure.
//
// Parameters:
//   - cause: the underlying error
//
// Returns:
//   - error: "git log: <cause>"
func GitLog(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrTraceGitLog), cause,
	)
}

// HookExists returns an error when a non-ctx hook already exists.
//
// Parameters:
//   - name: hook name (e.g. "prepare-commit-msg")
//   - path: path to the existing hook file
//
// Returns:
//   - error: "<name> hook already exists at <path> (not installed by ctx)..."
func HookExists(name, path string) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrTraceHookExists), name, path,
	)
}

// HookWrite wraps a hook file write failure.
//
// Parameters:
//   - name: hook name (e.g. "prepare-commit-msg")
//   - cause: the underlying error
//
// Returns:
//   - error: "write <name> hook: <cause>"
func HookWrite(name string, cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrTraceHookWrite), name, cause,
	)
}

// NoteRequired returns an error when --note is missing.
//
// Returns:
//   - error: "--note is required"
func NoteRequired() error {
	return errors.New(desc.Text(text.DescKeyErrTraceNoteRequired))
}

// ResolveCommit wraps a commit resolution failure.
//
// Parameters:
//   - ref: the commit ref that failed to resolve
//   - cause: the underlying error
//
// Returns:
//   - error: "resolve commit \"<ref>\": <cause>"
func ResolveCommit(ref string, cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrTraceResolveCommit), ref, cause,
	)
}

// UnknownAction returns an error for an unrecognized hook action.
//
// Parameters:
//   - action: the unknown action string
//
// Returns:
//   - error: "unknown action \"<action>\": use enable or disable"
func UnknownAction(action string) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrTraceUnknownAction), action,
	)
}

// WriteHistory wraps a history write failure.
//
// Parameters:
//   - cause: the underlying error
//
// Returns:
//   - error: "write history: <cause>"
func WriteHistory(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrTraceWriteHistory), cause,
	)
}

// WriteOverride wraps an override write failure.
//
// Parameters:
//   - cause: the underlying error
//
// Returns:
//   - error: "write override: <cause>"
func WriteOverride(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrTraceWriteOverride), cause,
	)
}
