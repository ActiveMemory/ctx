//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package serve

import (
	"net"

	"github.com/spf13/cobra"
)

// HubStarted prints the hub server address.
//
// Parameters:
//   - cmd: Cobra command for output
//   - addr: network address the server is listening on
func HubStarted(cmd *cobra.Command, addr net.Addr) {
	cmd.Println("Hub started on", addr)
}

// AdminToken prints the generated admin token.
//
// Parameters:
//   - cmd: Cobra command for output
//   - token: the generated admin token
func AdminToken(cmd *cobra.Command, token string) {
	cmd.Println("Admin token (save this):", token)
}
