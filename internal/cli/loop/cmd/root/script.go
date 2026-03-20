//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"fmt"
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/assets/tpl"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// GenerateLoopScript creates a bash script for running a Ralph loop.
//
// The generated script runs the specified AI tool repeatedly with the
// same prompt file until a completion signal is detected in the output.
//
// Parameters:
//   - promptFile: Path to the prompt file (converted to absolute path)
//   - tool: AI tool to use - "claude", "aider", or "generic"
//   - maxIterations: Maximum iterations before stopping (0 for unlimited)
//   - completionMsg: String to detect in output that signals completion
//
// Returns:
//   - string: Complete bash script content ready to write to file
func GenerateLoopScript(
	promptFile, tool string, maxIterations int, completionMsg string,
) string {
	// Get the absolute path for the prompt file
	absPrompt, _ := filepath.Abs(promptFile)

	var aiCommand string
	switch tool {
	case "claude":
		aiCommand = fmt.Sprintf(`claude --print "$(cat %s)"`, absPrompt)
	case "aider":
		aiCommand = fmt.Sprintf(`aider --message-file %s`, absPrompt)
	case "generic":
		aiCommand = fmt.Sprintf(`# Replace with your AI CLI command
    cat %s | your-ai-cli`, absPrompt)
	}

	maxIterCheck := ""
	if maxIterations > 0 {
		maxIterCheck = fmt.Sprintf(
			tpl.LoopMaxIter, maxIterations, maxIterations, tpl.LoopNotify)
	}

	script := fmt.Sprintf(tpl.LoopScript,
		absPrompt, completionMsg, maxIterCheck, aiCommand, desc.TextDesc(text.DescKeyLabelLoopComplete), tpl.LoopNotify)

	return script
}
