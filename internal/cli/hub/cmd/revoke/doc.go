//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package revoke wires the ctx hub revoke subcommand.
//
// # Overview
//
// [Cmd] builds the cobra command that revokes a client's token
// on the hub. It takes the target client ID as its single
// positional argument and the admin token from the --token flag
// or the CTX_HUB_ADMIN_TOKEN environment variable, then
// delegates to the core revoke package.
//
// # Behavior
//
// The command resolves the admin token (flag first, then
// environment) and fails early with a clear error if neither is
// set. Like the other hub subcommands it skips context init,
// since the hub lives at ~/.ctx/hub-data/ rather than .context/.
package revoke
