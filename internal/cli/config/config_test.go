//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package config

import (
	"testing"
)

// TestCmd_HasSubcommands verifies the config command
// includes expected subcommands.
func TestCmd_HasSubcommands(t *testing.T) {
	cmd := Cmd()

	expected := map[string]bool{
		"switch": false,
		"status": false,
		"schema": false,
	}

	for _, sub := range cmd.Commands() {
		if _, ok := expected[sub.Name()]; ok {
			expected[sub.Name()] = true
		}
	}

	for name, found := range expected {
		if !found {
			t.Errorf("config command should have %q subcommand", name)
		}
	}
}
