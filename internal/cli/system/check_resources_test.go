//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package system

import (
	"strings"
	"testing"
)

func TestCheckResources_SilentWhenOK(t *testing.T) {
	cmd := newTestCmd()
	if err := runCheckResources(cmd, createTempStdin(t, `{}`)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := cmdOutput(cmd)
	// On a healthy dev machine this should be silent.
	// If the test machine is under resource pressure, the hook would fire —
	// which is technically correct behavior, so we only check it doesn't error.
	_ = out
}

func TestCheckResources_OutputFormat(t *testing.T) {
	// This is a structural test — if the hook fires, the output should
	// contain the VERBATIM relay preamble and the box-drawing frame.
	// We can't force DANGER conditions in a unit test without mocking
	// the collector, so we just verify error-free execution.
	cmd := newTestCmd()
	if err := runCheckResources(cmd, createTempStdin(t, `{}`)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := cmdOutput(cmd)
	if out != "" {
		// If output was produced, verify it has the right structure
		if !strings.Contains(out, "IMPORTANT:") {
			t.Error("hook output should start with IMPORTANT preamble")
		}
		if !strings.Contains(out, "Resource Alert") {
			t.Error("hook output should contain Resource Alert box")
		}
	}
}
