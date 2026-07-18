//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package disclosure_test

import (
	"testing"

	"github.com/ActiveMemory/ctx/internal/config/disclosure"
)

// T01: the structural vocabulary is fixed and load-bearing — the layout
// proofs and Validate key on these exact strings.
func TestHeadingConstants(t *testing.T) {
	if disclosure.HeadingThemes != "## Themes" {
		t.Errorf("HeadingThemes = %q, want %q",
			disclosure.HeadingThemes, "## Themes")
	}
	if disclosure.HeadingRecent != "## Recent" {
		t.Errorf("HeadingRecent = %q, want %q",
			disclosure.HeadingRecent, "## Recent")
	}
}
