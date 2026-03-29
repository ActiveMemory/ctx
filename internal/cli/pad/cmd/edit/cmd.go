//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package edit

import (
	"strconv"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	errPad "github.com/ActiveMemory/ctx/internal/err/pad"
)

// Cmd returns the pad edit subcommand.
//
// Supports three modes:
//   - Replace: ctx pad edit N "text"
//   - Append:  ctx pad edit N --append "text"
//   - Prepend: ctx pad edit N --prepend "text"
//   - Blob file: ctx pad edit N --file ./v2.md
//   - Blob label: ctx pad edit N --label "new label"
//
// The --append and --prepend flags are mutually exclusive with each other
// and with the positional replacement text argument.
// The --file and --label flags conflict with positional/--append/--prepend.
//
// Returns:
//   - *cobra.Command: Configured edit subcommand
func Cmd() *cobra.Command {
	var appendText string
	var prependText string
	var filePath string
	var labelText string

	short, long := desc.Command(cmd.DescKeyPadEdit)
	c := &cobra.Command{
		Use:   cmd.UsePadEdit,
		Short: short,
		Long:  long,
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			n, err := strconv.Atoi(args[0])
			if err != nil {
				return errPad.InvalidIndex(args[0])
			}

			hasPositional := len(args) == 2
			hasAppend := appendText != ""
			hasPrepend := prependText != ""
			hasFile := filePath != ""
			hasLabel := labelText != ""

			// --file/--label conflict with positional/--append/--prepend.
			if (hasFile || hasLabel) && (hasPositional || hasAppend || hasPrepend) {
				return errPad.EditBlobTextConflict()
			}

			// Blob edit mode.
			if hasFile || hasLabel {
				return Run(cmd, Opts{
					N:         n,
					FilePath:  filePath,
					LabelText: labelText,
					Mode:      ModeBlob,
				})
			}

			// Validate mutual exclusivity of positional/--append/--prepend.
			flagCount := 0
			if hasPositional {
				flagCount++
			}
			if hasAppend {
				flagCount++
			}
			if hasPrepend {
				flagCount++
			}

			if flagCount == 0 {
				return errPad.EditNoMode()
			}
			if flagCount > 1 {
				return errPad.EditTextConflict()
			}

			switch {
			case hasAppend:
				return Run(cmd, Opts{N: n, Text: appendText, Mode: ModeAppend})
			case hasPrepend:
				return Run(cmd, Opts{N: n, Text: prependText, Mode: ModePrepend})
			default:
				return Run(cmd, Opts{N: n, Text: args[1], Mode: ModeReplace})
			}
		},
	}

	c.Flags().StringVar(&appendText,
		cFlag.Append, "", desc.Flag(flag.DescKeyPadEditAppend),
	)
	c.Flags().StringVar(&prependText,
		cFlag.Prepend, "", desc.Flag(flag.DescKeyPadEditPrepend),
	)
	c.Flags().StringVarP(&filePath,
		cFlag.File, cFlag.ShortFile, "", desc.Flag(flag.DescKeyPadEditFile),
	)
	c.Flags().StringVar(&labelText,
		cFlag.Label, "", desc.Flag(flag.DescKeyPadEditLabel),
	)

	return c
}
