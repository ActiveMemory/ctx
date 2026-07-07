//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package execute_test

import (
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/journal/core/execute"
	"github.com/ActiveMemory/ctx/internal/cli/journal/core/plan"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/entity"
	"github.com/ActiveMemory/ctx/internal/journal/state"
)

// TestImportCycle_SelfHealAndForeignEdit drives the real import
// pipeline (plan.Import → execute.Import → state round-trip) through
// the full self-healing lifecycle: a New import, an Unchanged no-op, a
// Grown re-render that self-heals, and a hand-edited entry that a later
// growth must never clobber. This is the integration seam the planner
// unit tests do not cover: execute's write, render-hash stamping, and
// state persistence across sweeps.
func TestImportCycle_SelfHealAndForeignEdit(t *testing.T) {
	jdir := t.TempDir()
	srcDir := t.TempDir()
	src := filepath.Join(srcDir, "sess.jsonl")
	writeAll(t, src, "line one\n")

	s := &entity.Session{
		ID:         "cycle-1",
		Tool:       "claude-code",
		Project:    "proj",
		SourceFile: src,
		StartTime:  time.Date(2026, 1, 15, 10, 30, 0, 0, time.UTC),
		EndTime:    time.Date(2026, 1, 15, 11, 0, 0, 0, time.UTC),
		Messages: []entity.Message{
			{Role: "user", Text: "first message"},
		},
	}
	opts := entity.ImportOpts{All: true, KeepFrontmatter: true}
	cmd := discardCmd()
	idx := map[string]string{}

	load := func() *state.State {
		st, err := state.Load(jdir)
		if err != nil {
			t.Fatalf("state load: %v", err)
		}
		return st
	}
	save := func(st *state.State) {
		if err := st.Save(jdir); err != nil {
			t.Fatalf("state save: %v", err)
		}
	}

	// --- 1. New import: renders the entry, records source + render hash.
	st := load()
	p := plan.Import([]*entity.Session{s}, jdir, idx, st, opts, false)
	if p.Actions[0].Action != entity.ActionNew {
		t.Fatalf("first run action = %v, want ActionNew", p.Actions[0].Action)
	}
	execute.Import(cmd, p, st, opts)
	entryPath := p.Actions[0].Path
	filename := p.Actions[0].Filename
	if !fileExists(entryPath) {
		t.Fatal("New import should have written the entry file")
	}
	if st.RenderHash(filename) == "" {
		t.Error("New import should record a render hash")
	}
	if _, ok := st.SessionSource("cycle-1"); !ok {
		t.Error("New import should record the source stats")
	}
	save(st)

	// --- 2. Unchanged: source untouched → skip, byte-identical file.
	before := readAll(t, entryPath)
	st = load()
	p = plan.Import([]*entity.Session{s}, jdir, idx, st, opts, false)
	if p.Actions[0].Action != entity.ActionSkip {
		t.Fatalf("unchanged action = %v, want ActionSkip", p.Actions[0].Action)
	}
	execute.Import(cmd, p, st, opts)
	if readAll(t, entryPath) != before {
		t.Error("Unchanged sweep must leave the entry byte-identical")
	}

	// --- 3. Grown, ctx-owned: source grows (size changes) and the
	// transcript now parses to more messages → re-render, self-heal.
	writeAll(t, src, "line one\nline two\nline three\n")
	s.Messages = append(s.Messages,
		entity.Message{Role: "assistant", Text: "a reply that grows the entry"},
	)
	st = load()
	p = plan.Import([]*entity.Session{s}, jdir, idx, st, opts, false)
	if p.Actions[0].Action != entity.ActionRegenerate {
		t.Fatalf("grown action = %v, want ActionRegenerate", p.Actions[0].Action)
	}
	// A grown re-render is counted as Grown, not Regen, so a routine
	// interactive sweep never trips the regenerate confirmation prompt.
	if p.GrownCount != 1 || p.RegenCount != 0 {
		t.Errorf("grown sweep: GrownCount=%d RegenCount=%d, want 1 and 0",
			p.GrownCount, p.RegenCount)
	}
	execute.Import(cmd, p, st, opts)
	grown := readAll(t, entryPath)
	if grown == before {
		t.Error("Grown re-render should update the entry body")
	}
	if st.RenderHash(filename) == "" {
		t.Error("Grown re-render should refresh the render hash")
	}
	save(st)

	// --- 4. Foreign edit: a human edits the body, then the source grows
	// again. The grown sweep must detect the edit and leave it untouched.
	edited := grown + "\n\n> hand-written note the next sweep must preserve\n"
	writeAll(t, entryPath, edited)
	writeAll(t, src, "line one\nline two\nline three\nline four\n")
	s.Messages = append(s.Messages,
		entity.Message{Role: "user", Text: "a third message"},
	)
	st = load()
	p = plan.Import([]*entity.Session{s}, jdir, idx, st, opts, false)
	if p.Actions[0].Action != entity.ActionForeignEdit {
		t.Fatalf("foreign-edit action = %v, want ActionForeignEdit",
			p.Actions[0].Action)
	}
	execute.Import(cmd, p, st, opts)
	if readAll(t, entryPath) != edited {
		t.Error("a hand-edited entry must survive a Grown sweep byte-identical")
	}
}

// TestImportCycle_FailedWriteDoesNotAdvanceSource locks in the M1
// invariant: if a grown session's write fails, its recorded source stat
// must NOT advance. Otherwise the next sweep sees grown=false and
// silently forgets the un-rendered growth (permanent for a completed
// session). The stat is committed by execute only after a clean write.
func TestImportCycle_FailedWriteDoesNotAdvanceSource(t *testing.T) {
	jdir := t.TempDir()
	srcDir := t.TempDir()
	src := filepath.Join(srcDir, "sess.jsonl")
	writeAll(t, src, "line one\n")

	s := &entity.Session{
		ID: "cycle-fail", Tool: "claude-code", Project: "proj",
		SourceFile: src,
		StartTime:  time.Date(2026, 1, 15, 10, 30, 0, 0, time.UTC),
		Messages:   []entity.Message{{Role: "user", Text: "first message"}},
	}
	opts := entity.ImportOpts{All: true, KeepFrontmatter: true}
	cmd := discardCmd()
	idx := map[string]string{}

	// New import establishes the entry and records the source stat.
	st, err := state.Load(jdir)
	if err != nil {
		t.Fatalf("state load: %v", err)
	}
	p := plan.Import([]*entity.Session{s}, jdir, idx, st, opts, false)
	execute.Import(cmd, p, st, opts)
	entryPath := p.Actions[0].Path
	recBefore, ok := st.SessionSource("cycle-fail")
	if !ok {
		t.Fatal("New import should record the source stat")
	}
	if saveErr := st.Save(jdir); saveErr != nil {
		t.Fatalf("state save: %v", saveErr)
	}

	// Grow the source so the next plan is a Regenerate.
	writeAll(t, src, "line one\nline two\nline three\n")
	s.Messages = append(s.Messages,
		entity.Message{Role: "assistant", Text: "a reply that grows the entry"},
	)
	st, err = state.Load(jdir)
	if err != nil {
		t.Fatalf("state reload: %v", err)
	}
	p = plan.Import([]*entity.Session{s}, jdir, idx, st, opts, false)
	if p.Actions[0].Action != entity.ActionRegenerate {
		t.Fatalf("grown action = %v, want ActionRegenerate", p.Actions[0].Action)
	}

	// Sabotage the write: replace the entry file with a directory so
	// SafeWriteFile fails deterministically.
	if rmErr := os.Remove(entryPath); rmErr != nil {
		t.Fatalf("remove entry: %v", rmErr)
	}
	if mkErr := os.Mkdir(entryPath, fs.PermExec); mkErr != nil {
		t.Fatalf("mkdir entry: %v", mkErr)
	}

	execute.Import(cmd, p, st, opts)

	recAfter, ok := st.SessionSource("cycle-fail")
	if !ok {
		t.Fatal("source record should still exist after a failed write")
	}
	if recAfter.SourceMtime != recBefore.SourceMtime ||
		recAfter.SourceSize != recBefore.SourceSize {
		t.Errorf("failed write advanced the source stat: before=%+v after=%+v",
			recBefore, recAfter)
	}
}

func discardCmd() *cobra.Command {
	cmd := &cobra.Command{}
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)
	return cmd
}

func writeAll(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), fs.PermFile); err != nil {
		t.Fatal(err)
	}
}

func readAll(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
