//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package err

import (
	"errors"
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets"
)

// FindSessions wraps a session-scanning failure.
//
// Parameters:
//   - cause: the underlying error from the parser.
//
// Returns:
//   - error: "failed to find sessions: <cause>"
func FindSessions(cause error) error {
	return fmt.Errorf(
		assets.TextDesc(assets.TextDescKeyErrSessionFindSessions), cause,
	)
}

// SessionNotFound returns an error for an unresolved session query.
//
// Parameters:
//   - query: the session ID or slug that was not found.
//
// Returns:
//   - error: "session not found: <query>"
func SessionNotFound(query string) error {
	return fmt.Errorf(
		assets.TextDesc(assets.TextDescKeyErrSessionNotFound), query,
	)
}

// AmbiguousQuery returns an error when a session query matches
// multiple results.
//
// Returns:
//   - error: "ambiguous query, use a more specific ID"
func AmbiguousQuery() error {
	return errors.New(assets.TextDesc(assets.TextDescKeyErrSessionAmbiguousQuery))
}

// NoSessionsFound returns an error when no sessions exist.
//
// Parameters:
//   - hint: additional guidance (e.g. "use --all-projects to search all").
//     Empty string omits the hint.
//
// Returns:
//   - error: "no sessions found" with optional hint
func NoSessionsFound(hint string) error {
	if hint != "" {
		return fmt.Errorf(
			assets.TextDesc(assets.TextDescKeyErrSessionNoSessionsFoundHint), hint,
		)
	}
	return errors.New(
		assets.TextDesc(assets.TextDescKeyErrSessionNoSessionsFound),
	)
}

// SessionIDRequired returns an error when no session ID was provided.
//
// Returns:
//   - error: "please provide a session ID or use --latest"
func SessionIDRequired() error {
	return errors.New(assets.TextDesc(assets.TextDescKeyErrSessionIDRequired))
}

// AllWithSessionID returns a validation error when --all is used with a session ID.
//
// Returns:
//   - error: "cannot use --all with a session ID; use one or the other"
func AllWithSessionID() error {
	return errors.New(
		assets.TextDesc(assets.TextDescKeyErrSessionAllWithSessionID),
	)
}

// AllWithPattern returns a validation error when --all is used with a pattern.
//
// Returns:
//   - error: "cannot use --all with a pattern; use one or the other"
func AllWithPattern() error {
	return errors.New(
		assets.TextDesc(assets.TextDescKeyErrSessionAllWithPattern),
	)
}
