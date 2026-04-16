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
// address, total entries, and project count.
//
// # Behavior
//
// [Run] dials the hub via gRPC, retrieves cluster metrics
// (connected clients, entry count, projects), and renders
// a summary showing node role, address, and totals.
//
// # Data Flow
//
// When [Run] is called it performs these steps:
//
//  1. Loads connection config to obtain the hub
//     address and authentication token.
//  2. Dials the hub via gRPC using hub.NewClient.
//  3. Calls the Status RPC to retrieve cluster
//     metrics including connected clients, total
//     entries, and per-project breakdowns.
//  4. Determines the node role: if there are
//     connected clients the node is marked active,
//     otherwise it is a follower.
//  5. Delegates to writeHub.ClusterStatus to render
//     the role, address, entry count, and project
//     count for the user.
package status
