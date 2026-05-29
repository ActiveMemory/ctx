//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package main_test

import (
	"os"
	"testing"

	"github.com/ActiveMemory/ctx/internal/assets/read/lookup"
)

// TestMain initializes the embedded-asset lookup tables once
// for the whole ctxctl test binary. The relocated audit
// logic no longer reads ctx's i18n descriptors, but the
// shared ctx packages it reuses (nudge, state, provenance)
// still resolve their own ctx-owned text through lookup.
func TestMain(m *testing.M) {
	lookup.Init()
	os.Exit(m.Run())
}
