//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package decision

import "testing"

func TestCmd(t *testing.T) {
	cmd := Cmd()

	if cmd == nil {
		t.Fatal("Cmd() returned nil")
	}

	if cmd.Use != "decision" {
		t.Errorf("Cmd().Use = %q, want %q", cmd.Use, "decision")
	}

	if cmd.Short == "" {
		t.Error("Cmd().Short is empty")
	}

	if cmd.Long == "" {
		t.Error("Cmd().Long is empty")
	}
}
