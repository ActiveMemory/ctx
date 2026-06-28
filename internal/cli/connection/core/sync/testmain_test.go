//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sync

import (
	"os"
	"testing"

	"github.com/ActiveMemory/ctx/internal/assets/read/lookup"
)

// TestMain initialises the embedded asset lookup so the
// lock-contention error (errHub.ConnectSyncLocked) renders its
// parsed format string rather than the empty default.
func TestMain(m *testing.M) {
	lookup.Init()
	os.Exit(m.Run())
}
