//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package session

import (
	"errors"
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// Find wraps a session-scanning failure.
//
// Parameters:
//   - cause: the underlying error from the parser.
//
// Returns:
//   - error: "failed to find sessions: <cause>"
func Find(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrSessionFindSessions), cause,
	)
}

// NotFound returns an error for an unresolved session query.
//
// Parameters:
//   - query: the session ID or slug that was not found.
//
// Returns:
//   - error: "session not found: <query>"
func NotFound(query string) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrSessionNotFound), query,
	)
}

// NoneFound returns an error when no sessions exist.
//
// Parameters:
//   - hint: additional guidance (e.g. "use --all-projects to search all").
//     Empty string omits the hint.
//
// Returns:
//   - error: "no sessions found" with an optional hint
func NoneFound(hint string) error {
	if hint != "" {
		return fmt.Errorf(
			desc.Text(text.DescKeyErrSessionNoSessionsFoundHint), hint,
		)
	}
	return errors.New(
		desc.Text(text.DescKeyErrSessionNoSessionsFound),
	)
}

// AmbiguousQuery returns an error when a session query matches
// multiple results.
//
// Returns:
//   - error: "ambiguous query, use a more specific ID"
func AmbiguousQuery() error {
	return errors.New(desc.Text(text.DescKeyErrSessionAmbiguousQuery))
}

// IDRequired returns an error when no session ID was provided.
//
// Returns:
//   - error: "please provide a session ID or use --latest"
func IDRequired() error {
	return errors.New(desc.Text(text.DescKeyErrSessionIDRequired))
}

// AllWithID returns a validation error when --all is used with a session ID.
//
// Returns:
//   - error: "cannot use --all with a session ID; use one or the other"
func AllWithID() error {
	return errors.New(
		desc.Text(text.DescKeyErrSessionAllWithSessionID),
	)
}

// AllWithPattern returns a validation error when --all is used with a pattern.
//
// Returns:
//   - error: "cannot use --all with a pattern; use one or the other"
func AllWithPattern() error {
	return errors.New(
		desc.Text(text.DescKeyErrSessionAllWithPattern),
	)
}
