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
	coreEdit "github.com/ActiveMemory/ctx/internal/cli/pad/core/edit"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	errPad "github.com/ActiveMemory/ctx/internal/err/pad"
	"github.com/ActiveMemory/ctx/internal/flagbind"
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
// The --append and --prepend flags are mutually exclusive
// with each other and with the positional replacement text.
// The --file and --label flags conflict with
// positional/--append/--prepend.
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
		Use:     cmd.UsePadEdit,
		Short:   short,
		Long:    long,
		Example: desc.Example(cmd.DescKeyPadEdit),
		Args:    cobra.RangeArgs(1, 2),
		RunE: func(
			cmd *cobra.Command, args []string,
		) error {
			n, err := strconv.Atoi(args[0])
			if err != nil {
				return errPad.InvalidIndex(args[0])
			}

			hasPositional := len(args) == 2
			hasAppend := appendText != ""
			hasPrepend := prependText != ""
			hasFile := filePath != ""
			hasLabel := labelText != ""

			// --file/--label conflict with text modes.
			if (hasFile || hasLabel) &&
				(hasPositional || hasAppend || hasPrepend) {
				return errPad.EditBlobTextConflict()
			}

			// Blob edit mode.
			if hasFile || hasLabel {
				return Run(cmd, coreEdit.Opts{
					N:         n,
					FilePath:  filePath,
					LabelText: labelText,
					Mode:      coreEdit.ModeBlob,
				})
			}

			// Validate mutual exclusivity.
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
				return Run(cmd, coreEdit.Opts{
					N:    n,
					Text: appendText,
					Mode: coreEdit.ModeAppend,
				})
			case hasPrepend:
				return Run(cmd, coreEdit.Opts{
					N:    n,
					Text: prependText,
					Mode: coreEdit.ModePrepend,
				})
			default:
				return Run(cmd, coreEdit.Opts{
					N:    n,
					Text: args[1],
					Mode: coreEdit.ModeReplace,
				})
			}
		},
	}

	flagbind.StringFlag(c, &appendText,
		cFlag.Append, flag.DescKeyPadEditAppend,
	)
	flagbind.StringFlag(c, &prependText,
		cFlag.Prepend, flag.DescKeyPadEditPrepend,
	)
	flagbind.StringFlagP(c, &filePath,
		cFlag.File, cFlag.ShortFile,
		flag.DescKeyPadEditFile,
	)
	flagbind.StringFlag(c, &labelText,
		cFlag.Label, flag.DescKeyPadEditLabel,
	)

	return c
}
