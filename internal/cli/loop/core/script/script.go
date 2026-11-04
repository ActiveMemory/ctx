//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package script

import (
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/assets/tpl"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	cfgLoop "github.com/ActiveMemory/ctx/internal/config/loop"
)

// Generate creates a bash script for running a Ralph loop.
//
// The generated script runs the specified AI tool repeatedly
// with the same prompt file until a completion signal is
// detected in the output. When opts.Regime is "eocc" or "sac",
// each iteration injects the assembled ctx agent state packet
// and execution-brief handover contract.
//
// Parameters:
//   - opts: Loop script generation options
//
// Returns:
//   - string: Complete bash script content
//   - error: non-nil if the prompt path or template rendering fails
func Generate(opts Opts) (string, error) {
	absPrompt, absErr := filepath.Abs(opts.PromptFile)
	if absErr != nil {
		return "", absErr
	}

	regime := opts.Regime
	if regime == "" {
		regime = cfgLoop.DefaultRegime
	}
	isEoCC := regime == cfgLoop.RegimeEoCC
	isSAC := regime == cfgLoop.RegimeSAC

	maxTurns := opts.MaxTurns
	if isSAC && maxTurns == 0 && opts.Tool == cfgLoop.DefaultTool {
		maxTurns = cfgLoop.DefaultSACMaxTurns
	}

	supportsMaxTurns := opts.Tool == cfgLoop.DefaultTool
	aiCommand := buildAICommand(
		opts.Tool, absPrompt, isEoCC || isSAC, maxTurns,
	)

	return tpl.Render(tpl.LoopScript, tpl.LoopData{
		PromptFile:       absPrompt,
		CompletionSignal: opts.CompletionMsg,
		MaxIter:          opts.MaxIterations,
		AICommand:        aiCommand,
		LoopComplete:     desc.Text(text.DescKeyLabelLoopComplete),
		Regime:           regime,
		IsEoCC:           isEoCC,
		IsSAC:            isSAC,
		MaxTurns:         maxTurns,
		SupportsMaxTurns: supportsMaxTurns,
	})
}
