//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package messages

// HookMessageInfo describes a single hook message template entry.
type HookMessageInfo struct {
	// Hook is the hook directory name (e.g., "qa-reminder").
	Hook string `yaml:"hook"`

	// Variant is the template file stem (e.g., "gate").
	Variant string `yaml:"variant"`

	// Category is "customizable" or "ctx-specific".
	Category string `yaml:"category"`

	// Description is a one-line human description of this message.
	Description string `yaml:"description"`

	// TemplateVars lists available Go template variables
	// (e.g., "PromptsSinceNudge").
	TemplateVars []string `yaml:"vars,omitempty"`
}
