//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package initialize

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/claude"
	"github.com/ActiveMemory/ctx/internal/config"
)

// mergeSettingsPermissions merges ctx permissions into settings.local.json.
//
// Only adds missing permissions to preserve user customizations. Does not
// manage hooks — hook configuration is now provided by the ctx Claude Code
// plugin.
//
// Parameters:
//   - cmd: Cobra command for output messages
//
// Returns:
//   - error: Non-nil if JSON parsing or file operations fail
func mergeSettingsPermissions(cmd *cobra.Command) error {
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	// Check if settings.local.json exists
	var settings claude.Settings
	existingContent, err := os.ReadFile(config.FileSettings)
	fileExists := err == nil

	if fileExists {
		if err := json.Unmarshal(existingContent, &settings); err != nil {
			return fmt.Errorf(
				"failed to parse existing %s: %w", config.FileSettings, err,
			)
		}
	}

	// Merge permissions - always additive, never removes existing permissions
	allowModified := mergePermissions(&settings.Permissions.Allow, assets.DefaultAllowPermissions())
	denyModified := mergePermissions(&settings.Permissions.Deny, assets.DefaultDenyPermissions())

	// Deduplicate after merge: remove exact dupes and FQ skill forms
	// subsumed by bare equivalents already in the list.
	allowDeduped := deduplicatePermissions(&settings.Permissions.Allow)
	denyDeduped := deduplicatePermissions(&settings.Permissions.Deny)

	if !allowModified && !denyModified && !allowDeduped && !denyDeduped {
		cmd.Println(fmt.Sprintf(
			"  %s %s (no changes needed)\n", yellow("○"), config.FileSettings,
		))
		return nil
	}

	// Create .claude/ directory if needed
	if err := os.MkdirAll(config.DirClaude, config.PermExec); err != nil {
		return fmt.Errorf("failed to create %s: %w", config.DirClaude, err)
	}

	// Write settings with pretty formatting
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(settings); err != nil {
		return fmt.Errorf("failed to marshal settings: %w", err)
	}

	if err := os.WriteFile(config.FileSettings, buf.Bytes(), config.PermFile); err != nil {
		return fmt.Errorf("failed to write %s: %w", config.FileSettings, err)
	}

	if fileExists {
		deduped := allowDeduped || denyDeduped
		merged := allowModified || denyModified
		switch {
		case merged && deduped:
			cmd.Println(fmt.Sprintf("  %s %s (added ctx permissions, removed duplicates)", green("✓"), config.FileSettings))
		case deduped:
			cmd.Println(fmt.Sprintf("  %s %s (removed duplicate permissions)", green("✓"), config.FileSettings))
		case allowModified && denyModified:
			cmd.Println(fmt.Sprintf("  %s %s (added ctx allow + deny permissions)", green("✓"), config.FileSettings))
		case denyModified:
			cmd.Println(fmt.Sprintf("  %s %s (added ctx deny permissions)", green("✓"), config.FileSettings))
		default:
			cmd.Println(fmt.Sprintf("  %s %s (added ctx permissions)", green("✓"), config.FileSettings))
		}
	} else {
		cmd.Println(fmt.Sprintf("  %s %s", green("✓"), config.FileSettings))
	}

	return nil
}

// mergePermissions adds missing entries to a string slice.
//
// Only adds entries that don't already exist. Never removes existing
// entries to preserve user customizations. Works on both allow and deny lists.
//
// Parameters:
//   - slice: Pointer to existing string slice to modify
//   - defaults: Default entries to add if missing
//
// Returns:
//   - bool: True if any entries were added
func mergePermissions(slice *[]string, defaults []string) bool {
	// Build a set of existing entries for fast lookup
	existing := make(map[string]bool)
	for _, p := range *slice {
		existing[p] = true
	}

	// Add missing entries
	added := false
	for _, p := range defaults {
		if !existing[p] {
			*slice = append(*slice, p)
			added = true
		}
	}

	return added
}

// pluginPrefix is the ctx plugin name used in fully-qualified skill forms.
const pluginPrefix = "ctx:"

// deduplicatePermissions removes redundant entries from a permission slice.
//
// Two kinds of redundancy are handled:
//  1. Exact duplicates — only the first occurrence is kept.
//  2. Fully-qualified skill forms subsumed by a bare equivalent already
//     in the list. When "Skill(foo)" exists, "Skill(ctx:foo)" and
//     "Skill(ctx:foo:*)" are redundant and removed. Only the ctx:
//     prefix is stripped (our plugin name), not arbitrary prefixes.
//
// The function preserves insertion order (stable dedup).
//
// Parameters:
//   - slice: Pointer to the string slice to deduplicate in place
//
// Returns:
//   - bool: True if any entries were removed
func deduplicatePermissions(slice *[]string) bool {
	if len(*slice) == 0 {
		return false
	}

	// First pass: collect bare Skill forms for subsumption checks.
	bareSkills := make(map[string]bool)
	for _, p := range *slice {
		if name, ok := skillName(p); ok {
			if !strings.Contains(name, ":") {
				bareSkills[name] = true
			}
		}
	}

	// Second pass: keep entries that are neither exact dupes nor
	// subsumed FQ forms.
	seen := make(map[string]bool)
	result := make([]string, 0, len(*slice))

	for _, p := range *slice {
		// Exact duplicate check.
		if seen[p] {
			continue
		}
		seen[p] = true

		// FQ skill subsumption check.
		if name, ok := skillName(p); ok && strings.HasPrefix(name, pluginPrefix) {
			bareName := strings.TrimPrefix(name, pluginPrefix)
			// Strip trailing :* variant as well.
			bareName = strings.TrimSuffix(bareName, ":*")
			if bareSkills[bareName] {
				continue
			}
		}

		result = append(result, p)
	}

	removed := len(*slice) != len(result)
	*slice = result
	return removed
}

// skillName extracts the inner name from a "Skill(name)" permission string.
// Returns the name and true if the string matches the Skill(...) pattern.
func skillName(perm string) (string, bool) {
	if !strings.HasPrefix(perm, "Skill(") || !strings.HasSuffix(perm, ")") {
		return "", false
	}
	return perm[len("Skill(") : len(perm)-1], true
}
