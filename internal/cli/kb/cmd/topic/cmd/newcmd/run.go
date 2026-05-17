//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package newcmd

import (
	"github.com/spf13/cobra"

	topicCore "github.com/ActiveMemory/ctx/internal/cli/kb/core/topic"
)

// Run scaffolds a new topic-page folder from the embedded
// template via [topicCore.Scaffold].
//
// Parameters:
//   - cobraCmd: cobra command for output.
//   - name: free-text topic name (slugified to kebab-case).
//
// Returns:
//   - error: scaffolding failure or refusal when topic exists.
func Run(cobraCmd *cobra.Command, name string) error {
	return topicCore.Scaffold(cobraCmd, name)
}
