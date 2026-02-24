//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package pad

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

// editCmd returns the pad edit subcommand.
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
func editCmd() *cobra.Command {
	var appendText string
	var prependText string
	var filePath string
	var labelText string

	cmd := &cobra.Command{
		Use:   "edit N [TEXT]",
		Short: "Replace, append to, or prepend to an entry by number",
		Long: `Replace, append to, or prepend to an entry by number.

By default, replaces the entire entry with the positional TEXT argument.
Use --append to add text to the end of an existing entry, or --prepend
to add text to the beginning.

For blob entries, use --file to replace file content and/or --label to
change the label.

Examples:
  ctx pad edit 2 "new text"           # replace entry 2
  ctx pad edit 2 --append "suffix"    # append to entry 2
  ctx pad edit 2 --prepend "prefix"   # prepend to entry 2
  ctx pad edit 2 --file ./v2.md       # replace blob file content
  ctx pad edit 2 --label "new name"   # rename blob label
  ctx pad edit 2 --file ./v2.md --label "new"  # replace both`,
		Args: cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			n, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid index: %s", args[0])
			}

			hasPositional := len(args) == 2
			hasAppend := appendText != ""
			hasPrepend := prependText != ""
			hasFile := filePath != ""
			hasLabel := labelText != ""

			// --file/--label conflict with positional/--append/--prepend.
			if (hasFile || hasLabel) && (hasPositional || hasAppend || hasPrepend) {
				return fmt.Errorf("--file/--label and positional text/--append/--prepend are mutually exclusive")
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
				return fmt.Errorf("provide replacement text, --append, or --prepend")
			}
			if flagCount > 1 {
				return fmt.Errorf("--append, --prepend, and positional text are mutually exclusive")
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

	cmd.Flags().StringVar(&appendText, "append", "", "append text to the end of the entry")
	cmd.Flags().StringVar(&prependText, "prepend", "", "prepend text to the beginning of the entry")
	cmd.Flags().StringVarP(&filePath, "file", "f", "", "replace blob file content")
	cmd.Flags().StringVar(&labelText, "label", "", "replace blob label")

	return cmd
}

// runEdit replaces entry at 1-based position n with new text.
func runEdit(cmd *cobra.Command, n int, text string) error {
	entries, err := readEntries()
	if err != nil {
		return err
	}

	if err := validateIndex(n, entries); err != nil {
		return err
	}

	entries[n-1] = text

	if err := writeEntries(entries); err != nil {
		return err
	}

	cmd.Printf("Updated entry %d.\n", n)
	return nil
}

// runEditAppend appends text to the entry at 1-based position n.
// The result is: old + " " + text.
func runEditAppend(cmd *cobra.Command, n int, text string) error {
	entries, err := readEntries()
	if err != nil {
		return err
	}

	if err := validateIndex(n, entries); err != nil {
		return err
	}

	if isBlob(entries[n-1]) {
		return fmt.Errorf("cannot append to a blob entry")
	}

	entries[n-1] = entries[n-1] + " " + text

	if err := writeEntries(entries); err != nil {
		return err
	}

	cmd.Printf("Updated entry %d.\n", n)
	return nil
}

// runEditPrepend prepends text to the entry at 1-based position n.
// The result is: text + " " + old.
func runEditPrepend(cmd *cobra.Command, n int, text string) error {
	entries, err := readEntries()
	if err != nil {
		return err
	}

	if err := validateIndex(n, entries); err != nil {
		return err
	}

	if isBlob(entries[n-1]) {
		return fmt.Errorf("cannot prepend to a blob entry")
	}

	entries[n-1] = text + " " + entries[n-1]

	if err := writeEntries(entries); err != nil {
		return err
	}

	cmd.Printf("Updated entry %d.\n", n)
	return nil
}

// runEditBlob replaces the file content and/or label of a blob entry.
func runEditBlob(cmd *cobra.Command, n int, filePath, labelText string) error {
	entries, err := readEntries()
	if err != nil {
		return err
	}

	if err := validateIndex(n, entries); err != nil {
		return err
	}

	oldLabel, oldData, ok := splitBlob(entries[n-1])
	if !ok {
		return fmt.Errorf("entry %d is not a blob entry", n)
	}

	newLabel := oldLabel
	newData := oldData

	if labelText != "" {
		newLabel = labelText
	}

	if filePath != "" {
		data, err := os.ReadFile(filePath) //nolint:gosec // user-provided path is intentional
		if err != nil {
			return fmt.Errorf("read file: %w", err)
		}
		if len(data) > MaxBlobSize {
			return fmt.Errorf("file too large: %d bytes (max %d)", len(data), MaxBlobSize)
		}
		newData = data
	}

	entries[n-1] = makeBlob(newLabel, newData)

	if err := writeEntries(entries); err != nil {
		return err
	}

	cmd.Printf("Updated entry %d.\n", n)
	return nil
}
