//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package stepdown

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	coreStep "github.com/ActiveMemory/ctx/internal/cli/hub/core/stepdown"
	"github.com/ActiveMemory/ctx/internal/config/cli"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

// Cmd returns the hub stepdown subcommand.
//
// Returns:
//   - *cobra.Command: The stepdown subcommand
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeyHubStepdown)

	return &cobra.Command{
		Use:     cmd.UseHubStepdown,
		Short:   short,
		Long:    long,
		Example: desc.Example(cmd.DescKeyHubStepdown),
		Args:    cobra.NoArgs,
		// Hub stores at ~/.ctx/hub-data/, not .context/.
		// Spec: specs/single-source-context-anchor.md.
		Annotations: map[string]string{cli.AnnotationSkipInit: cli.AnnotationTrue},
		RunE:        coreStep.Run,
	}
}
