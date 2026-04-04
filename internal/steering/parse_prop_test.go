//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package steering

import (
	cfgSteering "github.com/ActiveMemory/ctx/internal/config/steering"
	"math/rand"
	"reflect"
	"strings"
	"testing"
	"testing/quick"
)

// validSteeringFile is a wrapper around SteeringFile that implements
// quick.Generator to produce valid inputs for property testing.
type validSteeringFile struct {
	Name        string
	Description string
	Inclusion   cfgSteering.InclusionMode
	Tools       []string
	Priority    int
	Body        string
}

var inclusionModes = []cfgSteering.InclusionMode{cfgSteering.InclusionAlways, cfgSteering.InclusionAuto, cfgSteering.InclusionManual}

var validTools = []string{"claude", "cursor", "cline", "kiro", "codex"}

// Generate implements quick.Generator for validSteeringFile.
func (validSteeringFile) Generate(rand *rand.Rand, size int) reflect.Value {
	v := validSteeringFile{
		Name:      randAlphaName(rand, size),
		Inclusion: inclusionModes[rand.Intn(len(inclusionModes))],
		Priority:  rand.Intn(99) + 1, // 1–99, avoids 0 (which triggers default)
	}

	// Optional description.
	if rand.Intn(2) == 0 {
		v.Description = randSafeString(rand, size)
	}

	// Optional tools subset.
	if rand.Intn(2) == 0 {
		n := rand.Intn(len(validTools)) + 1
		perm := rand.Perm(len(validTools))
		v.Tools = make([]string, n)
		for i := 0; i < n; i++ {
			v.Tools[i] = validTools[perm[i]]
		}
	}

	// Optional body — must not contain frontmatter delimiter on its own line.
	if rand.Intn(2) == 0 {
		v.Body = randSafeBody(rand, size)
	}

	return reflect.ValueOf(v)
}

// TestProperty_RoundTripConsistency verifies the round-trip property:
// Parse(Print(Parse(data))) == Parse(data) for all valid inputs.
//
// **Validates: Requirements 1.8, 19.1**
func TestProperty_RoundTripConsistency(t *testing.T) {
	const filePath = "test.md"

	f := func(v validSteeringFile) bool {
		// Build a SteeringFile from the generated values.
		sf := &SteeringFile{
			Name:        v.Name,
			Description: v.Description,
			Inclusion:   v.Inclusion,
			Tools:       v.Tools,
			Priority:    v.Priority,
			Body:        v.Body,
		}

		// First trip: Print → Parse.
		printed := Print(sf)
		parsed1, err := Parse(printed, filePath)
		if err != nil {
			t.Logf("first Parse failed: %v\ninput bytes:\n%s", err, printed)
			return false
		}

		// Second trip: Print → Parse again.
		printed2 := Print(parsed1)
		parsed2, err := Parse(printed2, filePath)
		if err != nil {
			t.Logf("second Parse failed: %v\ninput bytes:\n%s", err, printed2)
			return false
		}

		// Structural equality (ignoring Path, which is set from filePath arg).
		if parsed1.Name != parsed2.Name {
			t.Logf("Name mismatch: %q vs %q", parsed1.Name, parsed2.Name)
			return false
		}
		if parsed1.Description != parsed2.Description {
			t.Logf("Description mismatch: %q vs %q", parsed1.Description, parsed2.Description)
			return false
		}
		if parsed1.Inclusion != parsed2.Inclusion {
			t.Logf("Inclusion mismatch: %q vs %q", parsed1.Inclusion, parsed2.Inclusion)
			return false
		}
		if parsed1.Priority != parsed2.Priority {
			t.Logf("Priority mismatch: %d vs %d", parsed1.Priority, parsed2.Priority)
			return false
		}
		if !toolsEqual(parsed1.Tools, parsed2.Tools) {
			t.Logf("Tools mismatch: %v vs %v", parsed1.Tools, parsed2.Tools)
			return false
		}
		if parsed1.Body != parsed2.Body {
			t.Logf("Body mismatch:\n  first:  %q\n  second: %q", parsed1.Body, parsed2.Body)
			return false
		}

		return true
	}

	cfg := &quick.Config{MaxCount: 200}
	if err := quick.Check(f, cfg); err != nil {
		t.Errorf("round-trip property failed: %v", err)
	}
}

// --- helpers ---

const alphaChars = "abcdefghijklmnopqrstuvwxyz"

// randAlphaName generates a non-empty lowercase alphabetic name.
func randAlphaName(r *rand.Rand, size int) string {
	n := r.Intn(max(size, 1)) + 1
	if n > 20 {
		n = 20
	}
	b := make([]byte, n)
	for i := range b {
		b[i] = alphaChars[r.Intn(len(alphaChars))]
	}
	return string(b)
}

// safeChars are characters that won't break YAML or frontmatter parsing.
const safeChars = "abcdefghijklmnopqrstuvwxyz ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789.,;:!?()-"

// randSafeString generates a string safe for YAML values (no newlines, no special YAML chars).
func randSafeString(r *rand.Rand, size int) string {
	n := r.Intn(max(size, 1)) + 1
	if n > 40 {
		n = 40
	}
	b := make([]byte, n)
	for i := range b {
		b[i] = safeChars[r.Intn(len(safeChars))]
	}
	return string(b)
}

// randSafeBody generates markdown body content that does not contain
// a frontmatter delimiter (---) on its own line.
func randSafeBody(r *rand.Rand, size int) string {
	lines := r.Intn(max(size, 1)) + 1
	if lines > 5 {
		lines = 5
	}
	var sb strings.Builder
	for i := 0; i < lines; i++ {
		line := randSafeString(r, size)
		// Ensure no line is exactly "---" which would break frontmatter.
		if strings.TrimSpace(line) == "---" {
			line = "safe content"
		}
		sb.WriteString(line)
		sb.WriteByte('\n')
	}
	return sb.String()
}

// toolsEqual compares two string slices for equality.
func toolsEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
