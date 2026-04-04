//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package text

// DescKeys for MCP prompt descriptions.
const (
	// DescKeyMCPPromptSessionStartDesc is the text key for mcp prompt session
	// start desc messages.
	DescKeyMCPPromptSessionStartDesc = "mcp.prompt-session-start-desc"
	// DescKeyMCPPromptAddDecisionDesc is the text key for mcp prompt add decision
	// desc messages.
	DescKeyMCPPromptAddDecisionDesc = "mcp.prompt-add-decision-desc"
	// DescKeyMCPPromptAddLearningDesc is the text key for mcp prompt add learning
	// desc messages.
	DescKeyMCPPromptAddLearningDesc = "mcp.prompt-add-learning-desc"
	// DescKeyMCPPromptReflectDesc is the text key for mcp prompt reflect desc
	// messages.
	DescKeyMCPPromptReflectDesc = "mcp.prompt-reflect-desc"
	// DescKeyMCPPromptCheckpointDesc is the text key for mcp prompt checkpoint
	// desc messages.
	DescKeyMCPPromptCheckpointDesc = "mcp.prompt-checkpoint-desc"
)

// DescKeys for MCP prompt arguments.
const (
	// DescKeyMCPPromptArgDecisionTitle is the text key for mcp prompt arg
	// decision title messages.
	DescKeyMCPPromptArgDecisionTitle = "mcp.prompt-arg-decision-title"
	// DescKeyMCPPromptArgDecisionCtx is the text key for mcp prompt arg decision
	// ctx messages.
	DescKeyMCPPromptArgDecisionCtx = "mcp.prompt-arg-decision-ctx"
	// DescKeyMCPPromptArgDecisionRat is the text key for mcp prompt arg decision
	// rat messages.
	DescKeyMCPPromptArgDecisionRat = "mcp.prompt-arg-decision-rationale"
	// DescKeyMCPPromptArgDecisionConseq is the text key for mcp prompt arg
	// decision conseq messages.
	DescKeyMCPPromptArgDecisionConseq = "mcp.prompt-arg-decision-consequence"
	// DescKeyMCPPromptArgLearningTitle is the text key for mcp prompt arg
	// learning title messages.
	DescKeyMCPPromptArgLearningTitle = "mcp.prompt-arg-learning-title"
	// DescKeyMCPPromptArgLearningCtx is the text key for mcp prompt arg learning
	// ctx messages.
	DescKeyMCPPromptArgLearningCtx = "mcp.prompt-arg-learning-ctx"
	// DescKeyMCPPromptArgLearningLesson is the text key for mcp prompt arg
	// learning lesson messages.
	DescKeyMCPPromptArgLearningLesson = "mcp.prompt-arg-learning-lesson"
	// DescKeyMCPPromptArgLearningApp is the text key for mcp prompt arg learning
	// app messages.
	DescKeyMCPPromptArgLearningApp = "mcp.prompt-arg-learning-app"
)

// DescKeys for MCP session-start prompt layout.
const (
	// DescKeyMCPPromptSessionStartHeader is the text key for mcp prompt session
	// start header messages.
	DescKeyMCPPromptSessionStartHeader = "mcp.prompt-session-start-header"
	// DescKeyMCPPromptSessionStartFooter is the text key for mcp prompt session
	// start footer messages.
	DescKeyMCPPromptSessionStartFooter = "mcp.prompt-session-start-footer"
	// DescKeyMCPPromptSessionStartResultD is the text key for mcp prompt session
	// start result d messages.
	DescKeyMCPPromptSessionStartResultD = "mcp.prompt-session-start-result-desc"
	// DescKeyMCPPromptSectionFormat is the text key for mcp prompt section format
	// messages.
	DescKeyMCPPromptSectionFormat = "mcp.prompt-section-format"
)

// DescKeys for MCP add-decision prompt.
const (
	// DescKeyMCPPromptAddDecisionHeader is the text key for mcp prompt add
	// decision header messages.
	DescKeyMCPPromptAddDecisionHeader = "mcp.prompt-add-decision-header"
	// DescKeyMCPPromptAddDecisionFieldFmt is the text key for mcp prompt add
	// decision field fmt messages.
	DescKeyMCPPromptAddDecisionFieldFmt = "mcp.prompt-add-decision-field-format"
)

// DescKeys for MCP prompt field labels.
const (
	// DescKeyMCPPromptLabelDecision is the text key for mcp prompt label decision
	// messages.
	DescKeyMCPPromptLabelDecision = "mcp.prompt-label-decision"
	// DescKeyMCPPromptLabelContext is the text key for mcp prompt label context
	// messages.
	DescKeyMCPPromptLabelContext = "mcp.prompt-label-context"
	// DescKeyMCPPromptLabelRationale is the text key for mcp prompt label
	// rationale messages.
	DescKeyMCPPromptLabelRationale = "mcp.prompt-label-rationale"
	// DescKeyMCPPromptLabelConsequence is the text key for mcp prompt label
	// consequence messages.
	DescKeyMCPPromptLabelConsequence = "mcp.prompt-label-consequence"
	// DescKeyMCPPromptLabelLearning is the text key for mcp prompt label learning
	// messages.
	DescKeyMCPPromptLabelLearning = "mcp.prompt-label-learning"
	// DescKeyMCPPromptLabelLesson is the text key for mcp prompt label lesson
	// messages.
	DescKeyMCPPromptLabelLesson = "mcp.prompt-label-lesson"
	// DescKeyMCPPromptLabelApplication is the text key for mcp prompt label
	// application messages.
	DescKeyMCPPromptLabelApplication = "mcp.prompt-label-application"
)

// DescKeys for MCP add-decision prompt result.
const (
	// DescKeyMCPPromptAddDecisionFooter is the text key for mcp prompt add
	// decision footer messages.
	DescKeyMCPPromptAddDecisionFooter = "mcp.prompt-add-decision-footer"
	// DescKeyMCPPromptAddDecisionResultD is the text key for mcp prompt add
	// decision result d messages.
	DescKeyMCPPromptAddDecisionResultD = "mcp.prompt-add-decision-result-desc"
)

// DescKeys for MCP add-learning prompt.
const (
	// DescKeyMCPPromptAddLearningHeader is the text key for mcp prompt add
	// learning header messages.
	DescKeyMCPPromptAddLearningHeader = "mcp.prompt-add-learning-header"
	// DescKeyMCPPromptAddLearningFieldFmt is the text key for mcp prompt add
	// learning field fmt messages.
	DescKeyMCPPromptAddLearningFieldFmt = "mcp.prompt-add-learning-field-format"
	// DescKeyMCPPromptAddLearningFooter is the text key for mcp prompt add
	// learning footer messages.
	DescKeyMCPPromptAddLearningFooter = "mcp.prompt-add-learning-footer"
	// DescKeyMCPPromptAddLearningResultD is the text key for mcp prompt add
	// learning result d messages.
	DescKeyMCPPromptAddLearningResultD = "mcp.prompt-add-learning-result-desc"
)

// DescKeys for MCP reflect prompt.
const (
	// DescKeyMCPPromptReflectBody is the text key for mcp prompt reflect body
	// messages.
	DescKeyMCPPromptReflectBody = "mcp.prompt-reflect-body"
	// DescKeyMCPPromptReflectResultD is the text key for mcp prompt reflect
	// result d messages.
	DescKeyMCPPromptReflectResultD = "mcp.prompt-reflect-result-desc"
)

// DescKeys for MCP checkpoint prompt.
const (
	// DescKeyMCPPromptCheckpointHeader is the text key for mcp prompt checkpoint
	// header messages.
	DescKeyMCPPromptCheckpointHeader = "mcp.prompt-checkpoint-header"
	// DescKeyMCPPromptCheckpointStatsFormat is the text key for mcp prompt
	// checkpoint stats format messages.
	DescKeyMCPPromptCheckpointStatsFormat = "mcp.prompt-checkpoint-stats-format"
	// DescKeyMCPPromptCheckpointSteps is the text key for mcp prompt checkpoint
	// steps messages.
	DescKeyMCPPromptCheckpointSteps = "mcp.prompt-checkpoint-steps"
	// DescKeyMCPPromptCheckpointResultD is the text key for mcp prompt checkpoint
	// result d messages.
	DescKeyMCPPromptCheckpointResultD = "mcp.prompt-checkpoint-result-desc"
)
