//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hub

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// clusterLines prints the leader and the combined entry/peer
// line for a hub running with a Raft node attached.
//
// Parameters:
//   - cmd: Cobra command for output
//   - info: status fields as reported by the Status RPC
func clusterLines(
	cmd *cobra.Command, info ClusterStatusInfo,
) {
	if info.Leader == "" {
		cmd.Println(desc.Text(
			text.DescKeyWriteHubLeaderUnknown,
		))
	} else {
		cmd.Println(fmt.Sprintf(
			desc.Text(text.DescKeyWriteHubLeader),
			info.Leader,
		))
	}

	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteHubClusterStats),
		info.Entries, info.Peers,
	))
}
