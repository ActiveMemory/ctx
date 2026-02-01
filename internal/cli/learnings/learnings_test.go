//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package learnings

import (
	"testing"
)

func TestCmd(t *testing.T) {
	cmd := Cmd()

	if cmd == nil {
		t.Fatal("Cmd() returned nil")
	}

	if cmd.Use != "learnings" {
		t.Errorf("Cmd().Use = %q, want %q", cmd.Use, "learnings")
	}

	if cmd.Short == "" {
		t.Error("Cmd().Short is empty")
	}

	if cmd.Long == "" {
		t.Error("Cmd().Long is empty")
	}
}

func TestCmd_HasReindexSubcommand(t *testing.T) {
	cmd := Cmd()

	var found bool
	for _, sub := range cmd.Commands() {
		if sub.Use == "reindex" {
			found = true
			if sub.Short == "" {
				t.Error("reindex subcommand has empty Short description")
			}
			if sub.RunE == nil {
				t.Error("reindex subcommand has no RunE function")
			}
			break
		}
	}

	if !found {
		t.Error("reindex subcommand not found")
	}
}
