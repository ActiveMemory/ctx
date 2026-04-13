//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package out

import (
	"encoding/json"
	"errors"
	"testing"

	cfgSchema "github.com/ActiveMemory/ctx/internal/config/mcp/schema"
	"github.com/ActiveMemory/ctx/internal/mcp/proto"
)

func TestOkResponse(t *testing.T) {
	id, _ := json.Marshal(1)
	resp := OkResponse(id, map[string]string{"k": "v"})
	if resp.JSONRPC != "2.0" {
		t.Errorf("jsonrpc = %q", resp.JSONRPC)
	}
	if resp.Error != nil {
		t.Error("unexpected error field")
	}
}

func TestErrResponse(t *testing.T) {
	id, _ := json.Marshal(1)
	resp := ErrResponse(id, cfgSchema.ErrCodeInternal, "boom")
	if resp.Error == nil {
		t.Fatal("expected error")
	}
	if resp.Error.Code != cfgSchema.ErrCodeInternal {
		t.Errorf("code = %d", resp.Error.Code)
	}
	if resp.Error.Message != "boom" {
		t.Errorf("msg = %q", resp.Error.Message)
	}
}

func TestToolOK(t *testing.T) {
	id, _ := json.Marshal(1)
	resp := ToolOK(id, "ok")
	raw, _ := json.Marshal(resp.Result)
	var r proto.CallToolResult
	_ = json.Unmarshal(raw, &r)
	if r.IsError {
		t.Error("unexpected isError")
	}
	if r.Content[0].Text != "ok" {
		t.Errorf("text = %q", r.Content[0].Text)
	}
}

func TestToolError(t *testing.T) {
	id, _ := json.Marshal(1)
	resp := ToolError(id, "fail")
	raw, _ := json.Marshal(resp.Result)
	var r proto.CallToolResult
	_ = json.Unmarshal(raw, &r)
	if !r.IsError {
		t.Error("expected isError")
	}
}

func TestToolResultSuccess(t *testing.T) {
	id, _ := json.Marshal(1)
	resp := ToolResult(id, "done", nil)
	raw, _ := json.Marshal(resp.Result)
	var r proto.CallToolResult
	_ = json.Unmarshal(raw, &r)
	if r.IsError {
		t.Error("unexpected isError")
	}
}

func TestToolResultError(t *testing.T) {
	id, _ := json.Marshal(1)
	resp := ToolResult(id, "", errors.New("bad"))
	raw, _ := json.Marshal(resp.Result)
	var r proto.CallToolResult
	_ = json.Unmarshal(raw, &r)
	if !r.IsError {
		t.Error("expected isError")
	}
}

func TestCallSuccess(t *testing.T) {
	id, _ := json.Marshal(1)
	resp := Call(id, func() (string, error) {
		return "ok", nil
	})
	raw, _ := json.Marshal(resp.Result)
	var r proto.CallToolResult
	_ = json.Unmarshal(raw, &r)
	if r.IsError {
		t.Error("unexpected isError")
	}
}

func TestCallError(t *testing.T) {
	id, _ := json.Marshal(1)
	resp := Call(id, func() (string, error) {
		return "", errors.New("oops")
	})
	raw, _ := json.Marshal(resp.Result)
	var r proto.CallToolResult
	_ = json.Unmarshal(raw, &r)
	if !r.IsError {
		t.Error("expected isError")
	}
}
