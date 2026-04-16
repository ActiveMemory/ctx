//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package register implements hub registration logic for
// the ctx connection register command.
//
// # Run
//
// [Run] registers the current project with a ctx Hub
// instance. It exchanges an admin token for a client
// token and persists the encrypted connection config.
//
// The execution flow is:
//
//  1. Dial the hub at the provided gRPC address using
//     hub.NewClient with an empty bearer token (the
//     admin token is sent as a registration parameter,
//     not as a connection credential).
//  2. Derive the project name from the context
//     directory path using filepath.Base.
//  3. Call client.Register with the admin token and
//     project name. The hub returns a client ID and a
//     client bearer token for future RPCs.
//  4. Build a connectCfg.Config with the hub address
//     and client token, then persist it via
//     connectCfg.Save. The config is encrypted at rest
//     in .context/.connect.enc.
//  5. Print a confirmation with the assigned client ID
//     via writeConnect.Registered.
//
// The function returns an error if dialing, registration,
// or config persistence fails. The gRPC connection is
// closed via a deferred Close call.
//
// # Data Flow
//
// The cmd/ layer extracts the hub address and admin
// token from flags or arguments and passes them to Run.
// After registration, the listen, publish, and status
// subpackages use the stored config for authenticated
// hub communication.
package register
