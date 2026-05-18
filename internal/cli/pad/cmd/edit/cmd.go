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
	"github.com/ActiveMemory/ctx/internal/config/pad"
	errPad "github.com/ActiveMemory/ctx/internal/err/pad"
	"github.com/ActiveMemory/ctx/internal/flagbind"
)

// Cmd returns the pad edit subcommand.
//
// Supports these modes:
//   - Replace: ctx pad edit N "text"
//   - Append:  ctx pad edit N --append "text"
//   - Prepend: ctx pad edit N --prepend "text"
//   - Tag:     ctx pad edit N --tag tagname
//   - Blob file: ctx pad edit N --file ./v2.md
//   - Blob label: ctx pad edit N --label "new label"
//
// The --tag flag can be used alone or combined with other text
// modes. The --append and --prepend flags are mutually exclusive
// with each other and with the positional replacement text.
// The --file and --label flags conflict with
// positional/--append/--prepend/--tag.
//
// Returns:
//   - *cobra.Command: Configured edit subcommand
func Cmd() *cobra.Command {
	var appendText string
	var prependText string
	var filePath string
	var labelText string
	var tagName string

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
			hasTag := tagName != ""

			// --file/--label conflict with text modes.
			if (hasFile || hasLabel) &&
				(hasPositional || hasAppend || hasPrepend || hasTag) {
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

			// Validate mutual exclusivity among text modes
			// (--tag is compatible with all, so excluded here).
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

			if flagCount > 1 {
				return errPad.EditTextConflict()
			}

			// --tag alone is valid (append " #tagname").
			if flagCount == 0 && !hasTag {
				return errPad.EditNoMode()
			}

			// Apply the primary text mode first.
			switch {
			case hasAppend:
				text := appendText
				if hasTag {
					text = text + pad.TagPrefixSpace + tagName
				}
				return Run(cmd, coreEdit.Opts{
					N:    n,
					Text: text,
					Mode: coreEdit.ModeAppend,
				})
			case hasPrepend:
				if hasTag {
					// Apply prepend, then tag via
					// sequential operations.
					if runErr := Run(cmd, coreEdit.Opts{
						N:    n,
						Text: prependText,
						Mode: coreEdit.ModePrepend,
					}); runErr != nil {
						return runErr
					}
					return Run(cmd, coreEdit.Opts{
						N:    n,
						Text: pad.TagPrefix + tagName,
						Mode: coreEdit.ModeAppend,
					})
				}
				return Run(cmd, coreEdit.Opts{
					N:    n,
					Text: prependText,
					Mode: coreEdit.ModePrepend,
				})
			case hasPositional:
				text := args[1]
				if hasTag {
					text = text + pad.TagPrefixSpace + tagName
				}
				return Run(cmd, coreEdit.Opts{
					N:    n,
					Text: text,
					Mode: coreEdit.ModeReplace,
				})
			default:
				// --tag only: append the tag.
				return Run(cmd, coreEdit.Opts{
					N:    n,
					Text: pad.TagPrefix + tagName,
					Mode: coreEdit.ModeAppend,
				})
			}
		},
	}

	flagbind.BindStringFlags(c,
		[]*string{&appendText, &prependText, &labelText},
		[]string{cFlag.Append, cFlag.Prepend, cFlag.Label},
		[]string{
			flag.DescKeyPadEditAppend,
			flag.DescKeyPadEditPrepend,
			flag.DescKeyPadEditLabel,
		},
	)
	flagbind.BindStringFlagsP(c,
		[]*string{&filePath, &tagName},
		[]string{cFlag.File, cFlag.Tag},
		[]string{cFlag.ShortFile, cFlag.ShortTag},
		[]string{flag.DescKeyPadEditFile, flag.DescKeyPadEditTag},
	)

	return c
}
