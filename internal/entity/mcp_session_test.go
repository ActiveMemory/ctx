//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package entity

import (
	"testing"
	"time"
)

func TestNewMCPSession(t *testing.T) {
	s := NewMCPSession()
	if s.ToolCalls != 0 {
		t.Errorf("ToolCalls = %d, want 0", s.ToolCalls)
	}
	if s.AddsPerformed == nil {
		t.Fatal("AddsPerformed should be initialized")
	}
	if len(s.AddsPerformed) != 0 {
		t.Errorf(
			"AddsPerformed length = %d, want 0",
			len(s.AddsPerformed),
		)
	}
	if s.SessionStartedAt.IsZero() {
		t.Error("SessionStartedAt should be set")
	}
	if len(s.PendingFlush) != 0 {
		t.Errorf(
			"PendingFlush length = %d, want 0",
			len(s.PendingFlush),
		)
	}
}

func TestRecordToolCall(t *testing.T) {
	s := NewMCPSession()
	s.RecordToolCall()
	if s.ToolCalls != 1 {
		t.Errorf("ToolCalls = %d, want 1", s.ToolCalls)
	}
	s.RecordToolCall()
	s.RecordToolCall()
	if s.ToolCalls != 3 {
		t.Errorf("ToolCalls = %d, want 3", s.ToolCalls)
	}
}

func TestRecordAdd(t *testing.T) {
	s := NewMCPSession()
	s.RecordAdd("task")
	s.RecordAdd("task")
	s.RecordAdd("decision")
	if s.AddsPerformed["task"] != 2 {
		t.Errorf(
			"task adds = %d, want 2",
			s.AddsPerformed["task"],
		)
	}
	if s.AddsPerformed["decision"] != 1 {
		t.Errorf(
			"decision adds = %d, want 1",
			s.AddsPerformed["decision"],
		)
	}
}

func TestQueuePendingUpdate(t *testing.T) {
	s := NewMCPSession()
	now := time.Now()
	s.QueuePendingUpdate(PendingUpdate{
		Type:     "task",
		Content:  "Build feature",
		QueuedAt: now,
	})
	if len(s.PendingFlush) != 1 {
		t.Fatalf(
			"PendingFlush length = %d, want 1",
			len(s.PendingFlush),
		)
	}
	pu := s.PendingFlush[0]
	if pu.Type != "task" {
		t.Errorf(
			"Type = %q, want %q",
			pu.Type, "task",
		)
	}
	if pu.Content != "Build feature" {
		t.Errorf(
			"Content = %q, want %q",
			pu.Content, "Build feature",
		)
	}
}

func TestPendingCount(t *testing.T) {
	s := NewMCPSession()
	if s.PendingCount() != 0 {
		t.Errorf("PendingCount = %d, want 0", s.PendingCount())
	}
	s.QueuePendingUpdate(PendingUpdate{Type: "task", Content: "A"})
	s.QueuePendingUpdate(PendingUpdate{Type: "decision", Content: "B"})
	if s.PendingCount() != 2 {
		t.Errorf("PendingCount = %d, want 2", s.PendingCount())
	}
}

func TestRecordSessionStart(t *testing.T) {
	s := NewMCPSession()
	if s.SessionStarted {
		t.Error("SessionStarted should be false initially")
	}
	before := time.Now()
	s.RecordSessionStart()
	if !s.SessionStarted {
		t.Error("SessionStarted should be true after RecordSessionStart")
	}
	if s.SessionStartedAt.Before(before) {
		t.Error("SessionStartedAt should be updated to now or later")
	}
}

func TestRecordContextLoaded(t *testing.T) {
	s := NewMCPSession()
	if s.ContextLoaded {
		t.Error("ContextLoaded should be false initially")
	}
	s.RecordContextLoaded()
	if !s.ContextLoaded {
		t.Error("ContextLoaded should be true after RecordContextLoaded")
	}
}

func TestRecordDriftCheck(t *testing.T) {
	s := NewMCPSession()
	if !s.LastDriftCheck.IsZero() {
		t.Error("LastDriftCheck should be zero initially")
	}
	before := time.Now()
	s.RecordDriftCheck()
	if s.LastDriftCheck.Before(before) {
		t.Error("LastDriftCheck should be updated to now or later")
	}
}

func TestRecordContextWrite(t *testing.T) {
	s := NewMCPSession()
	s.CallsSinceWrite = 5
	before := time.Now()
	s.RecordContextWrite()
	if s.LastContextWrite.Before(before) {
		t.Error("LastContextWrite should be updated to now or later")
	}
	if s.CallsSinceWrite != 0 {
		t.Errorf("CallsSinceWrite = %d, want 0 after RecordContextWrite", s.CallsSinceWrite)
	}
}

func TestIncrementCallsSinceWrite(t *testing.T) {
	s := NewMCPSession()
	if s.CallsSinceWrite != 0 {
		t.Errorf("CallsSinceWrite = %d, want 0 initially", s.CallsSinceWrite)
	}
	s.IncrementCallsSinceWrite()
	if s.CallsSinceWrite != 1 {
		t.Errorf("CallsSinceWrite = %d, want 1", s.CallsSinceWrite)
	}
	s.IncrementCallsSinceWrite()
	s.IncrementCallsSinceWrite()
	if s.CallsSinceWrite != 3 {
		t.Errorf("CallsSinceWrite = %d, want 3", s.CallsSinceWrite)
	}
}
