//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package system

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/assets/hooks/messages"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/spf13/cobra"
)

// messageCmd returns the "ctx system message" subcommand.
func messageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "message",
		Short: "Manage hook message templates",
		Long: `Manage hook message templates.

Hook messages control what text hooks emit. The hook logic (when to
fire, counting, state tracking) is universal. The messages are opinions
that can be customized per-project.

Subcommands:
  list     Show all hook messages with category and override status
  show     Print the effective message template for a hook/variant
  edit     Copy the embedded default to .context/ for editing
  reset    Delete a user override and revert to embedded default`,
	}

	cmd.AddCommand(
		messageListCmd(),
		messageShowCmd(),
		messageEditCmd(),
		messageResetCmd(),
	)

	return cmd
}

// messageListCmd returns the "ctx system message list" subcommand.
func messageListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Show all hook messages with category and override status",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runMessageList(cmd)
		},
	}
	cmd.Flags().Bool("json", false, "Output in JSON format")
	return cmd
}

type messageListEntry struct {
	Hook         string   `json:"hook"`
	Variant      string   `json:"variant"`
	Category     string   `json:"category"`
	Description  string   `json:"description"`
	TemplateVars []string `json:"template_vars"`
	HasOverride  bool     `json:"has_override"`
}

func runMessageList(cmd *cobra.Command) error {
	registry := messages.Registry()
	entries := make([]messageListEntry, 0, len(registry))

	for _, info := range registry {
		entry := messageListEntry{
			Hook:         info.Hook,
			Variant:      info.Variant,
			Category:     info.Category,
			Description:  info.Description,
			TemplateVars: info.TemplateVars,
			HasOverride:  hasOverride(info.Hook, info.Variant),
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
	cmd.Printf("%-24s %-20s %-16s %s\n", "Hook", "Variant", "Category", "Override")
	cmd.Printf("%-24s %-20s %-16s %s\n",
		strings.Repeat("\u2500", 22),
		strings.Repeat("\u2500", 18),
		strings.Repeat("\u2500", 14),
		strings.Repeat("\u2500", 8))

	for _, e := range entries {
		override := ""
		if e.HasOverride {
			override = "override"
		}
		cmd.Printf("%-24s %-20s %-16s %s\n", e.Hook, e.Variant, e.Category, override)
	}

	return nil
}

// messageShowCmd returns the "ctx system message show" subcommand.
func messageShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show <hook> <variant>",
		Short: "Print the effective message template for a hook/variant",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMessageShow(cmd, args[0], args[1])
		},
	}
}

func runMessageShow(cmd *cobra.Command, hook, variant string) error {
	info := messages.Lookup(hook, variant)
	if info == nil {
		return validationError(hook, variant)
	}

	// Check user override first
	overridePath := overridePath(hook, variant)
	if data, readErr := os.ReadFile(overridePath); readErr == nil { //nolint:gosec // project-local override path
		cmd.Printf("Source: user override (%s)\n", overridePath)
		printTemplateVars(cmd, info)
		cmd.Println()
		cmd.Print(string(data))
		if len(data) > 0 && data[len(data)-1] != '\n' {
			cmd.Println()
		}
		return nil
	}

	// Embedded default
	data, readErr := assets.HookMessage(hook, variant+".txt")
	if readErr != nil {
		return fmt.Errorf("embedded template not found for %s/%s", hook, variant)
	}

	cmd.Println("Source: embedded default")
	printTemplateVars(cmd, info)
	cmd.Println()
	cmd.Print(string(data))
	if len(data) > 0 && data[len(data)-1] != '\n' {
		cmd.Println()
	}
	return nil
}

// messageEditCmd returns the "ctx system message edit" subcommand.
func messageEditCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "edit <hook> <variant>",
		Short: "Copy the embedded default to .context/ for editing",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMessageEdit(cmd, args[0], args[1])
		},
	}
}

func runMessageEdit(cmd *cobra.Command, hook, variant string) error {
	info := messages.Lookup(hook, variant)
	if info == nil {
		return validationError(hook, variant)
	}

	oPath := overridePath(hook, variant)

	// Refuse if override already exists
	if _, statErr := os.Stat(oPath); statErr == nil {
		return fmt.Errorf("override already exists at %s\nEdit it directly or use `ctx system message reset %s %s` first",
			oPath, hook, variant)
	}

	// Warn for ctx-specific messages
	if info.Category == messages.CategoryCtxSpecific {
		cmd.Println("Warning: this message is ctx-specific (intended for ctx development).")
		cmd.Println("Customizing it may produce unexpected results.")
		cmd.Println()
	}

	// Read embedded default
	data, readErr := assets.HookMessage(hook, variant+".txt")
	if readErr != nil {
		return fmt.Errorf("embedded template not found for %s/%s", hook, variant)
	}

	// Create directories
	dir := filepath.Dir(oPath)
	if mkdirErr := os.MkdirAll(dir, 0o750); mkdirErr != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, mkdirErr)
	}

	// Write override file
	if writeErr := os.WriteFile(oPath, data, 0o600); writeErr != nil {
		return fmt.Errorf("failed to write override %s: %w", oPath, writeErr)
	}

	cmd.Printf("Override created at %s\n", oPath)
	cmd.Println("Edit this file to customize the message.")
	printTemplateVars(cmd, info)

	return nil
}

// messageResetCmd returns the "ctx system message reset" subcommand.
func messageResetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "reset <hook> <variant>",
		Short: "Delete a user override and revert to embedded default",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMessageReset(cmd, args[0], args[1])
		},
	}
}

func runMessageReset(cmd *cobra.Command, hook, variant string) error {
	info := messages.Lookup(hook, variant)
	if info == nil {
		return validationError(hook, variant)
	}

	oPath := overridePath(hook, variant)

	if removeErr := os.Remove(oPath); removeErr != nil {
		if os.IsNotExist(removeErr) {
			cmd.Printf("No override found for %s/%s. Already using embedded default.\n", hook, variant)
			return nil
		}
		return fmt.Errorf("failed to remove override %s: %w", oPath, removeErr)
	}

	// Clean up empty parent directories
	hookDir := filepath.Dir(oPath)
	_ = os.Remove(hookDir) // only succeeds if empty
	messagesDir := filepath.Dir(hookDir)
	_ = os.Remove(messagesDir) // only succeeds if empty

	cmd.Printf("Override removed for %s/%s. Using embedded default.\n", hook, variant)
	return nil
}

// overridePath returns the user override file path for a hook/variant.
func overridePath(hook, variant string) string {
	return filepath.Join(rc.ContextDir(), "hooks", "messages", hook, variant+".txt")
}

// hasOverride checks whether a user override file exists.
func hasOverride(hook, variant string) bool {
	_, statErr := os.Stat(overridePath(hook, variant))
	return statErr == nil
}

// validationError returns an error for an unknown hook/variant.
func validationError(hook, variant string) error {
	// Check if the hook exists at all
	if messages.Variants(hook) == nil {
		return fmt.Errorf("unknown hook: %s\nRun `ctx system message list` to see available hooks", hook)
	}
	return fmt.Errorf("unknown variant %q for hook %q\nRun `ctx system message list` to see available variants", variant, hook)
}

// printTemplateVars prints available template variables if any exist.
func printTemplateVars(cmd *cobra.Command, info *messages.HookMessageInfo) {
	if len(info.TemplateVars) == 0 {
		cmd.Println("Template variables: (none)")
		return
	}
	formatted := make([]string, len(info.TemplateVars))
	for i, v := range info.TemplateVars {
		formatted[i] = "{{." + v + "}}"
	}
	cmd.Printf("Template variables: %s\n", strings.Join(formatted, ", "))
}
