//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package insert_test

import (
	"os"
	"testing"

	"github.com/ActiveMemory/ctx/internal/assets/read/lookup"
)

// TestMain initializes the embedded asset lookup so desc.Text resolves
// real values. Without it every desc.Text call returns "", and
// strings.Index(s, "") == 0 silently rewrites the insert anchors into
// "match at offset 0" — tests would exercise a path production never
// takes, and pass for the wrong reason.
func TestMain(m *testing.M) {
	lookup.Init()
	os.Exit(m.Run())
}
