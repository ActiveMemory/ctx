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

func TestBlockNonPathCtx_RelativePath(t *testing.T) {
	tests := []struct {
		name    string
		command string
		blocked bool
	}{
		{"./ctx at start", "./ctx status", true},
		{"./ctx alone", "./ctx", true},
		{"./dist/ctx", "./dist/ctx status", true},
		{"go run ./cmd/ctx", "go run ./cmd/ctx status", true},
		{"absolute path", "/home/user/project/ctx status", true},
		{"absolute /tmp path", "/tmp/build/ctx status", true},
		{"after separator", "echo hello && ./ctx status", true},
		{"after pipe", "echo hello | ./ctx status", true},
		{"allowed: ctx from PATH", "ctx status", false},
		{"allowed: git with ctx path arg", "git -C ./ctx/path status", false},
		{"allowed: empty command", "", false},
		{"allowed: /tmp/ctx-test", "/tmp/ctx-test/bin run", false},
		{"absolute /var path", "/var/tmp/ctx run", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := newTestCmd()
			input := `{"tool_input":{"command":"` + tt.command + `"}}`
			stdin := createTempStdin(t, input)

			if err := runBlockNonPathCtx(cmd, stdin); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			out := cmdOutput(cmd)
			hasBlock := strings.Contains(out, `"decision":"block"`)

			if tt.blocked && !hasBlock {
				t.Errorf("expected block for %q, got: %s", tt.command, out)
			}
			if !tt.blocked && hasBlock {
				t.Errorf("expected allow for %q, got: %s", tt.command, out)
			}
		})
	}
}

func TestBlockNonPathCtx_JSONOutput(t *testing.T) {
	cmd := newTestCmd()
	stdin := createTempStdin(t, `{"tool_input":{"command":"./ctx status"}}`)

	if err := runBlockNonPathCtx(cmd, stdin); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := cmdOutput(cmd)
	if !strings.Contains(out, `"decision":"block"`) {
		t.Errorf("expected JSON block output, got: %s", out)
	}
	if !strings.Contains(out, `"reason"`) {
		t.Errorf("expected reason in output, got: %s", out)
	}
	if !strings.Contains(out, "CONSTITUTION.md") {
		t.Errorf("expected CONSTITUTION.md reference, got: %s", out)
	}
}
