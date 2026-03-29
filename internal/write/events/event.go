//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package events

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	writeIO "github.com/ActiveMemory/ctx/internal/write/line"
)

// JSON prints JSONL event lines. Nil cmd is a no-op.
//
// Parameters:
//   - cmd: Cobra command for output
//   - lines: pre-formatted JSON lines
func JSON(cmd *cobra.Command, lines []string) {
	writeIO.All(cmd, lines)
}

// Human prints formatted event lines. Nil cmd is a no-op.
//
// Parameters:
//   - cmd: Cobra command for output
//   - lines: pre-formatted human-readable lines
func Human(cmd *cobra.Command, lines []string) {
	writeIO.All(cmd, lines)
}

// Empty prints the "no events" message. Nil cmd is a no-op.
//
// Parameters:
//   - cmd: Cobra command for output
func Empty(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println(desc.Text(text.DescKeyEventsEmpty))
}
