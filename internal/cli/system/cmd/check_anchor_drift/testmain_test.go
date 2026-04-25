//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package check_anchor_drift

import (
	"os"
	"testing"

	"github.com/ActiveMemory/ctx/internal/assets/read/lookup"
)

// TestMain initializes the embedded text-asset lookup so the hook's
// box title and content templates resolve instead of returning
// empty strings.
func TestMain(m *testing.M) {
	lookup.Init()
	os.Exit(m.Run())
}
