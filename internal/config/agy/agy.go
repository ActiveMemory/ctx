//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package agy

// Antigravity CLI binary, directories, and plugin keys.
const (
	// Binary is the Antigravity CLI binary name.
	Binary = "agy"
	// DirConfig is the config directory name under Home (~/.gemini/config).
	DirConfig = "config"
	// DirPlugins is the plugins directory name under DirConfig.
	DirPlugins = "plugins"
	// PluginName is the ctx plugin name in the Antigravity manifest.
	PluginName = "ctx"
	// KeyCtx is the top-level hook group key in hooks.json.
	KeyCtx = "ctx"
	// KeyDisabled is the disabled flag key in plugin.json.
	KeyDisabled = "disabled"
	// FileConfigJSON is the Antigravity global config filename under DirConfig.
	FileConfigJSON = "config.json"
	// KeyPlugins is the plugins map key in config.json.
	KeyPlugins = "plugins"
	// KeyEnabled is the enabled flag key in config.json plugin entries.
	KeyEnabled = "enabled"
)

// Antigravity tool names reported in hook payloads.
const (
	// ToolRunCommand is the Antigravity tool name for shell commands.
	ToolRunCommand = "run_command"
	// ToolEditFile is the Antigravity tool name for editing files.
	ToolEditFile = "edit_file"
	// ToolWriteFile is the Antigravity tool name for writing files.
	ToolWriteFile = "write_file"
	// ToolReadFile is the Antigravity tool name for reading files.
	ToolReadFile = "read_file"
)

// Environment variable names for Antigravity integration.
const (
	// EnvHookTool is the environment variable that identifies the
	// invoking hook tool runner.
	EnvHookTool = "CTX_HOOK_TOOL"
	// EnvHome is the environment variable that overrides the
	// Antigravity home directory (~/.gemini).
	EnvHome = "AGY_HOME"
	// EnvConfigDir is the environment variable that overrides the
	// Antigravity config directory (~/.gemini/config).
	EnvConfigDir = "AGY_CONFIG_DIR"
)

// Antigravity hook decision values.
const (
	// DecisionAllow allows the tool execution to proceed.
	DecisionAllow = "allow"
	// DecisionDeny blocks the tool execution with a reason.
	DecisionDeny = "deny"
)

// Antigravity lifecycle hook event names.
const (
	// EventSessionStart fires at session start.
	EventSessionStart = "SessionStart"
	// EventPreInvocation fires before agent invocation.
	EventPreInvocation = "PreInvocation"
	// EventPreToolUse fires before tool execution.
	EventPreToolUse = "PreToolUse"
	// EventPostToolUse fires after tool execution.
	EventPostToolUse = "PostToolUse"
	// EventPostInvocation fires after agent invocation.
	EventPostInvocation = "PostInvocation"
	// EventStop fires when the session stops.
	EventStop = "Stop"
)

// Antigravity pending queue and file tokens.
const (
	// PrefixPendingQueue is the file prefix for pending nudges
	// queued during PreToolUse/PostToolUse.
	PrefixPendingQueue = "agy-pending-"
	// DirGemini is the Gemini home directory name under $HOME.
	DirGemini = ".gemini"
	// DirAntigravityCLI is the Antigravity CLI directory under
	// [DirGemini].
	DirAntigravityCLI = "antigravity-cli"
	// DirBrain is the brain directory under [DirAntigravityCLI].
	DirBrain = "brain"
	// DirSystemGenerated is the system-generated directory under
	// a brain session.
	DirSystemGenerated = ".system_generated"
	// DirLogs is the logs directory under [DirSystemGenerated].
	DirLogs = "logs"
	// FileTranscriptJSONL is the session transcript filename.
	FileTranscriptJSONL = "transcript.jsonl"
	// EmptyJSON is the empty JSON object representation.
	EmptyJSON = "{}"
	// PendingDelimiter is the separator between queued pending
	// messages in the state file.
	PendingDelimiter = "\n---\n"
)

// Antigravity session transcript files, step types, sources, and tags.
const (
	// FileTranscriptFullJSONL is the full session transcript filename.
	FileTranscriptFullJSONL = "transcript_full.jsonl"
	// FileHistoryJSONL is the workspace history filename.
	FileHistoryJSONL = "history.jsonl"
	// StepTypeUserInput identifies user input steps in transcripts.
	StepTypeUserInput = "USER_INPUT"
	// StepTypePlannerResponse identifies assistant planner response steps.
	StepTypePlannerResponse = "PLANNER_RESPONSE"
	// StepTypeGeneric identifies generic steps in transcripts.
	StepTypeGeneric = "GENERIC"
	// SourceUserExplicit identifies explicit user-originated steps.
	SourceUserExplicit = "USER_EXPLICIT"
	// SourceModel identifies model-generated steps.
	SourceModel = "MODEL"
	// SourceSystem identifies system-generated steps.
	SourceSystem = "SYSTEM"
	// TagUserRequestOpen is the opening tag for user request blocks.
	TagUserRequestOpen = "<USER_REQUEST>"
	// TagUserRequestClose is the closing tag for user request blocks.
	TagUserRequestClose = "</USER_REQUEST>"
	// TagAdditionalMetadataOpen is the opening tag for additional metadata.
	TagAdditionalMetadataOpen = "<ADDITIONAL_METADATA>"
	// TagAdditionalMetadataClose is the closing tag for additional metadata.
	TagAdditionalMetadataClose = "</ADDITIONAL_METADATA>"
	// TagSystemMessageOpen is the opening tag for internal system messages.
	TagSystemMessageOpen = "<SYSTEM_MESSAGE>"
)
