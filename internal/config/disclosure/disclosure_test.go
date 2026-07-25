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
//
// M4/T02 retired HeadingRecent ("## Recent") and ConventionLinePrefix
// ("### ") with the "### -under-## Recent" convention model they encoded.
func TestHeadingConstants(t *testing.T) {
	if disclosure.HeadingThemes != "## Themes" {
		t.Errorf("HeadingThemes = %q, want %q",
			disclosure.HeadingThemes, "## Themes")
	}
}

// T01/T02: the per-kind entry prefixes are the only structural difference
// between a timestamped root and a convention root, so their exact values
// are load-bearing for the single parse path.
func TestEntryPrefixConstants(t *testing.T) {
	if disclosure.EntryLinePrefix != "## [" {
		t.Errorf("EntryLinePrefix = %q, want %q",
			disclosure.EntryLinePrefix, "## [")
	}
	if disclosure.SectionLinePrefix != "## " {
		t.Errorf("SectionLinePrefix = %q, want %q",
			disclosure.SectionLinePrefix, "## ")
	}
}
