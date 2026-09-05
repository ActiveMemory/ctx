//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package codex

import (
	"strings"
	"testing"

	cfgCodex "github.com/ActiveMemory/ctx/internal/config/codex"
)

func TestMCPTable(t *testing.T) {
	want := "[mcp_servers.ctx]\ncommand = \"ctx\"\nargs = [\"mcp\", \"serve\"]\n"
	if got := MCPTable(); got != want {
		t.Fatalf("MCPTable() =\n%s\nwant\n%s", got, want)
	}
}

func TestEnsureMCPTable_Empty(t *testing.T) {
	for _, in := range []string{"", "  \n\t"} {
		out, outcome := EnsureMCPTable([]byte(in))
		if outcome != OutcomeCreated {
			t.Fatalf("outcome = %v, want Created", outcome)
		}
		if string(out) != MCPTable() {
			t.Fatalf("out = %q, want table", out)
		}
	}
}

func TestEnsureMCPTable_AppendsPreservingBytes(t *testing.T) {
	user := "# my config\nmodel = \"gpt-5\"\n\n[mcp_servers.other]\ncommand = \"other\""
	out, outcome := EnsureMCPTable([]byte(user))
	if outcome != OutcomeMerged {
		t.Fatalf("outcome = %v, want Merged", outcome)
	}
	got := string(out)
	if !strings.HasPrefix(got, user) {
		t.Fatalf("existing bytes rewritten:\n%s", got)
	}
	if !strings.HasSuffix(got, "\n\n"+MCPTable()) {
		t.Fatalf("table not appended after a blank line:\n%s", got)
	}
}

func TestEnsureMCPTable_AppendsAfterTrailingNewline(t *testing.T) {
	user := "model = \"gpt-5\"\n"
	out, _ := EnsureMCPTable([]byte(user))
	if string(out) != user+"\n"+MCPTable() {
		t.Fatalf("unexpected layout:\n%q", out)
	}
}

func TestEnsureMCPTable_SkipsWhenHeaderPresent(t *testing.T) {
	user := "model = \"x\"\n\n  " + cfgCodex.TOMLHeaderMCPCtx + "  \ncommand = \"custom\"\n"
	out, outcome := EnsureMCPTable([]byte(user))
	if outcome != OutcomeSkipped {
		t.Fatalf("outcome = %v, want Skipped", outcome)
	}
	if string(out) != user {
		t.Fatalf("skipped output must be the input bytes")
	}
}

// TestEnsureMCPTable_HeaderVariants guards the tolerant header
// matching: trailing comments and quoted key segments must all
// count as an existing registration (no duplicate table appended).
func TestEnsureMCPTable_HeaderVariants(t *testing.T) {
	variants := []string{
		"[mcp_servers.ctx]\n",
		"  [mcp_servers.ctx]  \n",
		"[mcp_servers.ctx] # registered by hand\n",
		"[mcp_servers.\"ctx\"]\n",
		"[mcp_servers.'ctx']\n",
	}
	for _, v := range variants {
		if _, outcome := EnsureMCPTable([]byte(v)); outcome != OutcomeSkipped {
			t.Errorf("variant %q: outcome %v, want Skipped", v, outcome)
		}
	}
	foreign := "[mcp_servers.ctxother]\n"
	if _, outcome := EnsureMCPTable([]byte(foreign)); outcome == OutcomeSkipped {
		t.Errorf("unrelated header %q wrongly detected", foreign)
	}
}
