//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/claude"
	"github.com/ActiveMemory/ctx/internal/config"
	ctxerr "github.com/ActiveMemory/ctx/internal/err"
	"github.com/ActiveMemory/ctx/internal/write"
)

// MergeSettingsPermissions merges ctx permissions into settings.local.json.
//
// Parameters:
//   - cmd: Cobra command for output
//
// Returns:
//   - error: Non-nil if file operations fail
func MergeSettingsPermissions(cmd *cobra.Command) error {
	var settings claude.Settings
	existingContent, err := os.ReadFile(config.FileSettings)
	fileExists := err == nil
	if fileExists {
		if err := json.Unmarshal(existingContent, &settings); err != nil {
			return ctxerr.ParseFile(config.FileSettings, err)
		}
	}
	allowModified := MergePermissions(&settings.Permissions.Allow, assets.DefaultAllowPermissions())
	denyModified := MergePermissions(&settings.Permissions.Deny, assets.DefaultDenyPermissions())
	allowDeduped := DeduplicatePermissions(&settings.Permissions.Allow)
	denyDeduped := DeduplicatePermissions(&settings.Permissions.Deny)
	if !allowModified && !denyModified && !allowDeduped && !denyDeduped {
		write.InitNoChanges(cmd, config.FileSettings)
		return nil
	}
	if err := os.MkdirAll(config.DirClaude, config.PermExec); err != nil {
		return ctxerr.Mkdir(config.DirClaude, err)
	}
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(settings); err != nil {
		return ctxerr.MarshalSettings(err)
	}
	if err := os.WriteFile(config.FileSettings, buf.Bytes(), config.PermFile); err != nil {
		return ctxerr.FileWrite(config.FileSettings, err)
	}
	if fileExists {
		deduped := allowDeduped || denyDeduped
		merged := allowModified || denyModified
		switch {
		case merged && deduped:
			write.InitPermsMergedDeduped(cmd, config.FileSettings)
		case deduped:
			write.InitPermsDeduped(cmd, config.FileSettings)
		case allowModified && denyModified:
			write.InitPermsAllowDeny(cmd, config.FileSettings)
		case denyModified:
			write.InitPermsDeny(cmd, config.FileSettings)
		default:
			write.InitPermsAllow(cmd, config.FileSettings)
		}
	} else {
		write.InitCreated(cmd, config.FileSettings)
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
