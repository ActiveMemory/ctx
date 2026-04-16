//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package peer implements runtime peer management for
// the ctx hub peer command.
//
// # Overview
//
// This package provides the business logic for adding
// and removing peers from a hub cluster. When a user
// runs ctx hub peer add or ctx hub peer remove, the
// command layer delegates to [Run], which dispatches
// on the action argument and reports the result.
//
// # Behavior
//
// [Run] dispatches on the action argument ("add" or "remove")
// to register or deregister a peer address in the cluster.
//
// # Data Flow
//
// The peer management pipeline works as follows:
//
//  1. The cmd layer invokes [Run] with cobra args
//     containing [action, address].
//  2. [Run] switches on the action string, matching
//     against the configured add and remove constants
//     from the hub config package.
//  3. For "add", a confirmation message is printed
//     via writeHub.PeerAdded.
//  4. For "remove", a confirmation message is printed
//     via writeHub.PeerRemoved.
//  5. An invalid action returns an error from the
//     hub error package.
package peer
