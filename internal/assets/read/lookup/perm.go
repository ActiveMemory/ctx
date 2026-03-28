//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package lookup

var (
	allowPerms []string
	denyPerms  []string
)

// PermAllowListDefault returns the default allow permissions for ctx
// commands and skills, parsed from the embedded permissions/allow.txt.
//
// Returns:
//   - []string: Allow permission patterns from the embedded allow list
func PermAllowListDefault() []string {
	return allowPerms
}

// PermDenyListDefault returns the default deny permissions that block
// dangerous operations, parsed from the embedded permissions/deny.txt.
//
// Returns:
//   - []string: Deny permission patterns from the embedded deny list
func PermDenyListDefault() []string {
	return denyPerms
}
