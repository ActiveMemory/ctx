//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package stats

// Check result status constants used by doctor, drift, and other health checks.
const (
	StatusOK      = "ok"
	StatusWarning = "warning"
	StatusError   = "error"
	StatusInfo    = "info"
)

// Unicode icons for status display.
const (
	IconOK      = "✓"
	IconWarning = "⚠"
	IconError   = "✗"
	IconInfo    = "○"
	IconUnknown = "?"
)

// StatusIcon returns the Unicode icon for a given status string.
//
// Parameters:
//   - status: One of StatusOK, StatusWarning, StatusError, or StatusInfo
//
// Returns:
//   - string: A single Unicode character representing the status
func StatusIcon(status string) string {
	switch status {
	case StatusOK:
		return IconOK
	case StatusWarning:
		return IconWarning
	case StatusError:
		return IconError
	case StatusInfo:
		return IconInfo
	default:
		return IconUnknown
	}
}
