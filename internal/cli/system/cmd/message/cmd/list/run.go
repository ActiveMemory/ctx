//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package list

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/hooks/messages"
	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/cli/system/core"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	cflag "github.com/ActiveMemory/ctx/internal/config/flag"
	"github.com/ActiveMemory/ctx/internal/config/msg"
)

// Run executes the message list logic.
//
// Parameters:
//   - cmd: Cobra command for output and flag access
//
// Returns:
//   - error: Non-nil on JSON encoding failure
func Run(cmd *cobra.Command) error {
	registry := messages.Registry()
	entries := make([]core.MessageListEntry, 0, len(registry))

	for _, info := range registry {
		entry := core.MessageListEntry{
			Hook:         info.Hook,
			Variant:      info.Variant,
			Category:     info.Category,
			Description:  info.Description,
			TemplateVars: info.TemplateVars,
			HasOverride:  core.HasOverride(info.Hook, info.Variant),
		}
		if entry.TemplateVars == nil {
			entry.TemplateVars = []string{}
		}
		entries = append(entries, entry)
	}

	jsonFlag, _ := cmd.Flags().GetBool(cflag.JSON)
	if jsonFlag {
		enc := json.NewEncoder(cmd.OutOrStdout())
		enc.SetIndent("", "  ")
		return enc.Encode(entries)
	}

	headerFmt := fmt.Sprintf("%%-%ds %%-%ds %%-%ds %%s",
		msg.MessageColHook, msg.MessageColVariant, msg.MessageColCategory)
	cmd.Println(fmt.Sprintf(headerFmt,
		desc.Text(text.DescKeyMessageListHeaderHook),
		desc.Text(text.DescKeyMessageListHeaderVariant),
		desc.Text(text.DescKeyMessageListHeaderCategory),
		desc.Text(text.DescKeyMessageListHeaderOverride)))
	cmd.Println(fmt.Sprintf(headerFmt,
		strings.Repeat("\u2500", msg.MessageSepHook),
		strings.Repeat("\u2500", msg.MessageSepVariant),
		strings.Repeat("\u2500", msg.MessageSepCategory),
		strings.Repeat("\u2500", msg.MessageSepOverride)))

	for _, e := range entries {
		override := ""
		if e.HasOverride {
			override = desc.Text(text.DescKeyMessageOverrideLabel)
		}
		cmd.Println(fmt.Sprintf(headerFmt, e.Hook, e.Variant, e.Category, override))
	}

	return nil
}
