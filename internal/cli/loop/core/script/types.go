//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package script

// Opts holds the configuration for generating a Ralph loop script.
type Opts struct {
	// PromptFile is the path to the prompt file.
	PromptFile string
	// Tool is the AI CLI tool identifier ("claude", "aider", or "generic").
	Tool string
	// Regime is the context regime ("standard", "eocc", or "sac").
	Regime string
	// MaxTurns is the per-iteration turn cap (0 for tool default).
	MaxTurns int
	// MaxIterations is the maximum loop iterations (0 for unlimited).
	MaxIterations int
	// CompletionMsg is the signal string that terminates the loop.
	CompletionMsg string
}
