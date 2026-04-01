//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package trace

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWorkingRefsInProgressTasks(t *testing.T) {
	contextDir := t.TempDir()

	tasks := `# Tasks

- [ ] First pending task
- [x] Completed task
- [ ] Second pending task
`
	if err := os.WriteFile(filepath.Join(contextDir, "TASKS.md"), []byte(tasks), 0o600); err != nil {
		t.Fatalf("WriteFile() error: %v", err)
	}

	refs := WorkingRefs(contextDir)

	found := map[string]bool{}
	for _, r := range refs {
		found[r] = true
	}

	if !found["task:1"] {
		t.Errorf("expected task:1 in refs %v", refs)
	}
	if !found["task:2"] {
		t.Errorf("expected task:2 in refs %v", refs)
	}
	if found["task:3"] {
		t.Errorf("did not expect task:3 in refs %v (completed tasks should not appear)", refs)
	}
}

func TestWorkingRefsSessionEnv(t *testing.T) {
	contextDir := t.TempDir()

	t.Setenv("CTX_SESSION_ID", "test-session-42")

	refs := WorkingRefs(contextDir)

	found := false
	for _, r := range refs {
		if r == "session:test-session-42" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected session:test-session-42 in refs %v", refs)
	}
}

func TestWorkingRefsNoTasksFile(t *testing.T) {
	contextDir := t.TempDir()

	// Should not panic when TASKS.md is absent.
	refs := WorkingRefs(contextDir)
	_ = refs
}
