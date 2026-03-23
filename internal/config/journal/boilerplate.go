//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package journal

// Boilerplate detection patterns for tool output filtering.
// These are parser vocabulary (not i18n text) - they match
// Claude Code's fixed output strings.
const (
	BoilerplateNoMatch    = "No matches found"
	BoilerplateFilePrefix = "The file "
	BoilerplateFileSuffix = "has been updated successfully."
	BoilerplateDenied     = "denied this tool"
)

// PartPrefix is the Markdown bold prefix for multipart navigation labels.
const PartPrefix = "**Part "
