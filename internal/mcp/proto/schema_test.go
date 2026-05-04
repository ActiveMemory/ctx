//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package proto_test

import (
	"encoding/json"
	"testing"

	cfgSchema "github.com/ActiveMemory/ctx/internal/config/mcp/schema"
	"github.com/ActiveMemory/ctx/internal/mcp/proto"
)

func roundTrip(
	t *testing.T, v interface{}, dst interface{},
) {
	t.Helper()
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if err := json.Unmarshal(data, dst); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
}

func TestRequestRoundTrip(t *testing.T) {
	orig := proto.Request{
		JSONRPC: "2.0",
		ID:      json.RawMessage(`1`),
		Method:  "tools/call",
		Params:  json.RawMessage(`{"name":"ctx_status"}`),
	}
	var got proto.Request
	roundTrip(t, orig, &got)
	if got.JSONRPC != orig.JSONRPC {
		t.Errorf("JSONRPC = %q, want %q",
			got.JSONRPC, orig.JSONRPC)
	}
	if got.Method != orig.Method {
		t.Errorf("Method = %q, want %q",
			got.Method, orig.Method)
	}
	if string(got.ID) != string(orig.ID) {
		t.Errorf("ID = %s, want %s", got.ID, orig.ID)
	}
}

func TestResponseSuccessRoundTrip(t *testing.T) {
	orig := proto.Response{
		JSONRPC: "2.0",
		ID:      json.RawMessage(`1`),
		Result:  map[string]string{"key": "value"},
	}
	var got proto.Response
	roundTrip(t, orig, &got)
	if got.JSONRPC != "2.0" {
		t.Errorf("JSONRPC = %q, want %q",
			got.JSONRPC, "2.0")
	}
	if got.Error != nil {
		t.Errorf("unexpected error: %v", got.Error)
	}
}

func TestResponseErrorRoundTrip(t *testing.T) {
	orig := proto.Response{
		JSONRPC: "2.0",
		ID:      json.RawMessage(`1`),
		Error: &proto.RPCError{
			Code:    cfgSchema.ErrCodeNotFound,
			Message: "method not found",
		},
	}
	var got proto.Response
	roundTrip(t, orig, &got)
	if got.Error == nil {
		t.Fatal("expected error in response")
	}
	if got.Error.Code != cfgSchema.ErrCodeNotFound {
		t.Errorf("Code = %d, want %d",
			got.Error.Code, cfgSchema.ErrCodeNotFound)
	}
	if got.Error.Message != "method not found" {
		t.Errorf("Message = %q, want %q",
			got.Error.Message, "method not found")
	}
}

func TestNotificationRoundTrip(t *testing.T) {
	orig := proto.Notification{
		JSONRPC: "2.0",
		Method:  "notifications/initialized",
	}
	var got proto.Notification
	roundTrip(t, orig, &got)
	if got.Method != "notifications/initialized" {
		t.Errorf("Method = %q, want %q",
			got.Method, "notifications/initialized")
	}
}

func TestRPCErrorWithData(t *testing.T) {
	orig := proto.RPCError{
		Code:    cfgSchema.ErrCodeInvalidArg,
		Message: "invalid",
		Data:    map[string]string{"field": "name"},
	}
	var got proto.RPCError
	roundTrip(t, orig, &got)
	if got.Code != cfgSchema.ErrCodeInvalidArg {
		t.Errorf("Code = %d, want %d",
			got.Code, cfgSchema.ErrCodeInvalidArg)
	}
}

func TestInitializeParamsRoundTrip(t *testing.T) {
	orig := proto.InitializeParams{
		ProtocolVersion: cfgSchema.ProtocolVersion,
		ClientInfo: proto.AppInfo{
			Name:    "test-client",
			Version: "1.0.0",
		},
	}
	var got proto.InitializeParams
	roundTrip(t, orig, &got)
	if got.ProtocolVersion != cfgSchema.ProtocolVersion {
		t.Errorf("ProtocolVersion = %q, want %q",
			got.ProtocolVersion, cfgSchema.ProtocolVersion)
	}
	if got.ClientInfo.Name != "test-client" {
		t.Errorf("ClientInfo.Name = %q, want %q",
			got.ClientInfo.Name, "test-client")
	}
}

func TestInitializeResultRoundTrip(t *testing.T) {
	orig := proto.InitializeResult{
		ProtocolVersion: cfgSchema.ProtocolVersion,
		Capabilities: proto.ServerCaps{
			Resources: &proto.ResourcesCap{
				Subscribe:   true,
				ListChanged: true,
			},
			Tools:   &proto.ToolsCap{ListChanged: true},
			Prompts: &proto.PromptsCap{ListChanged: false},
		},
		ServerInfo: proto.AppInfo{
			Name:    "ctx",
			Version: "0.3.0",
		},
	}
	var got proto.InitializeResult
	roundTrip(t, orig, &got)
	if got.Capabilities.Resources == nil {
		t.Fatal("expected Resources capability")
	}
	if !got.Capabilities.Resources.Subscribe {
		t.Error("expected Subscribe = true")
	}
}

func TestResourceRoundTrip(t *testing.T) {
	orig := proto.Resource{
		URI:      "ctx://context/tasks",
		Name:     "tasks",
		MimeType: "text/markdown",
	}
	var got proto.Resource
	roundTrip(t, orig, &got)
	if got.URI != orig.URI {
		t.Errorf("URI = %q, want %q",
			got.URI, orig.URI)
	}
}

func TestToolRoundTrip(t *testing.T) {
	orig := proto.Tool{
		Name: "ctx_status",
		InputSchema: proto.InputSchema{
			Type: "object",
			Properties: map[string]proto.Property{
				"verbose": {
					Type:        "boolean",
					Description: "Verbose",
				},
			},
			Required: []string{"verbose"},
		},
		Annotations: &proto.ToolAnnotations{
			ReadOnlyHint: true,
		},
	}
	var got proto.Tool
	roundTrip(t, orig, &got)
	if got.Name != "ctx_status" {
		t.Errorf("Name = %q, want %q",
			got.Name, "ctx_status")
	}
	if got.Annotations == nil ||
		!got.Annotations.ReadOnlyHint {
		t.Error("expected ReadOnlyHint = true")
	}
}

func TestCallToolParamsRoundTrip(t *testing.T) {
	orig := proto.CallToolParams{
		Name: "ctx_add",
		Arguments: map[string]interface{}{
			"type":    "task",
			"content": "Test",
		},
	}
	var got proto.CallToolParams
	roundTrip(t, orig, &got)
	if got.Name != "ctx_add" {
		t.Errorf("Name = %q, want %q",
			got.Name, "ctx_add")
	}
}

func TestCallToolResultRoundTrip(t *testing.T) {
	orig := proto.CallToolResult{
		Content: []proto.ToolContent{
			{Type: "text", Text: "Done"},
		},
	}
	var got proto.CallToolResult
	roundTrip(t, orig, &got)
	if len(got.Content) != 1 {
		t.Fatalf("Content count = %d, want 1",
			len(got.Content))
	}
	if got.Content[0].Text != "Done" {
		t.Errorf("Text = %q, want %q",
			got.Content[0].Text, "Done")
	}
	if got.IsError {
		t.Error("expected IsError = false")
	}
}

func TestCallToolResultErrorRoundTrip(t *testing.T) {
	orig := proto.CallToolResult{
		Content: []proto.ToolContent{
			{Type: "text", Text: "failed"},
		},
		IsError: true,
	}
	var got proto.CallToolResult
	roundTrip(t, orig, &got)
	if !got.IsError {
		t.Error("expected IsError = true")
	}
}

func TestPromptRoundTrip(t *testing.T) {
	orig := proto.Prompt{
		Name: "ctx-session-start",
		Arguments: []proto.PromptArgument{
			{Name: "content", Required: true},
		},
	}
	var got proto.Prompt
	roundTrip(t, orig, &got)
	if got.Name != "ctx-session-start" {
		t.Errorf("Name = %q, want %q",
			got.Name, "ctx-session-start")
	}
	if len(got.Arguments) != 1 ||
		!got.Arguments[0].Required {
		t.Error("expected 1 required argument")
	}
}

func TestGetPromptResultRoundTrip(t *testing.T) {
	orig := proto.GetPromptResult{
		Description: "Test",
		Messages: []proto.PromptMessage{
			{
				Role: "user",
				Content: proto.ToolContent{
					Type: "text",
					Text: "Hi",
				},
			},
		},
	}
	var got proto.GetPromptResult
	roundTrip(t, orig, &got)
	if len(got.Messages) != 1 {
		t.Fatalf("Messages count = %d, want 1",
			len(got.Messages))
	}
	if got.Messages[0].Role != "user" {
		t.Errorf("Role = %q, want %q",
			got.Messages[0].Role, "user")
	}
}

func TestSubscribeParamsRoundTrip(t *testing.T) {
	orig := proto.SubscribeParams{
		URI: "ctx://context/tasks",
	}
	var got proto.SubscribeParams
	roundTrip(t, orig, &got)
	if got.URI != orig.URI {
		t.Errorf("URI = %q, want %q",
			got.URI, orig.URI)
	}
}

func TestUnsubscribeParamsRoundTrip(t *testing.T) {
	orig := proto.UnsubscribeParams{
		URI: "ctx://context/decisions",
	}
	var got proto.UnsubscribeParams
	roundTrip(t, orig, &got)
	if got.URI != orig.URI {
		t.Errorf("URI = %q, want %q",
			got.URI, orig.URI)
	}
}

func TestResourceUpdatedParamsRoundTrip(t *testing.T) {
	orig := proto.ResourceUpdatedParams{
		URI: "ctx://context/tasks",
	}
	var got proto.ResourceUpdatedParams
	roundTrip(t, orig, &got)
	if got.URI != orig.URI {
		t.Errorf("URI = %q, want %q",
			got.URI, orig.URI)
	}
}

func TestErrorCodeConstants(t *testing.T) {
	if cfgSchema.ErrCodeParse != -32700 {
		t.Errorf("ErrCodeParse = %d, want -32700",
			cfgSchema.ErrCodeParse)
	}
	if cfgSchema.ErrCodeNotFound != -32601 {
		t.Errorf("ErrCodeNotFound = %d, want -32601",
			cfgSchema.ErrCodeNotFound)
	}
	if cfgSchema.ErrCodeInvalidArg != -32602 {
		t.Errorf("ErrCodeInvalidArg = %d, want -32602",
			cfgSchema.ErrCodeInvalidArg)
	}
	if cfgSchema.ErrCodeInternal != -32603 {
		t.Errorf("ErrCodeInternal = %d, want -32603",
			cfgSchema.ErrCodeInternal)
	}
}

func TestProtocolVersionValue(t *testing.T) {
	if cfgSchema.ProtocolVersion != "2024-11-05" {
		t.Errorf("ProtocolVersion = %q, want %q",
			cfgSchema.ProtocolVersion, "2024-11-05")
	}
}

func TestRequestNilParams(t *testing.T) {
	orig := proto.Request{
		JSONRPC: "2.0",
		ID:      json.RawMessage(`"abc"`),
		Method:  "ping",
	}
	data, err := json.Marshal(orig)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var got proto.Request
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Params != nil {
		t.Errorf("expected nil Params, got %s",
			got.Params)
	}
}

func TestResponseNilID(t *testing.T) {
	orig := proto.Response{
		JSONRPC: "2.0",
		Error: &proto.RPCError{
			Code:    cfgSchema.ErrCodeParse,
			Message: "parse error",
		},
	}
	var got proto.Response
	roundTrip(t, orig, &got)
	if got.ID != nil {
		t.Errorf("expected nil ID, got %s", got.ID)
	}
}
