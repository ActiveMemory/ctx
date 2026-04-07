//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package connect

import "github.com/spf13/cobra"

// Registered confirms a successful hub registration.
//
// Parameters:
//   - cmd: Cobra command for output
//   - clientID: assigned client identifier
func Registered(cmd *cobra.Command, clientID string) {
	cmd.Println("Registered as", clientID)
}
