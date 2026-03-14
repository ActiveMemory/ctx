//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package edit

import (
	"strconv"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
	ctxerr "github.com/ActiveMemory/ctx/internal/err"
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

	short, long := assets.CommandDesc(assets.CmdDescKeyPadEdit)
	cmd := &cobra.Command{
		Use:   "edit N [TEXT]",
		Short: short,
		Long:  long,
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			n, err := strconv.Atoi(args[0])
			if err != nil {
				return ctxerr.InvalidIndex(args[0])
			}

			hasPositional := len(args) == 2
			hasAppend := appendText != ""
			hasPrepend := prependText != ""
			hasFile := filePath != ""
			hasLabel := labelText != ""

			// --file/--label conflict with positional/--append/--prepend.
			if (hasFile || hasLabel) && (hasPositional || hasAppend || hasPrepend) {
				return ctxerr.EditBlobTextConflict()
			}

			// Blob edit mode.
			if hasFile || hasLabel {
				return runEditBlob(cmd, n, filePath, labelText)
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
				return ctxerr.EditNoMode()
			}
			if flagCount > 1 {
				return ctxerr.EditTextConflict()
			}

			switch {
			case hasAppend:
				return runEditAppend(cmd, n, appendText)
			case hasPrepend:
				return runEditPrepend(cmd, n, prependText)
			default:
				return runEdit(cmd, n, args[1])
			}
		},
	}

	cmd.Flags().StringVar(&appendText,
		"append", "", assets.FlagDesc(assets.FlagDescKeyPadEditAppend),
	)
	cmd.Flags().StringVar(&prependText,
		"prepend", "", assets.FlagDesc(assets.FlagDescKeyPadEditPrepend),
	)
	cmd.Flags().StringVarP(&filePath,
		"file", "f", "", assets.FlagDesc(assets.FlagDescKeyPadEditFile),
	)
	cmd.Flags().StringVar(&labelText,
		"label", "", assets.FlagDesc(assets.FlagDescKeyPadEditLabel),
	)

	return cmd
}
