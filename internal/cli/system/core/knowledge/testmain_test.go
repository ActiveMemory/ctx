//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package knowledge

import (
	"os"
	"testing"

	"github.com/ActiveMemory/ctx/internal/assets/read/lookup"
)

// TestMain initializes the embedded asset lookup so desc.Text resolves
// in FormatWarnings/Report — otherwise the finding format and remedy
// lines come back empty. (See the "uninitialized desc.Text" learning.)
func TestMain(m *testing.M) {
	lookup.Init()
	os.Exit(m.Run())
}
