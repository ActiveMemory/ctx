//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package collapse

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// TestToolOutputsMatchesLegacy asserts the details-template rewrite of
// the long-output wrap path reproduces the legacy paired-tag output
// byte-for-byte. The fixture was captured from the legacy code path.
func TestToolOutputsMatchesLegacy(t *testing.T) {
	header := turnHeader(
		1, desc.Text(text.DescKeyLabelToolOutput), "10:00:00",
	)
	input := header + "\n\n" + bodyLines(12) + "\n"

	want, readErr := os.ReadFile(filepath.Join("testdata", "wrapped.golden"))
	if readErr != nil {
		t.Fatal(readErr)
	}
	got := ToolOutputs(input)
	if got != string(want) {
		t.Errorf(
			"drift:\n--- want ---\n%q\n--- got ---\n%q", string(want), got,
		)
	}
}
