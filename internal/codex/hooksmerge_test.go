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

	"github.com/ActiveMemory/ctx/internal/assets/read/agent"
	cfgCodex "github.com/ActiveMemory/ctx/internal/config/codex"
)

func embeddedHooks(t *testing.T) []byte {
	t.Helper()
	data, err := agent.CodexHooksJSON()
	if err != nil {
		t.Fatalf("CodexHooksJSON: %v", err)
	}
	return data
}

type manifest struct {
	Description string                       `json:"description"`
	Hooks       map[string][]json.RawMessage `json:"hooks"`
}

func parse(t *testing.T, data []byte) manifest {
	t.Helper()
	var m manifest
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("parse: %v\n%s", err, data)
	}
	return m
}

const foreignGroup = `{"matcher":"Bash","hooks":[{"type":"command","command":"echo foreign"}]}`

func ctxGroup(command string) string {
	quoted, _ := json.Marshal(cfgCodex.HookAnchor + command)
	return `{"matcher":"Bash","hooks":[{"type":"command","command":` +
		string(quoted) + `}]}`
}

func TestMergeHooks_EmptyExistingIsCreated(t *testing.T) {
	embedded := embeddedHooks(t)
	for _, in := range [][]byte{nil, []byte("  \n")} {
		out, outcome, err := MergeHooks(in, embedded)
		if err != nil {
			t.Fatal(err)
		}
		if outcome != OutcomeCreated {
			t.Fatalf("outcome = %v, want Created", outcome)
		}
		if !bytes.Equal(out, embedded) {
			t.Fatal("created output must be the embedded manifest verbatim")
		}
	}
}

func TestMergeHooks_SecondRunSkips(t *testing.T) {
	embedded := embeddedHooks(t)
	out, outcome, err := MergeHooks(embedded, embedded)
	if err != nil {
		t.Fatal(err)
	}
	if outcome != OutcomeSkipped {
		t.Fatalf("outcome = %v, want Skipped", outcome)
	}
	if !bytes.Equal(out, embedded) {
		t.Fatal("skipped output must be the existing bytes")
	}
}

func TestMergeHooks_PreservesForeignReplacesStale(t *testing.T) {
	embedded := embeddedHooks(t)
	existing := []byte(`{
  "description": "user description",
  "custom": {"a": 1},
  "hooks": {
    "PreToolUse": [
      ` + foreignGroup + `,
      ` + ctxGroup("ctx system stale-command") + `
    ],
    "PreCompact": [` + foreignGroup + `],
    "SessionEnd": [` + ctxGroup("ctx journal import --all -y") + `]
  }
}
`)
	out, outcome, err := MergeHooks(existing, embedded)
	if err != nil {
		t.Fatal(err)
	}
	if outcome != OutcomeMerged {
		t.Fatalf("outcome = %v, want Merged", outcome)
	}
	if bytes.Contains(out, []byte("stale-command")) {
		t.Fatal("stale ctx group survived the merge")
	}
	if !bytes.Contains(out, []byte("echo foreign")) {
		t.Fatal("foreign group lost")
	}
	if !bytes.Contains(out, []byte(`"custom": {`)) {
		t.Fatal("unrelated top-level key lost")
	}
	got := parse(t, out)
	if got.Description != "user description" {
		t.Fatalf("description overwritten: %q", got.Description)
	}
	if len(got.Hooks[cfgCodex.EventPreCompact]) != 1 {
		t.Fatal("event only in existing file was not preserved")
	}
	want := parse(t, embedded)
	for event, groups := range want.Hooks {
		foreign := 0
		if event == cfgCodex.EventPreToolUse {
			foreign = 1
		}
		if len(got.Hooks[event]) != len(groups)+foreign {
			t.Fatalf("%s: %d groups, want %d", event, len(got.Hooks[event]), len(groups)+foreign)
		}
	}
	// Foreign group is first (original order), ctx groups follow.
	first := string(got.Hooks[cfgCodex.EventPreToolUse][0])
	if !strings.Contains(first, "echo foreign") {
		t.Fatalf("foreign group not first: %s", first)
	}
	if !strings.HasSuffix(string(out), "\n") {
		t.Fatal("missing trailing newline")
	}
	if bytes.Contains(out, []byte(`\u0026`)) || bytes.Contains(out, []byte(`\u003e`)) {
		t.Fatal("HTML escaping leaked into output")
	}
	if !bytes.Contains(out, []byte(`&& ctx`)) {
		t.Fatal("anchor not written verbatim")
	}

	// Idempotent: merging the merged output again is a no-op.
	again, outcome2, err := MergeHooks(out, embedded)
	if err != nil {
		t.Fatal(err)
	}
	if outcome2 != OutcomeSkipped || !bytes.Equal(again, out) {
		t.Fatalf("second merge not a no-op: %v", outcome2)
	}
}

func TestMergeHooks_TakesEmbeddedDescriptionWhenMissing(t *testing.T) {
	embedded := embeddedHooks(t)
	out, _, err := MergeHooks([]byte(`{"hooks":{}}`), embedded)
	if err != nil {
		t.Fatal(err)
	}
	if parse(t, out).Description != parse(t, embedded).Description {
		t.Fatal("embedded description not adopted")
	}
}

func TestMergeHooks_InvalidJSON(t *testing.T) {
	embedded := embeddedHooks(t)
	for _, in := range []string{"{not json", `{"hooks": []}`, `{"hooks": {"PreToolUse": {}}}`} {
		if _, _, err := MergeHooks([]byte(in), embedded); err == nil {
			t.Fatalf("expected error for %q", in)
		}
	}
}

func TestMergeHooks_UnparseableGroupIsForeign(t *testing.T) {
	embedded := embeddedHooks(t)
	existing := []byte(`{"hooks":{"PreToolUse":[{"hooks":"not-a-list"},{"hooks":[]}]}}`)
	out, _, err := MergeHooks(existing, embedded)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(out, []byte("not-a-list")) {
		t.Fatal("odd group dropped")
	}
}

// TestMergeHooks_ForeignAnchoredGroupSurvives guards the ownership
// rule: a user group whose commands copy ctx's git-root anchor but
// do not invoke ctx must survive the merge untouched.
func TestMergeHooks_ForeignAnchoredGroupSurvives(t *testing.T) {
	existing := []byte(`{
  "hooks": {
    "PreToolUse": [
      {
        "matcher": "Bash",
        "hooks": [
          {
            "type": "command",
            "command": "cd \"$(git rev-parse --show-toplevel)\" && make lint"
          }
        ]
      }
    ]
  }
}`)
	embedded := []byte(`{
  "hooks": {
    "PreToolUse": [
      {
        "matcher": ".*",
        "hooks": [
          {
            "type": "command",
            "command": "cd \"$(git rev-parse --show-toplevel)\" && ctx system context-load-gate"
          }
        ]
      }
    ]
  }
}`)
	out, outcome, mergeErr := MergeHooks(existing, embedded)
	if mergeErr != nil {
		t.Fatalf("MergeHooks: %v", mergeErr)
	}
	if outcome != OutcomeMerged {
		t.Fatalf("outcome = %v, want OutcomeMerged", outcome)
	}
	if !strings.Contains(string(out), "make lint") {
		t.Fatalf("foreign anchored group was dropped:\n%s", out)
	}
	if !strings.Contains(string(out), "context-load-gate") {
		t.Fatalf("embedded ctx group missing:\n%s", out)
	}
}

// TestMergeHooks_LegacyAnchorGroupMigrates guards migration: a
// ctx group deployed by an earlier build (bare git-root anchor,
// no non-repo fallback) is recognized as ctx-managed and replaced
// by the current manifest instead of being duplicated.
func TestMergeHooks_LegacyAnchorGroupMigrates(t *testing.T) {
	existing := []byte(`{
  "hooks": {
    "PreToolUse": [
      {
        "matcher": ".*",
        "hooks": [
          {
            "type": "command",
            "command": "cd \"$(git rev-parse --show-toplevel)\" && ctx system context-load-gate"
          }
        ]
      }
    ]
  }
}`)
	embedded := []byte(`{
  "hooks": {
    "PreToolUse": [
      {
        "matcher": ".*",
        "hooks": [
          {
            "type": "command",
            "command": "cd \"$(git rev-parse --show-toplevel 2>/dev/null || pwd)\" && ctx system context-load-gate"
          }
        ]
      }
    ]
  }
}`)
	out, outcome, mergeErr := MergeHooks(existing, embedded)
	if mergeErr != nil {
		t.Fatalf("MergeHooks: %v", mergeErr)
	}
	if outcome != OutcomeMerged {
		t.Fatalf("outcome = %v, want OutcomeMerged", outcome)
	}
	if strings.Count(string(out), "context-load-gate") != 1 {
		t.Fatalf("legacy group not migrated (duplicated or dropped):\n%s", out)
	}
	if !strings.Contains(string(out), "|| pwd") {
		t.Fatalf("migrated group lacks tolerant anchor:\n%s", out)
	}
}

// TestMergeHooks_MixedGroupKeepsUserHandlers guards the mixed-group
// rule: a user handler added inside a deployed ctx group survives,
// the ctx handlers in that group are stripped (the fresh embedded
// groups replace them), and nothing duplicates.
func TestMergeHooks_MixedGroupKeepsUserHandlers(t *testing.T) {
	ctxHandler := map[string]string{
		"type":    "command",
		"command": cfgCodex.HookCommandPrefix + "system qa-reminder",
	}
	userHandler := map[string]string{
		"type":    "command",
		"command": "echo mine",
	}
	build := func(handlers ...map[string]string) []byte {
		doc := map[string]any{
			"hooks": map[string]any{
				"PreToolUse": []any{
					map[string]any{
						"matcher": "Bash",
						"hooks":   handlers,
					},
				},
			},
		}
		out, marshalErr := json.Marshal(doc)
		if marshalErr != nil {
			t.Fatal(marshalErr)
		}
		return out
	}

	existing := build(ctxHandler, userHandler)
	embedded := build(ctxHandler)

	out, _, mergeErr := MergeHooks(existing, embedded)
	if mergeErr != nil {
		t.Fatalf("MergeHooks: %v", mergeErr)
	}
	if strings.Count(string(out), "qa-reminder") != 1 {
		t.Fatalf("ctx handler duplicated or lost:\n%s", out)
	}
	if !strings.Contains(string(out), "echo mine") {
		t.Fatalf("user handler lost:\n%s", out)
	}
}
