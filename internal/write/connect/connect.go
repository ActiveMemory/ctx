//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package connect

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Registered confirms a successful hub registration.
//
// Parameters:
//   - cmd: Cobra command for output
//   - clientID: assigned client identifier
func Registered(cmd *cobra.Command, clientID string) {
	cmd.Println("Registered as", clientID)
}

// Subscribed confirms subscription types were updated.
//
// Parameters:
//   - cmd: Cobra command for output
//   - types: subscribed entry types
func Subscribed(cmd *cobra.Command, types []string) {
	cmd.Println("Subscribed to", types)
}

// Synced confirms entries were synced from the hub.
//
// Parameters:
//   - cmd: Cobra command for output
//   - count: number of entries synced
func Synced(cmd *cobra.Command, count int) {
	cmd.Println(fmt.Sprintf("Synced %d entries", count))
}

// Published confirms entries were published to the hub.
//
// Parameters:
//   - cmd: Cobra command for output
//   - count: number of entries published
func Published(cmd *cobra.Command, count int) {
	cmd.Println(fmt.Sprintf("Published %d entries", count))
}

// Listening confirms the listen stream is active.
//
// Parameters:
//   - cmd: Cobra command for output
func Listening(cmd *cobra.Command) {
	cmd.Println("Listening for new entries (Ctrl-C to stop)")
}

// EntryReceived confirms a single entry was received via
// the Listen stream.
//
// Parameters:
//   - cmd: Cobra command for output
//   - entryType: type of the received entry
func EntryReceived(cmd *cobra.Command, entryType string) {
	cmd.Println("Received", entryType)
}

// PublishFailed warns that hub publish failed.
//
// Parameters:
//   - cmd: Cobra command for output
//   - pubErr: the publish error
func PublishFailed(cmd *cobra.Command, pubErr error) {
	cmd.Println(fmt.Sprintf(
		"Warning: hub publish failed: %v", pubErr,
	))
}

// Status prints hub connection information.
//
// Parameters:
//   - cmd: Cobra command for output
//   - addr: hub address
//   - total: total entries on hub
//   - clients: connected client count
func Status(
	cmd *cobra.Command,
	addr string,
	total uint64,
	clients uint32,
) {
	cmd.Println("Hub:", addr)
	cmd.Println(fmt.Sprintf(
		"Entries: %d  Clients: %d", total, clients,
	))
}
