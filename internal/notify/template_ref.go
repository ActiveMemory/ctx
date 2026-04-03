//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package notify

import "github.com/ActiveMemory/ctx/internal/entity"

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
//   - *entity.TemplateRef: Populated reference
func NewTemplateRef(
	hook, variant string, vars map[string]any,
) *entity.TemplateRef {
	return entity.NewTemplateRef(hook, variant, vars)
}
