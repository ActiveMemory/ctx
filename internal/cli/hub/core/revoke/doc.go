//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package revoke implements client token revocation for the
// ctx hub revoke command.
//
// # Overview
//
// This package revokes a registered client's token on a remote
// hub, invalidating it immediately. It is the operator-side
// counterpart to registration: where registration mints a
// token, revocation destroys one.
//
// # Behavior
//
// [Run] dials the hub via gRPC and calls the admin-gated Revoke
// RPC. Authentication uses the hub admin token (resolved from
// the --token flag or the CTX_HUB_ADMIN_TOKEN environment
// variable by the command layer), not the stored bearer token.
//
// # Data Flow
//
// When [Run] is called it performs these steps:
//
//  1. Loads connection config to obtain the hub address
//     (same source as ctx hub status).
//  2. Dials the hub via gRPC using hub.NewClient with no
//     bearer token; the Revoke RPC is admin-gated instead.
//  3. Calls the Revoke RPC with the admin token and the
//     target client ID.
//  4. On success, delegates to writeHub.Revoked to confirm
//     the revocation to the user.
package revoke
