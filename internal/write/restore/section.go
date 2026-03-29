//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package restore

import (
	"fmt"

	"github.com/spf13/cobra"
)

// printSection prints a header and list items if the list is non-empty.
//
// Parameters:
//   - cmd: Cobra command for output
//   - headerTpl: Format string for the section header (receives item count)
//   - itemTpl: Format string for each item line (receives item name)
//   - items: Entries to display; skips output entirely when empty
func printSection(
	cmd *cobra.Command,
	headerTpl, itemTpl string,
	items []string,
) {
	if len(items) == 0 {
		return
	}
	cmd.Println(fmt.Sprintf(headerTpl, len(items)))
	for _, item := range items {
		cmd.Println(fmt.Sprintf(itemTpl, item))
	}
}
