//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hub_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	writeHub "github.com/ActiveMemory/ctx/internal/write/hub"
)

// clusterStatus renders ClusterStatus into a buffer.
func clusterStatus(
	info writeHub.ClusterStatusInfo,
) string {
	var buf bytes.Buffer
	cmd := &cobra.Command{}
	cmd.SetOut(&buf)
	writeHub.ClusterStatus(cmd, info)
	return buf.String()
}

// clustered returns the info a healthy three-node leader
// reports.
func clustered() writeHub.ClusterStatusInfo {
	return writeHub.ClusterStatusInfo{
		Role:      "Leader",
		Clustered: true,
		Leader:    "127.0.0.1:9901",
		Entries:   42,
		Peers:     2,
	}
}

// TestClusterStatus_Cluster pins the clustered shape: the role,
// the leader address the hub reported, and the peer count.
func TestClusterStatus_Cluster(t *testing.T) {
	out := clusterStatus(clustered())

	for _, want := range []string{
		"Role: Leader",
		"Leader: 127.0.0.1:9901",
		"Entries: 42  Peers: 2",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("want %q, got:\n%s", want, out)
		}
	}
}

// TestClusterStatus_Standalone pins the standalone shape. A hub
// with no Raft node has no leader and no peers; printing either
// would be inventing them, which is what the old renderer did.
func TestClusterStatus_Standalone(t *testing.T) {
	out := clusterStatus(writeHub.ClusterStatusInfo{
		Role:    "Standalone",
		Entries: 42,
	})

	if !strings.Contains(out, "Role: Standalone") {
		t.Errorf("want the standalone role, got:\n%s", out)
	}
	if !strings.Contains(out, "Entries: 42") {
		t.Errorf("want the entry count, got:\n%s", out)
	}
	if strings.Contains(out, "Leader") {
		t.Errorf("want no leader line, got:\n%s", out)
	}
	if strings.Contains(out, "Peers") {
		t.Errorf("want no peer count, got:\n%s", out)
	}
}

// TestClusterStatus_LeaderUnknown covers the clustered node with
// no leader: mid-election, or after quorum loss. The line has to
// say so rather than render an empty address.
func TestClusterStatus_LeaderUnknown(t *testing.T) {
	info := clustered()
	info.Role = "Follower"
	info.Leader = ""

	out := clusterStatus(info)

	if !strings.Contains(out, "Leader: unknown") {
		t.Errorf("want the unknown-leader line, got:\n%s", out)
	}
	if !strings.Contains(out, "Entries: 42  Peers: 2") {
		t.Errorf("want the stats line intact, got:\n%s", out)
	}
}

// TestClusterStatus_DroppedListeners pins the conditional
// slow-listener line. desc.Text returns "" for an unknown key, so a
// renamed text key would silently blank the line; asserting on the
// rendered count catches that.
func TestClusterStatus_DroppedListeners(t *testing.T) {
	info := clustered()
	info.Dropped = 3

	out := clusterStatus(info)

	if !strings.Contains(out, "Dropped listeners: 3") {
		t.Errorf("want dropped-listener line with count, got:\n%s", out)
	}
}

// TestClusterStatus_NoDroppedListeners pins the omission at zero so
// a healthy hub's output stays what it was.
func TestClusterStatus_NoDroppedListeners(t *testing.T) {
	out := clusterStatus(clustered())

	if strings.Contains(out, "Dropped listeners") {
		t.Errorf("want no dropped-listener line at zero, got:\n%s", out)
	}
	if !strings.Contains(out, "Entries: 42") {
		t.Errorf("want the existing stats line intact, got:\n%s", out)
	}
}
