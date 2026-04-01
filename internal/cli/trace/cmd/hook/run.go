//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hook

import (
	"fmt"

	"github.com/spf13/cobra"

	coreHook "github.com/ActiveMemory/ctx/internal/cli/trace/core/hook"
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
	case "enable":
		return coreHook.Enable(cmd)
	case "disable":
		return coreHook.Disable(cmd)
	default:
		return fmt.Errorf("unknown action %q: use enable or disable", action)
	}
}
