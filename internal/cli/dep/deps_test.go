//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package dep

import (
	"bytes"
	"strings"
	"testing"
)

func TestCmd_UnknownFormat(t *testing.T) {
	cmd := Cmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	cmd.SetArgs([]string{"--format", "xml"})
	execErr := cmd.Execute()
	if execErr == nil {
		t.Fatal("expected error for --format xml, got nil")
	}
	if !strings.Contains(execErr.Error(), "unknown format") {
		t.Errorf("unexpected error message: %v", execErr)
	}
}

func TestCmd_UnknownType(t *testing.T) {
	cmd := Cmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	cmd.SetArgs([]string{"--type", "invalid"})
	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for invalid --type, got nil")
	}
	if !strings.Contains(err.Error(), "unknown project type") {
		t.Errorf("unexpected error: %v", err)
	}
}
