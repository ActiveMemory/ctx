//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package session

// Event type constants for session lifecycle events.
const (
	// EventStart marks the beginning of a workspace session.
	EventStart = "start"
	// EventEnd marks the end of a workspace session.
	EventEnd = "end"
)

// Session and template constants.
const (
	// IDUnknown is the fallback session ID when input lacks one.
	IDUnknown = "unknown"
	// IDSuffixSummary is appended to session IDs for summary messages.
	IDSuffixSummary = "-summary"
	// IDSuffixTopic is appended to session IDs for topic messages.
	IDSuffixTopic = "-topic"
	// PreviewMaxLen is the maximum character length for first-message previews.
	PreviewMaxLen = 100
	// TemplateName is the name used for Go text/template instances.
	TemplateName = "msg"
)
