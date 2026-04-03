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
func Generate(
	promptFile, tool string,
	maxIterations int,
	completionMsg string,
) string {
	// Get the absolute path for the prompt file
	absPrompt, _ := filepath.Abs(promptFile)

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

	maxIterCheck := ""
	if maxIterations > 0 {
		maxIterCheck = fmt.Sprintf(
			tpl.LoopMaxIter,
			maxIterations, maxIterations, tpl.LoopNotify,
		)
	}

	script := fmt.Sprintf(tpl.LoopScript,
		absPrompt, completionMsg, maxIterCheck, aiCommand,
		desc.Text(text.DescKeyLabelLoopComplete),
		tpl.LoopNotify,
	)

	return script
}
