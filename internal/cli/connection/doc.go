//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package connection provides the ctx connection command
// group for ctx Hub client operations.
//
// The connection command manages the relationship between
// a local project and a remote ctx Hub instance. It
// handles device registration, topic subscription,
// context publishing, and real-time event listening.
//
// # Subcommands
//
//   - register: register this device with a Hub instance
//   - subscribe: subscribe to context topics on the Hub
//   - sync: pull latest context from subscribed topics
//   - publish: push local context entries to the Hub
//   - listen: stream real-time events from the Hub
//   - status: show connection state and subscription info
//
// # Subpackages
//
//	cmd/register: device registration flow
//	cmd/subscribe: topic subscription management
//	cmd/sync: context pull from Hub
//	cmd/publish: context push to Hub
//	cmd/listen: real-time event streaming
//	cmd/status: connection status display
//	core: shared Hub client helpers
package connection
