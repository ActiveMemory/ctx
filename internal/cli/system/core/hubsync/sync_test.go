//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hubsync_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	connectCfg "github.com/ActiveMemory/ctx/internal/cli/connection/core/config"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/hubsync"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/crypto"
	"github.com/ActiveMemory/ctx/internal/log/warn"
	"github.com/ActiveMemory/ctx/internal/testutil/testctx"
)

// TestSync_WarnsOnLoadError verifies that the session-start
// hub sync hook emits a warning via [warn.Warn] when the
// connection config cannot be loaded, instead of swallowing
// the error and producing an empty status. The "warn but do
// not block" contract is documented in the package doc; this
// test is the regression guard for ActiveMemory/ctx#100.
func TestSync_WarnsOnLoadError(t *testing.T) {
	tempDir := t.TempDir()
	testctx.Declare(t, tempDir)

	var buf bytes.Buffer
	restore := warn.SetSinkForTesting(&buf)
	defer restore()

	msg := hubsync.Sync("session-id-ignored")
	if msg != "" {
		t.Errorf("Sync should return empty on load error, got %q", msg)
	}

	got := buf.String()
	if !strings.Contains(got, "hubsync: load connection config:") {
		t.Errorf(
			"warning output missing hubsync load prefix; got %q",
			got,
		)
	}
}

// TestSync_WarnsOnDialError verifies that hubsync emits a
// warning when grpc.NewClient rejects a malformed hub
// address. "%" reliably trips the URL parser inside
// grpc.NewClient ("invalid URL escape"), exercising the
// dial-error branch without needing a network listener.
// Anything reachable enough to satisfy NewClient (like
// "dns:///") would defer the failure to the Sync RPC and
// hit a different warn path; this test is specifically the
// regression guard for the dial-error branch.
func TestSync_WarnsOnDialError(t *testing.T) {
	tempDir := t.TempDir()
	ctxDir := testctx.Declare(t, tempDir)
	if mkErr := os.Mkdir(ctxDir, fs.PermKeyDir); mkErr != nil {
		t.Fatal(mkErr)
	}
	writeTestKey(t, tempDir)
	if saveErr := connectCfg.Save(connectCfg.Config{
		HubAddr: "%",
		Token:   "test-token",
	}); saveErr != nil {
		t.Fatal(saveErr)
	}

	var buf bytes.Buffer
	restore := warn.SetSinkForTesting(&buf)
	defer restore()

	msg := hubsync.Sync("session-id-ignored")
	if msg != "" {
		t.Errorf("Sync should return empty on dial error, got %q", msg)
	}

	got := buf.String()
	if !strings.Contains(got, "hubsync: dial %:") {
		t.Errorf(
			"warning output missing hubsync dial prefix; got %q",
			got,
		)
	}
}

// TestSync_NonBlockingOnLoadError verifies the second half of
// the contract: even when the load fails and a warning is
// emitted, Sync returns without panicking or propagating an
// error to the caller. The check-hub-sync hook depends on
// this invariant to never block session start.
func TestSync_NonBlockingOnLoadError(t *testing.T) {
	tempDir := t.TempDir()
	testctx.Declare(t, tempDir)

	restore := warn.SetSinkForTesting(&bytes.Buffer{})
	defer restore()

	// The bare fact that Sync returns at all (rather than
	// panicking or hanging) is the assertion. The returned
	// string is checked separately in TestSync_WarnsOnLoadError.
	_ = hubsync.Sync("")
}

func writeTestKey(t *testing.T, home string) {
	t.Helper()
	key, genErr := crypto.GenerateKey()
	if genErr != nil {
		t.Fatal(genErr)
	}
	keyDir := filepath.Join(home, dir.CtxData)
	if mkErr := os.Mkdir(keyDir, fs.PermKeyDir); mkErr != nil {
		t.Fatal(mkErr)
	}
	if saveErr := crypto.SaveKey(crypto.GlobalKeyPath(), key); saveErr != nil {
		t.Fatal(saveErr)
	}
}
