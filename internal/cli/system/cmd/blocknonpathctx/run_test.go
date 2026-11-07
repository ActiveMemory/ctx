//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package blocknonpathctx

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	coreSession "github.com/ActiveMemory/ctx/internal/cli/system/core/session"
	"github.com/ActiveMemory/ctx/internal/config/agy"
)

func writeTempFile(t *testing.T, content string) *os.File {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "hook-stdin-*.json")
	if err != nil {
		t.Fatalf("CreateTemp: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("WriteString: %v", err)
	}
	if _, err := f.Seek(0, 0); err != nil {
		t.Fatalf("Seek: %v", err)
	}
	return f
}

func TestRun_AgyAllowAndDeny(t *testing.T) {
	defer coreSession.SetAgyMode(false, "", "", "")

	// 1. Deny case: relative path invocation "./ctx task list"
	denyJSON := `{
		"conversationId": "sess-deny",
		"toolCall": {
			"name": "run_command",
			"args": {
				"CommandLine": "./ctx task list"
			}
		}
	}`
	denyFile := writeTempFile(t, denyJSON)
	defer denyFile.Close()

	cmd := &cobra.Command{}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := Run(cmd, denyFile)
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	out := strings.TrimSpace(buf.String())
	var denyResp struct {
		Decision string `json:"decision"`
		Reason   string `json:"reason"`
	}
	if unmarshalErr := json.Unmarshal([]byte(out), &denyResp); unmarshalErr != nil {
		t.Fatalf("expected valid JSON deny output, got %q (err: %v)", out, unmarshalErr)
	}
	if denyResp.Decision != agy.DecisionDeny {
		t.Errorf("decision = %q, want %q", denyResp.Decision, agy.DecisionDeny)
	}
	if denyResp.Reason == "" {
		t.Error("expected non-empty reason on deny")
	}

	// 2. Allow case: legitimate PATH invocation "ctx task list"
	coreSession.SetAgyMode(false, "", "", "")
	allowJSON := `{
		"conversationId": "sess-allow",
		"toolCall": {
			"name": "run_command",
			"args": {
				"CommandLine": "ctx task list"
			}
		}
	}`
	allowFile := writeTempFile(t, allowJSON)
	defer allowFile.Close()

	buf.Reset()
	cmd.SetOut(&buf)

	err = Run(cmd, allowFile)
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	out = strings.TrimSpace(buf.String())
	if out != `{"decision":"allow"}` {
		t.Errorf("allow output = %q, want %q", out, `{"decision":"allow"}`)
	}

	// 3. Bail case: empty command -> emits {"decision":"allow"}
	coreSession.SetAgyMode(false, "", "", "")
	emptyJSON := `{
		"conversationId": "sess-empty",
		"toolCall": {
			"name": "run_command",
			"args": {
				"CommandLine": ""
			}
		}
	}`
	emptyFile := writeTempFile(t, emptyJSON)
	defer emptyFile.Close()

	buf.Reset()
	cmd.SetOut(&buf)

	err = Run(cmd, emptyFile)
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	out = strings.TrimSpace(buf.String())
	if out != `{"decision":"allow"}` {
		t.Errorf("empty command output = %q, want %q", out, `{"decision":"allow"}`)
	}
}
