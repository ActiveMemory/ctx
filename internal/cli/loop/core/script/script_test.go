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
			got, genErr := Generate("/tmp/prompt.md", c.tool, c.maxIter, "DONE")
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
