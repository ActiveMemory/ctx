//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package subscribe implements subscription type management
// for the ctx connection subscribe command.
//
// # Overview
//
// This package provides the business logic for updating
// which entry types a connection subscribes to. When a
// user runs ctx connection subscribe, the command layer
// delegates to [Run], which persists the new subscription
// list to the connection configuration.
//
// # Behavior
//
// [Run] replaces the subscribed entry types in the
// connection config file and persists the change to disk.
//
// # Data Flow
//
// The subscribe pipeline works as follows:
//
//  1. The cmd layer invokes [Run] with cobra args
//     containing the desired entry types.
//  2. [Run] loads the current connection config via
//     the config sub-package.
//  3. The Types field is replaced with the new list.
//  4. The updated config is saved back to disk.
//  5. A confirmation message is printed via the
//     write/connect layer.
//
// If the config file cannot be loaded or saved, the
// error propagates back to the cmd layer for display.
package subscribe
