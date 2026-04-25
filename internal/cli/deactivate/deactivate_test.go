//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package deactivate_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/cli/deactivate"
)

// runDeactivate invokes `ctx deactivate` with the given args and
// returns (stdout, error).
func runDeactivate(t *testing.T, args []string) (string, error) {
	t.Helper()
	c := deactivate.Cmd()
	c.SetArgs(args)
	var out bytes.Buffer
	c.SetOut(&out)
	c.SetErr(&out)
	err := c.Execute()
	return out.String(), err
}

// TestDeactivate_DefaultShell: no --shell flag → autodetect from
// $SHELL → bash emitter → `unset CTX_DIR`.
func TestDeactivate_DefaultShell(t *testing.T) {
	t.Setenv("SHELL", "/bin/bash")

	stdout, err := runDeactivate(t, nil)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if strings.TrimSpace(stdout) != "unset CTX_DIR" {
		t.Errorf("stdout = %q, want 'unset CTX_DIR\\n'", stdout)
	}
}

// TestDeactivate_ExplicitZsh: --shell zsh → same POSIX unset
// statement (v1 bash/zsh/sh share syntax).
func TestDeactivate_ExplicitZsh(t *testing.T) {
	stdout, err := runDeactivate(t, []string{"--shell", "zsh"})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if !strings.Contains(stdout, "unset CTX_DIR") {
		t.Errorf("stdout missing unset: %q", stdout)
	}
}

// TestDeactivate_UnknownShell: unknown shell → POSIX unset fallback.
func TestDeactivate_UnknownShell(t *testing.T) {
	stdout, err := runDeactivate(t, []string{"--shell", "rc"})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if !strings.Contains(stdout, "unset CTX_DIR") {
		t.Errorf("stdout missing unset fallback: %q", stdout)
	}
}

// TestDeactivate_RejectsPositionalArgs: deactivate takes no args.
func TestDeactivate_RejectsPositionalArgs(t *testing.T) {
	_, err := runDeactivate(t, []string{"unexpected-arg"})
	if err == nil {
		t.Fatalf("expected error for positional arg, got nil")
	}
}
