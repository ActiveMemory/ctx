//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package text

// DescKeys for hub error messages.
const (
	// DescKeyErrHubGenerateToken is the text key for token
	// generation failures.
	DescKeyErrHubGenerateToken = "err.hub.generate-token"
	// DescKeyErrHubInternal is the text key for internal hub
	// errors.
	DescKeyErrHubInternal = "err.hub.internal"
	// DescKeyErrHubDuplicateProject is the text key for
	// duplicate project registration errors.
	DescKeyErrHubDuplicateProject = "err.hub.duplicate-project"
	// DescKeyErrHubUnknownClient is the text key for a
	// revocation targeting a client ID that is not registered.
	DescKeyErrHubUnknownClient = "err.hub.unknown-client"
	// DescKeyErrHubAdminTokenRequired is the text key for an
	// admin-gated command invoked with no admin token supplied.
	DescKeyErrHubAdminTokenRequired = "err.hub.admin-token-required"
	// DescKeyErrHubInvalidPeerAction is the text key for
	// unrecognized peer action errors.
	DescKeyErrHubInvalidPeerAction = "err.hub.invalid-peer-action"
)
