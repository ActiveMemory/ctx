//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package server

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ActiveMemory/ctx/internal/config/ctx"
	cfgSchema "github.com/ActiveMemory/ctx/internal/config/mcp/schema"
	"github.com/ActiveMemory/ctx/internal/mcp/proto"
	mcpIO "github.com/ActiveMemory/ctx/internal/mcp/server/io"
	"github.com/ActiveMemory/ctx/internal/rc"
)

func newTestServer(t *testing.T) (*Server, string) {
	t.Helper()
	dir := t.TempDir()

	// Change CWD to the temp dir so ValidateBoundary passes.
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(origDir) })

	contextDir := filepath.Join(dir, ".context")
	if err := os.MkdirAll(contextDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	// Tools dispatched through the MCP server call rc.ContextDir()
	// for paths under .context/; declare it so they resolve without
	// the "context directory not declared" error.
	t.Setenv("CTX_DIR", contextDir)
	rc.Reset()
	t.Cleanup(rc.Reset)
	files := map[string]string{
		ctx.Constitution:  "# Constitution\n\n- Rule 1: Never break things\n",
		ctx.Task:          "# Tasks\n\n- [ ] Build MCP server\n- [ ] Write tests\n",
		ctx.Decision:      "# Decisions\n",
		ctx.Convention:    "# Conventions\n\n- Use Go idioms\n",
		ctx.Learning:      "# Learnings\n",
		ctx.Architecture:  "# Architecture\n",
		ctx.Glossary:      "# Glossary\n",
		ctx.AgentPlaybook: "# Agent Playbook\n\nRead context files first.\n",
	}
	for name, content := range files {
		p := filepath.Join(contextDir, name)
		if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
			t.Fatalf("write %s: %v", name, err)
		}
	}
	srv := New(contextDir, "test")
	return srv, contextDir
}

func request(
	t *testing.T, srv *Server,
	method string, params interface{},
) *proto.Response {
	t.Helper()
	var rawParams json.RawMessage
	if params != nil {
		b, err := json.Marshal(params)
		if err != nil {
			t.Fatalf("marshal params: %v", err)
		}
		rawParams = b
	}
	idBytes, _ := json.Marshal(1)
	req := proto.Request{
		JSONRPC: "2.0",
		ID:      idBytes,
		Method:  method,
		Params:  rawParams,
	}
	line, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}
	var out bytes.Buffer
	srv.in = bytes.NewReader(append(line, '\n'))
	srv.out = mcpIO.NewWriter(&out)
	if serveErr := srv.Serve(); serveErr != nil {
		t.Fatalf("serve: %v", serveErr)
	}
	var resp proto.Response
	if err := json.Unmarshal(out.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v (raw: %s)", err, out.String())
	}
	return &resp
}

func TestInitialize(t *testing.T) {
	srv, _ := newTestServer(t)
	resp := request(t, srv, "initialize", proto.InitializeParams{
		ProtocolVersion: cfgSchema.ProtocolVersion,
		ClientInfo:      proto.AppInfo{Name: "test", Version: "1.0"},
	})
	if resp.Error != nil {
		t.Fatalf("unexpected error: %v", resp.Error.Message)
	}
	raw, _ := json.Marshal(resp.Result)
	var result proto.InitializeResult
	if err := json.Unmarshal(raw, &result); err != nil {
		t.Fatalf("unmarshal result: %v", err)
	}
	if result.ProtocolVersion != cfgSchema.ProtocolVersion {
		t.Errorf(
			"protocol version = %q, want %q",
			result.ProtocolVersion, cfgSchema.ProtocolVersion,
		)
	}
	if result.ServerInfo.Name != "ctx" {
		t.Errorf("server name = %q, want %q", result.ServerInfo.Name, "ctx")
	}
	if result.Capabilities.Resources == nil {
		t.Error("expected resources capability")
	}
	hasRes := result.Capabilities.Resources != nil
	if hasRes && !result.Capabilities.Resources.Subscribe {
		t.Error("expected resources subscribe capability")
	}
	if result.Capabilities.Tools == nil {
		t.Error("expected tools capability")
	}
	if result.Capabilities.Prompts == nil {
		t.Error("expected prompts capability")
	}
}

func TestPing(t *testing.T) {
	srv, _ := newTestServer(t)
	resp := request(t, srv, "ping", nil)
	if resp.Error != nil {
		t.Fatalf("unexpected error: %v", resp.Error.Message)
	}
}

func TestMethodNotFound(t *testing.T) {
	srv, _ := newTestServer(t)
	resp := request(t, srv, "nonexistent/method", nil)
	if resp.Error == nil {
		t.Fatal("expected error for unknown method")
	}
	if resp.Error.Code != cfgSchema.ErrCodeNotFound {
		t.Errorf("error code = %d, want %d", resp.Error.Code, cfgSchema.ErrCodeNotFound)
	}
}

func TestResourcesList(t *testing.T) {
	srv, _ := newTestServer(t)
	resp := request(t, srv, "resources/list", nil)
	if resp.Error != nil {
		t.Fatalf("unexpected error: %v", resp.Error.Message)
	}
	raw, _ := json.Marshal(resp.Result)
	var result proto.ResourceListResult
	if err := json.Unmarshal(raw, &result); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(result.Resources) != 9 {
		t.Errorf("resource count = %d, want 9", len(result.Resources))
	}
	found := false
	for _, r := range result.Resources {
		if r.URI == "ctx://context/agent" {
			found = true
			break
		}
	}
	if !found {
		t.Error("agent resource not found in list")
	}
}

func TestResourcesRead(t *testing.T) {
	srv, _ := newTestServer(t)
	resp := request(t, srv, "resources/read", proto.ReadResourceParams{
		URI: "ctx://context/tasks",
	})
	if resp.Error != nil {
		t.Fatalf("unexpected error: %v", resp.Error.Message)
	}
	raw, _ := json.Marshal(resp.Result)
	var result proto.ReadResourceResult
	if err := json.Unmarshal(raw, &result); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(result.Contents) != 1 {
		t.Fatalf("contents count = %d, want 1", len(result.Contents))
	}
	if !strings.Contains(result.Contents[0].Text, "Build MCP server") {
		t.Errorf("expected tasks content, got: %s", result.Contents[0].Text)
	}
}

func TestResourcesReadAgent(t *testing.T) {
	srv, _ := newTestServer(t)
	resp := request(t, srv, "resources/read", proto.ReadResourceParams{
		URI: "ctx://context/agent",
	})
	if resp.Error != nil {
		t.Fatalf("unexpected error: %v", resp.Error.Message)
	}
	raw, _ := json.Marshal(resp.Result)
	var result proto.ReadResourceResult
	if err := json.Unmarshal(raw, &result); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	text := result.Contents[0].Text
	if !strings.Contains(text, "Context Packet") {
		t.Error("expected Context Packet header in agent resource")
	}
}

func TestResourcesReadUnknown(t *testing.T) {
	srv, _ := newTestServer(t)
	resp := request(t, srv, "resources/read", proto.ReadResourceParams{
		URI: "ctx://context/nonexistent",
	})
	if resp.Error == nil {
		t.Fatal("expected error for unknown resource")
	}
}

func TestToolsList(t *testing.T) {
	srv, _ := newTestServer(t)
	resp := request(t, srv, "tools/list", nil)
	if resp.Error != nil {
		t.Fatalf("unexpected error: %v", resp.Error.Message)
	}
	raw, _ := json.Marshal(resp.Result)
	var result proto.ToolListResult
	if err := json.Unmarshal(raw, &result); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(result.Tools) != 15 {
		t.Errorf("tool count = %d, want 15", len(result.Tools))
	}
	names := make(map[string]bool)
	for _, tool := range result.Tools {
		names[tool.Name] = true
	}
	for _, want := range []string{
		"ctx_status", "ctx_add", "ctx_complete", "ctx_drift",
		"ctx_journal_source", "ctx_watch_update", "ctx_compact",
		"ctx_next", "ctx_check_task_completion",
		"ctx_session_event", "ctx_remind",
		"ctx_steering_get", "ctx_search",
		"ctx_session_start", "ctx_session_end",
	} {
		if !names[want] {
			t.Errorf("missing tool: %s", want)
		}
	}
}

func TestToolStatus(t *testing.T) {
	srv, _ := newTestServer(t)
	resp := request(t, srv, "tools/call", proto.CallToolParams{
		Name: "ctx_status",
	})
	if resp.Error != nil {
		t.Fatalf("unexpected error: %v", resp.Error.Message)
	}
	raw, _ := json.Marshal(resp.Result)
	var result proto.CallToolResult
	if err := json.Unmarshal(raw, &result); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if result.IsError {
		t.Fatalf("unexpected tool error: %s", result.Content[0].Text)
	}
	text := result.Content[0].Text
	if !strings.Contains(text, "TASKS.md") {
		t.Errorf("expected TASKS.md in status output, got: %s", text)
	}
}

func TestToolComplete(t *testing.T) {
	srv, contextDir := newTestServer(t)
	resp := request(t, srv, "tools/call", proto.CallToolParams{
		Name:      "ctx_complete",
		Arguments: map[string]interface{}{"query": "1"},
	})
	if resp.Error != nil {
		t.Fatalf("unexpected error: %v", resp.Error.Message)
	}
	raw, _ := json.Marshal(resp.Result)
	var result proto.CallToolResult
	if err := json.Unmarshal(raw, &result); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if result.IsError {
		t.Fatalf("unexpected tool error: %s", result.Content[0].Text)
	}
	if !strings.Contains(result.Content[0].Text, "Build MCP server") {
		t.Errorf("expected completed task name, got: %s", result.Content[0].Text)
	}
	content, err := os.ReadFile(filepath.Join(contextDir, ctx.Task))
	if err != nil {
		t.Fatalf("read tasks: %v", err)
	}
	if !strings.Contains(string(content), "- [x] Build MCP server") {
		t.Errorf("task not marked complete in file: %s", string(content))
	}
}

func TestToolDrift(t *testing.T) {
	srv, _ := newTestServer(t)
	resp := request(t, srv, "tools/call", proto.CallToolParams{
		Name: "ctx_drift",
	})
	if resp.Error != nil {
		t.Fatalf("unexpected error: %v", resp.Error.Message)
	}
	raw, _ := json.Marshal(resp.Result)
	var result proto.CallToolResult
	if err := json.Unmarshal(raw, &result); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if result.IsError {
		t.Fatalf("unexpected tool error: %s", result.Content[0].Text)
	}
	if !strings.Contains(result.Content[0].Text, "Status:") {
		t.Errorf("expected Status in drift output, got: %s", result.Content[0].Text)
	}
}

func TestToolAdd(t *testing.T) {
	tests := []struct {
		name         string
		args         map[string]interface{}
		wantErr      bool
		wantFile     string
		wantContains string
	}{
		{
			name: "add task",
			args: map[string]interface{}{
				"type": "task", "content": "Test task", "section": "Misc",
				"session_id": "test1234", "branch": "main", "commit": "abc123",
			},
			wantFile:     ctx.Task,
			wantContains: "Test task",
		},
		{
			name: "add convention",
			args: map[string]interface{}{
				"type": "convention", "content": "Use tabs",
			},
			wantFile:     ctx.Convention,
			wantContains: "Use tabs",
		},
		{
			name: "add decision",
			args: map[string]interface{}{
				"type":        "decision",
				"content":     "Use Redis",
				"session_id":  "test1234",
				"branch":      "main",
				"commit":      "abc123",
				"context":     "Need caching",
				"rationale":   "Fast and simple",
				"consequence": "Ops must manage Redis",
			},
			wantFile:     ctx.Decision,
			wantContains: "Use Redis",
		},
		{
			name: "add learning",
			args: map[string]interface{}{
				"type":        "learning",
				"content":     "Go embed requires same package",
				"session_id":  "test1234",
				"branch":      "main",
				"commit":      "abc123",
				"context":     "Tried parent dir",
				"lesson":      "Only same or child dirs",
				"application": "Keep files in internal",
			},
			wantFile:     ctx.Learning,
			wantContains: "Go embed",
		},
		{
			name: "decision missing rationale",
			args: map[string]interface{}{
				"type": "decision", "content": "X",
				"session_id": "test1234", "branch": "main", "commit": "abc123",
				"context": "Y",
			},
			wantErr: true,
		},
		{
			name: "learning missing lesson",
			args: map[string]interface{}{
				"type": "learning", "content": "X",
				"session_id": "test1234", "branch": "main", "commit": "abc123",
				"context": "Y",
			},
			wantErr: true,
		},
		{
			name:    "missing content",
			args:    map[string]interface{}{"type": "task"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv, contextDir := newTestServer(t)
			resp := request(t, srv, "tools/call", proto.CallToolParams{
				Name:      "ctx_add",
				Arguments: tt.args,
			})
			if resp.Error != nil {
				t.Fatalf("unexpected error: %v", resp.Error.Message)
			}
			raw, _ := json.Marshal(resp.Result)
			var result proto.CallToolResult
			if err := json.Unmarshal(raw, &result); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}

			if tt.wantErr {
				if !result.IsError {
					t.Fatalf("expected tool error, got success: %s", result.Content[0].Text)
				}
				return
			}

			if result.IsError {
				t.Fatalf("unexpected tool error: %s", result.Content[0].Text)
			}

			content, err := os.ReadFile(filepath.Join(contextDir, tt.wantFile))
			if err != nil {
				t.Fatalf("read %s: %v", tt.wantFile, err)
			}
			if !strings.Contains(string(content), tt.wantContains) {
				t.Errorf(
					"expected %q in %s, got: %s",
					tt.wantContains, tt.wantFile,
					string(content),
				)
			}
		})
	}
}

func TestToolUnknown(t *testing.T) {
	srv, _ := newTestServer(t)
	resp := request(t, srv, "tools/call", proto.CallToolParams{
		Name: "nonexistent_tool",
	})
	if resp.Error == nil {
		t.Fatal("expected error for unknown tool")
	}
}

func TestNotification(t *testing.T) {
	srv, _ := newTestServer(t)
	req := proto.Request{
		JSONRPC: "2.0",
		Method:  "notifications/initialized",
	}
	line, _ := json.Marshal(req)
	var out bytes.Buffer
	srv.in = bytes.NewReader(append(line, '\n'))
	srv.out = mcpIO.NewWriter(&out)
	if err := srv.Serve(); err != nil {
		t.Fatalf("serve: %v", err)
	}
	if out.Len() != 0 {
		t.Errorf("expected no output for notification, got: %s", out.String())
	}
}

func TestParseError(t *testing.T) {
	srv, _ := newTestServer(t)
	var out bytes.Buffer
	srv.in = bytes.NewReader([]byte("not json\n"))
	srv.out = mcpIO.NewWriter(&out)
	if err := srv.Serve(); err != nil {
		t.Fatalf("serve: %v", err)
	}
	var resp proto.Response
	if err := json.Unmarshal(out.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if resp.Error == nil || resp.Error.Code != cfgSchema.ErrCodeParse {
		t.Errorf("expected parse error, got: %+v", resp.Error)
	}
}

// --- New tool tests (v0.2) ---

func TestToolRecall(t *testing.T) {
	srv, _ := newTestServer(t)
	resp := request(t, srv, "tools/call", proto.CallToolParams{
		Name:      "ctx_journal_source",
		Arguments: map[string]interface{}{"limit": float64(3)},
	})
	if resp.Error != nil {
		t.Fatalf("unexpected error: %v", resp.Error.Message)
	}
	raw, _ := json.Marshal(resp.Result)
	var result proto.CallToolResult
	if err := json.Unmarshal(raw, &result); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if result.IsError {
		t.Fatalf("unexpected tool error: %s", result.Content[0].Text)
	}
	// Should return something (either sessions or "No sessions found.")
	if len(result.Content) == 0 {
		t.Error("expected content in recall response")
	}
}

func TestToolRecallInvalidDate(t *testing.T) {
	srv, _ := newTestServer(t)
	resp := request(t, srv, "tools/call", proto.CallToolParams{
		Name:      "ctx_journal_source",
		Arguments: map[string]interface{}{"since": "not-a-date"},
	})
	if resp.Error != nil {
		t.Fatalf("unexpected error: %v", resp.Error.Message)
	}
	raw, _ := json.Marshal(resp.Result)
	var result proto.CallToolResult
	if err := json.Unmarshal(raw, &result); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !result.IsError {
		t.Error("expected tool error for invalid date")
	}
}

func TestToolWatchUpdate(t *testing.T) {
	srv, contextDir := newTestServer(t)
	resp := request(t, srv, "tools/call", proto.CallToolParams{
		Name: "ctx_watch_update",
		Arguments: map[string]interface{}{
			"type":       "task",
			"content":    "New MCP task from watch",
			"section":    "Misc",
			"session_id": "test1234",
			"branch":     "main",
			"commit":     "abc123",
		},
	})
	if resp.Error != nil {
		t.Fatalf("unexpected error: %v", resp.Error.Message)
	}
	raw, _ := json.Marshal(resp.Result)
	var result proto.CallToolResult
	if err := json.Unmarshal(raw, &result); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if result.IsError {
		t.Fatalf("unexpected tool error: %s", result.Content[0].Text)
	}
	text := result.Content[0].Text
	if !strings.Contains(text, "Wrote task") {
		t.Errorf("expected advisory text, got: %s", text)
	}

	// Verify the entry was written.
	content, err := os.ReadFile(filepath.Join(contextDir, ctx.Task))
	if err != nil {
		t.Fatalf("read tasks: %v", err)
	}
	if !strings.Contains(string(content), "New MCP task from watch") {
		t.Errorf("task not found in file: %s", string(content))
	}
}

func TestToolWatchUpdateDecision(t *testing.T) {
	srv, contextDir := newTestServer(t)
	resp := request(t, srv, "tools/call", proto.CallToolParams{
		Name: "ctx_watch_update",
		Arguments: map[string]interface{}{
			"type":        "decision",
			"content":     "Use MCP protocol",
			"session_id":  "test1234",
			"branch":      "main",
			"commit":      "abc123",
			"context":     "Need AI tool integration",
			"rationale":   "Standard protocol",
			"consequence": "Must maintain compatibility",
		},
	})
	if resp.Error != nil {
		t.Fatalf("unexpected error: %v", resp.Error.Message)
	}
	raw, _ := json.Marshal(resp.Result)
	var result proto.CallToolResult
	if err := json.Unmarshal(raw, &result); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if result.IsError {
		t.Fatalf("unexpected tool error: %s", result.Content[0].Text)
	}

	content, err := os.ReadFile(filepath.Join(contextDir, ctx.Decision))
	if err != nil {
		t.Fatalf("read decisions: %v", err)
	}
	if !strings.Contains(string(content), "Use MCP protocol") {
		t.Errorf("decision not found in file: %s", string(content))
	}
}

func TestToolWatchUpdateValidationError(t *testing.T) {
	srv, _ := newTestServer(t)
	resp := request(t, srv, "tools/call", proto.CallToolParams{
		Name: "ctx_watch_update",
		Arguments: map[string]interface{}{
			"type":    "decision",
			"content": "Missing fields",
			// Missing context, rationale, consequences.
		},
	})
	if resp.Error != nil {
		t.Fatalf("unexpected error: %v", resp.Error.Message)
	}
	raw, _ := json.Marshal(resp.Result)
	var result proto.CallToolResult
	if err := json.Unmarshal(raw, &result); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !result.IsError {
		t.Error("expected validation error for decision missing required fields")
	}
}

func TestToolWatchUpdateComplete(t *testing.T) {
	srv, contextDir := newTestServer(t)
	resp := request(t, srv, "tools/call", proto.CallToolParams{
		Name: "ctx_watch_update",
		Arguments: map[string]interface{}{
			"type":    "complete",
			"content": "1",
		},
	})
	if resp.Error != nil {
		t.Fatalf("unexpected error: %v", resp.Error.Message)
	}
	raw, _ := json.Marshal(resp.Result)
	var result proto.CallToolResult
	if err := json.Unmarshal(raw, &result); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if result.IsError {
		t.Fatalf("unexpected tool error: %s", result.Content[0].Text)
	}
	if !strings.Contains(result.Content[0].Text, "Build MCP server") {
		t.Errorf("expected completed task name, got: %s", result.Content[0].Text)
	}

	content, err := os.ReadFile(filepath.Join(contextDir, ctx.Task))
	if err != nil {
		t.Fatalf("read tasks: %v", err)
	}
	if !strings.Contains(string(content), "- [x] Build MCP server") {
		t.Errorf("task not marked complete: %s", string(content))
	}
}

func TestToolCompact(t *testing.T) {
	srv, contextDir := newTestServer(t)

	// Set up TASKS.md with a completed task and a Completed section.
	tasksContent := "# Tasks\n\n" +
		"- [x] Done task\n- [ ] Pending task\n\n" +
		"## Completed\n\n"
	if err := os.WriteFile(
		filepath.Join(contextDir, ctx.Task),
		[]byte(tasksContent), 0o644,
	); err != nil {
		t.Fatalf("write tasks: %v", err)
	}

	resp := request(t, srv, "tools/call", proto.CallToolParams{
		Name:      "ctx_compact",
		Arguments: map[string]interface{}{},
	})
	if resp.Error != nil {
		t.Fatalf("unexpected error: %v", resp.Error.Message)
	}
	raw, _ := json.Marshal(resp.Result)
	var result proto.CallToolResult
	if err := json.Unmarshal(raw, &result); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if result.IsError {
		t.Fatalf("unexpected tool error: %s", result.Content[0].Text)
	}
	text := result.Content[0].Text
	if !strings.Contains(text, "Compacted") {
		t.Errorf("expected compacted message, got: %s", text)
	}

	// Verify task was moved.
	content, err := os.ReadFile(filepath.Join(contextDir, ctx.Task))
	if err != nil {
		t.Fatalf("read tasks: %v", err)
	}
	if strings.Contains(string(content), "- [x] Done task\n- [ ] Pending task") {
		t.Error("completed task should have been moved to Completed section")
	}
}

func TestToolCompactClean(t *testing.T) {
	srv, _ := newTestServer(t)
	resp := request(t, srv, "tools/call", proto.CallToolParams{
		Name:      "ctx_compact",
		Arguments: map[string]interface{}{},
	})
	if resp.Error != nil {
		t.Fatalf("unexpected error: %v", resp.Error.Message)
	}
	raw, _ := json.Marshal(resp.Result)
	var result proto.CallToolResult
	if err := json.Unmarshal(raw, &result); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if result.IsError {
		t.Fatalf("unexpected tool error: %s", result.Content[0].Text)
	}
	// No completed tasks to move - should report clean.
	text := result.Content[0].Text
	if !strings.Contains(text, "clean") && !strings.Contains(text, "Compacted") {
		t.Errorf("expected clean or compacted message, got: %s", text)
	}
}

func TestToolNext(t *testing.T) {
	srv, _ := newTestServer(t)
	resp := request(t, srv, "tools/call", proto.CallToolParams{
		Name: "ctx_next",
	})
	if resp.Error != nil {
		t.Fatalf("unexpected error: %v", resp.Error.Message)
	}
	raw, _ := json.Marshal(resp.Result)
	var result proto.CallToolResult
	if err := json.Unmarshal(raw, &result); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if result.IsError {
		t.Fatalf("unexpected tool error: %s", result.Content[0].Text)
	}
	text := result.Content[0].Text
	if !strings.Contains(text, "Build MCP server") {
		t.Errorf("expected first pending task, got: %s", text)
	}
}

func TestToolNextAllComplete(t *testing.T) {
	srv, contextDir := newTestServer(t)

	tasksContent := "# Tasks\n\n- [x] Done 1\n- [x] Done 2\n"
	if err := os.WriteFile(
		filepath.Join(contextDir, ctx.Task),
		[]byte(tasksContent), 0o644,
	); err != nil {
		t.Fatalf("write tasks: %v", err)
	}

	resp := request(t, srv, "tools/call", proto.CallToolParams{
		Name: "ctx_next",
	})
	if resp.Error != nil {
		t.Fatalf("unexpected error: %v", resp.Error.Message)
	}
	raw, _ := json.Marshal(resp.Result)
	var result proto.CallToolResult
	if err := json.Unmarshal(raw, &result); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !strings.Contains(result.Content[0].Text, "All tasks completed") {
		t.Errorf("expected all complete message, got: %s", result.Content[0].Text)
	}
}

func TestToolCheckTaskCompletion(t *testing.T) {
	srv, _ := newTestServer(t)
	resp := request(t, srv, "tools/call", proto.CallToolParams{
		Name: "ctx_check_task_completion",
		Arguments: map[string]interface{}{
			"recent_action": "Finished build of the MCP server",
		},
	})
	if resp.Error != nil {
		t.Fatalf("unexpected error: %v", resp.Error.Message)
	}
	raw, _ := json.Marshal(resp.Result)
	var result proto.CallToolResult
	if err := json.Unmarshal(raw, &result); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	text := result.Content[0].Text
	// Should find overlap with "Build MCP server".
	if !strings.Contains(text, "Build MCP server") {
		t.Errorf("expected task match nudge, got: %s", text)
	}
}

func TestToolCheckTaskCompletionNoMatch(t *testing.T) {
	srv, _ := newTestServer(t)

	// Prime session state to avoid governance warnings in response.
	request(t, srv, "tools/call", proto.CallToolParams{
		Name:      "ctx_session_event",
		Arguments: map[string]interface{}{"type": "start"},
	})
	request(t, srv, "tools/call", proto.CallToolParams{
		Name: "ctx_status",
	})

	resp := request(t, srv, "tools/call", proto.CallToolParams{
		Name: "ctx_check_task_completion",
		Arguments: map[string]interface{}{
			"recent_action": "Updated CSS styles",
		},
	})
	if resp.Error != nil {
		t.Fatalf("unexpected error: %v", resp.Error.Message)
	}
	raw, _ := json.Marshal(resp.Result)
	var result proto.CallToolResult
	if err := json.Unmarshal(raw, &result); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	// Should not match.
	if result.Content[0].Text != "" {
		t.Errorf(
			"expected empty response for no match, got: %s",
			result.Content[0].Text,
		)
	}
}

func TestToolSessionEventStart(t *testing.T) {
	srv, _ := newTestServer(t)
	resp := request(t, srv, "tools/call", proto.CallToolParams{
		Name: "ctx_session_event",
		Arguments: map[string]interface{}{
			"type":   "start",
			"caller": "vscode",
		},
	})
	if resp.Error != nil {
		t.Fatalf("unexpected error: %v", resp.Error.Message)
	}
	raw, _ := json.Marshal(resp.Result)
	var result proto.CallToolResult
	if err := json.Unmarshal(raw, &result); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if result.IsError {
		t.Fatalf("unexpected tool error: %s", result.Content[0].Text)
	}
	text := result.Content[0].Text
	if !strings.Contains(text, "Session started") {
		t.Errorf("expected session start message, got: %s", text)
	}
	if !strings.Contains(text, "vscode") {
		t.Errorf("expected caller in message, got: %s", text)
	}
}

func TestToolSessionEventEnd(t *testing.T) {
	srv, _ := newTestServer(t)
	resp := request(t, srv, "tools/call", proto.CallToolParams{
		Name:      "ctx_session_event",
		Arguments: map[string]interface{}{"type": "end"},
	})
	if resp.Error != nil {
		t.Fatalf("unexpected error: %v", resp.Error.Message)
	}
	raw, _ := json.Marshal(resp.Result)
	var result proto.CallToolResult
	if err := json.Unmarshal(raw, &result); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if result.IsError {
		t.Fatalf("unexpected tool error: %s", result.Content[0].Text)
	}
	text := result.Content[0].Text
	if !strings.Contains(text, "Session ending") {
		t.Errorf("expected session end message, got: %s", text)
	}
}

func TestToolSessionEventInvalid(t *testing.T) {
	srv, _ := newTestServer(t)
	resp := request(t, srv, "tools/call", proto.CallToolParams{
		Name:      "ctx_session_event",
		Arguments: map[string]interface{}{"type": "pause"},
	})
	if resp.Error != nil {
		t.Fatalf("unexpected error: %v", resp.Error.Message)
	}
	raw, _ := json.Marshal(resp.Result)
	var result proto.CallToolResult
	if err := json.Unmarshal(raw, &result); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !result.IsError {
		t.Error("expected error for invalid event type")
	}
}

func TestToolRemind(t *testing.T) {
	srv, _ := newTestServer(t)
	resp := request(t, srv, "tools/call", proto.CallToolParams{
		Name: "ctx_remind",
	})
	if resp.Error != nil {
		t.Fatalf("unexpected error: %v", resp.Error.Message)
	}
	raw, _ := json.Marshal(resp.Result)
	var result proto.CallToolResult
	if err := json.Unmarshal(raw, &result); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if result.IsError {
		t.Fatalf("unexpected tool error: %s", result.Content[0].Text)
	}
	// No reminders file in test setup - should return "No reminders."
	if !strings.Contains(result.Content[0].Text, "No reminders") {
		t.Errorf("expected no reminders message, got: %s", result.Content[0].Text)
	}
}

// --- Prompt tests ---

func TestPromptsList(t *testing.T) {
	srv, _ := newTestServer(t)
	resp := request(t, srv, "prompts/list", nil)
	if resp.Error != nil {
		t.Fatalf("unexpected error: %v", resp.Error.Message)
	}
	raw, _ := json.Marshal(resp.Result)
	var result proto.PromptListResult
	if err := json.Unmarshal(raw, &result); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(result.Prompts) != 5 {
		t.Errorf("prompt count = %d, want 5", len(result.Prompts))
	}
	names := make(map[string]bool)
	for _, p := range result.Prompts {
		names[p.Name] = true
	}
	for _, want := range []string{
		"ctx-session-start", "ctx-decision-add", "ctx-learning-add",
		"ctx-reflect", "ctx-checkpoint",
	} {
		if !names[want] {
			t.Errorf("missing prompt: %s", want)
		}
	}
}

func TestPromptSessionStart(t *testing.T) {
	srv, _ := newTestServer(t)
	resp := request(t, srv, "prompts/get", proto.GetPromptParams{
		Name: "ctx-session-start",
	})
	if resp.Error != nil {
		t.Fatalf("unexpected error: %v", resp.Error.Message)
	}
	raw, _ := json.Marshal(resp.Result)
	var result proto.GetPromptResult
	if err := json.Unmarshal(raw, &result); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(result.Messages) == 0 {
		t.Fatal("expected at least one message in session-start prompt")
	}
	text := result.Messages[0].Content.Text
	if !strings.Contains(text, "session") {
		t.Errorf("expected session orientation text, got: %s", text)
	}
}

func TestPromptAddDecision(t *testing.T) {
	srv, _ := newTestServer(t)
	resp := request(t, srv, "prompts/get", proto.GetPromptParams{
		Name: "ctx-decision-add",
		Arguments: map[string]string{
			"content":     "Use Go",
			"context":     "Need compiled language",
			"rationale":   "Fast",
			"consequence": "Team needs Go skills",
		},
	})
	if resp.Error != nil {
		t.Fatalf("unexpected error: %v", resp.Error.Message)
	}
	raw, _ := json.Marshal(resp.Result)
	var result proto.GetPromptResult
	if err := json.Unmarshal(raw, &result); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(result.Messages) == 0 {
		t.Fatal("expected message in decision prompt")
	}
	text := result.Messages[0].Content.Text
	if !strings.Contains(text, "Use Go") {
		t.Errorf("expected decision content in text, got: %s", text)
	}
}

func TestPromptReflect(t *testing.T) {
	srv, _ := newTestServer(t)
	resp := request(t, srv, "prompts/get", proto.GetPromptParams{
		Name: "ctx-reflect",
	})
	if resp.Error != nil {
		t.Fatalf("unexpected error: %v", resp.Error.Message)
	}
	raw, _ := json.Marshal(resp.Result)
	var result proto.GetPromptResult
	if err := json.Unmarshal(raw, &result); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(result.Messages) == 0 {
		t.Fatal("expected message in reflect prompt")
	}
	text := result.Messages[0].Content.Text
	if !strings.Contains(text, "Reflect") {
		t.Errorf("expected reflect text, got: %s", text)
	}
}

func TestPromptCheckpoint(t *testing.T) {
	srv, _ := newTestServer(t)
	resp := request(t, srv, "prompts/get", proto.GetPromptParams{
		Name: "ctx-checkpoint",
	})
	if resp.Error != nil {
		t.Fatalf("unexpected error: %v", resp.Error.Message)
	}
	raw, _ := json.Marshal(resp.Result)
	var result proto.GetPromptResult
	if err := json.Unmarshal(raw, &result); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(result.Messages) == 0 {
		t.Fatal("expected message in checkpoint prompt")
	}
	text := result.Messages[0].Content.Text
	if !strings.Contains(text, "checkpoint") {
		t.Errorf("expected checkpoint text, got: %s", text)
	}
}

func TestPromptUnknown(t *testing.T) {
	srv, _ := newTestServer(t)
	resp := request(t, srv, "prompts/get", proto.GetPromptParams{
		Name: "nonexistent",
	})
	if resp.Error == nil {
		t.Fatal("expected error for unknown prompt")
	}
}

// --- Session state tests ---

func TestSessionStateTracking(t *testing.T) {
	srv, _ := newTestServer(t)

	// Start session.
	request(t, srv, "tools/call", proto.CallToolParams{
		Name:      "ctx_session_event",
		Arguments: map[string]interface{}{"type": "start"},
	})

	// Call a few tools.
	request(t, srv, "tools/call", proto.CallToolParams{Name: "ctx_status"})
	request(t, srv, "tools/call", proto.CallToolParams{Name: "ctx_next"})

	// End session - should report tool call count.
	resp := request(t, srv, "tools/call", proto.CallToolParams{
		Name:      "ctx_session_event",
		Arguments: map[string]interface{}{"type": "end"},
	})
	raw, _ := json.Marshal(resp.Result)
	var result proto.CallToolResult
	if err := json.Unmarshal(raw, &result); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	text := result.Content[0].Text
	// After start, status, next, end = 4 calls
	// (start resets, so status + next + end = 3)
	if !strings.Contains(text, "tool calls") {
		t.Errorf("expected tool call stats, got: %s", text)
	}
}

func TestResourcesSubscribe(t *testing.T) {
	srv, _ := newTestServer(t)
	resp := request(t, srv, "resources/subscribe", proto.SubscribeParams{
		URI: "ctx://context/tasks",
	})
	if resp.Error != nil {
		t.Fatalf("unexpected error: %v", resp.Error.Message)
	}
	// Cleanup: stop poller.
	srv.poller.Stop()
}

func TestResourcesUnsubscribe(t *testing.T) {
	srv, _ := newTestServer(t)
	// Subscribe first.
	request(t, srv, "resources/subscribe", proto.SubscribeParams{
		URI: "ctx://context/tasks",
	})
	// Then unsubscribe.
	resp := request(t, srv, "resources/unsubscribe", proto.UnsubscribeParams{
		URI: "ctx://context/tasks",
	})
	if resp.Error != nil {
		t.Fatalf("unexpected error: %v", resp.Error.Message)
	}
	srv.poller.Stop()
}

func TestResourcePollerNotification(t *testing.T) {
	srv, contextDir := newTestServer(t)

	var mu sync.Mutex
	var notifications []proto.Notification
	srv.poller.SetNotifyFunc(func(n proto.Notification) {
		mu.Lock()
		notifications = append(notifications, n)
		mu.Unlock()
	})

	// Subscribe to tasks.
	request(t, srv, "resources/subscribe", proto.SubscribeParams{
		URI: "ctx://context/tasks",
	})

	// Modify the tasks file.
	time.Sleep(10 * time.Millisecond) // Ensure mtime differs.
	taskFile := filepath.Join(contextDir, ctx.Task)
	if err := os.WriteFile(
		taskFile,
		[]byte("# Tasks\n\n- [ ] Modified task\n"),
		0o644,
	); err != nil {
		t.Fatalf("write: %v", err)
	}

	// Manually trigger a poll check instead of waiting for the timer.
	srv.poller.CheckChanges()

	mu.Lock()
	count := len(notifications)
	mu.Unlock()

	if count != 1 {
		t.Fatalf("notification count = %d, want 1", count)
	}

	params, ok := notifications[0].Params.(proto.ResourceUpdatedParams)
	if !ok {
		t.Fatalf("unexpected params type: %T", notifications[0].Params)
	}
	if params.URI != "ctx://context/tasks" {
		t.Errorf("notification URI = %q, want %q", params.URI, "ctx://context/tasks")
	}

	srv.poller.Stop()
}

// --- Steering and session hook tool tests ---

func TestToolSteeringGetWithPrompt(t *testing.T) {
	srv, contextDir := newTestServer(t)

	// Create steering directory with test files.
	steeringDir := filepath.Join(contextDir, "steering")
	if err := os.MkdirAll(steeringDir, 0o755); err != nil {
		t.Fatalf("mkdir steering: %v", err)
	}
	alwaysFile := "---\nname: always-rules\ndescription: Always included\ninclusion: always\npriority: 10\n---\n\nAlways body content.\n"
	autoFile := "---\nname: api-rules\ndescription: API design\ninclusion: auto\npriority: 20\n---\n\nAPI body content.\n"
	if err := os.WriteFile(filepath.Join(steeringDir, "always-rules.md"), []byte(alwaysFile), 0o644); err != nil {
		t.Fatalf("write always: %v", err)
	}
	if err := os.WriteFile(filepath.Join(steeringDir, "api-rules.md"), []byte(autoFile), 0o644); err != nil {
		t.Fatalf("write auto: %v", err)
	}

	resp := request(t, srv, "tools/call", proto.CallToolParams{
		Name:      "ctx_steering_get",
		Arguments: map[string]interface{}{"prompt": "API design"},
	})
	if resp.Error != nil {
		t.Fatalf("unexpected error: %v", resp.Error.Message)
	}
	raw, _ := json.Marshal(resp.Result)
	var result proto.CallToolResult
	if err := json.Unmarshal(raw, &result); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if result.IsError {
		t.Fatalf("unexpected tool error: %s", result.Content[0].Text)
	}
	text := result.Content[0].Text
	// Should include both always and auto-matched files.
	if !strings.Contains(text, "always-rules") {
		t.Errorf("expected always-rules in response, got: %s", text)
	}
	if !strings.Contains(text, "api-rules") {
		t.Errorf("expected api-rules in response, got: %s", text)
	}
}

func TestToolSteeringGetWithoutPrompt(t *testing.T) {
	srv, contextDir := newTestServer(t)

	// Create steering directory with test files.
	steeringDir := filepath.Join(contextDir, "steering")
	if err := os.MkdirAll(steeringDir, 0o755); err != nil {
		t.Fatalf("mkdir steering: %v", err)
	}
	alwaysFile := "---\nname: always-rules\ndescription: Always included\ninclusion: always\npriority: 10\n---\n\nAlways body.\n"
	autoFile := "---\nname: api-rules\ndescription: API design\ninclusion: auto\npriority: 20\n---\n\nAPI body.\n"
	if err := os.WriteFile(filepath.Join(steeringDir, "always-rules.md"), []byte(alwaysFile), 0o644); err != nil {
		t.Fatalf("write always: %v", err)
	}
	if err := os.WriteFile(filepath.Join(steeringDir, "api-rules.md"), []byte(autoFile), 0o644); err != nil {
		t.Fatalf("write auto: %v", err)
	}

	// No prompt — should return only always files.
	resp := request(t, srv, "tools/call", proto.CallToolParams{
		Name: "ctx_steering_get",
	})
	if resp.Error != nil {
		t.Fatalf("unexpected error: %v", resp.Error.Message)
	}
	raw, _ := json.Marshal(resp.Result)
	var result proto.CallToolResult
	if err := json.Unmarshal(raw, &result); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if result.IsError {
		t.Fatalf("unexpected tool error: %s", result.Content[0].Text)
	}
	text := result.Content[0].Text
	if !strings.Contains(text, "always-rules") {
		t.Errorf("expected always-rules in response, got: %s", text)
	}
	if strings.Contains(text, "api-rules") {
		t.Errorf("auto file should not be included without matching prompt, got: %s", text)
	}
}

func TestToolSessionStartNoHooks(t *testing.T) {
	srv, contextDir := newTestServer(t)

	// Create empty hooks directory so discovery succeeds.
	hooksDir := filepath.Join(contextDir, "hooks")
	if err := os.MkdirAll(filepath.Join(hooksDir, "session-start"), 0o755); err != nil {
		t.Fatalf("mkdir hooks: %v", err)
	}

	resp := request(t, srv, "tools/call", proto.CallToolParams{
		Name: "ctx_session_start",
	})
	if resp.Error != nil {
		t.Fatalf("unexpected error: %v", resp.Error.Message)
	}
	raw, _ := json.Marshal(resp.Result)
	var result proto.CallToolResult
	if err := json.Unmarshal(raw, &result); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if result.IsError {
		t.Fatalf("unexpected tool error: %s", result.Content[0].Text)
	}
	// No hooks exist — should return success message.
	text := result.Content[0].Text
	if !strings.Contains(text, "Session start hooks executed") &&
		!strings.Contains(text, "No additional context") {
		t.Errorf("expected success message for no hooks, got: %s", text)
	}
}

func TestToolSessionEndWithSummary(t *testing.T) {
	srv, contextDir := newTestServer(t)

	// Create empty hooks directory.
	hooksDir := filepath.Join(contextDir, "hooks")
	if err := os.MkdirAll(filepath.Join(hooksDir, "session-end"), 0o755); err != nil {
		t.Fatalf("mkdir hooks: %v", err)
	}

	resp := request(t, srv, "tools/call", proto.CallToolParams{
		Name: "ctx_session_end",
		Arguments: map[string]interface{}{
			"summary": "Completed MCP server implementation",
		},
	})
	if resp.Error != nil {
		t.Fatalf("unexpected error: %v", resp.Error.Message)
	}
	raw, _ := json.Marshal(resp.Result)
	var result proto.CallToolResult
	if err := json.Unmarshal(raw, &result); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if result.IsError {
		t.Fatalf("unexpected tool error: %s", result.Content[0].Text)
	}
	// No hooks exist — should return success.
	text := result.Content[0].Text
	if !strings.Contains(text, "Session end hooks executed") {
		t.Errorf("expected session end success message, got: %s", text)
	}
}

func TestToolSearch(t *testing.T) {
	srv, _ := newTestServer(t)

	resp := request(t, srv, "tools/call", proto.CallToolParams{
		Name:      "ctx_search",
		Arguments: map[string]interface{}{"query": "Rule 1"},
	})
	if resp.Error != nil {
		t.Fatalf("unexpected error: %v", resp.Error.Message)
	}
	raw, _ := json.Marshal(resp.Result)
	var result proto.CallToolResult
	if err := json.Unmarshal(raw, &result); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if result.IsError {
		t.Fatalf("unexpected tool error: %s", result.Content[0].Text)
	}
	text := result.Content[0].Text
	if !strings.Contains(text, "CONSTITUTION.md") {
		t.Errorf("expected match in CONSTITUTION.md, got: %s", text)
	}
}

func TestToolSearchNoQuery(t *testing.T) {
	srv, _ := newTestServer(t)

	resp := request(t, srv, "tools/call", proto.CallToolParams{
		Name: "ctx_search",
	})
	if resp.Error != nil {
		t.Fatalf("unexpected error: %v", resp.Error.Message)
	}
	raw, _ := json.Marshal(resp.Result)
	var result proto.CallToolResult
	if err := json.Unmarshal(raw, &result); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !result.IsError {
		t.Error("expected error when query is missing")
	}
}
