//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package journal

import (
	"testing"
)

func TestCmd(t *testing.T) {
	cmd := Cmd()

	if cmd == nil {
		t.Fatal("Cmd() returned nil")
	}

	if cmd.Use != "journal" {
		t.Errorf("Cmd().Use = %q, want %q", cmd.Use, "journal")
	}

	if cmd.Short == "" {
		t.Error("Cmd().Short is empty")
	}

	if cmd.Long == "" {
		t.Error("Cmd().Long is empty")
	}
}

func TestCmd_HasSiteSubcommand(t *testing.T) {
	cmd := Cmd()

	var found bool
	for _, sub := range cmd.Commands() {
		if sub.Use == "site" {
			found = true
			if sub.Short == "" {
				t.Error("site subcommand has empty Short description")
			}
			if sub.RunE == nil {
				t.Error("site subcommand has no RunE function")
			}

			// Check flags
			outputFlag := sub.Flags().Lookup("output")
			if outputFlag == nil {
				t.Error("site subcommand missing --output flag")
			}

			buildFlag := sub.Flags().Lookup("build")
			if buildFlag == nil {
				t.Error("site subcommand missing --build flag")
			}

			serveFlag := sub.Flags().Lookup("serve")
			if serveFlag == nil {
				t.Error("site subcommand missing --serve flag")
			}

			break
		}
	}

	if !found {
		t.Error("site subcommand not found")
	}
}

func TestCmd_HasObsidianSubcommand(t *testing.T) {
	cmd := Cmd()

	var found bool
	for _, sub := range cmd.Commands() {
		if sub.Use == "obsidian" {
			found = true
			if sub.Short == "" {
				t.Error("obsidian subcommand has empty Short description")
			}
			if sub.RunE == nil {
				t.Error("obsidian subcommand has no RunE function")
			}

			outputFlag := sub.Flags().Lookup("output")
			if outputFlag == nil {
				t.Error("obsidian subcommand missing --output flag")
			}

			break
		}
	}

	if !found {
		t.Error("obsidian subcommand not found")
	}
}
