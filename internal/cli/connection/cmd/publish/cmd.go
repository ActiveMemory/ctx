//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package publish

import (
	"time"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	corePub "github.com/ActiveMemory/ctx/internal/cli/connection/core/publish"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/hub"
)

// Cmd returns the connect publish subcommand.
//
// Returns:
//   - *cobra.Command: The publish subcommand
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeyConnectionPublish)

	return &cobra.Command{
		Use:     cmd.UseConnectionPublish,
		Short:   short,
		Long:    long,
		Example: desc.Example(cmd.DescKeyConnectionPublish),
		Args:    cobra.MinimumNArgs(2),
		RunE: func(
			cobraCmd *cobra.Command, args []string,
		) error {
			entry := hub.PublishEntry{
				Type:      args[0],
				Content:   args[1],
				Timestamp: time.Now().Unix(),
			}
			return corePub.Run(
				cobraCmd, []hub.PublishEntry{entry},
			)
		},
	}
}
