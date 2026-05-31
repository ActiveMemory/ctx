//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package tpl

// Agent load and loop script templates.
const (
	// LoadBudget formats the token budget summary line.
	// Args: budget, totalTokens.
	LoadBudget = "Token Budget: %d | Available: %d"

	// LoadTruncated formats the truncation notice for budget-limited output.
	// Args: fileName.
	LoadTruncated = "*[Truncated: %s and remaining" +
		" files excluded due to token budget]*"

	// LoadSectionHeading formats a file section heading in assembled output.
	// Args: title.
	LoadSectionHeading = "## %s"

	// LoopCmdClaude is the shell command template for Claude Code.
	// Args: promptFile.
	LoopCmdClaude = `claude --print "$(cat %s)"`

	// LoopCmdAider is the shell command template for Aider.
	// Args: promptFile.
	LoopCmdAider = `aider --message-file %s`

	// LoopCmdGeneric is the shell command placeholder for custom tools.
	// Args: promptFile.
	LoopCmdGeneric = `# Replace with your AI CLI command
    cat %s | your-ai-cli`
)
