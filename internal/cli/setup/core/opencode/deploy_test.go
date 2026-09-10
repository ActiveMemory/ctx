//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package opencode

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/ActiveMemory/ctx/internal/assets/read/agent"
)

func withTempProjectDir(t *testing.T) string {
	t.Helper()
	tmp := t.TempDir()
	origDir, _ := os.Getwd()
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(origDir) })
	return tmp
}

func TestDeployPlugin_RefreshesStalePlugin(t *testing.T) {
	withTempProjectDir(t)
	target := filepath.Join(".opencode", "plugins", "ctx.ts")
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(target, []byte("stale plugin"), 0o644); err != nil {
		t.Fatalf("seed plugin: %v", err)
	}

	var buf bytes.Buffer
	if err := deployPlugin(testCmd(&buf)); err != nil {
		t.Fatalf("deployPlugin: %v", err)
	}

	files, err := agent.OpenCodePlugin()
	if err != nil {
		t.Fatalf("OpenCodePlugin: %v", err)
	}
	got, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("read plugin: %v", err)
	}
	if !bytes.Equal(got, files["index.ts"]) {
		t.Fatalf("plugin not refreshed")
	}
	if bytes.Contains(buf.Bytes(), []byte("skipped")) {
		t.Fatalf("expected refresh, got skipped output %q", buf.String())
	}
}

func TestDeploySkills_RefreshesStaleSkill(t *testing.T) {
	withTempProjectDir(t)
	target := filepath.Join(".opencode", "skills", "ctx-agent", "SKILL.md")
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(target, []byte("stale skill"), 0o644); err != nil {
		t.Fatalf("seed skill: %v", err)
	}

	var buf bytes.Buffer
	if err := deploySkills(testCmd(&buf)); err != nil {
		t.Fatalf("deploySkills: %v", err)
	}

	skills, err := agent.OpenCodeSkills()
	if err != nil {
		t.Fatalf("OpenCodeSkills: %v", err)
	}
	got, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("read skill: %v", err)
	}
	if !bytes.Equal(got, skills["ctx-agent"]) {
		t.Fatalf("skill not refreshed")
	}
	if bytes.Contains(buf.Bytes(), []byte("skipped")) {
		t.Fatalf("expected refresh, got skipped output %q", buf.String())
	}
}

func TestDeploySkills_DeploysReferences(t *testing.T) {
	withTempProjectDir(t)

	var buf bytes.Buffer
	if err := deploySkills(testCmd(&buf)); err != nil {
		t.Fatalf("deploySkills: %v", err)
	}

	refs, err := agent.OpenCodeSkillReferences()
	if err != nil {
		t.Fatalf("OpenCodeSkillReferences: %v", err)
	}
	if len(refs) == 0 {
		t.Fatal("no embedded reference files — embed glob regressed")
	}
	for name, files := range refs {
		for refName, want := range files {
			target := filepath.Join(
				".opencode", "skills", name, "references", refName,
			)
			got, readErr := os.ReadFile(target)
			if readErr != nil {
				t.Errorf("%s: not deployed: %v", target, readErr)
				continue
			}
			if !bytes.Equal(got, want) {
				t.Errorf("%s: deployed content differs from embedded", target)
			}
		}
	}
}

func TestDeploySkills_RefreshesStaleReferenceWhenSkillUnchanged(t *testing.T) {
	withTempProjectDir(t)

	skills, err := agent.OpenCodeSkills()
	if err != nil {
		t.Fatalf("OpenCodeSkills: %v", err)
	}
	refs, err := agent.OpenCodeSkillReferences()
	if err != nil {
		t.Fatalf("OpenCodeSkillReferences: %v", err)
	}
	const skill = "ctx-humanize"
	if refs[skill] == nil {
		t.Fatalf("%s has no embedded references — pick another fixture", skill)
	}

	// Seed the skill's SKILL.md with current content (deploy will skip
	// it) and one stale reference: the skip path must still refresh
	// references, or a re-run after a reference-only upstream change
	// would leave deployed skills citing stale files.
	skillDir := filepath.Join(".opencode", "skills", skill)
	refDir := filepath.Join(skillDir, "references")
	if mkErr := os.MkdirAll(refDir, 0o755); mkErr != nil {
		t.Fatalf("mkdir: %v", mkErr)
	}
	target := filepath.Join(skillDir, "SKILL.md")
	if seedErr := os.WriteFile(target, skills[skill], 0o644); seedErr != nil {
		t.Fatalf("seed skill: %v", seedErr)
	}
	var refName string
	for n := range refs[skill] {
		refName = n
		break
	}
	refTarget := filepath.Join(refDir, refName)
	if seedRefErr := os.WriteFile(refTarget, []byte("stale reference"), 0o644); seedRefErr != nil {
		t.Fatalf("seed reference: %v", seedRefErr)
	}

	var buf bytes.Buffer
	if deployErr := deploySkills(testCmd(&buf)); deployErr != nil {
		t.Fatalf("deploySkills: %v", deployErr)
	}

	got, readErr := os.ReadFile(refTarget)
	if readErr != nil {
		t.Fatalf("read reference: %v", readErr)
	}
	if !bytes.Equal(got, refs[skill][refName]) {
		t.Fatalf("reference not refreshed on the skill-skipped path")
	}
}

func TestDeployPlugin_RejectsSymlinkTarget(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink behavior varies on Windows in this environment")
	}
	withTempProjectDir(t)
	target := filepath.Join(".opencode", "plugins", "ctx.ts")
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	realFile := filepath.Join(t.TempDir(), "outside.ts")
	if err := os.WriteFile(realFile, []byte("secret"), 0o644); err != nil {
		t.Fatalf("seed real file: %v", err)
	}
	if err := os.Symlink(realFile, target); err != nil {
		t.Fatalf("symlink: %v", err)
	}

	if err := deployPlugin(testCmd(&bytes.Buffer{})); err == nil {
		t.Fatal("expected symlink rejection, got nil")
	}
}

func TestDeploySkills_RejectsNonRegularTarget(t *testing.T) {
	withTempProjectDir(t)
	target := filepath.Join(".opencode", "skills", "ctx-agent", "SKILL.md")
	if err := os.MkdirAll(target, 0o755); err != nil {
		t.Fatalf("mkdir target dir: %v", err)
	}

	if err := deploySkills(testCmd(&bytes.Buffer{})); err == nil {
		t.Fatal("expected non-regular target rejection, got nil")
	}
}
