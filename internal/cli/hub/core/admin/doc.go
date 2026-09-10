//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package admin resolves the hub admin credential for the
// operator-facing ctx hub subcommands.
//
// # Overview
//
// Revoke, peer and stepdown are admin-token-gated on the hub:
// they are operator actions, not client ones. Each command
// accepts the token the same way, and [Token] is that one
// way, so the three cannot drift apart.
//
// # Behavior
//
// [Token] takes the --token flag value and falls back to the
// CTX_HUB_ADMIN_TOKEN environment variable. An empty result
// is an error naming both, rather than an RPC that fails with
// PermissionDenied at the far end.
//
// # Usage
//
//	token, tokenErr := admin.Token(flagValue)
//	if tokenErr != nil {
//		return tokenErr
//	}
package admin
