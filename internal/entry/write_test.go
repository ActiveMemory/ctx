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

// TestWrite_RefusesEntriesTrappedInIndexBlock is the regression guard for the
// data-loss bug: when entry bodies live between the INDEX markers, Write must
// refuse and leave the file byte-identical rather than regenerate the index
// and delete them.
func TestWrite_RefusesEntriesTrappedInIndexBlock(t *testing.T) {
	malformed := "# Learnings\n\n<!-- INDEX:START -->\n\n" +
		"## [2026-01-01-090000] First\n\n**Lesson:** alpha must survive.\n\n" +
		"## [2026-01-02-090000] Second\n\n**Lesson:** beta must survive.\n\n" +
		"<!-- INDEX:END -->\n"
	dir, path := seedLearnings(t, malformed)

	if err := Write(learningParams(dir)); err == nil {
		t.Fatal("Write() must refuse a LEARNINGS.md with entries trapped in the index block")
	}

	got, readErr := os.ReadFile(path) //nolint:gosec // path is test-controlled
	if readErr != nil {
		t.Fatalf("read back: %v", readErr)
	}
	if string(got) != malformed {
		t.Errorf("Write() must not modify a refused file\nGot:\n%s", got)
	}
}

// TestWrite_PreservesBodiesWellFormed confirms the guard does not regress the
// happy path: a well-formed file gains the new entry and keeps prior bodies
// and exactly one marker pair.
func TestWrite_PreservesBodiesWellFormed(t *testing.T) {
	wellFormed := "# Learnings\n\n<!-- INDEX:START -->\n<!-- INDEX:END -->\n\n" +
		"## [2026-01-01-090000] First\n\n**Lesson:** alpha must survive.\n"
	dir, path := seedLearnings(t, wellFormed)

	if err := Write(learningParams(dir)); err != nil {
		t.Fatalf("Write() on a well-formed file: %v", err)
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
	if n := strings.Count(body, "<!-- INDEX:START -->"); n != 1 {
		t.Errorf("Write() left %d INDEX:START markers, want 1\nGot:\n%s", n, body)
	}
}
