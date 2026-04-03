//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package text

// DescKeys for MCP prompt descriptions.
const (
	DescKeyMCPPromptSessionStartDesc = "mcp.prompt-session-start-desc"
	DescKeyMCPPromptAddDecisionDesc  = "mcp.prompt-add-decision-desc"
	DescKeyMCPPromptAddLearningDesc  = "mcp.prompt-add-learning-desc"
	DescKeyMCPPromptReflectDesc      = "mcp.prompt-reflect-desc"
	DescKeyMCPPromptCheckpointDesc   = "mcp.prompt-checkpoint-desc"
)

// DescKeys for MCP prompt arguments.
const (
	DescKeyMCPPromptArgDecisionTitle  = "mcp.prompt-arg-decision-title"
	DescKeyMCPPromptArgDecisionCtx    = "mcp.prompt-arg-decision-ctx"
	DescKeyMCPPromptArgDecisionRat    = "mcp.prompt-arg-decision-rationale"
	DescKeyMCPPromptArgDecisionConseq = "mcp.prompt-arg-decision-consequence"
	DescKeyMCPPromptArgLearningTitle  = "mcp.prompt-arg-learning-title"
	DescKeyMCPPromptArgLearningCtx    = "mcp.prompt-arg-learning-ctx"
	DescKeyMCPPromptArgLearningLesson = "mcp.prompt-arg-learning-lesson"
	DescKeyMCPPromptArgLearningApp    = "mcp.prompt-arg-learning-app"
)

// DescKeys for MCP session-start prompt layout.
const (
	DescKeyMCPPromptSessionStartHeader  = "mcp.prompt-session-start-header"
	DescKeyMCPPromptSessionStartFooter  = "mcp.prompt-session-start-footer"
	DescKeyMCPPromptSessionStartResultD = "mcp.prompt-session-start-result-desc"
	DescKeyMCPPromptSectionFormat       = "mcp.prompt-section-format"
)

// DescKeys for MCP add-decision prompt.
const (
	DescKeyMCPPromptAddDecisionHeader   = "mcp.prompt-add-decision-header"
	DescKeyMCPPromptAddDecisionFieldFmt = "mcp.prompt-add-decision-field-format"
)

// DescKeys for MCP prompt field labels.
const (
	DescKeyMCPPromptLabelDecision    = "mcp.prompt-label-decision"
	DescKeyMCPPromptLabelContext     = "mcp.prompt-label-context"
	DescKeyMCPPromptLabelRationale   = "mcp.prompt-label-rationale"
	DescKeyMCPPromptLabelConsequence = "mcp.prompt-label-consequence"
	DescKeyMCPPromptLabelLearning    = "mcp.prompt-label-learning"
	DescKeyMCPPromptLabelLesson      = "mcp.prompt-label-lesson"
	DescKeyMCPPromptLabelApplication = "mcp.prompt-label-application"
)

// DescKeys for MCP add-decision prompt result.
const (
	DescKeyMCPPromptAddDecisionFooter  = "mcp.prompt-add-decision-footer"
	DescKeyMCPPromptAddDecisionResultD = "mcp.prompt-add-decision-result-desc"
)

// DescKeys for MCP add-learning prompt.
const (
	DescKeyMCPPromptAddLearningHeader   = "mcp.prompt-add-learning-header"
	DescKeyMCPPromptAddLearningFieldFmt = "mcp.prompt-add-learning-field-format"
	DescKeyMCPPromptAddLearningFooter   = "mcp.prompt-add-learning-footer"
	DescKeyMCPPromptAddLearningResultD  = "mcp.prompt-add-learning-result-desc"
)

// DescKeys for MCP reflect prompt.
const (
	DescKeyMCPPromptReflectBody    = "mcp.prompt-reflect-body"
	DescKeyMCPPromptReflectResultD = "mcp.prompt-reflect-result-desc"
)

// DescKeys for MCP checkpoint prompt.
const (
	DescKeyMCPPromptCheckpointHeader      = "mcp.prompt-checkpoint-header"
	DescKeyMCPPromptCheckpointStatsFormat = "mcp.prompt-checkpoint-stats-format"
	DescKeyMCPPromptCheckpointSteps       = "mcp.prompt-checkpoint-steps"
	DescKeyMCPPromptCheckpointResultD     = "mcp.prompt-checkpoint-result-desc"
)
