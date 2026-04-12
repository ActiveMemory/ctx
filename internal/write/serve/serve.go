//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package serve

import (
	"fmt"
	"net"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// HubStarted prints the hub server address.
//
// Parameters:
//   - cmd: Cobra command for output
//   - addr: network address the server is listening on
func HubStarted(cmd *cobra.Command, addr net.Addr) {
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteServeHubStarted), addr,
	))
}

// AdminToken prints the generated admin token.
//
// Parameters:
//   - cmd: Cobra command for output
//   - token: the generated admin token
func AdminToken(cmd *cobra.Command, token string) {
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteServeAdminToken), token,
	))
}

// Daemonized confirms the hub started in background.
//
// Parameters:
//   - cmd: Cobra command for output
//   - pid: process ID of the daemon
func Daemonized(cmd *cobra.Command, pid int) {
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteServeHubBackground), pid,
	))
}

// Stopped confirms the hub daemon was killed.
//
// Parameters:
//   - cmd: Cobra command for output
//   - pid: process ID that was stopped
func Stopped(cmd *cobra.Command, pid int) {
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteServeHubStopped), pid,
	))
}
