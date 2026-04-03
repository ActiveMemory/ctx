//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package budget

// packet represents the JSON output format for the agent command.
type packet struct {
	Generated    string   `json:"generated"`
	Budget       int      `json:"budget"`
	TokensUsed   int      `json:"tokens_used"`
	ReadOrder    []string `json:"read_order"`
	Constitution []string `json:"constitution"`
	Tasks        []string `json:"tasks"`
	Conventions  []string `json:"conventions"`
	Decisions    []string `json:"decisions"`
	Learnings    []string `json:"learnings,omitempty"`
	Summaries    []string `json:"summaries,omitempty"`
	Steering     []string `json:"steering,omitempty"`
	Skill        string   `json:"skill,omitempty"`
	Instruction  string   `json:"instruction"`
}

// AssembledPacket holds the budget-aware output sections ready for rendering.
//
// Fields:
//   - ReadOrder: Files in priority order
//   - Constitution: Inviolable rules
//   - Tasks: Active task items
//   - Conventions: Code conventions
//   - Decisions: Architectural decisions (scored)
//   - Learnings: Gotchas and tips (scored)
//   - Summaries: Title-only overflow entries
//   - Steering: Applicable steering file bodies
//   - Skill: Named skill content (from --skill flag)
//   - Instruction: Behavioral instruction text
//   - Budget: Token budget limit
//   - TokensUsed: Estimated tokens consumed
type AssembledPacket struct {
	ReadOrder    []string
	Constitution []string
	Tasks        []string
	Conventions  []string
	Decisions    []string
	Learnings    []string
	Summaries    []string
	Steering     []string
	Skill        string
	Instruction  string
	Budget       int
	TokensUsed   int
}
