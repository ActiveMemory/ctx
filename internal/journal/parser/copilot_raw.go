//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package parser

import "encoding/json"

// Copilot Chat JSONL raw types.
//
// Copilot Chat stores sessions as JSONL files in VS Code's workspaceStorage.
// Each file contains one session. The first line (kind=0) is the full session
// snapshot, subsequent lines are incremental patches (kind=1 for scalar
// replacements, kind=2 for array/object replacements).

// copilotRawLine represents a single JSONL line from a Copilot Chat session.
//
// Kind discriminates the line type:
//   - 0: Full session snapshot (V contains copilotRawSession)
//   - 1: Scalar property patch (K is the JSON path, V is the new value)
//   - 2: Array/object patch (K is the JSON path, V is the new value)
type copilotRawLine struct {
	Kind int               `json:"kind"`
	K    []json.RawMessage `json:"k,omitempty"`
	V    json.RawMessage   `json:"v"`
}

// copilotRawSession is the full session snapshot from a copilotKindSnapshot line.
type copilotRawSession struct {
	// Version is the schema version of the session format.
	Version int `json:"version"`
	// CreationDate is the session creation time as Unix milliseconds.
	CreationDate int64 `json:"creationDate"`
	// CustomTitle is the user-set session title, if any.
	CustomTitle string `json:"customTitle,omitempty"`
	// SessionID is the unique identifier for this session.
	SessionID string `json:"sessionId"`
	// ResponderUsername is the Copilot responder identity.
	ResponderUsername string `json:"responderUsername,omitempty"`
	// InitialLocation is the VS Code view that started the session.
	InitialLocation string `json:"initialLocation,omitempty"`
	// Requests holds the ordered request-response pairs.
	Requests []copilotRawRequest `json:"requests"`
}

// copilotRawRequest represents a single request-response pair.
type copilotRawRequest struct {
	RequestID         string               `json:"requestId"`
	Timestamp         int64                `json:"timestamp"`
	ModelID           string               `json:"modelId,omitempty"`
	Message           copilotRawMessage    `json:"message"`
	Response          []copilotRawRespItem `json:"response,omitempty"`
	Result            *copilotRawResult    `json:"result,omitempty"`
	ContentReferences []json.RawMessage    `json:"contentReferences,omitempty"`
}

// copilotRawMessage is the user's input message.
type copilotRawMessage struct {
	Text string `json:"text"`
}

// copilotRawRespItem is a single item in the response array.
//
// The Kind field discriminates the type:
//   - "thinking": Extended thinking (Value contains the text)
//   - "toolInvocationSerialized": Tool call
//   - "textEditGroup": File edit
//   - "": Plain markdown text (Value field only)
type copilotRawRespItem struct {
	Kind              string          `json:"kind,omitempty"`
	Value             json.RawMessage `json:"value,omitempty"`
	ID                string          `json:"id,omitempty"`
	InvocationMessage json.RawMessage `json:"invocationMessage,omitempty"`
	ToolID            string          `json:"toolId,omitempty"`
	ToolCallID        string          `json:"toolCallId,omitempty"`
	IsComplete        json.RawMessage `json:"isComplete,omitempty"`
}

// copilotRawResult contains completion metadata for a request.
type copilotRawResult struct {
	Timings  copilotRawTimings  `json:"timings"`
	Metadata copilotRawMetadata `json:"metadata,omitempty"`
}

// copilotRawTimings contains timing information.
type copilotRawTimings struct {
	FirstProgress int64 `json:"firstProgress"`
	TotalElapsed  int64 `json:"totalElapsed"`
}

// copilotRawMetadata contains token usage and other metadata.
type copilotRawMetadata struct {
	PromptTokens int `json:"promptTokens,omitempty"`
	OutputTokens int `json:"outputTokens,omitempty"`
}

// copilotRawWorkspace is the workspace.json file in workspaceStorage.
type copilotRawWorkspace struct {
	Folder string `json:"folder,omitempty"`
}
