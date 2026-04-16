//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package setup provides terminal output for the setup
// command (ctx setup) and hook lifecycle output.
//
// This package covers two distinct output surfaces:
// tool deployment results and hook runtime messages.
//
// # Tool Deployment
//
// Setup deploys integration files for various AI
// tools. Copilot output uses [InfoCopilotCreated],
// [InfoCopilotMerged], [InfoCopilotSkipped], and
// [InfoCopilotSummary]. Copilot-CLI uses
// [InfoCopilotCLICreated], [InfoCopilotCLISkipped],
// and [InfoCopilotCLISummary]. AGENTS.md uses
// [InfoAgentsCreated], [InfoAgentsMerged],
// [InfoAgentsSkipped], and [InfoAgentsSummary].
//
// Generic deploy functions handle file creation
// across tools: [DeployComplete], [DeployFileExists],
// [DeployFileCreated], [DeploySteeringSynced],
// [DeploySteeringSkipped], [DeployNoSteering].
//
// Editor-specific integration instructions are
// printed by [InfoCursorIntegration],
// [InfoKiroIntegration], and [InfoClineIntegration].
//
// # Hook Runtime
//
// [Nudge] emits a single-line relay directive.
// [NudgeBlock] emits a multi-line boxed message
// followed by a blank line. Both are consumed by
// the agent as VERBATIM relay directives.
//
// [Context] and [BlockResponse] print JSON hook
// response lines. [Content] prints raw hook
// content. [Separator] prints blank-line dividers.
//
// [InfoTool] prints a pre-formatted tool section.
// [InfoUnknownTool] reports an unrecognized tool.
package setup
