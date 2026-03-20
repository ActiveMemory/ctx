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

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/assets/read/hook"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/msg"
	"github.com/ActiveMemory/ctx/internal/err/fs"
	ctxerr "github.com/ActiveMemory/ctx/internal/err/hook"
	"github.com/ActiveMemory/ctx/internal/io"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/hooks/messages"
	"github.com/ActiveMemory/ctx/internal/cli/system/core"
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
		desc.TextDesc(text.DescKeyMessageListHeaderHook),
		desc.TextDesc(text.DescKeyMessageListHeaderVariant),
		desc.TextDesc(text.DescKeyMessageListHeaderCategory),
		desc.TextDesc(text.DescKeyMessageListHeaderOverride)))
	cmd.Println(fmt.Sprintf(headerFmt,
		strings.Repeat("\u2500", msg.MessageSepHook),
		strings.Repeat("\u2500", msg.MessageSepVariant),
		strings.Repeat("\u2500", msg.MessageSepCategory),
		strings.Repeat("\u2500", msg.MessageSepOverride)))

	for _, e := range entries {
		override := ""
		if e.HasOverride {
			override = desc.TextDesc(text.DescKeyMessageOverrideLabel)
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
func RunMessageShow(cmd *cobra.Command, hk, variant string) error {
	info := messages.Lookup(hk, variant)
	if info == nil {
		return core.ValidationError(hk, variant)
	}

	// Check user override first
	oPath := core.OverridePath(hk, variant)
	if data, readErr := io.SafeReadUserFile(oPath); readErr == nil {
		cmd.Println(fmt.Sprintf(desc.TextDesc(text.DescKeyMessageSourceOverride), oPath))
		core.PrintTemplateVars(cmd, info)
		cmd.Println()
		cmd.Print(string(data))
		if len(data) > 0 && data[len(data)-1] != '\n' {
			cmd.Println()
		}
		return nil
	}

	// Embedded default
	data, readErr := hook.Message(hk, variant+file.ExtTxt)
	if readErr != nil {
		return ctxerr.EmbeddedTemplateNotFound(hk, variant)
	}

	cmd.Println(desc.TextDesc(text.DescKeyMessageSourceDefault))
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
func RunMessageEdit(cmd *cobra.Command, hk, variant string) error {
	info := messages.Lookup(hk, variant)
	if info == nil {
		return core.ValidationError(hk, variant)
	}

	oPath := core.OverridePath(hk, variant)

	// Refuse if override already exists
	if _, statErr := os.Stat(oPath); statErr == nil {
		return ctxerr.OverrideExists(oPath, hk, variant)
	}

	// Warn for ctx-specific messages
	if info.Category == messages.CategoryCtxSpecific {
		cmd.Println(desc.TextDesc(text.DescKeyMessageCtxSpecificWarning))
		cmd.Println()
	}

	// Read embedded default
	data, readErr := hook.Message(hk, variant+file.ExtTxt)
	if readErr != nil {
		return ctxerr.EmbeddedTemplateNotFound(hk, variant)
	}

	// Create directories
	dir := filepath.Dir(oPath)
	if mkdirErr := os.MkdirAll(dir, 0o750); mkdirErr != nil {
		return fs.CreateDir(dir, mkdirErr)
	}

	// Write override file
	if writeErr := os.WriteFile(oPath, data, 0o600); writeErr != nil {
		return ctxerr.WriteOverride(oPath, writeErr)
	}

	cmd.Println(fmt.Sprintf(desc.TextDesc(text.DescKeyMessageOverrideCreated), oPath))
	cmd.Println(desc.TextDesc(text.DescKeyMessageEditHint))
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
			cmd.Println(fmt.Sprintf(desc.TextDesc(text.DescKeyMessageNoOverride), hook, variant))
			return nil
		}
		return ctxerr.RemoveOverride(oPath, removeErr)
	}

	// Clean up empty parent directories
	hookDir := filepath.Dir(oPath)
	_ = os.Remove(hookDir) // only succeeds if empty
	messagesDir := filepath.Dir(hookDir)
	_ = os.Remove(messagesDir) // only succeeds if empty

	cmd.Println(fmt.Sprintf(desc.TextDesc(text.DescKeyMessageOverrideRemoved), hook, variant))
	return nil
}
