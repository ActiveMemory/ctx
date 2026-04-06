//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hub

import (
	"testing"
	"time"
)

func TestNewStore(t *testing.T) {
	dir := t.TempDir()
	s, err := NewStore(dir)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	total, _, _ := s.Stats()
	if total != 0 {
		t.Errorf("new store should have 0 entries, got %d", total)
	}
}

func TestStoreAppendAndQuery(t *testing.T) {
	dir := t.TempDir()
	s, err := NewStore(dir)
	if err != nil {
		t.Fatal(err)
	}

	entries := []Entry{
		{ID: "a", Type: "decision", Content: "Use Go", Origin: "alpha", Timestamp: time.Now()},
		{ID: "b", Type: "learning", Content: "Avoid mocks", Origin: "beta", Timestamp: time.Now()},
		{ID: "c", Type: "decision", Content: "Use UTC", Origin: "alpha", Timestamp: time.Now()},
	}

	seqs, appendErr := s.Append(entries)
	if appendErr != nil {
		t.Fatalf("Append: %v", appendErr)
	}

	if len(seqs) != 3 {
		t.Fatalf("expected 3 sequences, got %d", len(seqs))
	}
	if seqs[0] != 1 || seqs[1] != 2 || seqs[2] != 3 {
		t.Errorf("sequences should be 1,2,3, got %v", seqs)
	}

	// Query all
	all := s.Query(nil, 0)
	if len(all) != 3 {
		t.Errorf("expected 3 entries, got %d", len(all))
	}

	// Query by type
	decisions := s.Query([]string{"decision"}, 0)
	if len(decisions) != 2 {
		t.Errorf("expected 2 decisions, got %d", len(decisions))
	}

	// Query since sequence
	since2 := s.Query(nil, 2)
	if len(since2) != 1 {
		t.Errorf("expected 1 entry after seq 2, got %d", len(since2))
	}
	if since2[0].ID != "c" {
		t.Errorf("expected entry 'c', got %q", since2[0].ID)
	}

	// Query by type + since
	decisionsSince1 := s.Query([]string{"decision"}, 1)
	if len(decisionsSince1) != 1 {
		t.Errorf("expected 1 decision after seq 1, got %d", len(decisionsSince1))
	}
}

func TestStorePersistence(t *testing.T) {
	dir := t.TempDir()

	// Write entries
	s1, err := NewStore(dir)
	if err != nil {
		t.Fatal(err)
	}
	_, appendErr := s1.Append([]Entry{
		{ID: "x", Type: "learning", Content: "Persist works", Origin: "proj", Timestamp: time.Now()},
	})
	if appendErr != nil {
		t.Fatal(appendErr)
	}

	// Reopen and verify
	s2, err := NewStore(dir)
	if err != nil {
		t.Fatal(err)
	}
	all := s2.Query(nil, 0)
	if len(all) != 1 {
		t.Fatalf("expected 1 entry after reopen, got %d", len(all))
	}
	if all[0].Content != "Persist works" {
		t.Errorf("wrong content: %q", all[0].Content)
	}
	if all[0].Sequence != 1 {
		t.Errorf("sequence should be 1, got %d", all[0].Sequence)
	}
}

func TestStoreRegisterAndValidate(t *testing.T) {
	dir := t.TempDir()
	s, err := NewStore(dir)
	if err != nil {
		t.Fatal(err)
	}

	client := ClientInfo{ID: "c1", ProjectName: "alpha", Token: "tok_abc"}
	if regErr := s.RegisterClient(client); regErr != nil {
		t.Fatal(regErr)
	}

	// Valid token
	found := s.ValidateToken("tok_abc")
	if found == nil {
		t.Fatal("expected to find client")
	}
	if found.ProjectName != "alpha" {
		t.Errorf("wrong project: %q", found.ProjectName)
	}

	// Invalid token
	if s.ValidateToken("invalid") != nil {
		t.Error("should not find client with invalid token")
	}

	// Persistence
	s2, err := NewStore(dir)
	if err != nil {
		t.Fatal(err)
	}
	found2 := s2.ValidateToken("tok_abc")
	if found2 == nil {
		t.Fatal("client should persist across reopens")
	}
}

func TestStoreStats(t *testing.T) {
	dir := t.TempDir()
	s, err := NewStore(dir)
	if err != nil {
		t.Fatal(err)
	}

	_, _ = s.Append([]Entry{
		{ID: "1", Type: "decision", Origin: "a", Timestamp: time.Now()},
		{ID: "2", Type: "decision", Origin: "b", Timestamp: time.Now()},
		{ID: "3", Type: "learning", Origin: "a", Timestamp: time.Now()},
	})

	total, byType, byProject := s.Stats()
	if total != 3 {
		t.Errorf("total: want 3, got %d", total)
	}
	if byType["decision"] != 2 {
		t.Errorf("decisions: want 2, got %d", byType["decision"])
	}
	if byType["learning"] != 1 {
		t.Errorf("learnings: want 1, got %d", byType["learning"])
	}
	if byProject["a"] != 2 {
		t.Errorf("project a: want 2, got %d", byProject["a"])
	}
	if byProject["b"] != 1 {
		t.Errorf("project b: want 1, got %d", byProject["b"])
	}
}
