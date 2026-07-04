//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package merge

import (
	"encoding/json"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/lookup"
	"github.com/ActiveMemory/ctx/internal/claude"
	cfgClaude "github.com/ActiveMemory/ctx/internal/config/claude"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/err/config"
	errParser "github.com/ActiveMemory/ctx/internal/err/parser"
	"github.com/ActiveMemory/ctx/internal/write/initialize"
)

// SettingsPermissions merges ctx permissions into settings.local.json.
//
// Only the permissions section is rewritten: all other top-level keys
// in the settings file (hooks, statusLine, env, anything ctx does not
// model) survive byte-for-byte via the raw-map read-modify-write in
// settings.go.
//
// Parameters:
//   - cmd: Cobra command for output
//
// Returns:
//   - error: Non-nil if file operations fail
func SettingsPermissions(cmd *cobra.Command) error {
	raw, fileExists, readErr := readSettingsRaw()
	if readErr != nil {
		return readErr
	}
	var perms claude.PermissionsConfig
	if rawPerms, exists := raw[cfgClaude.FieldPermissions]; exists {
		if unmarshalErr := json.Unmarshal(rawPerms, &perms); unmarshalErr != nil {
			return errParser.ParseFile(cfgClaude.Settings, unmarshalErr)
		}
	}
	allowModified := Permissions(
		&perms.Allow, lookup.PermAllowListDefault(),
	)
	denyModified := Permissions(
		&perms.Deny, lookup.PermDenyListDefault(),
	)
	allowDeduped := DeduplicatePermissions(&perms.Allow)
	denyDeduped := DeduplicatePermissions(&perms.Deny)
	if !allowModified && !denyModified && !allowDeduped && !denyDeduped {
		initialize.NoChanges(cmd, cfgClaude.Settings)
		return nil
	}
	section, marshalErr := marshalSettingsSection(perms)
	if marshalErr != nil {
		return config.MarshalSettings(marshalErr)
	}
	raw[cfgClaude.FieldPermissions] = section
	if writeErr := writeSettingsRaw(raw); writeErr != nil {
		return writeErr
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
		name, ok := skillName(p)
		if ok && strings.HasPrefix(name, cfgClaude.PluginScope) {
			bareName := strings.TrimPrefix(name, cfgClaude.PluginScope)
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
