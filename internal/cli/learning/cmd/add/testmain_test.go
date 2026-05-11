//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package add

import (
	"os"
	"testing"

	"github.com/ActiveMemory/ctx/internal/assets/read/lookup"
)

// TestMain initialises the embedded asset lookup before any
// test runs. Tests that exercise error paths through
// internal/err/cli depend on desc.Text returning the parsed
// format strings rather than the empty default.
func TestMain(m *testing.M) {
	lookup.Init()
	os.Exit(m.Run())
}
