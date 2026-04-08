//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package pad

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// EntryAdded prints confirmation that a pad entry was added.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - n: entry number (1-based).
func EntryAdded(cmd *cobra.Command, n int) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWritePadEntryAdded), n))
}

// EntryUpdated prints confirmation that a pad entry was updated.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - n: entry number (1-based).
func EntryUpdated(cmd *cobra.Command, n int) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWritePadEntryUpdated), n))
}

// EntryRemoved prints confirmation that a pad entry was removed.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - n: entry number (1-based).
func EntryRemoved(cmd *cobra.Command, n int) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWritePadEntryRemoved), n))
}

// Normalized prints confirmation that IDs were normalized.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - count: number of entries renumbered.
func Normalized(cmd *cobra.Command, count int) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWritePadNormalized),
		count))
}

// EntryMoved prints confirmation that a pad entry was moved.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - from: source position (1-based).
//   - to: destination position (1-based).
func EntryMoved(cmd *cobra.Command, from, to int) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWritePadEntryMoved), from, to))
}

// EntryShow prints a pad entry with a trailing newline.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - entry: Entry text.
func EntryShow(cmd *cobra.Command, entry string) {
	if cmd == nil {
		return
	}
	cmd.Println(entry)
}

// EntryList prints a formatted pad list item.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - line: Pre-formatted list item string.
func EntryList(cmd *cobra.Command, line string) {
	if cmd == nil {
		return
	}
	cmd.Println(line)
}
