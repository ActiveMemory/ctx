//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package loop

// Loop script configuration.
const (
	// DefaultCompletionSignal is the default loop completion signal string.
	DefaultCompletionSignal = "SYSTEM_CONVERGED"
	// DefaultOutput is the default output filename for generated loop scripts.
	DefaultOutput = "loop.sh"
	// DefaultTool is the default AI tool for loop scripts.
	DefaultTool = "claude"
	// ToolAider is the aider tool identifier for loop scripts.
	ToolAider = "aider"
	// ToolGeneric is the generic tool identifier for loop scripts.
	ToolGeneric = "generic"

	// RegimeStandard is the default loop context regime (stateless prompt).
	RegimeStandard = "standard"
	// RegimeEoCC is the End-of-Conversation Compaction loop regime.
	RegimeEoCC = "eocc"
	// RegimeSAC is the Single-Action Compaction (SSSM cycle) loop regime.
	RegimeSAC = "sac"
	// DefaultRegime is the default loop context regime.
	DefaultRegime = RegimeStandard
	// DefaultSACMaxTurns is the default per-iteration turn cap for SAC on Claude.
	DefaultSACMaxTurns = 2
)

// ValidTools is the set of supported tool identifiers for loop scripts.
var ValidTools = map[string]bool{
	DefaultTool: true,
	ToolAider:   true,
	ToolGeneric: true,
}

// ValidRegimes is the set of supported context regimes for loop scripts.
var ValidRegimes = map[string]bool{
	RegimeStandard: true,
	RegimeEoCC:     true,
	RegimeSAC:      true,
}
