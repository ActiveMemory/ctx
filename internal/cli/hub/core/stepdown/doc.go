//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package stepdown implements graceful leadership
// transfer for the ctx hub stepdown command.
//
// # Overview
//
// This package provides the business logic for
// requesting that the current hub node relinquish
// its leader role. When a user runs ctx hub stepdown,
// the command layer delegates to [Run], which signals
// the transfer and reports the result.
//
// # Behavior
//
// [Run] asks the hub to hand leadership to a follower and
// prints a confirmation once the transfer returns. Only a
// leader can transfer leadership: a follower answers
// FailedPrecondition rather than reporting a handoff that
// did not happen.
//
// # Data Flow
//
// The stepdown pipeline works as follows:
//
//  1. The cmd layer resolves the admin token and invokes
//     [Run].
//  2. The connection config supplies the hub address; the
//     client dials without a bearer token, since the RPC is
//     authenticated by the admin credential.
//  3. The hub calls raft LeadershipTransfer and returns
//     when it completes or fails.
//  4. writeHub.SteppedDown reports the handoff. Which node
//     won is answered by ctx hub status.
package stepdown
