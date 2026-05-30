//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package jsonpayload

import (
	"os"
	"testing"

	"github.com/ActiveMemory/ctx/internal/assets/read/lookup"
)

// TestMain initialises the embedded asset lookup so that the
// error helpers (errAdd.JSONParse, errFs.FileRead) render their
// parsed format strings rather than the empty default.
func TestMain(m *testing.M) {
	lookup.Init()
	os.Exit(m.Run())
}
