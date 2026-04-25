//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package status

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	coreStatus "github.com/ActiveMemory/ctx/internal/cli/hub/core/status"
	"github.com/ActiveMemory/ctx/internal/config/cli"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

// Cmd returns the hub status subcommand.
//
// Returns:
//   - *cobra.Command: The status subcommand
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeyHubStatus)

	return &cobra.Command{
		Use:     cmd.UseHubStatus,
		Short:   short,
		Long:    long,
		Example: desc.Example(cmd.DescKeyHubStatus),
		Args:    cobra.NoArgs,
		// Hub stores at ~/.ctx/hub-data/, not .context/.
		// Spec: specs/single-source-context-anchor.md.
		Annotations: map[string]string{cli.AnnotationSkipInit: cli.AnnotationTrue},
		RunE:        coreStatus.Run,
	}
}
