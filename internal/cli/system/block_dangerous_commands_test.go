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

func TestBlockDangerousCommands(t *testing.T) {
	tests := []struct {
		name    string
		command string
		blocked bool
		reason  string
	}{
		{
			name:    "mid-command sudo blocked",
			command: "echo ok && sudo rm -rf /",
			blocked: true,
			reason:  "sudo",
		},
		{
			name:    "mid-command sudo after semicolon",
			command: "echo ok; sudo apt install foo",
			blocked: true,
			reason:  "sudo",
		},
		{
			name:    "mid-command sudo after or",
			command: "test -f foo || sudo mkdir /bar",
			blocked: true,
			reason:  "sudo",
		},
		{
			name:    "mid-command git push blocked",
			command: "make && git push",
			blocked: true,
			reason:  "git push",
		},
		{
			name:    "mid-command git push after semicolon",
			command: "echo done; git push origin main",
			blocked: true,
			reason:  "git push",
		},
		{
			name:    "cp to /usr/local/bin blocked",
			command: "cp ./ctx /usr/local/bin/",
			blocked: true,
			reason:  "bin directories",
		},
		{
			name:    "mv to ~/go/bin blocked",
			command: "mv ctx ~/go/bin/ctx",
			blocked: true,
			reason:  "bin directories",
		},
		{
			name:    "cp to /usr/bin blocked",
			command: "cp ctx /usr/bin/ctx",
			blocked: true,
			reason:  "bin directories",
		},
		{
			name:    "mv to /home/user/go/bin blocked",
			command: "mv ctx /home/jose/go/bin/ctx",
			blocked: true,
			reason:  "bin directories",
		},
		{
			name:    "mv to /home/user/.local/bin blocked",
			command: "mv ctx /home/jose/.local/bin/ctx",
			blocked: true,
			reason:  "bin directories",
		},
		{
			name:    "install to ~/.local/bin blocked",
			command: "install ctx ~/.local/bin/",
			blocked: true,
			reason:  "~/.local/bin",
		},
		{
			name:    "cp to ~/.local/bin blocked",
			command: "cp ctx ~/.local/bin/ctx",
			blocked: true,
			reason:  "bin directories",
		},
		{
			name:    "clean command allowed",
			command: "go test ./...",
			blocked: false,
		},
		{
			name:    "prefix sudo not caught (deny-list job)",
			command: "sudo make install",
			blocked: false,
		},
		{
			name:    "prefix git push not caught (deny-list job)",
			command: "git push origin main",
			blocked: false,
		},
		{
			name:    "empty command silent",
			command: "",
			blocked: false,
		},
		{
			name:    "make build allowed",
			command: "make build && make test",
			blocked: false,
		},
		{
			name:    "cp to normal dir allowed",
			command: "cp file.txt /tmp/backup/",
			blocked: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := newTestCmd()
			input := `{"tool_input":{"command":"` + tt.command + `"}}`
			stdin := createTempStdin(t, input)

			if err := runBlockDangerousCommands(cmd, stdin); err != nil {
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
			if tt.blocked && tt.reason != "" && !strings.Contains(out, tt.reason) {
				t.Errorf("expected reason containing %q, got: %s", tt.reason, out)
			}
		})
	}
}

func TestBlockDangerousCommands_JSONOutput(t *testing.T) {
	cmd := newTestCmd()
	stdin := createTempStdin(t, `{"tool_input":{"command":"echo ok && sudo rm -rf /"}}`)

	if err := runBlockDangerousCommands(cmd, stdin); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := cmdOutput(cmd)
	if !strings.Contains(out, `"decision":"block"`) {
		t.Errorf("expected JSON block output, got: %s", out)
	}
	if !strings.Contains(out, `"reason"`) {
		t.Errorf("expected reason in output, got: %s", out)
	}
}
