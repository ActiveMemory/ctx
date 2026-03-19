//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package lookup

import (
	"strings"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config/token"
)

var (
	allowPerms []string
	denyPerms  []string
)

// loadPermissions reads an embedded permission file and splits it into entries.
func loadPermissions(path string) []string {
	data, readErr := assets.FS.ReadFile(path)
	if readErr != nil {
		return nil
	}
	var result []string
	for _, line := range strings.Split(string(data), token.NewlineLF) {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, token.PrefixHeading) {
			continue
		}
		result = append(result, line)
	}
	return result
}

// PermAllowListDefault returns the default allow permissions for ctx
// commands and skills, parsed from the embedded permissions/allow.txt.
func PermAllowListDefault() []string {
	return allowPerms
}

// PermDenyListDefault returns the default deny permissions that block
// dangerous operations, parsed from the embedded permissions/deny.txt.
func PermDenyListDefault() []string {
	return denyPerms
}
