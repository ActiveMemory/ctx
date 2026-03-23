//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package merge

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"

	"github.com/ActiveMemory/ctx/internal/assets/read/lookup"
	cfgClaude "github.com/ActiveMemory/ctx/internal/config/claude"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/err/config"
	errFs "github.com/ActiveMemory/ctx/internal/err/fs"
	errParser "github.com/ActiveMemory/ctx/internal/err/parser"
	"github.com/ActiveMemory/ctx/internal/write/initialize"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/claude"
)

// SettingsPermissions merges ctx permissions into settings.local.json.
//
// Parameters:
//   - cmd: Cobra command for output
//
// Returns:
//   - error: Non-nil if file operations fail
func SettingsPermissions(cmd *cobra.Command) error {
	var settings claude.Settings
	existingContent, err := os.ReadFile(cfgClaude.Settings)
	fileExists := err == nil
	if fileExists {
		if err := json.Unmarshal(existingContent, &settings); err != nil {
			return errParser.ParseFile(cfgClaude.Settings, err)
		}
	}
	allowModified := Permissions(
		&settings.Permissions.Allow, lookup.PermAllowListDefault(),
	)
	denyModified := Permissions(
		&settings.Permissions.Deny, lookup.PermDenyListDefault(),
	)
	allowDeduped := DeduplicatePermissions(&settings.Permissions.Allow)
	denyDeduped := DeduplicatePermissions(&settings.Permissions.Deny)
	if !allowModified && !denyModified && !allowDeduped && !denyDeduped {
		initialize.NoChanges(cmd, cfgClaude.Settings)
		return nil
	}
	if err := os.MkdirAll(dir.Claude, fs.PermExec); err != nil {
		return errFs.Mkdir(dir.Claude, err)
	}
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(settings); err != nil {
		return config.MarshalSettings(err)
	}
	if err := os.WriteFile(
		cfgClaude.Settings, buf.Bytes(), fs.PermFile,
	); err != nil {
		return errFs.FileWrite(cfgClaude.Settings, err)
	}
	if fileExists {
		deduped := allowDeduped || denyDeduped
		merged := allowModified || denyModified
		switch {
		case merged && deduped:
			initialize.PermsMergedDeduped(cmd, cfgClaude.Settings)
		case deduped:
			initialize.PermsDeduped(cmd, cfgClaude.Settings)
		case allowModified && denyModified:
			initialize.PermsAllowDeny(cmd, cfgClaude.Settings)
		case denyModified:
			initialize.PermsDeny(cmd, cfgClaude.Settings)
		default:
			initialize.PermsAllow(cmd, cfgClaude.Settings)
		}
	} else {
		initialize.Created(cmd, cfgClaude.Settings)
	}
	return nil
}

// Permissions adds default permissions that are not already present.
//
// Parameters:
//   - slice: Existing permissions slice to modify
//   - defaults: Default permissions to merge in
//
// Returns:
//   - bool: True if any permissions were added
func Permissions(slice *[]string, defaults []string) bool {
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
const PluginPrefix = cfgClaude.PluginScope

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
		if name, ok := skillName(p); ok {
			if !strings.Contains(name, token.Colon) {
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
		if name, ok := skillName(p); ok && strings.HasPrefix(name, PluginPrefix) {
			bareName := strings.TrimPrefix(name, PluginPrefix)
			bareName = strings.TrimSuffix(bareName, cfgClaude.PluginScopeWildcard)
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
