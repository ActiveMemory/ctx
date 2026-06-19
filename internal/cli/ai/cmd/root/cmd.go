//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	"github.com/ActiveMemory/ctx/internal/flagbind"
)

// Cmd returns the ctx ai command tree.
//
// Returns:
//   - *cobra.Command: configured ai command with subcommands
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeyAI)
	c := &cobra.Command{
		Use:   cmd.UseAI,
		Short: short,
		Long:  long,
	}
	var backendName string
	pingShort, pingLong := desc.Command(cmd.DescKeyAIPing)
	ping := &cobra.Command{
		Use:   cmd.UseAIPing,
		Short: pingShort,
		Long:  pingLong,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunPing(cmd, backendName)
		},
	}
	flagbind.StringFlag(ping, &backendName, cFlag.Backend, flag.DescKeyAIBackend)
	var proposeBackend string
	var emit string
	proposeShort, proposeLong := desc.Command(cmd.DescKeyAIPropose)
	propose := &cobra.Command{
		Use:   cmd.UseAIPropose,
		Short: proposeShort,
		Long:  proposeLong,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunPropose(cmd, args[0], proposeBackend, emit)
		},
	}
	flagbind.StringFlag(
		propose,
		&proposeBackend,
		cFlag.Backend,
		flag.DescKeyAIBackend,
	)
	flagbind.StringFlag(propose, &emit, cFlag.Emit, flag.DescKeyAIEmit)
	c.AddCommand(ping)
	c.AddCommand(propose)
	return c
}
