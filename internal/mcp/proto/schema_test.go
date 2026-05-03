//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package proto

import (
	"encoding/json"
	"testing"
)

// TestToolContentTextFieldAlwaysPresent locks the wire format that
// strict MCP clients (notably OpenCode, whose Zod schema requires
// `text` present on every `type:"text"` content) depend on.
//
// The MCP spec defines `text` as required for text content. If this
// test ever flips back to omitting `text` on the empty value, OpenCode
// rejects the response with a Zod validation error before the agent
// ever sees it. See PR #72 for the live verification against Claude
// Code and Copilot CLI MCP clients.
func TestToolContentTextFieldAlwaysPresent(t *testing.T) {
	cases := []struct {
		name string
		in   ToolContent
		want string
	}{
		{
			name: "empty text still emits the key",
			in:   ToolContent{Type: "text"},
			want: `{"type":"text","text":""}`,
		},
		{
			name: "non-empty text round-trips",
			in:   ToolContent{Type: "text", Text: "hello"},
			want: `{"type":"text","text":"hello"}`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := json.Marshal(tc.in)
			if err != nil {
				t.Fatalf("json.Marshal: %v", err)
			}
			if string(got) != tc.want {
				t.Errorf("got %s, want %s", got, tc.want)
			}
		})
	}
}
