//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package connect provides terminal output for the hub
// client commands (ctx connect).
//
// # Registration and Subscription
//
// [Registered] confirms a successful hub registration
// and prints the assigned client ID. [Subscribed]
// confirms which entry types the client is subscribed
// to receive.
//
// # Data Transfer
//
// [Synced] reports how many entries were pulled from
// the hub. [Published] reports how many entries were
// pushed to the hub. [PublishFailed] warns when a
// publish operation fails without aborting.
//
// # Live Stream
//
// [Listening] confirms the listen stream is active.
// [EntryReceived] reports each entry received via the
// live stream with its type.
//
// # Hub Status
//
// [Status] prints the hub connection dashboard: the
// hub address, total entry count, and connected client
// count.
//
// # Message Categories
//
//   - Info: registration, sync, publish confirmations
//   - Warning: publish failures
//   - Status: hub connection dashboard
package connect
