//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package register

import (
	"context"
	"path/filepath"

	"github.com/spf13/cobra"

	connectCfg "github.com/ActiveMemory/ctx/internal/cli/connect/core/config"
	"github.com/ActiveMemory/ctx/internal/hub"
	"github.com/ActiveMemory/ctx/internal/rc"
	writeConnect "github.com/ActiveMemory/ctx/internal/write/connect"
)

// Run registers this project with a shared context hub.
//
// Connects to the hub, sends the admin token and project
// name, receives a client token, and stores the encrypted
// connection config in .context/.connect.enc.
//
// Parameters:
//   - cmd: cobra command for output
//   - hubAddr: hub gRPC address (host:port)
//   - adminToken: admin token from hub startup
//
// Returns:
//   - error: non-nil if registration or storage fails
func Run(
	cmd *cobra.Command,
	hubAddr string,
	adminToken string,
) error {
	client, dialErr := hub.NewClient(hubAddr, "")
	if dialErr != nil {
		return dialErr
	}
	defer func() { _ = client.Close() }()

	projectName := filepath.Base(rc.ContextDir())

	resp, regErr := client.Register(
		context.Background(),
		adminToken,
		projectName,
	)
	if regErr != nil {
		return regErr
	}

	cfg := connectCfg.Config{
		HubAddr: hubAddr,
		Token:   resp.ClientToken,
	}
	if saveErr := connectCfg.Save(cfg); saveErr != nil {
		return saveErr
	}

	writeConnect.Registered(cmd, resp.ClientID)
	return nil
}
