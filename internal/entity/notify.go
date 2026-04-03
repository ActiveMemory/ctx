//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package entity

// TemplateRef identifies the hook template and variables that produced a
// notification, allowing receivers to filter, re-render, or aggregate
// without parsing opaque rendered text.
//
// Fields:
//   - Hook: Hook name that produced this notification
//   - Variant: Template variant within the hook
//   - Variables: Template variables used for rendering
type TemplateRef struct {
	Hook      string         `json:"hook"`
	Variant   string         `json:"variant"`
	Variables map[string]any `json:"variables,omitempty"`
}

// NewTemplateRef constructs a TemplateRef.
//
// Nil variables are omitted from JSON.
//
// Parameters:
//   - hook: Hook name that triggered the notification
//   - variant: Template variant within the hook
//   - vars: Template variables; nil is omitted from JSON
//
// Returns:
//   - *TemplateRef: Populated reference
func NewTemplateRef(hook, variant string, vars map[string]any) *TemplateRef {
	return &TemplateRef{Hook: hook, Variant: variant, Variables: vars}
}

// NotifyPayload is the JSON body sent to the webhook endpoint.
//
// Fields:
//   - Event: Event type (loop, nudge, relay, heartbeat)
//   - Message: Rendered notification text
//   - Detail: Template reference for re-rendering
//   - SessionID: Claude Code session ID
//   - Timestamp: ISO 8601 send time
//   - Project: Project directory name
type NotifyPayload struct {
	Event     string       `json:"event"`
	Message   string       `json:"message"`
	Detail    *TemplateRef `json:"detail,omitempty"`
	SessionID string       `json:"session_id,omitempty"`
	Timestamp string       `json:"timestamp"`
	Project   string       `json:"project"`
}
