//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package server

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/mcp/proto"
)

// callTool sends one tools/call request and returns its tool result.
func callTool(
	t *testing.T, srv *Server, name string, args map[string]interface{},
) proto.CallToolResult {
	t.Helper()
	resp := request(t, srv, "tools/call", proto.CallToolParams{
		Name: name, Arguments: args,
	})
	if resp.Error != nil {
		t.Fatalf("unexpected error: %v", resp.Error.Message)
	}
	raw, _ := json.Marshal(resp.Result)
	var result proto.CallToolResult
	if err := json.Unmarshal(raw, &result); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	return result
}

// readContext returns the content of one context file.
func readContext(t *testing.T, contextDir, name string) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(contextDir, name))
	if err != nil {
		t.Fatalf("read %s: %v", name, err)
	}
	return string(b)
}

// runGit runs git in dir and returns its trimmed output.
func runGit(t *testing.T, dir string, args ...string) string {
	t.Helper()
	//nolint:gosec // G204: test fixture, args are hardcoded by callers
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
	return strings.TrimSpace(string(out))
}

// gitInit makes dir a repository with one commit, skipping the test
// when git is not installed.
func gitInit(t *testing.T, dir string) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skipf("git not on PATH: %v", err)
	}
	runGit(t, dir, "init", "-q")
	runGit(t, dir, "config", "user.email", "test@example.com")
	runGit(t, dir, "config", "user.name", "Test User")
	// A global commit.gpgsign must never reach for the GPG agent.
	runGit(t, dir, "config", "commit.gpgsign", "false")
	runGit(t, dir, "commit", "--allow-empty", "-q", "-m", "init")
}

// Branch and commit come from git when the caller omits them or sends
// "unknown". Expected values are read from git, not assumed: the
// default branch name depends on the machine's git config.
func TestToolEntryResolvesGitProvenance(t *testing.T) {
	for _, tool := range []string{"ctx_add", "ctx_watch_update"} {
		t.Run(tool, func(t *testing.T) {
			srv, contextDir := newTestServer(t)
			root := filepath.Dir(contextDir)
			gitInit(t, root)
			result := callTool(t, srv, tool, map[string]interface{}{
				"type": "task", "content": "Resolved task",
				"section": "Misc", "session_id": "unknown",
				"branch": "unknown",
			})
			if result.IsError {
				t.Fatalf("unexpected tool error: %s", result.Content[0].Text)
			}
			want := "Resolved task #session:unknown" +
				" #branch:" + runGit(t, root, "branch", "--show-current") +
				" #commit:" + runGit(t, root, "rev-parse", "--short", "HEAD")
			tasks := readContext(t, contextDir, ctx.Task)
			if !strings.Contains(tasks, want) {
				t.Errorf("want %q in TASKS.md, got:\n%s", want, tasks)
			}
		})
	}
}

// An explicit value wins over the repository's own branch and commit.
func TestToolAddKeepsExplicitProvenance(t *testing.T) {
	srv, contextDir := newTestServer(t)
	gitInit(t, filepath.Dir(contextDir))
	result := callTool(t, srv, "ctx_add", map[string]interface{}{
		"type": "task", "content": "Explicit task", "section": "Misc",
		"session_id": "test1234", "branch": "feature/x",
		"commit": "abc1234",
	})
	if result.IsError {
		t.Fatalf("unexpected tool error: %s", result.Content[0].Text)
	}
	want := "Explicit task #session:test1234" +
		" #branch:feature/x #commit:abc1234"
	tasks := readContext(t, contextDir, ctx.Task)
	if !strings.Contains(tasks, want) {
		t.Errorf("want %q in TASKS.md, got:\n%s", want, tasks)
	}
}

// The probe from issue #190, from a client with no session ID to relay.
func TestToolAddLearningUnknownSession(t *testing.T) {
	srv, contextDir := newTestServer(t)
	result := callTool(t, srv, "ctx_add", map[string]interface{}{
		"type": "learning", "content": "Probe learning",
		"context": "c", "lesson": "l", "application": "a",
		"session_id": "unknown",
	})
	if result.IsError {
		t.Fatalf("unexpected tool error: %s", result.Content[0].Text)
	}
	if !strings.Contains(
		readContext(t, contextDir, ctx.Learning), "Probe learning",
	) {
		t.Error("learning not written to LEARNINGS.md")
	}
}

// The server resolves branch and commit, never the session ID.
func TestToolAddRequiresSessionID(t *testing.T) {
	srv, _ := newTestServer(t)
	result := callTool(t, srv, "ctx_add", map[string]interface{}{
		"type": "task", "content": "No session", "section": "Misc",
	})
	if !result.IsError {
		t.Fatal("expected a missing-field error without session_id")
	}
	text := result.Content[0].Text
	if !strings.Contains(text, "--session-id") ||
		strings.Contains(text, "--branch") ||
		strings.Contains(text, "--commit") {
		t.Errorf("want only --session-id reported missing, got: %s", text)
	}
}
