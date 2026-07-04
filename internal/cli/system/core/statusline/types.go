//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package statusline

// Payload models the subset of Claude Code's status line stdin JSON
// that the ctx status line renders. Pointer fields distinguish absent
// or null values (segment dropped) from legitimate zeroes.
//
// Schema reference: https://code.claude.com/docs/en/statusline
//
// Fields:
//   - Model: model identity; DisplayName feeds the model segment
//   - Workspace: working directory (preferred over the legacy Cwd)
//   - Cwd: legacy duplicate of Workspace.CurrentDir
//   - Cost: session cost figures; TotalCostUSD feeds the $ segment
//   - ContextWindow: context usage; UsedPercentage feeds the ctx%
//     segment and may be null early in a session
type Payload struct {
	Model struct {
		DisplayName string `json:"display_name"`
	} `json:"model"`
	Workspace struct {
		CurrentDir string `json:"current_dir"`
	} `json:"workspace"`
	Cwd  string `json:"cwd"`
	Cost struct {
		TotalCostUSD *float64 `json:"total_cost_usd"`
	} `json:"cost"`
	ContextWindow struct {
		UsedPercentage *float64 `json:"used_percentage"`
	} `json:"context_window"`
}
