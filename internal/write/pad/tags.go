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

// TagsItem prints a single tag with its count.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - name: Tag name without the # prefix.
//   - count: Number of entries containing the tag.
func TagsItem(cmd *cobra.Command, name string, count int) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWritePadTagsItem), name, count))
}

// TagsJSON prints the JSON-encoded tag data.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - data: Pre-marshaled JSON bytes.
func TagsJSON(cmd *cobra.Command, data []byte) {
	if cmd == nil {
		return
	}
	cmd.Println(string(data))
}

// TagsNone prints the message when no tags are found.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
func TagsNone(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println(desc.Text(text.DescKeyWritePadTagsNone))
}
