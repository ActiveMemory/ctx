//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hub

import (
	"errors"
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
		desc.Text(text.DescKeyErrHubGenerateFailed), cause,
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

// UnknownClient returns an error when no registered client
// matches the given ID (e.g. an operator tries to revoke a
// client that was never registered or was already revoked).
//
// Parameters:
//   - id: the client ID that could not be found
//
// Returns:
//   - error: "unknown client: <id>"
func UnknownClient(id string) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrHubUnknownClient), id,
	)
}

// AdminTokenRequired returns an error when an admin-gated
// command is run without an admin token supplied via either the
// --token flag or the CTX_HUB_ADMIN_TOKEN environment variable.
//
// Returns:
//   - error: guidance on how to supply the admin token
func AdminTokenRequired() error {
	return errors.New(
		desc.Text(text.DescKeyErrHubAdminRequired),
	)
}

// RaftBindRequired returns an error when cluster mode is
// requested without a Raft bind address. Raft refuses to
// advertise a wildcard address, so there is no address to
// derive: the operator has to name the one their peers dial.
//
// Returns:
//   - error: guidance on supplying --raft-bind
func RaftBindRequired() error {
	return errors.New(
		desc.Text(text.DescKeyErrHubRaftBindRequired),
	)
}

// JoinWithPeers returns an error when a start asks both to
// join an existing cluster and to bootstrap one. A joining
// node has no configuration of its own by design.
//
// Returns:
//   - error: guidance on the join flow
func JoinWithPeers() error {
	return errors.New(
		desc.Text(text.DescKeyErrHubJoinWithPeers),
	)
}

// RaftBindUnroutable returns an error for a Raft bind address
// no peer could dial: a bare port, an empty host, or a
// wildcard such as 0.0.0.0 or [::].
//
// Parameters:
//   - addr: the rejected address
//
// Returns:
//   - error: "--raft-bind <addr> is not an address a peer can
//     dial: ..."
func RaftBindUnroutable(addr string) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrHubRaftBindUnroutable), addr,
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
