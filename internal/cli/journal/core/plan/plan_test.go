//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package plan

import (
	"os"
	"testing"
	"time"

	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/entity"
	"github.com/ActiveMemory/ctx/internal/journal/state"
)

// session builds a minimal, non-empty session pointing at a real
// source transcript so growth detection can stat it.
func buildSession(t *testing.T, id, sourceFile string) *entity.Session {
	t.Helper()
	return &entity.Session{
		ID:           id,
		Tool:         "claude-code",
		Project:      "proj",
		SourceFile:   sourceFile,
		StartTime:    time.Date(2026, 1, 15, 10, 30, 0, 0, time.UTC),
		EndTime:      time.Date(2026, 1, 15, 11, 0, 0, 0, time.UTC),
		Duration:     30 * time.Minute,
		TurnCount:    1,
		FirstUserMsg: "Hello there",
		Messages: []entity.Message{
			{
				Role: "user", Text: "Hello there",
				Timestamp: time.Date(2026, 1, 15, 10, 30, 0, 0, time.UTC),
			},
		},
	}
}

// writeSource writes a transcript file and returns its path and stat.
func writeSource(t *testing.T, path, content string) (mtime, size int64) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), fs.PermFile); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	return info.ModTime().Unix(), info.Size()
}

func emptyIndex() map[string]string { return map[string]string{} }

func keepOpts() entity.ImportOpts {
	return entity.ImportOpts{All: true, KeepFrontmatter: true}
}

// planOnce runs Import and returns the single-part plan (sessions here
// have one message, so exactly one action).
func planOnce(
	s *entity.Session, journalDir string, st *state.State,
) entity.FileAction {
	p := Import(
		[]*entity.Session{s}, journalDir, emptyIndex(), st, keepOpts(), false,
	)
	return p.Actions[0]
}

func TestImport_NewSession(t *testing.T) {
	dir := t.TempDir()
	src := dir + "/sess-new.jsonl"
	writeSource(t, src, "line one\n")
	s := buildSession(t, "sess-new", src)

	st := &state.State{
		Version:  state.CurrentVersion,
		Entries:  map[string]state.File{},
		Sessions: map[string]state.Source{},
	}

	p := Import(
		[]*entity.Session{s}, dir, emptyIndex(), st, keepOpts(), false,
	)
	fa := p.Actions[0]
	if fa.Action != entity.ActionNew {
		t.Errorf("Action = %v, want ActionNew", fa.Action)
	}
	// Plan stashes the source observation so execute can commit it after
	// the write succeeds; it must NOT be committed to state at plan time
	// (a failed write would then strand the un-rendered growth).
	if _, ok := p.Sources["sess-new"]; !ok {
		t.Error("source observation should be stashed in the plan")
	}
	if _, ok := st.SessionSource("sess-new"); ok {
		t.Error("plan must not commit source stats to state; execute does")
	}
}

func TestImport_Unchanged(t *testing.T) {
	dir := t.TempDir()
	src := dir + "/sess-u.jsonl"
	mtime, size := writeSource(t, src, "line one\n")
	s := buildSession(t, "sess-u", src)

	// First: discover the output filename/path via a New plan.
	st := &state.State{
		Version:  state.CurrentVersion,
		Entries:  map[string]state.File{},
		Sessions: map[string]state.Source{},
	}
	fa := planOnce(s, dir, st)

	// Simulate the entry existing and the source already recorded.
	if err := os.WriteFile(fa.Path, []byte("body"), fs.PermFile); err != nil {
		t.Fatal(err)
	}
	st.Sessions["sess-u"] = state.Source{
		SourceFile: src, SourceMtime: mtime, SourceSize: size,
	}

	got := planOnce(s, dir, st)
	if got.Action != entity.ActionSkip {
		t.Errorf("Action = %v, want ActionSkip (unchanged)", got.Action)
	}
}

func TestImport_GrownCtxOwned(t *testing.T) {
	dir := t.TempDir()
	src := dir + "/sess-g.jsonl"
	mtime, size := writeSource(t, src, "line one\nline two\n")
	s := buildSession(t, "sess-g", src)

	st := &state.State{
		Version:  state.CurrentVersion,
		Entries:  map[string]state.File{},
		Sessions: map[string]state.Source{},
	}
	fa := planOnce(s, dir, st)

	// The file exists with a body ctx wrote (matching render hash), and
	// the recorded source stats are older than the current source.
	body := "the rendered body"
	if err := os.WriteFile(fa.Path, []byte(body), fs.PermFile); err != nil {
		t.Fatal(err)
	}
	st.Entries[fa.Filename] = state.File{
		Exported: "2026-01-15", RenderHash: state.HashRender(body),
	}
	st.Sessions["sess-g"] = state.Source{
		SourceFile: src, SourceMtime: mtime - 100, SourceSize: size - 5,
	}

	got := planOnce(s, dir, st)
	if got.Action != entity.ActionRegenerate {
		t.Errorf("Action = %v, want ActionRegenerate (grown, ctx-owned)",
			got.Action)
	}
}

func TestImport_GrownForeignEdit(t *testing.T) {
	dir := t.TempDir()
	src := dir + "/sess-f.jsonl"
	mtime, size := writeSource(t, src, "line one\nline two\n")
	s := buildSession(t, "sess-f", src)

	st := &state.State{
		Version:  state.CurrentVersion,
		Entries:  map[string]state.File{},
		Sessions: map[string]state.Source{},
	}
	fa := planOnce(s, dir, st)

	// The body on disk does NOT match the recorded hash: a human edited it.
	if err := os.WriteFile(
		fa.Path, []byte("HAND EDITED body"), fs.PermFile,
	); err != nil {
		t.Fatal(err)
	}
	st.Entries[fa.Filename] = state.File{
		Exported: "2026-01-15", RenderHash: state.HashRender("original body"),
	}
	st.Sessions["sess-f"] = state.Source{
		SourceFile: src, SourceMtime: mtime - 100, SourceSize: size - 5,
	}

	got := planOnce(s, dir, st)
	if got.Action != entity.ActionForeignEdit {
		t.Errorf("Action = %v, want ActionForeignEdit (grown, edited)",
			got.Action)
	}
}

func TestImport_GrownNoHashTreatedAsEdited(t *testing.T) {
	dir := t.TempDir()
	src := dir + "/sess-p.jsonl"
	mtime, size := writeSource(t, src, "line one\nline two\n")
	s := buildSession(t, "sess-p", src)

	st := &state.State{
		Version:  state.CurrentVersion,
		Entries:  map[string]state.File{},
		Sessions: map[string]state.Source{},
	}
	fa := planOnce(s, dir, st)

	// Existing file, but no recorded render hash (pre-v2 entry).
	if err := os.WriteFile(fa.Path, []byte("body"), fs.PermFile); err != nil {
		t.Fatal(err)
	}
	st.Entries[fa.Filename] = state.File{Exported: "2026-01-15"}
	st.Sessions["sess-p"] = state.Source{
		SourceFile: src, SourceMtime: mtime - 100, SourceSize: size - 5,
	}

	got := planOnce(s, dir, st)
	if got.Action != entity.ActionForeignEdit {
		t.Errorf("Action = %v, want ActionForeignEdit (no hash → treat edited)",
			got.Action)
	}
}

func TestImport_AdoptV1Session(t *testing.T) {
	dir := t.TempDir()
	src := dir + "/sess-a.jsonl"
	writeSource(t, src, "line one\n")
	s := buildSession(t, "sess-a", src)

	st := &state.State{
		Version:  state.CurrentVersion,
		Entries:  map[string]state.File{},
		Sessions: map[string]state.Source{},
	}
	fa := planOnce(s, dir, st)

	// File exists (imported under v1) but the session is untracked.
	if err := os.WriteFile(fa.Path, []byte("body"), fs.PermFile); err != nil {
		t.Fatal(err)
	}
	// Fresh state: session NOT in Sessions map.
	st2 := &state.State{
		Version:  state.CurrentVersion,
		Entries:  map[string]state.File{},
		Sessions: map[string]state.Source{},
	}

	p := Import(
		[]*entity.Session{s}, dir, emptyIndex(), st2, keepOpts(), false,
	)
	got := p.Actions[0]
	if got.Action != entity.ActionSkip {
		t.Errorf("Action = %v, want ActionSkip (adopt)", got.Action)
	}
	// Adoption stashes the source observation in the plan (execute commits
	// it) without re-rendering.
	if _, ok := p.Sources["sess-a"]; !ok {
		t.Error("adoption should stash source stats in the plan")
	}
	// A pre-v2 entry is adopted as ctx-owned: its render hash is
	// backfilled now so a later growth sweep can self-heal it instead of
	// mistaking the absent hash for a hand edit.
	if st2.RenderHash(got.Filename) == "" {
		t.Error("adoption should backfill the render hash")
	}
}

func TestImport_UnreadableSourceDegradesToSkip(t *testing.T) {
	dir := t.TempDir()
	src := dir + "/gone.jsonl" // never created
	s := buildSession(t, "sess-x", src)

	st := &state.State{
		Version:  state.CurrentVersion,
		Entries:  map[string]state.File{},
		Sessions: map[string]state.Source{},
	}
	// Discover filename with a throwaway source so the plan can run.
	tmpSrc := dir + "/tmp.jsonl"
	writeSource(t, tmpSrc, "x\n")
	probe := buildSession(t, "sess-x", tmpSrc)
	fa := planOnce(probe, dir, &state.State{
		Version: state.CurrentVersion, Entries: map[string]state.File{},
		Sessions: map[string]state.Source{},
	})

	// The entry exists and the session is "seen", but the source is gone.
	if err := os.WriteFile(fa.Path, []byte("body"), fs.PermFile); err != nil {
		t.Fatal(err)
	}
	st.Sessions["sess-x"] = state.Source{
		SourceFile: src, SourceMtime: 1, SourceSize: 1,
	}

	got := planOnce(s, dir, st)
	if got.Action != entity.ActionSkip {
		t.Errorf("Action = %v, want ActionSkip (unreadable source)", got.Action)
	}
}
