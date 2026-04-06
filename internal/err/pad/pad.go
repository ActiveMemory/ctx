//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package pad

import (
	"errors"
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// EntryRange returns an error for an out-of-range scratchpad entry.
//
// Parameters:
//   - n: the requested entry number.
//   - total: the total number of entries.
//
// Returns:
//   - error: "entry <n> does not exist, scratchpad has <total> entries"
func EntryRange(n, total int) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrPadEntryRange), n, total,
	)
}

// EditBlobTextConflict returns an error when --file/--label and text
// editing flags are used together.
//
// Returns:
//   - error: describing the mutual exclusivity
func EditBlobTextConflict() error {
	return errors.New(
		desc.Text(text.DescKeyErrPadEditBlobTextConflict),
	)
}

// EditTextConflict returns an error when multiple text editing modes
// are used together.
//
// Returns:
//   - error: describing the mutual exclusivity
func EditTextConflict() error {
	return errors.New(
		desc.Text(text.DescKeyErrPadEditTextConflict),
	)
}

// EditNoMode returns an error when no editing mode was specified.
//
// Returns:
//   - error: prompting for a mode
func EditNoMode() error {
	return errors.New(
		desc.Text(text.DescKeyErrPadEditNoMode),
	)
}

// NotBlobEntry returns an error when a blob operation targets a non-blob.
//
// Parameters:
//   - n: the 1-based entry index.
//
// Returns:
//   - error: "entry <n> is not a blob entry"
func NotBlobEntry(n int) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrPadNotBlobEntry), n,
	)
}

// ResolveNotEncrypted returns an error when resolve is used on an
// unencrypted scratchpad.
//
// Returns:
//   - error: "resolve is only needed for encrypted scratchpads"
func ResolveNotEncrypted() error {
	return errors.New(
		desc.Text(text.DescKeyErrPadResolveNotEncrypted),
	)
}

// NoConflictFiles returns an error when no merge conflict files are found.
//
// Parameters:
//   - filename: the base scratchpad filename.
//
// Returns:
//   - error: "no conflict files found (<filename>.ours / <filename>.theirs)"
func NoConflictFiles(filename string) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrPadNoConflictFiles),
		filename, filename,
	)
}

// OutFlagRequiresBlob returns an error when --out is used on a non-blob entry.
//
// Returns:
//   - error: "--out can only be used with blob entries"
func OutFlagRequiresBlob() error {
	return errors.New(
		desc.Text(text.DescKeyErrPadOutFlagRequiresBlob),
	)
}

// Read wraps a scratchpad read failure.
//
// Parameters:
//   - cause: the underlying read error.
//
// Returns:
//   - error: "read scratchpad: <cause>"
func Read(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrPadReadScratchpad), cause,
	)
}

// InvalidIndex returns an error for a non-numeric entry index.
//
// Parameters:
//   - value: the invalid index string.
//
// Returns:
//   - error: "invalid index: <value>"
func InvalidIndex(value string) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrPadInvalidIndex), value,
	)
}

// FileTooLarge returns an error for a file exceeding the size limit.
//
// Parameters:
//   - size: actual file size in bytes.
//   - max: maximum allowed size in bytes.
//
// Returns:
//   - error: "file too large: <size> bytes (max <max>)"
func FileTooLarge(size, max int) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrPadFileTooLarge), size, max,
	)
}
