//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package system

import (
	"strings"
	"testing"
)

func TestPostCommit_GitCommit(t *testing.T) {
	cmd := newTestCmd()
	stdin := createTempStdin(t, `{"tool_input":{"command":"git commit -m 'test'"}}`)

	if err := runPostCommit(cmd, stdin); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := cmdOutput(cmd)
	if !strings.Contains(out, "Post-Commit") {
		t.Errorf("expected post-commit message, got: %s", out)
	}
	if !strings.Contains(out, "Offer context capture") {
		t.Errorf("expected context capture prompt, got: %s", out)
	}
}

func TestPostCommit_AmendSkipped(t *testing.T) {
	cmd := newTestCmd()
	stdin := createTempStdin(t, `{"tool_input":{"command":"git commit --amend"}}`)

	if err := runPostCommit(cmd, stdin); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := cmdOutput(cmd)
	if strings.Contains(out, "Post-Commit") {
		t.Errorf("expected silence for amend, got: %s", out)
	}
}

func TestPostCommit_NonGitCommand(t *testing.T) {
	cmd := newTestCmd()
	stdin := createTempStdin(t, `{"tool_input":{"command":"ls -la"}}`)

	if err := runPostCommit(cmd, stdin); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := cmdOutput(cmd)
	if strings.Contains(out, "Post-Commit") {
		t.Errorf("expected silence for non-git command, got: %s", out)
	}
}

func TestPostCommit_EmptyCommand(t *testing.T) {
	cmd := newTestCmd()
	stdin := createTempStdin(t, `{"tool_input":{"command":""}}`)

	if err := runPostCommit(cmd, stdin); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := cmdOutput(cmd)
	if strings.Contains(out, "Post-Commit") {
		t.Errorf("expected silence for empty command, got: %s", out)
	}
}

func TestPostCommit_GitCommitWithHeredoc(t *testing.T) {
	cmd := newTestCmd()
	stdin := createTempStdin(t, `{"tool_input":{"command":"git commit -m \"$(cat <<'EOF'\nFix bug\nEOF\n)\""}}`)

	if err := runPostCommit(cmd, stdin); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := cmdOutput(cmd)
	if !strings.Contains(out, "Post-Commit") {
		t.Errorf("expected post-commit for heredoc commit, got: %s", out)
	}
}
