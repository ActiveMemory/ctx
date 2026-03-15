//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package notify

// TemplateRef identifies the hook template and variables that produced a
// notification, allowing receivers to filter, re-render, or aggregate
// without parsing opaque rendered text.
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
