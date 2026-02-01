//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package bootstrap

import (
	"testing"
)

func TestRootCmd(t *testing.T) {
	cmd := RootCmd()

	if cmd == nil {
		t.Fatal("RootCmd() returned nil")
	}

	if cmd.Use != "ctx" {
		t.Errorf("RootCmd().Use = %q, want %q", cmd.Use, "ctx")
	}

	if cmd.Short == "" {
		t.Error("RootCmd().Short is empty")
	}

	if cmd.Long == "" {
		t.Error("RootCmd().Long is empty")
	}

	// Check global flags exist
	contextDirFlag := cmd.PersistentFlags().Lookup("context-dir")
	if contextDirFlag == nil {
		t.Error("--context-dir flag not found")
	}

	noColorFlag := cmd.PersistentFlags().Lookup("no-color")
	if noColorFlag == nil {
		t.Error("--no-color flag not found")
	}
}

func TestInitialize(t *testing.T) {
	root := RootCmd()
	cmd := Initialize(root)

	if cmd == nil {
		t.Fatal("Initialize() returned nil")
	}

	// Verify all expected subcommands are registered
	expectedCommands := []string{
		"init",
		"status",
		"load",
		"add",
		"complete",
		"agent",
		"drift",
		"sync",
		"compact",
		"decisions",
		"watch",
		"hook",
		"learnings",
		"session",
		"tasks",
		"loop",
		"recall",
		"journal",
		"serve",
	}

	commands := make(map[string]bool)
	for _, c := range cmd.Commands() {
		commands[c.Use] = true
		// Handle commands with args in Use (e.g., "serve [directory]")
		for _, exp := range expectedCommands {
			if c.Name() == exp {
				commands[exp] = true
			}
		}
	}

	for _, exp := range expectedCommands {
		if !commands[exp] {
			t.Errorf("missing subcommand: %s", exp)
		}
	}
}

func TestRootCmdVersion(t *testing.T) {
	cmd := RootCmd()

	if cmd.Version == "" {
		t.Error("RootCmd().Version is empty")
	}
}
