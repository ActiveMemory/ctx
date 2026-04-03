//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package steering

import (
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"testing/quick"
)

// syncInput bundles one or more valid steering files for property testing
// the sync idempotence property.
type syncInput struct {
	Files []validSteeringFile
}

// Generate implements quick.Generator for syncInput.
// It produces 1–3 steering files with distinct names.
func (syncInput) Generate(r *rand.Rand, size int) reflect.Value {
	n := r.Intn(3) + 1
	seen := make(map[string]bool)
	var files []validSteeringFile
	for len(files) < n {
		v := validSteeringFile{}.Generate(r, size).Interface().(validSteeringFile)
		if seen[v.Name] {
			continue
		}
		seen[v.Name] = true
		files = append(files, v)
	}
	return reflect.ValueOf(syncInput{Files: files})
}

// TestProperty_SyncIdempotence verifies the sync idempotence property:
// running SyncTool twice produces identical output files and the second
// run reports zero written files (all skipped as unchanged).
//
// **Validates: Requirements 19.4**
func TestProperty_SyncIdempotence(t *testing.T) {
	f := func(input syncInput) bool {
		root := t.TempDir()
		steeringDir := filepath.Join(root, ".context", "steering")
		if err := os.MkdirAll(steeringDir, 0o755); err != nil {
			t.Logf("mkdir steering dir: %v", err)
			return false
		}

		// Write generated steering files to disk.
		for _, v := range input.Files {
			sf := &SteeringFile{
				Name:        v.Name,
				Description: v.Description,
				Inclusion:   v.Inclusion,
				Tools:       v.Tools,
				Priority:    v.Priority,
				Body:        v.Body,
			}
			data := Print(sf)
			path := filepath.Join(steeringDir, v.Name+".md")
			if err := os.WriteFile(path, data, 0o644); err != nil {
				t.Logf("write steering file %s: %v", v.Name, err)
				return false
			}
		}

		tools := []string{"cursor", "cline", "kiro"}
		for _, tool := range tools {
			// First sync — writes files.
			r1, err := SyncTool(steeringDir, root, tool)
			if err != nil {
				t.Logf("first SyncTool(%s): %v", tool, err)
				return false
			}
			if len(r1.Errors) > 0 {
				t.Logf("first SyncTool(%s) errors: %v", tool, r1.Errors)
				return false
			}

			// Capture output file contents after first sync.
			snapshot := make(map[string][]byte)
			for _, name := range r1.Written {
				outPath := nativePath(root, tool, name)
				data, readErr := os.ReadFile(outPath)
				if readErr != nil {
					t.Logf("read output %s: %v", outPath, readErr)
					return false
				}
				snapshot[name] = data
			}

			// Second sync — should skip all files.
			r2, err := SyncTool(steeringDir, root, tool)
			if err != nil {
				t.Logf("second SyncTool(%s): %v", tool, err)
				return false
			}
			if len(r2.Errors) > 0 {
				t.Logf("second SyncTool(%s) errors: %v", tool, r2.Errors)
				return false
			}

			// Verify second run wrote zero files.
			if len(r2.Written) != 0 {
				t.Logf("tool %s: second sync wrote %d files, expected 0: %v",
					tool, len(r2.Written), r2.Written)
				return false
			}

			// Verify output files are byte-identical after second sync.
			for name, before := range snapshot {
				outPath := nativePath(root, tool, name)
				after, readErr := os.ReadFile(outPath)
				if readErr != nil {
					t.Logf("re-read output %s: %v", outPath, readErr)
					return false
				}
				if string(before) != string(after) {
					t.Logf("tool %s, file %s: content changed between syncs", tool, name)
					return false
				}
			}
		}

		return true
	}

	cfg := &quick.Config{MaxCount: 100}
	if err := quick.Check(f, cfg); err != nil {
		t.Errorf("sync idempotence property failed: %v", err)
	}
}
