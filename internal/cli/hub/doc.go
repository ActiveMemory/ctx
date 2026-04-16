//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package hub provides the ctx hub command group for
// cluster management operations.
//
// The Hub is a lightweight server that synchronizes
// context across multiple machines or team members. The
// hub command group manages the server lifecycle and
// cluster topology from the command line.
//
// # Subcommands
//
//   - start: launch a Hub instance on the local machine
//   - stop: gracefully shut down a running Hub
//   - status: display the Hub's health, peer list, and
//     subscription counts
//   - peer: add or remove peer Hub nodes for multi-node
//     replication
//   - stepdown: ask the current leader to yield its role
//     to another node
//
// # Subpackages
//
//	cmd/start: server startup logic
//	cmd/stop: graceful shutdown
//	cmd/status: health and topology display
//	cmd/peer: peer management
//	cmd/stepdown: leader yield
//	core: shared Hub client and config helpers
package hub
