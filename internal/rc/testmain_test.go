//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package rc

import (
	"os"
	"testing"

	"github.com/ActiveMemory/ctx/internal/assets/read/lookup"
)

// TestMain initializes the embedded text-asset lookup so that error
// factories (internal/err/context.NotDeclared, etc.) resolve their
// DescKey-based messages instead of returning empty strings.
func TestMain(m *testing.M) {
	lookup.Init()
	os.Exit(m.Run())
}
