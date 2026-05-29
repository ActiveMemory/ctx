//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package unknown

import (
	"os"
	"testing"

	"github.com/ActiveMemory/ctx/internal/assets/read/lookup"
)

// TestMain loads the embedded description maps before running tests;
// desc.Text returns empty strings until lookup.Init has run.
func TestMain(m *testing.M) {
	lookup.Init()
	os.Exit(m.Run())
}
