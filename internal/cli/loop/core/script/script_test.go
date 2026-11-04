//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package script

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ActiveMemory/ctx/internal/assets/read/lookup"
	cfgLoop "github.com/ActiveMemory/ctx/internal/config/loop"
)

func TestMain(m *testing.M) {
	lookup.Init()
	os.Exit(m.Run())
}

// TestGenerateMatchesLegacy asserts the text/template-based Generate
// reproduces, byte-for-byte, the output of the pre-migration
// fmt.Sprintf composition. The golden fixtures were captured from the
// legacy code path (see git history) and cover both the iteration-cap
// on/off branch and each tool's command.
func TestGenerateMatchesLegacy(t *testing.T) {
	cases := []struct {
		name    string
		tool    string
		maxIter int
	}{
		{"claude-nomax", cfgLoop.DefaultTool, 0},
		{"claude-max", cfgLoop.DefaultTool, 5},
		{"aider-nomax", cfgLoop.ToolAider, 0},
		{"generic-max", cfgLoop.ToolGeneric, 5},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			want, readErr := os.ReadFile(
				filepath.Join("testdata", c.name+".golden"),
			)
			if readErr != nil {
				t.Fatal(readErr)
			}
			got, genErr := Generate(Opts{
				PromptFile:    "/tmp/prompt.md",
				Tool:          c.tool,
				MaxIterations: c.maxIter,
				CompletionMsg: "DONE",
			})
			if genErr != nil {
				t.Fatalf("Generate: %v", genErr)
			}
			if got != string(want) {
				t.Errorf(
					"drift:\n--- want ---\n%q\n--- got ---\n%q",
					string(want), got,
				)
			}
		})
	}
}

func TestGenerate_EoCCAndSACRegimes(t *testing.T) {
	eoccOut, err := Generate(Opts{
		PromptFile:    "/tmp/prompt.md",
		Tool:          cfgLoop.DefaultTool,
		Regime:        cfgLoop.RegimeEoCC,
		MaxIterations: 10,
		CompletionMsg: "SYSTEM_CONVERGED",
	})
	if err != nil {
		t.Fatalf("Generate EoCC: %v", err)
	}
	for _, wantSub := range []string{
		`echo "Regime: eocc"`,
		`STATE_PACKET=$(ctx agent --budget 4000 2>/dev/null || true)`,
		`# End-of-Conversation Compaction (EoCC) Contract`,
		`ctx handover write --no-fold --summary`,
		`OUTPUT=$(claude --print "$ITER_PROMPT" 2>&1) || true`,
	} {
		if !containsSubstring(eoccOut, wantSub) {
			t.Errorf("EoCC script missing %q:\n%s", wantSub, eoccOut)
		}
	}

	sacOut, err := Generate(Opts{
		PromptFile:    "/tmp/prompt.md",
		Tool:          cfgLoop.DefaultTool,
		Regime:        cfgLoop.RegimeSAC,
		MaxIterations: 25,
		CompletionMsg: "SYSTEM_CONVERGED",
	})
	if err != nil {
		t.Fatalf("Generate SAC: %v", err)
	}
	for _, wantSub := range []string{
		`echo "Regime: sac"`,
		`STATE_PACKET=$(ctx agent --budget 4000 2>/dev/null || true)`,
		`# Single-Action Compaction (SAC) Contract`,
		`Execute at most ONE environment action`,
		`ctx handover write --no-fold --summary`,
		`OUTPUT=$(claude --print --max-turns 2 "$ITER_PROMPT" 2>&1) || true`,
	} {
		if !containsSubstring(sacOut, wantSub) {
			t.Errorf("SAC script missing %q:\n%s", wantSub, sacOut)
		}
	}

	// Verify aider in SAC regime uses prompt envelope without unsupported --max-turns CLI flag.
	aiderSAC, err := Generate(Opts{
		PromptFile:    "/tmp/prompt.md",
		Tool:          cfgLoop.ToolAider,
		Regime:        cfgLoop.RegimeSAC,
		CompletionMsg: "SYSTEM_CONVERGED",
	})
	if err != nil {
		t.Fatalf("Generate aider SAC: %v", err)
	}
	if containsSubstring(aiderSAC, "--max-turns") {
		t.Errorf("aider SAC script should not contain --max-turns:\n%s", aiderSAC)
	}
	if !containsSubstring(aiderSAC, `# Single-Action Compaction (SAC) Contract`) {
		t.Errorf("aider SAC script missing SAC contract:\n%s", aiderSAC)
	}
}

func containsSubstring(haystack, needle string) bool {
	for i := 0; i+len(needle) <= len(haystack); i++ {
		if haystack[i:i+len(needle)] == needle {
			return true
		}
	}
	return false
}
