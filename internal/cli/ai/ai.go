//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package ai

import (
	"github.com/spf13/cobra"

	aiRoot "github.com/ActiveMemory/ctx/internal/cli/ai/cmd/root"
)

// Cmd returns the ctx ai command.
//
// Returns:
//   - *cobra.Command: configured ai command
func Cmd() *cobra.Command {
	return aiRoot.Cmd()
}
