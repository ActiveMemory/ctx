//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hook

import (
	"github.com/spf13/cobra"

	coreHook "github.com/ActiveMemory/ctx/internal/cli/trace/core/hook"
	cfgTrace "github.com/ActiveMemory/ctx/internal/config/trace"
	errTrace "github.com/ActiveMemory/ctx/internal/err/trace"
)

// Run executes the hook enable or disable action.
//
// Parameters:
//   - cmd: Cobra command for output stream
//   - action: "enable" or "disable"
//
// Returns:
//   - error: non-nil on unknown action or execution failure
func Run(cmd *cobra.Command, action string) error {
	switch action {
	case cfgTrace.ActionEnable:
		return coreHook.Enable(cmd)
	case cfgTrace.ActionDisable:
		return coreHook.Disable(cmd)
	default:
		return errTrace.UnknownAction(action)
	}
}
