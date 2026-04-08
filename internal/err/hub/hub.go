//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package hub provides error constructors for the hub subsystem.
package hub

import "fmt"

// GenerateToken wraps a token generation failure.
//
// Parameters:
//   - cause: the underlying error from crypto/rand
//
// Returns:
//   - error: "generate token: <cause>"
func GenerateToken(cause error) error {
	return fmt.Errorf("generate token: %w", cause)
}

// InternalErr wraps an internal server error.
//
// Parameters:
//   - cause: the underlying error
//
// Returns:
//   - error: "internal: <cause>"
func InternalErr(cause error) error {
	return fmt.Errorf("internal: %w", cause)
}

// InvalidPeerAction returns an error for an unrecognized
// peer action.
//
// Parameters:
//   - action: the unrecognized action string
//
// Returns:
//   - error: "action must be 'add' or 'remove', got <action>"
func InvalidPeerAction(action string) error {
	return fmt.Errorf(
		"action must be 'add' or 'remove', got %q",
		action,
	)
}
