//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package deactivate_test

import (
	"os"
	"testing"

	"github.com/ActiveMemory/ctx/internal/assets/read/lookup"
)

// TestMain initializes the embedded text-asset lookup so deactivate's
// command metadata (Use/Short/Long from cmd/root) resolves correctly.
func TestMain(m *testing.M) {
	lookup.Init()
	os.Exit(m.Run())
}
