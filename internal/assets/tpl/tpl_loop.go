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

	// LoopCmdClaudeMaxTurns is the shell command template for Claude Code
	// with a per-iteration turn cap in standard regime.
	// Args: maxTurns, promptFile.
	LoopCmdClaudeMaxTurns = `claude --print --max-turns %d "$(cat %s)"`

	// LoopCmdClaudeIterPrompt is the shell command template for Claude Code
	// in EoCC/SAC regimes using the composite $ITER_PROMPT variable.
	LoopCmdClaudeIterPrompt = `claude --print "$ITER_PROMPT"`

	// LoopCmdClaudeIterPromptMaxTurns is the shell command template for
	// Claude Code in EoCC/SAC regimes with --max-turns.
	// Args: maxTurns.
	LoopCmdClaudeIterPromptMaxTurns = `claude --print ` +
		`--max-turns %d "$ITER_PROMPT"`

	// LoopCmdAider is the shell command template for Aider.
	// Args: promptFile.
	LoopCmdAider = `aider --message-file %s`

	// LoopCmdAiderIterPrompt is the shell command template for Aider in
	// EoCC/SAC regimes using a temporary message file for $ITER_PROMPT.
	LoopCmdAiderIterPrompt = `ITER_PROMPT_FILE=$(mktemp) && ` +
		`printf '%s\n' "$ITER_PROMPT" > "$ITER_PROMPT_FILE" && ` +
		`aider --message-file "$ITER_PROMPT_FILE"; ` +
		`rm -f "$ITER_PROMPT_FILE"`

	// LoopCmdGeneric is the shell command placeholder for custom tools.
	// Args: promptFile.
	LoopCmdGeneric = `# Replace with your AI CLI command
    cat %s | your-ai-cli`

	// LoopCmdGenericIterPrompt is the shell command placeholder for custom
	// tools in EoCC/SAC regimes.
	LoopCmdGenericIterPrompt = `# Replace with your AI CLI command
    printf '%s\n' "$ITER_PROMPT" | your-ai-cli`
)
