//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

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

// MergeSettingsPermissions merges ctx permissions into settings.local.json.
//
// Parameters:
//   - cmd: Cobra command for output
//
// Returns:
//   - error: Non-nil if file operations fail
func MergeSettingsPermissions(cmd *cobra.Command) error {
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	var settings claude.Settings
	existingContent, err := os.ReadFile(config.FileSettings)
	fileExists := err == nil
	if fileExists {
		if err := json.Unmarshal(existingContent, &settings); err != nil {
			return fmt.Errorf("failed to parse existing %s: %w", config.FileSettings, err)
		}
	}
	allowModified := MergePermissions(&settings.Permissions.Allow, assets.DefaultAllowPermissions())
	denyModified := MergePermissions(&settings.Permissions.Deny, assets.DefaultDenyPermissions())
	allowDeduped := DeduplicatePermissions(&settings.Permissions.Allow)
	denyDeduped := DeduplicatePermissions(&settings.Permissions.Deny)
	if !allowModified && !denyModified && !allowDeduped && !denyDeduped {
		cmd.Println(fmt.Sprintf("  %s %s (no changes needed)\n", yellow("○"), config.FileSettings))
		return nil
	}
	if err := os.MkdirAll(config.DirClaude, config.PermExec); err != nil {
		return fmt.Errorf("failed to create %s: %w", config.DirClaude, err)
	}
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

// MergePermissions adds default permissions that are not already present.
//
// Parameters:
//   - slice: Existing permissions slice to modify
//   - defaults: Default permissions to merge in
//
// Returns:
//   - bool: True if any permissions were added
func MergePermissions(slice *[]string, defaults []string) bool {
	existing := make(map[string]bool)
	for _, p := range *slice {
		existing[p] = true
	}
	added := false
	for _, p := range defaults {
		if !existing[p] {
			*slice = append(*slice, p)
			added = true
		}
	}
	return added
}

// PluginPrefix is the prefix for plugin-scoped skill permissions.
const PluginPrefix = "ctx:"

// DeduplicatePermissions removes duplicate and redundant FQ-form permissions.
//
// Parameters:
//   - slice: Permissions slice to deduplicate
//
// Returns:
//   - bool: True if any duplicates were removed
func DeduplicatePermissions(slice *[]string) bool {
	if len(*slice) == 0 {
		return false
	}
	bareSkills := make(map[string]bool)
	for _, p := range *slice {
		if name, ok := SkillName(p); ok {
			if !strings.Contains(name, ":") {
				bareSkills[name] = true
			}
		}
	}
	seen := make(map[string]bool)
	result := make([]string, 0, len(*slice))
	for _, p := range *slice {
		if seen[p] {
			continue
		}
		seen[p] = true
		if name, ok := SkillName(p); ok && strings.HasPrefix(name, PluginPrefix) {
			bareName := strings.TrimPrefix(name, PluginPrefix)
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

// SkillName extracts the skill name from a permission string like "Skill(name)".
//
// Parameters:
//   - perm: Permission string to parse
//
// Returns:
//   - string: The skill name
//   - bool: True if perm matches the Skill(...) format
func SkillName(perm string) (string, bool) {
	if !strings.HasPrefix(perm, "Skill(") || !strings.HasSuffix(perm, ")") {
		return "", false
	}
	return perm[len("Skill(") : len(perm)-1], true
}
