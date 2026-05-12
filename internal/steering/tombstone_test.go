//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package steering

import "testing"

func TestHasTombstone_PresentReturnsTrue(t *testing.T) {
	body := "# Product Context\n\n" + Tombstone + "\n\nbody text"
	if !HasTombstone(body) {
		t.Errorf("expected HasTombstone to return true for body containing the marker")
	}
}

func TestHasTombstone_AbsentReturnsFalse(t *testing.T) {
	body := "# Product Context\n\nWe build a thing for AI coding sessions.\n"
	if HasTombstone(body) {
		t.Errorf("expected HasTombstone to return false for body without the marker")
	}
}

func TestHasTombstone_EmptyReturnsFalse(t *testing.T) {
	if HasTombstone("") {
		t.Errorf("expected HasTombstone to return false for empty body")
	}
}

func TestHasTombstone_PartialMatchReturnsFalse(t *testing.T) {
	// A prefix substring of the tombstone must NOT trigger detection.
	// This guards against accidental matches when a user happens to
	// write a similar but distinct comment in their own steering body.
	body := "# Project Structure\n\n<!-- remove this -->\n\nbody text"
	if HasTombstone(body) {
		t.Errorf("expected HasTombstone to return false for a partial-match comment")
	}
}

func TestHasTombstone_TombstoneAnywhereInBodyReturnsTrue(t *testing.T) {
	// The marker may appear anywhere in the body; detection is a
	// strings.Contains check, not a position-anchored one.
	body := "# Workflow\n\nProse line.\n\nMore prose.\n\n" + Tombstone + "\n"
	if !HasTombstone(body) {
		t.Errorf("expected HasTombstone to return true regardless of marker position")
	}
}
