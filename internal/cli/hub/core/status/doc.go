//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package status implements cluster status display for
// the ctx hub status command.
//
// # Overview
//
// This package queries a remote hub for its cluster
// state and renders a summary showing the node role,
// the leader, the entry count and the peer count.
//
// # Behavior
//
// [Run] dials the hub via gRPC and renders what the
// Status RPC reports. The role, the leader and the
// peer count are the hub's answers, read from its
// Raft node: a hub started without peers reports
// Standalone and names no leader.
//
// # Data Flow
//
// When [Run] is called it performs these steps:
//
//  1. Loads connection config to obtain the hub
//     address and authentication token.
//  2. Dials the hub via gRPC using hub.NewClient.
//  3. Calls the Status RPC, whose response carries
//     the entry count, the listener counts and the
//     cluster leadership fields.
//  4. Maps the response onto writeHub's render
//     fields: Standalone when no Raft node is
//     attached, otherwise Leader or Follower from
//     the reported leadership state.
//  5. Delegates to writeHub.ClusterStatus to render
//     the result.
package status
