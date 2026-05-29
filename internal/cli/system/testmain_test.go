//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package system_test

import (
	"os"
	"testing"

	"github.com/ActiveMemory/ctx/internal/assets/read/lookup"
)

// TestMain loads the embedded description maps so the system command
// tree renders real Short/Long text (desc lookups return empty until
// lookup.Init has run).
func TestMain(m *testing.M) {
	lookup.Init()
	os.Exit(m.Run())
}
