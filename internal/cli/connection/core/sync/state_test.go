//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sync

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/ActiveMemory/ctx/internal/config/fs"
	cfgHub "github.com/ActiveMemory/ctx/internal/config/hub"
	"github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/testutil/testctx"
)

// declareContext positions the test in a temp project with a
// materialized .context/ directory and returns its path.
func declareContext(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	ctxDir := testctx.Declare(t, dir)
	if mkErr := os.MkdirAll(ctxDir, fs.PermExec); mkErr != nil {
		t.Fatal(mkErr)
	}
	return ctxDir
}

func TestLoadState_RejectsConcurrentSyncs(t *testing.T) {
	declareContext(t)

	const attempts = 16

	type outcome struct {
		release func()
		err     error
	}

	start := make(chan struct{})
	results := make(chan outcome, attempts)
	for i := 0; i < attempts; i++ {
		go func() {
			<-start
			_, release, lockErr := loadState()
			results <- outcome{release: release, err: lockErr}
		}()
	}
	close(start)

	// Winners must not release until every attempt has
	// finished, or a late goroutine could legitimately
	// re-acquire the freed lock and skew the count.
	var winners []func()
	var contended int
	for i := 0; i < attempts; i++ {
		got := <-results
		switch {
		case got.err == nil:
			winners = append(winners, got.release)
		case errors.Is(got.err, os.ErrExist):
			contended++
		default:
			t.Fatalf("unexpected error: %v", got.err)
		}
	}

	if len(winners) != 1 {
		t.Errorf("winners = %d, want exactly 1", len(winners))
	}
	if contended != attempts-1 {
		t.Errorf(
			"contended = %d, want %d", contended, attempts-1,
		)
	}
	for _, release := range winners {
		release()
	}
}

func TestLoadState_ReleaseRemovesLock(t *testing.T) {
	ctxDir := declareContext(t)
	lockPath := filepath.Join(
		ctxDir, cfgHub.DirHub, cfgHub.FileSyncLock,
	)

	_, release, lockErr := loadState()
	if lockErr != nil {
		t.Fatalf("first loadState: %v", lockErr)
	}
	if _, statErr := os.Stat(lockPath); statErr != nil {
		t.Fatalf("lock file should exist while held: %v", statErr)
	}

	release()

	if _, statErr := os.Stat(lockPath); !os.IsNotExist(statErr) {
		t.Errorf(
			"lock file should be gone after release: %v", statErr,
		)
	}

	// The next sync must be able to proceed.
	_, release, lockErr = loadState()
	if lockErr != nil {
		t.Fatalf("loadState after release: %v", lockErr)
	}
	release()
}

func TestLoadState_ReleasesLockOnCorruptState(t *testing.T) {
	ctxDir := declareContext(t)
	hubDir := filepath.Join(ctxDir, cfgHub.DirHub)
	if mkErr := io.SafeMkdirAll(hubDir, fs.PermKeyDir); mkErr != nil {
		t.Fatal(mkErr)
	}
	statePath := filepath.Join(hubDir, cfgHub.FileSyncState)
	if writeErr := os.WriteFile(
		statePath, []byte("{not json"), fs.PermFile,
	); writeErr != nil {
		t.Fatal(writeErr)
	}

	_, _, loadErr := loadState()
	if loadErr == nil {
		t.Fatal("loadState should fail on corrupt state")
	}

	lockPath := filepath.Join(hubDir, cfgHub.FileSyncLock)
	if _, statErr := os.Stat(lockPath); !os.IsNotExist(statErr) {
		t.Errorf(
			"lock must not leak after a failed load: %v", statErr,
		)
	}
}

func TestLoadState_LockFileLocation(t *testing.T) {
	ctxDir := declareContext(t)

	_, release, lockErr := loadState()
	if lockErr != nil {
		t.Fatalf("loadState: %v", lockErr)
	}
	defer release()

	want := filepath.Join(ctxDir, cfgHub.DirHub, cfgHub.FileSyncLock)
	if _, statErr := os.Stat(want); statErr != nil {
		t.Errorf("lock file not at %s: %v", want, statErr)
	}
}
