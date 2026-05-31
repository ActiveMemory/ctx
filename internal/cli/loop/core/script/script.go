//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package script

import (
	"fmt"
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
// detected in the output.
//
// Parameters:
//   - promptFile: Path to the prompt file (absolute path)
//   - tool: AI tool - "claude", "aider", or "generic"
//   - maxIterations: Max iterations (0 for unlimited)
//   - completionMsg: Signal string for completion
//
// Returns:
//   - string: Complete bash script content
//   - error: non-nil if the prompt path or template rendering fails
func Generate(
	promptFile, tool string,
	maxIterations int,
	completionMsg string,
) (string, error) {
	// Get the absolute path for the prompt file
	absPrompt, absErr := filepath.Abs(promptFile)
	if absErr != nil {
		return "", absErr
	}

	var aiCommand string
	switch tool {
	case cfgLoop.DefaultTool:
		aiCommand = fmt.Sprintf(tpl.LoopCmdClaude, absPrompt)
	case cfgLoop.ToolAider:
		aiCommand = fmt.Sprintf(tpl.LoopCmdAider, absPrompt)
	case cfgLoop.ToolGeneric:
		aiCommand = fmt.Sprintf(
			tpl.LoopCmdGeneric, absPrompt,
		)
	}

	return tpl.Render(tpl.LoopScript, tpl.LoopData{
		PromptFile:       absPrompt,
		CompletionSignal: completionMsg,
		MaxIter:          maxIterations,
		AICommand:        aiCommand,
		LoopComplete:     desc.Text(text.DescKeyLabelLoopComplete),
	})
}
