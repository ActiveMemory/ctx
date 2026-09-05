//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package codex

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	cfgCodex "github.com/ActiveMemory/ctx/internal/config/codex"
	cfgSetup "github.com/ActiveMemory/ctx/internal/config/setup"
)

func TestDeployHooks_MergesForeignAndStaleGroups(t *testing.T) {
	withTempProjectDir(t)
	stale, _ := json.Marshal(cfgCodex.HookAnchor + "ctx system old-hook")
	seedFile(t, cfgSetup.HooksPathCodex, []byte(`{
  "description": "team hooks",
  "hooks": {
    "PreToolUse": [
      {"matcher": "Bash", "hooks": [{"type": "command", "command": "echo team"}]},
      {"matcher": "Bash", "hooks": [{"type": "command", "command": `+string(stale)+`}]}
    ],
    "Stop": [
      {"hooks": [{"type": "command", "command": "echo bye"}]}
    ]
  }
}
`))

	var buf bytes.Buffer
	if _, err := deployHooks(testCmd(&buf)); err != nil {
		t.Fatalf("deployHooks: %v", err)
	}
	got := readFile(t, cfgSetup.HooksPathCodex)
	var m struct {
		Description string                       `json:"description"`
		Hooks       map[string][]json.RawMessage `json:"hooks"`
	}
	if err := json.Unmarshal(got, &m); err != nil {
		t.Fatalf("merged file invalid: %v\n%s", err, got)
	}
	if m.Description != "team hooks" {
		t.Fatalf("description lost: %q", m.Description)
	}
	if !bytes.Contains(got, []byte("echo team")) || !bytes.Contains(got, []byte("echo bye")) {
		t.Fatal("foreign groups lost")
	}
	if bytes.Contains(got, []byte("old-hook")) {
		t.Fatal("stale ctx group survived")
	}
	if len(m.Hooks[cfgCodex.EventSessionStart]) == 0 ||
		len(m.Hooks[cfgCodex.EventUserPromptSubmit]) == 0 {
		t.Fatal("embedded events missing after merge")
	}
	if !strings.Contains(buf.String(), cfgSetup.HooksPathCodex+" (merged)") {
		t.Fatalf("expected merged line:\n%s", buf.String())
	}

	// Second run: skipped and byte-identical.
	buf.Reset()
	if _, err := deployHooks(testCmd(&buf)); err != nil {
		t.Fatalf("second deployHooks: %v", err)
	}
	if !bytes.Equal(got, readFile(t, cfgSetup.HooksPathCodex)) {
		t.Fatal("merged file rewritten on second run")
	}
	if !strings.Contains(buf.String(), "(up to date, skipped)") {
		t.Fatalf("expected skip line:\n%s", buf.String())
	}
}

func TestDeployHooks_InvalidJSONLeftUntouched(t *testing.T) {
	withTempProjectDir(t)
	broken := []byte("{not json")
	seedFile(t, cfgSetup.HooksPathCodex, broken)

	var buf bytes.Buffer
	if _, err := deployHooks(testCmd(&buf)); err != nil {
		t.Fatalf("deployHooks must warn, not fail: %v", err)
	}
	if !bytes.Equal(readFile(t, cfgSetup.HooksPathCodex), broken) {
		t.Fatal("invalid hooks.json was modified")
	}
	out := buf.String()
	if !strings.Contains(out, "! "+cfgSetup.HooksPathCodex+": ") {
		t.Fatalf("expected a warning naming the file:\n%s", out)
	}
}

func TestDeploy_InvalidHooksStillDeploysRest(t *testing.T) {
	withTempProjectDir(t)
	seedFile(t, cfgSetup.HooksPathCodex, []byte("{not json"))

	if err := Deploy(testCmd(&bytes.Buffer{})); err != nil {
		t.Fatalf("Deploy: %v", err)
	}
	if _, err := readFile(t, cfgSetup.MCPConfigPathCodex), error(nil); err != nil {
		t.Fatal(err)
	}
	readFile(t, skillPath("ctx-agent"))
}

func TestDeployHooks_RejectsDirectoryTarget(t *testing.T) {
	withTempProjectDir(t)
	if err := mkdirAll(cfgSetup.HooksPathCodex); err != nil {
		t.Fatal(err)
	}
	if _, err := deployHooks(testCmd(&bytes.Buffer{})); err == nil {
		t.Fatal("expected non-regular target rejection")
	}
}
