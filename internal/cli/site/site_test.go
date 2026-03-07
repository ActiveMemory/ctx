//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package site

import (
	"testing"
)

// TestCmd_HasFeedSubcommand verifies the site command includes "feed".
func TestCmd_HasFeedSubcommand(t *testing.T) {
	cmd := Cmd()

	found := false
	for _, sub := range cmd.Commands() {
		if sub.Name() == "feed" {
			found = true
			break
		}
	}

	if !found {
		t.Error("site command should have a 'feed' subcommand")
	}
}
