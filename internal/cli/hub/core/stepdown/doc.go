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
// [Run] signals the current hub node to relinquish its
// leader role and prints a confirmation once the transfer
// is initiated.
//
// # Data Flow
//
// The stepdown pipeline works as follows:
//
//  1. The cmd layer invokes [Run] with the cobra
//     command and unused args.
//  2. [Run] calls writeHub.SteppedDown to print a
//     confirmation message indicating the node has
//     initiated leadership transfer.
//  3. The function returns nil on success. Future
//     implementations may add gRPC calls to
//     coordinate the transfer with the cluster.
package stepdown
