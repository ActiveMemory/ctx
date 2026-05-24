//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hub

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/log/warn"
)

// TestReplicateOnce_WarnsOnDialError verifies that
// replicateOnce emits a warning via [warn.Warn] when the
// gRPC NewClient call fails on a malformed target. The "%"
// literal trips grpc.NewClient's URL parser ("invalid URL
// escape") before any network I/O - the most reliable way
// to hit the dial-error branch without standing up a
// listener. Regression guard for ActiveMemory/ctx#100:
// pre-fix, dial failures returned silently and the
// follower would loop forever with no operator-visible
// signal that the master address was unreachable or
// malformed.
func TestReplicateOnce_WarnsOnDialError(t *testing.T) {
	store, storeErr := NewStore(t.TempDir())
	if storeErr != nil {
		t.Fatal(storeErr)
	}

	var buf bytes.Buffer
	restore := warn.SetSinkForTesting(&buf)
	defer restore()

	replicateOnce(context.Background(), "%", store, "token")

	got := buf.String()
	if !strings.Contains(got, "replication: dial %:") {
		t.Errorf(
			"warning output missing replication dial prefix; got %q",
			got,
		)
	}
}
