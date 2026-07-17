//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package disclosure_test

import (
	"os"
	"testing"

	"github.com/ActiveMemory/ctx/internal/assets/read/lookup"
)

// TestMain initializes the embedded asset lookup so sentinel .Error()
// text resolves — otherwise a failing assertion prints an empty error.
// (See the "uninitialized desc.Text" learning.)
func TestMain(m *testing.M) {
	lookup.Init()
	os.Exit(m.Run())
}
