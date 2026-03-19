//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package assets

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Cannot import lookup.Init() here (cycle: assets → lookup → assets).
	// The assets test suite does not exercise desc lookups directly, so
	// no YAML loading is needed. Tests that require desc lookups live in
	// their own packages with proper TestMain calls.
	os.Exit(m.Run())
}
