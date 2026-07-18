//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package insert_test

import (
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/cli/add/core/insert"
	cfgEntry "github.com/ActiveMemory/ctx/internal/config/entry"
)

// Regression tests for specs/fix-afterheader-tail-truncation.md.
//
// AfterHeader returned content[:insertPoint] + entry — truncating the
// file at the insert point and discarding everything after it. Its sole
// caller is beforeFirstEntry's fallback, taken when a knowledge file has
// no "## [" entry, so `ctx learning add` against an entry-less file
// destroyed any content below the H1 header and its comment block.

const preamble = `# Learnings

<!--
UPDATE WHEN:
- Discover a gotcha worth recording
-->

`

// tailSection is content that lives below the preamble of a file which
// has no "## [" entries — exactly the shape that triggers the fallback.
const tailSection = `## Notes

- a hand-written section that must survive an add
`

const addedEntry = "## [2026-07-17-120000] A brand new learning\n\n" +
	"**Context**: freshly added.\n"

const addedEntryHeader = "## [2026-07-17-120000]"

// TestAfterHeader_PreservesTail is the bug: content below the preamble
// of an entry-less file must survive, with the new entry above it.
func TestAfterHeader_PreservesTail(t *testing.T) {
	out := string(insert.AppendEntry(
		[]byte(preamble+tailSection), addedEntry, cfgEntry.Learning, "",
	))

	if !strings.Contains(out, "## Notes") {
		t.Fatalf("tail section was destroyed by add; got:\n%s", out)
	}
	if !strings.Contains(out, "a hand-written section that must survive an add") {
		t.Fatalf("tail body was destroyed by add; got:\n%s", out)
	}

	entryIdx := strings.Index(out, addedEntryHeader)
	tailIdx := strings.Index(out, "## Notes")
	if entryIdx == -1 {
		t.Fatalf("added entry missing; got:\n%s", out)
	}
	if entryIdx > tailIdx {
		t.Errorf("entry landed below the tail (entry=%d tail=%d); "+
			"it must land directly after the preamble", entryIdx, tailIdx)
	}
}

// TestAfterHeader_EmptyTailUnchanged pins the shape that works today:
// an entry-less file with nothing after the comment block must produce
// byte-identical output to the pre-fix implementation (content + entry).
func TestAfterHeader_EmptyTailUnchanged(t *testing.T) {
	out := string(insert.AppendEntry(
		[]byte(preamble), addedEntry, cfgEntry.Learning, "",
	))

	want := preamble + addedEntry
	if out != want {
		t.Errorf("empty-tail output changed.\n got: %q\nwant: %q", out, want)
	}
}

// TestBeforeFirstEntry_PrimaryPathUnchanged guards the common path: a
// file that already has entries still inserts before the first one and
// keeps every existing entry.
func TestBeforeFirstEntry_PrimaryPathUnchanged(t *testing.T) {
	existing := preamble +
		"## [2026-07-15-141726] an existing learning\n\n**Context**: old.\n"

	out := string(insert.AppendEntry(
		[]byte(existing), addedEntry, cfgEntry.Learning, "",
	))

	if !strings.Contains(out, "an existing learning") {
		t.Fatalf("existing entry destroyed; got:\n%s", out)
	}

	newIdx := strings.Index(out, addedEntryHeader)
	oldIdx := strings.Index(out, "## [2026-07-15-141726]")
	if newIdx == -1 || oldIdx == -1 {
		t.Fatalf("an entry is missing; got:\n%s", out)
	}
	if newIdx > oldIdx {
		t.Errorf("new entry must precede the existing one "+
			"(new=%d old=%d)", newIdx, oldIdx)
	}
}
