//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package script

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets/tpl"
	cfgLoop "github.com/ActiveMemory/ctx/internal/config/loop"
)

// buildAICommand constructs the shell command invocation for the
// selected AI tool, context regime, and per-iteration turn cap.
//
// Parameters:
//   - tool: AI tool identifier ("claude", "aider", or "generic")
//   - absPrompt: Absolute path to the prompt file
//   - stateful: True when regime is EoCC or SAC ($ITER_PROMPT)
//   - maxTurns: Per-iteration turn cap (0 for tool default)
//
// Returns:
//   - string: Formatted shell command string
func buildAICommand(
	tool, absPrompt string,
	stateful bool,
	maxTurns int,
) string {
	switch tool {
	case cfgLoop.DefaultTool:
		if stateful {
			if maxTurns > 0 {
				return fmt.Sprintf(
					tpl.LoopCmdClaudeIterPromptMaxTurns, maxTurns,
				)
			}
			return tpl.LoopCmdClaudeIterPrompt
		}
		if maxTurns > 0 {
			return fmt.Sprintf(
				tpl.LoopCmdClaudeMaxTurns, maxTurns, absPrompt,
			)
		}
		return fmt.Sprintf(tpl.LoopCmdClaude, absPrompt)
	case cfgLoop.ToolAider:
		if stateful {
			return tpl.LoopCmdAiderIterPrompt
		}
		return fmt.Sprintf(tpl.LoopCmdAider, absPrompt)
	case cfgLoop.ToolGeneric:
		if stateful {
			return tpl.LoopCmdGenericIterPrompt
		}
		return fmt.Sprintf(tpl.LoopCmdGeneric, absPrompt)
	default:
		return ""
	}
}
