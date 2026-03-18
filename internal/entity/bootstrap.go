//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package entity

// BootstrapOutput is the JSON output structure for the bootstrap command.
//
// Fields:
//   - ContextDir: absolute path to the context directory
//   - Files: list of context file names
//   - Rules: list of rule strings
//   - NextSteps: list of next-step strings
//   - Warnings: optional warning strings
type BootstrapOutput struct {
	ContextDir string   `json:"context_dir"`
	Files      []string `json:"files"`
	Rules      []string `json:"rules"`
	NextSteps  []string `json:"next_steps"`
	Warnings   []string `json:"warnings,omitempty"`
}
