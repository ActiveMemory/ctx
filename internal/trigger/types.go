//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package trigger

import (
	cfgTrigger "github.com/ActiveMemory/ctx/internal/config/trigger"
	"github.com/ActiveMemory/ctx/internal/entity"
)

// HookType is an alias for the trigger event type.
type HookType = cfgTrigger.TriggerType

// ValidTypes returns all valid trigger type strings.
func ValidTypes() []HookType {
	return []HookType{
		cfgTrigger.PreToolUse,
		cfgTrigger.PostToolUse,
		cfgTrigger.SessionStart,
		cfgTrigger.SessionEnd,
		cfgTrigger.FileSave,
		cfgTrigger.ContextAdd,
	}
}

// HookSession is an alias for entity.TriggerSession.
type HookSession = entity.TriggerSession

// HookInput is an alias for entity.TriggerInput.
type HookInput = entity.TriggerInput

// HookOutput is the JSON object returned by trigger scripts via stdout.
//
// Fields:
//   - Cancel: If true, halt execution of subsequent triggers
//   - Context: Optional text to append to AI conversation context
//   - Message: Optional user-visible message
type HookOutput struct {
	Cancel  bool   `json:"cancel"`
	Context string `json:"context,omitempty"`
	Message string `json:"message,omitempty"`
}

// HookInfo describes a discovered trigger script.
//
// Fields:
//   - Name: Script filename without extension
//   - Type: Lifecycle event category
//   - Path: Filesystem path to the script
//   - Enabled: True if the executable permission bit is set
type HookInfo struct {
	Name    string
	Type    HookType
	Path    string
	Enabled bool
}

// AggregatedOutput collects results from all triggers in a run.
//
// Fields:
//   - Cancelled: True if a trigger returned cancel:true
//   - Message: Cancellation or summary message
//   - Context: Concatenated context from all triggers
//   - Errors: Warnings from failed triggers
type AggregatedOutput struct {
	Cancelled bool
	Message   string
	Context   string
	Errors    []string
}
