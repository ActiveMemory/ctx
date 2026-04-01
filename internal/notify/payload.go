//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package notify

// Payload is the JSON body sent to the webhook endpoint.
//
// Fields:
//   - Event: Event type (loop, nudge, relay, heartbeat)
//   - Message: Rendered notification text
//   - Detail: Template reference for re-rendering
//   - SessionID: Claude Code session ID
//   - Timestamp: ISO 8601 send time
//   - Project: Project directory name
type Payload struct {
	Event     string       `json:"event"`
	Message   string       `json:"message"`
	Detail    *TemplateRef `json:"detail,omitempty"`
	SessionID string       `json:"session_id,omitempty"`
	Timestamp string       `json:"timestamp"`
	Project   string       `json:"project"`
}
