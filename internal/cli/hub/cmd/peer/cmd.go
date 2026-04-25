//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package peer

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	corePeer "github.com/ActiveMemory/ctx/internal/cli/hub/core/peer"
	"github.com/ActiveMemory/ctx/internal/config/cli"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

// Cmd returns the hub peer subcommand.
//
// Returns:
//   - *cobra.Command: The peer subcommand
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeyHubPeer)

	return &cobra.Command{
		Use:     cmd.UseHubPeer,
		Short:   short,
		Long:    long,
		Example: desc.Example(cmd.DescKeyHubPeer),
		Args:    cobra.ExactArgs(2),
		// Hub stores at ~/.ctx/hub-data/, not .context/.
		// Spec: specs/single-source-context-anchor.md.
		Annotations: map[string]string{cli.AnnotationSkipInit: cli.AnnotationTrue},
		RunE:        corePeer.Run,
	}
}
