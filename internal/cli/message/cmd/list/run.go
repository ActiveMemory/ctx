//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package list

import (
	"encoding/json"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/hooks/messages"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/message"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/entity"
	"github.com/ActiveMemory/ctx/internal/rc"
	writeMessage "github.com/ActiveMemory/ctx/internal/write/message"
)

// Run executes the message list logic.
//
// Parameters:
//   - cmd: Cobra command for output and flag access
//
// Returns:
//   - error: Non-nil on JSON encoding failure
func Run(cmd *cobra.Command) error {
	if _, ctxErr := rc.RequireContextDir(); ctxErr != nil {
		cmd.SilenceUsage = true
		return ctxErr
	}
	registry := messages.Registry()
	entries := make([]entity.MessageListEntry, 0, len(registry))

	for _, info := range registry {
		hasOverride, overrideErr := message.HasOverride(info.Hook, info.Variant)
		if overrideErr != nil {
			return overrideErr
		}
		entry := entity.MessageListEntry{
			Hook:         info.Hook,
			Variant:      info.Variant,
			Category:     info.Category,
			Description:  info.Description,
			TemplateVars: info.TemplateVars,
			HasOverride:  hasOverride,
		}
		if entry.TemplateVars == nil {
			entry.TemplateVars = []string{}
		}
		entries = append(entries, entry)
	}

	jsonFlag, _ := cmd.Flags().GetBool(cFlag.JSON)
	if jsonFlag {
		enc := json.NewEncoder(cmd.OutOrStdout())
		enc.SetIndent("", token.Indent2)
		return enc.Encode(entries)
	}

	writeMessage.ListHeader(cmd)
	for _, e := range entries {
		writeMessage.ListRow(cmd, e.Hook, e.Variant, e.Category, e.HasOverride)
	}

	return nil
}
