//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package message

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/msg"
	"github.com/ActiveMemory/ctx/internal/io"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/assets/hooks/messages"
	"github.com/ActiveMemory/ctx/internal/cli/system/core"
	ctxerr "github.com/ActiveMemory/ctx/internal/err"
)

// RunMessageList executes the message list logic.
//
// Collects all registered hook messages from the registry and outputs
// them as either a JSON array or a formatted table.
//
// Parameters:
//   - cmd: Cobra command for output and flag access
//
// Returns:
//   - error: Non-nil on JSON encoding failure
func RunMessageList(cmd *cobra.Command) error {
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

	jsonFlag, _ := cmd.Flags().GetBool("json")
	if jsonFlag {
		enc := json.NewEncoder(cmd.OutOrStdout())
		enc.SetIndent("", "  ")
		return enc.Encode(entries)
	}

	// Table output
	headerFmt := fmt.Sprintf("%%-%ds %%-%ds %%-%ds %%s",
		msg.MessageColHook, msg.MessageColVariant, msg.MessageColCategory)
	cmd.Println(fmt.Sprintf(headerFmt,
		assets.TextDesc(assets.TextDescKeyMessageListHeaderHook),
		assets.TextDesc(assets.TextDescKeyMessageListHeaderVariant),
		assets.TextDesc(assets.TextDescKeyMessageListHeaderCategory),
		assets.TextDesc(assets.TextDescKeyMessageListHeaderOverride)))
	cmd.Println(fmt.Sprintf(headerFmt,
		strings.Repeat("\u2500", msg.MessageSepHook),
		strings.Repeat("\u2500", msg.MessageSepVariant),
		strings.Repeat("\u2500", msg.MessageSepCategory),
		strings.Repeat("\u2500", msg.MessageSepOverride)))

	for _, e := range entries {
		override := ""
		if e.HasOverride {
			override = assets.TextDesc(assets.TextDescKeyMessageOverrideLabel)
		}
		cmd.Println(fmt.Sprintf(headerFmt, e.Hook, e.Variant, e.Category, override))
	}

	return nil
}

// RunMessageShow executes the message show logic.
//
// Displays the content of a hook message template, checking for a user
// override first and falling back to the embedded default.
//
// Parameters:
//   - cmd: Cobra command for output
//   - hook: hook name
//   - variant: template variant name
//
// Returns:
//   - error: Non-nil if the hook/variant is unknown or template is missing
func RunMessageShow(cmd *cobra.Command, hook, variant string) error {
	info := messages.Lookup(hook, variant)
	if info == nil {
		return core.ValidationError(hook, variant)
	}

	// Check user override first
	oPath := core.OverridePath(hook, variant)
	if data, readErr := io.SafeReadUserFile(oPath); readErr == nil {
		cmd.Println(fmt.Sprintf(assets.TextDesc(assets.TextDescKeyMessageSourceOverride), oPath))
		core.PrintTemplateVars(cmd, info)
		cmd.Println()
		cmd.Print(string(data))
		if len(data) > 0 && data[len(data)-1] != '\n' {
			cmd.Println()
		}
		return nil
	}

	// Embedded default
	data, readErr := assets.HookMessage(hook, variant+file.ExtTxt)
	if readErr != nil {
		return ctxerr.EmbeddedTemplateNotFound(hook, variant)
	}

	cmd.Println(assets.TextDesc(assets.TextDescKeyMessageSourceDefault))
	core.PrintTemplateVars(cmd, info)
	cmd.Println()
	cmd.Print(string(data))
	if len(data) > 0 && data[len(data)-1] != '\n' {
		cmd.Println()
	}
	return nil
}

// RunMessageEdit executes the message edit logic.
//
// Creates a user override file by copying the embedded default template
// to the project's .context/hooks/messages/ directory.
//
// Parameters:
//   - cmd: Cobra command for output
//   - hook: hook name
//   - variant: template variant name
//
// Returns:
//   - error: Non-nil if the hook/variant is unknown, override exists,
//     or file operations fail
func RunMessageEdit(cmd *cobra.Command, hook, variant string) error {
	info := messages.Lookup(hook, variant)
	if info == nil {
		return core.ValidationError(hook, variant)
	}

	oPath := core.OverridePath(hook, variant)

	// Refuse if override already exists
	if _, statErr := os.Stat(oPath); statErr == nil {
		return ctxerr.OverrideExists(oPath, hook, variant)
	}

	// Warn for ctx-specific messages
	if info.Category == messages.CategoryCtxSpecific {
		cmd.Println(assets.TextDesc(assets.TextDescKeyMessageCtxSpecificWarning))
		cmd.Println()
	}

	// Read embedded default
	data, readErr := assets.HookMessage(hook, variant+file.ExtTxt)
	if readErr != nil {
		return ctxerr.EmbeddedTemplateNotFound(hook, variant)
	}

	// Create directories
	dir := filepath.Dir(oPath)
	if mkdirErr := os.MkdirAll(dir, 0o750); mkdirErr != nil {
		return ctxerr.CreateDir(dir, mkdirErr)
	}

	// Write override file
	if writeErr := os.WriteFile(oPath, data, 0o600); writeErr != nil {
		return ctxerr.WriteOverride(oPath, writeErr)
	}

	cmd.Println(fmt.Sprintf(assets.TextDesc(assets.TextDescKeyMessageOverrideCreated), oPath))
	cmd.Println(assets.TextDesc(assets.TextDescKeyMessageEditHint))
	core.PrintTemplateVars(cmd, info)

	return nil
}

// RunMessageReset executes the message reset logic.
//
// Removes a user override file, reverting to the embedded default.
// Cleans up empty parent directories after removal.
//
// Parameters:
//   - cmd: Cobra command for output
//   - hook: hook name
//   - variant: template variant name
//
// Returns:
//   - error: Non-nil if the hook/variant is unknown or removal fails
func RunMessageReset(cmd *cobra.Command, hook, variant string) error {
	info := messages.Lookup(hook, variant)
	if info == nil {
		return core.ValidationError(hook, variant)
	}

	oPath := core.OverridePath(hook, variant)

	if removeErr := os.Remove(oPath); removeErr != nil {
		if os.IsNotExist(removeErr) {
			cmd.Println(fmt.Sprintf(assets.TextDesc(assets.TextDescKeyMessageNoOverride), hook, variant))
			return nil
		}
		return ctxerr.RemoveOverride(oPath, removeErr)
	}

	// Clean up empty parent directories
	hookDir := filepath.Dir(oPath)
	_ = os.Remove(hookDir) // only succeeds if empty
	messagesDir := filepath.Dir(hookDir)
	_ = os.Remove(messagesDir) // only succeeds if empty

	cmd.Println(fmt.Sprintf(assets.TextDesc(assets.TextDescKeyMessageOverrideRemoved), hook, variant))
	return nil
}
