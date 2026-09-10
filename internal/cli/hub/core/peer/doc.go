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
// [Run] validates the action, dials the hub named by the
// saved connection config, and calls the admin-gated Peer
// RPC. The confirmation prints only after the hub has
// committed the configuration change; a follower answers
// FailedPrecondition, because only the leader can change
// cluster membership.
//
// # Data Flow
//
// The peer management pipeline works as follows:
//
//  1. The cmd layer resolves the admin token and invokes
//     [Run] with the action and the peer's Raft address.
//  2. [Run] rejects an action that is neither add nor
//     remove, using the configured constants from the hub
//     config package.
//  3. The connection config supplies the hub address; the
//     client dials without a bearer token, since the RPC is
//     authenticated by the admin credential.
//  4. The hub applies the change through raft AddVoter or
//     RemoveServer and returns.
//  5. writeHub.PeerAdded or writeHub.PeerRemoved reports
//     what the cluster accepted.
package peer
