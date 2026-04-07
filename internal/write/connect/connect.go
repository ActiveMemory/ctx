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
