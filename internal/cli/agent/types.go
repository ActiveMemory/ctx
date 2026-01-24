//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package agent

// Packet represents the JSON output format for the agent command.
//
// This struct is serialized when using --format json and contains
// all extracted context information for AI consumption.
//
// Fields:
//   - Generated: RFC3339 timestamp of when the packet was created
//   - Budget: Token budget specified by the user
//   - TokensUsed: Actual token count of loaded context files
//   - ReadOrder: File paths in recommended reading order
//   - Constitution: Rules from CONSTITUTION.md
//   - Tasks: Active (unchecked) tasks from TASKS.md
//   - Conventions: Key conventions from CONVENTIONS.md
//   - Decisions: Recent decision titles from DECISIONS.md
type Packet struct {
	Generated    string   `json:"generated"`
	Budget       int      `json:"budget"`
	TokensUsed   int      `json:"tokens_used"`
	ReadOrder    []string `json:"read_order"`
	Constitution []string `json:"constitution"`
	Tasks        []string `json:"tasks"`
	Conventions  []string `json:"conventions"`
	Decisions    []string `json:"decisions"`
}
