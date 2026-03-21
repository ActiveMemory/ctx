//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package message

import (
	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	"github.com/spf13/cobra"
)

// Cmd returns the "ctx system message" subcommand.
//
// Returns:
//   - *cobra.Command: Configured message subcommand with sub-subcommands
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeySystemMessage)

	cmd := &cobra.Command{
		Use:   cmd.UseSystemMessage,
		Short: short,
		Long:  long,
	}

	cmd.AddCommand(
		messageListCmd(),
		messageShowCmd(),
		messageEditCmd(),
		messageResetCmd(),
	)

	return cmd
}

// messageListCmd returns the "ctx system message list" subcommand.
func messageListCmd() *cobra.Command {
	short, _ := desc.Command(cmd.DescKeySystemMessageList)

	cmd := &cobra.Command{
		Use:   cmd.UseSystemMessageList,
		Short: short,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return RunMessageList(cmd)
		},
	}
	cmd.Flags().Bool("json", false, desc.Flag(flag.DescKeySystemMessageJson))
	return cmd
}

// messageShowCmd returns the "ctx system message show" subcommand.
func messageShowCmd() *cobra.Command {
	short, _ := desc.Command(cmd.DescKeySystemMessageShow)

	return &cobra.Command{
		Use:   cmd.UseSystemMessageShow,
		Short: short,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunMessageShow(cmd, args[0], args[1])
		},
	}
}

// messageEditCmd returns the "ctx system message edit" subcommand.
func messageEditCmd() *cobra.Command {
	short, _ := desc.Command(cmd.DescKeySystemMessageEdit)

	return &cobra.Command{
		Use:   cmd.UseSystemMessageEdit,
		Short: short,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunMessageEdit(cmd, args[0], args[1])
		},
	}
}

// messageResetCmd returns the "ctx system message reset" subcommand.
func messageResetCmd() *cobra.Command {
	short, _ := desc.Command(cmd.DescKeySystemMessageReset)

	return &cobra.Command{
		Use:   cmd.UseSystemMessageReset,
		Short: short,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunMessageReset(cmd, args[0], args[1])
		},
	}
}
