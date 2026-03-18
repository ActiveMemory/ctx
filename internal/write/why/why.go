//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package why

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
)

// Banner prints the ctx ASCII art banner for the why menu.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
func Banner(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println(assets.TextDesc(assets.TextDescKeyWhyBanner))
}

// MenuItem prints a numbered menu item.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - index: 1-based menu index.
//   - label: display label for the document.
func MenuItem(cmd *cobra.Command, index int, label string) {
	if cmd == nil {
		return
	}
	cmd.Println(
		fmt.Sprintf(
			assets.TextDesc(assets.TextDescKeyWhyMenuItemFormat), index, label,
		),
	)
}

// MenuPrompt prints the selection prompt.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
func MenuPrompt(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Print(assets.TextDesc(assets.TextDescKeyWhyMenuPrompt))
}
