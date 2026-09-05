//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package codex

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/codex"
	cfgCodex "github.com/ActiveMemory/ctx/internal/config/codex"
	cfgSetup "github.com/ActiveMemory/ctx/internal/config/setup"
)

func mkdirAll(path string) error {
	return os.MkdirAll(path, 0o755)
}

func TestEnsureMCPConfig_CreatesFile(t *testing.T) {
	withTempProjectDir(t)

	var buf bytes.Buffer
	if err := ensureMCPConfig(testCmd(&buf)); err != nil {
		t.Fatalf("ensureMCPConfig: %v", err)
	}
	if string(readFile(t, cfgSetup.MCPConfigPathCodex)) != codex.MCPTable() {
		t.Fatal("created config.toml is not the MCP table")
	}
	if !strings.Contains(buf.String(), "✓ "+cfgSetup.MCPConfigPathCodex) {
		t.Fatalf("expected created line:\n%s", buf.String())
	}
}

func TestEnsureMCPConfig_AppendsPreservingUserBytes(t *testing.T) {
	withTempProjectDir(t)
	user := "# team config\nmodel = \"gpt-5\"\n[mcp_servers.other]\ncommand = \"other\"\n"
	seedFile(t, cfgSetup.MCPConfigPathCodex, []byte(user))

	var buf bytes.Buffer
	if err := ensureMCPConfig(testCmd(&buf)); err != nil {
		t.Fatalf("ensureMCPConfig: %v", err)
	}
	got := string(readFile(t, cfgSetup.MCPConfigPathCodex))
	if !strings.HasPrefix(got, user) {
		t.Fatalf("user bytes rewritten:\n%s", got)
	}
	if !strings.HasSuffix(got, "\n"+codex.MCPTable()) {
		t.Fatalf("table not appended:\n%s", got)
	}
	if !strings.Contains(buf.String(), cfgSetup.MCPConfigPathCodex+" (merged)") {
		t.Fatalf("expected merged line:\n%s", buf.String())
	}
}

func TestEnsureMCPConfig_SkipsWhenHeaderPresent(t *testing.T) {
	withTempProjectDir(t)
	user := cfgCodex.TOMLHeaderMCPCtx + "\ncommand = \"/custom/ctx\"\n"
	seedFile(t, cfgSetup.MCPConfigPathCodex, []byte(user))

	var buf bytes.Buffer
	if err := ensureMCPConfig(testCmd(&buf)); err != nil {
		t.Fatalf("ensureMCPConfig: %v", err)
	}
	if string(readFile(t, cfgSetup.MCPConfigPathCodex)) != user {
		t.Fatal("user-owned table body was touched")
	}
	if !strings.Contains(buf.String(), "(up to date, skipped)") {
		t.Fatalf("expected skip line:\n%s", buf.String())
	}
}

func TestEnsureMCPConfig_RejectsDirectoryTarget(t *testing.T) {
	withTempProjectDir(t)
	if err := mkdirAll(cfgSetup.MCPConfigPathCodex); err != nil {
		t.Fatal(err)
	}
	if err := ensureMCPConfig(testCmd(&bytes.Buffer{})); err == nil {
		t.Fatal("expected non-regular target rejection")
	}
}
