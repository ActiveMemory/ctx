//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package collect

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	coreCollect "github.com/ActiveMemory/ctx/internal/cli/trace/core/collect"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	"github.com/ActiveMemory/ctx/internal/flagbind"
)

// Cmd returns the trace collect subcommand.
//
// Returns:
//   - *cobra.Command: Configured trace collect command with flags registered
func Cmd() *cobra.Command {
	var record string
	short, long := desc.Command(cmd.DescKeyTraceCollect)
	c := &cobra.Command{
		Use:    cmd.UseTraceCollect,
		Short:  short,
		Long:   long,
		Hidden: true,
		Args:   cobra.ExactArgs(0),
		RunE: func(cobraCmd *cobra.Command, _ []string) error {
			if record != "" {
				return coreCollect.RecordCommit(record)
			}
			return Run(cobraCmd)
		},
	}
	flagbind.StringFlag(c, &record, cFlag.Record, flag.DescKeyTraceCollectRecord)
	return c
}
