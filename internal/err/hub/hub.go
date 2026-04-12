//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package hub provides error constructors for the hub subsystem.
package hub

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// GenerateToken wraps a token generation failure.
//
// Parameters:
//   - cause: the underlying error from crypto/rand
//
// Returns:
//   - error: "generate token: <cause>"
func GenerateToken(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrHubGenerateToken), cause,
	)
}

// InternalErr wraps an internal server error.
//
// Parameters:
//   - cause: the underlying error
//
// Returns:
//   - error: "internal: <cause>"
func InternalErr(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrHubInternal), cause,
	)
}

// DuplicateProject returns an error when a project is
// already registered.
//
// Parameters:
//   - name: the duplicate project name
//
// Returns:
//   - error: "project already registered: <name>"
func DuplicateProject(name string) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrHubDuplicateProject), name,
	)
}

// InvalidPeerAction returns an error for an unrecognized
// peer action.
//
// Parameters:
//   - action: the unrecognized action string
//
// Returns:
//   - error: formatted error with the invalid action
func InvalidPeerAction(action string) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrHubInvalidPeerAction),
		action,
	)
}
