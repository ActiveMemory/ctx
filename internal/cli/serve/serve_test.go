//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package serve

import (
	"testing"
)

func TestCmd(t *testing.T) {
	cmd := Cmd()

	if cmd == nil {
		t.Fatal("Cmd() returned nil")
	}

	if cmd.Use != "serve [directory]" {
		t.Errorf("Cmd().Use = %q, want %q", cmd.Use, "serve [directory]")
	}

	if cmd.Short == "" {
		t.Error("Cmd().Short is empty")
	}

	if cmd.Long == "" {
		t.Error("Cmd().Long is empty")
	}

	if cmd.RunE == nil {
		t.Error("Cmd().RunE is nil")
	}
}

func TestCmd_AcceptsArgs(t *testing.T) {
	cmd := Cmd()

	// Should accept 0 or 1 args
	if err := cmd.Args(cmd, []string{}); err != nil {
		t.Errorf("should accept 0 args: %v", err)
	}

	if err := cmd.Args(cmd, []string{"./docs"}); err != nil {
		t.Errorf("should accept 1 arg: %v", err)
	}

	if err := cmd.Args(cmd, []string{"a", "b"}); err == nil {
		t.Error("should reject 2 args")
	}
}
