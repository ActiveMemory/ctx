//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package entry

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/config/ctx"
	cfgEntry "github.com/ActiveMemory/ctx/internal/config/entry"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/entity"
)

// seedLearnings writes content to a temp LEARNINGS.md and returns its dir/path.
func seedLearnings(t *testing.T, content string) (dir, path string) {
	t.Helper()
	dir = t.TempDir()
	path = filepath.Join(dir, ctx.Learning)
	if err := os.WriteFile(path, []byte(content), fs.PermFile); err != nil {
		t.Fatalf("seed LEARNINGS.md: %v", err)
	}
	return dir, path
}

func learningParams(dir string) entity.EntryParams {
	return entity.EntryParams{
		Type:        cfgEntry.Learning,
		Content:     "New learning",
		Context:     "ctx",
		Lesson:      "lesson",
		Application: "apply",
		ContextDir:  dir,
	}
}

// TestWrite_AppendsWithoutIndexBlock confirms Write appends the new entry,
// preserves prior bodies, and never writes an INDEX block: the index is now
// projected on demand by `ctx index`, not stored in the file.
func TestWrite_AppendsWithoutIndexBlock(t *testing.T) {
	seed := "# Learnings\n\n" +
		"## [2026-01-01-090000] First\n\n**Lesson:** alpha must survive.\n"
	dir, path := seedLearnings(t, seed)

	if err := Write(learningParams(dir)); err != nil {
		t.Fatalf("Write() on a plain file: %v", err)
	}

	got, readErr := os.ReadFile(path) //nolint:gosec // path is test-controlled
	if readErr != nil {
		t.Fatalf("read back: %v", readErr)
	}
	body := string(got)
	if !strings.Contains(body, "alpha must survive.") {
		t.Errorf("Write() dropped an existing body\nGot:\n%s", body)
	}
	if !strings.Contains(body, "New learning") {
		t.Errorf("Write() did not add the new entry\nGot:\n%s", body)
	}
	if strings.Contains(body, "<!-- INDEX:START -->") {
		t.Errorf("Write() wrote an INDEX block; none should be created\nGot:\n%s", body)
	}
}
