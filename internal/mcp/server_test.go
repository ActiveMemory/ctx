//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package mcp

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/config/ctx"
)

func newTestServer(t *testing.T) (*Server, string) {
	t.Helper()
	dir := t.TempDir()
	contextDir := filepath.Join(dir, ".context")
	if err := os.MkdirAll(contextDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
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
	srv := NewServer(contextDir, "test")
	return srv, contextDir
}

func request(t *testing.T, srv *Server, method string, params interface{}) *Response {
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
	req := Request{
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
	srv.out = &out
	if serveErr := srv.Serve(); serveErr != nil {
		t.Fatalf("serve: %v", serveErr)
	}
	var resp Response
	if err := json.Unmarshal(out.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v (raw: %s)", err, out.String())
	}
	return &resp
}

func TestInitialize(t *testing.T) {
	srv, _ := newTestServer(t)
	resp := request(t, srv, "initialize", InitializeParams{
		ProtocolVersion: protocolVersion,
		ClientInfo:      AppInfo{Name: "test", Version: "1.0"},
	})
	if resp.Error != nil {
		t.Fatalf("unexpected error: %v", resp.Error.Message)
	}
	raw, _ := json.Marshal(resp.Result)
	var result InitializeResult
	if err := json.Unmarshal(raw, &result); err != nil {
		t.Fatalf("unmarshal result: %v", err)
	}
	if result.ProtocolVersion != protocolVersion {
		t.Errorf("protocol version = %q, want %q", result.ProtocolVersion, protocolVersion)
	}
	if result.ServerInfo.Name != "ctx" {
		t.Errorf("server name = %q, want %q", result.ServerInfo.Name, "ctx")
	}
	if result.Capabilities.Resources == nil {
		t.Error("expected resources capability")
	}
	if result.Capabilities.Tools == nil {
		t.Error("expected tools capability")
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
	if resp.Error.Code != errCodeNotFound {
		t.Errorf("error code = %d, want %d", resp.Error.Code, errCodeNotFound)
	}
}

func TestResourcesList(t *testing.T) {
	srv, _ := newTestServer(t)
	resp := request(t, srv, "resources/list", nil)
	if resp.Error != nil {
		t.Fatalf("unexpected error: %v", resp.Error.Message)
	}
	raw, _ := json.Marshal(resp.Result)
	var result ResourceListResult
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
	resp := request(t, srv, "resources/read", ReadResourceParams{
		URI: "ctx://context/tasks",
	})
	if resp.Error != nil {
		t.Fatalf("unexpected error: %v", resp.Error.Message)
	}
	raw, _ := json.Marshal(resp.Result)
	var result ReadResourceResult
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
	resp := request(t, srv, "resources/read", ReadResourceParams{
		URI: "ctx://context/agent",
	})
	if resp.Error != nil {
		t.Fatalf("unexpected error: %v", resp.Error.Message)
	}
	raw, _ := json.Marshal(resp.Result)
	var result ReadResourceResult
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
	resp := request(t, srv, "resources/read", ReadResourceParams{
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
	var result ToolListResult
	if err := json.Unmarshal(raw, &result); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(result.Tools) != 4 {
		t.Errorf("tool count = %d, want 4", len(result.Tools))
	}
	names := make(map[string]bool)
	for _, tool := range result.Tools {
		names[tool.Name] = true
	}
	for _, want := range []string{"ctx_status", "ctx_add", "ctx_complete", "ctx_drift"} {
		if !names[want] {
			t.Errorf("missing tool: %s", want)
		}
	}
}

func TestToolStatus(t *testing.T) {
	srv, _ := newTestServer(t)
	resp := request(t, srv, "tools/call", CallToolParams{
		Name: "ctx_status",
	})
	if resp.Error != nil {
		t.Fatalf("unexpected error: %v", resp.Error.Message)
	}
	raw, _ := json.Marshal(resp.Result)
	var result CallToolResult
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
	resp := request(t, srv, "tools/call", CallToolParams{
		Name:      "ctx_complete",
		Arguments: map[string]interface{}{"query": "1"},
	})
	if resp.Error != nil {
		t.Fatalf("unexpected error: %v", resp.Error.Message)
	}
	raw, _ := json.Marshal(resp.Result)
	var result CallToolResult
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
	resp := request(t, srv, "tools/call", CallToolParams{
		Name: "ctx_drift",
	})
	if resp.Error != nil {
		t.Fatalf("unexpected error: %v", resp.Error.Message)
	}
	raw, _ := json.Marshal(resp.Result)
	var result CallToolResult
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
			name:         "add task",
			args:         map[string]interface{}{"type": "task", "content": "Test task"},
			wantFile:     ctx.Task,
			wantContains: "Test task",
		},
		{
			name:         "add convention",
			args:         map[string]interface{}{"type": "convention", "content": "Use tabs"},
			wantFile:     ctx.Convention,
			wantContains: "Use tabs",
		},
		{
			name: "add decision",
			args: map[string]interface{}{
				"type":         "decision",
				"content":      "Use Redis",
				"context":      "Need caching",
				"rationale":    "Fast and simple",
				"consequences": "Ops must manage Redis",
			},
			wantFile:     ctx.Decision,
			wantContains: "Use Redis",
		},
		{
			name: "add learning",
			args: map[string]interface{}{
				"type":        "learning",
				"content":     "Go embed requires same package",
				"context":     "Tried parent dir",
				"lesson":      "Only same or child dirs",
				"application": "Keep files in internal",
			},
			wantFile:     ctx.Learning,
			wantContains: "Go embed",
		},
		{
			name:    "decision missing rationale",
			args:    map[string]interface{}{"type": "decision", "content": "X", "context": "Y"},
			wantErr: true,
		},
		{
			name:    "learning missing lesson",
			args:    map[string]interface{}{"type": "learning", "content": "X", "context": "Y"},
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
			resp := request(t, srv, "tools/call", CallToolParams{
				Name:      "ctx_add",
				Arguments: tt.args,
			})
			if resp.Error != nil {
				t.Fatalf("unexpected error: %v", resp.Error.Message)
			}
			raw, _ := json.Marshal(resp.Result)
			var result CallToolResult
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
				t.Errorf("expected %q in %s, got: %s", tt.wantContains, tt.wantFile, string(content))
			}
		})
	}
}

func TestToolUnknown(t *testing.T) {
	srv, _ := newTestServer(t)
	resp := request(t, srv, "tools/call", CallToolParams{
		Name: "nonexistent_tool",
	})
	if resp.Error == nil {
		t.Fatal("expected error for unknown tool")
	}
}

func TestNotification(t *testing.T) {
	srv, _ := newTestServer(t)
	req := Request{
		JSONRPC: "2.0",
		Method:  "notifications/initialized",
	}
	line, _ := json.Marshal(req)
	var out bytes.Buffer
	srv.in = bytes.NewReader(append(line, '\n'))
	srv.out = &out
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
	srv.out = &out
	if err := srv.Serve(); err != nil {
		t.Fatalf("serve: %v", err)
	}
	var resp Response
	if err := json.Unmarshal(out.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if resp.Error == nil || resp.Error.Code != errCodeParse {
		t.Errorf("expected parse error, got: %+v", resp.Error)
	}
}
