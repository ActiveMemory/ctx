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
func clusterStatus(dropped uint64) string {
	var buf bytes.Buffer
	cmd := &cobra.Command{}
	cmd.SetOut(&buf)
	writeHub.ClusterStatus(
		cmd, "leader", "127.0.0.1:9901", 42, 2, dropped,
	)
	return buf.String()
}

// TestClusterStatus_DroppedListeners pins the conditional
// slow-listener line. desc.Text returns "" for an unknown key, so a
// renamed text key would silently blank the line; asserting on the
// rendered count catches that.
func TestClusterStatus_DroppedListeners(t *testing.T) {
	out := clusterStatus(3)

	if !strings.Contains(out, "Dropped listeners: 3") {
		t.Errorf("want dropped-listener line with count, got:\n%s", out)
	}
}

// TestClusterStatus_NoDroppedListeners pins the omission at zero so
// a healthy hub's output stays what it was.
func TestClusterStatus_NoDroppedListeners(t *testing.T) {
	out := clusterStatus(0)

	if strings.Contains(out, "Dropped listeners") {
		t.Errorf("want no dropped-listener line at zero, got:\n%s", out)
	}
	if !strings.Contains(out, "Entries: 42") {
		t.Errorf("want the existing stats line intact, got:\n%s", out)
	}
}
