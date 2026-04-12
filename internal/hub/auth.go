//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hub

import cfgHub "github.com/ActiveMemory/ctx/internal/config/hub"

// GenerateAdminToken creates a new admin token.
//
// The token is a hex-encoded 32-byte random value prefixed
// with "ctx_adm_". It is generated once on first hub startup
// and must be saved by the operator.
//
// Returns:
//   - string: the admin token
//   - error: non-nil if the system random source fails
func GenerateAdminToken() (string, error) {
	return generateToken(cfgHub.AdminTokenPrefix)
}

// GenerateClientToken creates a new client token.
//
// The token is a hex-encoded 32-byte random value prefixed
// with "ctx_cli_". It is issued to a client during
// registration.
//
// Returns:
//   - string: the client token
//   - error: non-nil if the system random source fails
func GenerateClientToken() (string, error) {
	return generateToken(cfgHub.ClientTokenPrefix)
}
