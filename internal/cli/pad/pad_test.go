//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package pad

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/crypto"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// setupEncrypted creates a temp dir with a .context/ directory and encryption key.
// It sets the RC context dir override and returns a cleanup function.
func setupEncrypted(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(origDir)
		rc.Reset()
	})

	rc.Reset()
	rc.OverrideContextDir(config.DirContext)

	ctxDir := filepath.Join(dir, config.DirContext)
	if err := os.MkdirAll(ctxDir, 0750); err != nil {
		t.Fatal(err)
	}

	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	keyFile := filepath.Join(ctxDir, config.FileScratchpadKey)
	if err := crypto.SaveKey(keyFile, key); err != nil {
		t.Fatal(err)
	}

	return dir
}

// setupPlaintext creates a temp dir with a .context/ directory and
// scratchpad_encrypt: false in .contextrc.
func setupPlaintext(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(origDir)
		rc.Reset()
	})

	// Write .contextrc with encryption disabled
	rcContent := "scratchpad_encrypt: false\n"
	if err := os.WriteFile(filepath.Join(dir, ".contextrc"), []byte(rcContent), 0600); err != nil {
		t.Fatal(err)
	}

	rc.Reset()

	ctxDir := filepath.Join(dir, config.DirContext)
	if err := os.MkdirAll(ctxDir, 0750); err != nil {
		t.Fatal(err)
	}

	return dir
}

// runCmd executes a cobra command and captures its output.
func runCmd(cmd *cobra.Command) (string, error) {
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	err := cmd.Execute()
	return buf.String(), err
}

// newPadCmd builds a fresh pad command with the given args.
func newPadCmd(args ...string) *cobra.Command {
	cmd := Cmd()
	cmd.SetArgs(args)
	return cmd
}

func TestList_Empty(t *testing.T) {
	setupEncrypted(t)

	out, err := runCmd(newPadCmd())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, msgEmpty) {
		t.Errorf("output = %q, want %q", out, msgEmpty)
	}
}

func TestAdd_Encrypted(t *testing.T) {
	setupEncrypted(t)

	out, err := runCmd(newPadCmd("add", "check DNS config"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Added entry 1.") {
		t.Errorf("output = %q, want 'Added entry 1.'", out)
	}

	// Verify listing
	out, err = runCmd(newPadCmd())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "1. check DNS config") {
		t.Errorf("list output = %q, want entry listed", out)
	}
}

func TestAdd_Plaintext(t *testing.T) {
	setupPlaintext(t)

	out, err := runCmd(newPadCmd("add", "plaintext note"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Added entry 1.") {
		t.Errorf("output = %q, want 'Added entry 1.'", out)
	}

	// Verify the file is plain text
	path := filepath.Join(config.DirContext, config.FileScratchpadMd)
	data, err := os.ReadFile(path) //nolint:gosec // test reads a known test file path
	if err != nil {
		t.Fatalf("ReadFile() error: %v", err)
	}
	if string(data) != "plaintext note\n" {
		t.Errorf("file contents = %q, want %q", string(data), "plaintext note\n")
	}
}

func TestMultipleAdd_List(t *testing.T) {
	setupEncrypted(t)

	entries := []string{"first", "second", "third"}
	for _, e := range entries {
		if _, err := runCmd(newPadCmd("add", e)); err != nil {
			t.Fatalf("add %q: %v", e, err)
		}
	}

	out, err := runCmd(newPadCmd())
	if err != nil {
		t.Fatalf("list error: %v", err)
	}

	for i, e := range entries {
		expected := strings.TrimSpace(
			strings.Repeat(" ", 2) + strings.Join(
				[]string{""}, "",
			),
		)
		_ = expected
		line := strings.TrimSpace(out)
		_ = line
		if !strings.Contains(out, e) {
			t.Errorf("list missing entry %d: %q", i+1, e)
		}
	}
}

func TestRm(t *testing.T) {
	setupEncrypted(t)

	for _, e := range []string{"one", "two", "three"} {
		if _, err := runCmd(newPadCmd("add", e)); err != nil {
			t.Fatal(err)
		}
	}

	out, err := runCmd(newPadCmd("rm", "2"))
	if err != nil {
		t.Fatalf("rm error: %v", err)
	}
	if !strings.Contains(out, "Removed entry 2.") {
		t.Errorf("output = %q, want 'Removed entry 2.'", out)
	}

	// Verify remaining entries
	out, err = runCmd(newPadCmd())
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out, "two") {
		t.Error("entry 'two' should have been removed")
	}
	if !strings.Contains(out, "one") || !strings.Contains(out, "three") {
		t.Error("entries 'one' and 'three' should remain")
	}
}

func TestRm_OutOfRange(t *testing.T) {
	setupEncrypted(t)

	if _, err := runCmd(newPadCmd("add", "only")); err != nil {
		t.Fatal(err)
	}

	_, err := runCmd(newPadCmd("rm", "5"))
	if err == nil {
		t.Fatal("expected error for out-of-range index")
	}
	if !strings.Contains(err.Error(), "does not exist") {
		t.Errorf("error = %q, want 'does not exist'", err.Error())
	}
}

func TestEdit(t *testing.T) {
	setupEncrypted(t)

	if _, err := runCmd(newPadCmd("add", "original")); err != nil {
		t.Fatal(err)
	}

	out, err := runCmd(newPadCmd("edit", "1", "updated"))
	if err != nil {
		t.Fatalf("edit error: %v", err)
	}
	if !strings.Contains(out, "Updated entry 1.") {
		t.Errorf("output = %q, want 'Updated entry 1.'", out)
	}

	// Verify
	out, err = runCmd(newPadCmd())
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out, "original") {
		t.Error("old entry should be gone")
	}
	if !strings.Contains(out, "updated") {
		t.Error("new entry should be present")
	}
}

func TestEdit_OutOfRange(t *testing.T) {
	setupEncrypted(t)

	_, err := runCmd(newPadCmd("edit", "1", "text"))
	if err == nil {
		t.Fatal("expected error for empty scratchpad")
	}
}

func TestEdit_Append(t *testing.T) {
	setupEncrypted(t)

	if _, err := runCmd(newPadCmd("add", "check DNS")); err != nil {
		t.Fatal(err)
	}

	out, err := runCmd(newPadCmd("edit", "1", "--append", "on staging"))
	if err != nil {
		t.Fatalf("edit --append error: %v", err)
	}
	if !strings.Contains(out, "Updated entry 1.") {
		t.Errorf("output = %q, want 'Updated entry 1.'", out)
	}

	// Verify the entry was appended
	out, err = runCmd(newPadCmd())
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "check DNS on staging") {
		t.Errorf("list output = %q, want 'check DNS on staging'", out)
	}
}

func TestEdit_Prepend(t *testing.T) {
	setupEncrypted(t)

	if _, err := runCmd(newPadCmd("add", "check DNS")); err != nil {
		t.Fatal(err)
	}

	out, err := runCmd(newPadCmd("edit", "1", "--prepend", "URGENT:"))
	if err != nil {
		t.Fatalf("edit --prepend error: %v", err)
	}
	if !strings.Contains(out, "Updated entry 1.") {
		t.Errorf("output = %q, want 'Updated entry 1.'", out)
	}

	// Verify the entry was prepended
	out, err = runCmd(newPadCmd())
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "URGENT: check DNS") {
		t.Errorf("list output = %q, want 'URGENT: check DNS'", out)
	}
}

func TestEdit_AppendAndPrependMutuallyExclusive(t *testing.T) {
	setupEncrypted(t)

	if _, err := runCmd(newPadCmd("add", "note")); err != nil {
		t.Fatal(err)
	}

	_, err := runCmd(newPadCmd("edit", "1", "--append", "suffix", "--prepend", "prefix"))
	if err == nil {
		t.Fatal("expected error for --append + --prepend")
	}
	if !strings.Contains(err.Error(), "mutually exclusive") {
		t.Errorf("error = %q, want 'mutually exclusive'", err.Error())
	}
}

func TestEdit_PositionalAndFlagMutuallyExclusive(t *testing.T) {
	setupEncrypted(t)

	if _, err := runCmd(newPadCmd("add", "note")); err != nil {
		t.Fatal(err)
	}

	_, err := runCmd(newPadCmd("edit", "1", "replacement", "--append", "suffix"))
	if err == nil {
		t.Fatal("expected error for positional + --append")
	}
	if !strings.Contains(err.Error(), "mutually exclusive") {
		t.Errorf("error = %q, want 'mutually exclusive'", err.Error())
	}
}

func TestEdit_NoTextProvided(t *testing.T) {
	setupEncrypted(t)

	if _, err := runCmd(newPadCmd("add", "note")); err != nil {
		t.Fatal(err)
	}

	_, err := runCmd(newPadCmd("edit", "1"))
	if err == nil {
		t.Fatal("expected error when no text or flag provided")
	}
	if !strings.Contains(err.Error(), "provide replacement text") {
		t.Errorf("error = %q, want 'provide replacement text'", err.Error())
	}
}

func TestShow_Valid(t *testing.T) {
	setupEncrypted(t)

	for _, e := range []string{"alpha", "beta", "gamma"} {
		if _, err := runCmd(newPadCmd("add", e)); err != nil {
			t.Fatal(err)
		}
	}

	out, err := runCmd(newPadCmd("show", "2"))
	if err != nil {
		t.Fatalf("show error: %v", err)
	}

	// Should output raw text with a single trailing newline, no numbering prefix.
	if out != "beta\n" {
		t.Errorf("output = %q, want %q", out, "beta\n")
	}
}

func TestShow_OutOfRange(t *testing.T) {
	setupEncrypted(t)

	if _, err := runCmd(newPadCmd("add", "only")); err != nil {
		t.Fatal(err)
	}

	_, err := runCmd(newPadCmd("show", "5"))
	if err == nil {
		t.Fatal("expected error for out-of-range index")
	}
	if !strings.Contains(err.Error(), "does not exist") {
		t.Errorf("error = %q, want 'does not exist'", err.Error())
	}
}

func TestShow_EmptyScratchpad(t *testing.T) {
	setupEncrypted(t)

	_, err := runCmd(newPadCmd("show", "1"))
	if err == nil {
		t.Fatal("expected error for empty scratchpad")
	}
	if !strings.Contains(err.Error(), "does not exist") {
		t.Errorf("error = %q, want 'does not exist'", err.Error())
	}
}

func TestMv(t *testing.T) {
	setupEncrypted(t)

	for _, e := range []string{"A", "B", "C"} {
		if _, err := runCmd(newPadCmd("add", e)); err != nil {
			t.Fatal(err)
		}
	}

	// Move entry 3 to position 1
	out, err := runCmd(newPadCmd("mv", "3", "1"))
	if err != nil {
		t.Fatalf("mv error: %v", err)
	}
	if !strings.Contains(out, "Moved entry 3 to 1.") {
		t.Errorf("output = %q, want 'Moved entry 3 to 1.'", out)
	}

	// Verify order: C, A, B
	out, err = runCmd(newPadCmd())
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d: %q", len(lines), out)
	}
	if !strings.Contains(lines[0], "C") {
		t.Errorf("line 1 = %q, want 'C'", lines[0])
	}
	if !strings.Contains(lines[1], "A") {
		t.Errorf("line 2 = %q, want 'A'", lines[1])
	}
	if !strings.Contains(lines[2], "B") {
		t.Errorf("line 3 = %q, want 'B'", lines[2])
	}
}

func TestMv_OutOfRange(t *testing.T) {
	setupEncrypted(t)

	if _, err := runCmd(newPadCmd("add", "only")); err != nil {
		t.Fatal(err)
	}

	_, err := runCmd(newPadCmd("mv", "1", "5"))
	if err == nil {
		t.Fatal("expected error for out-of-range destination")
	}
}

func TestNoKey_EncryptedFileExists(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(origDir)
		rc.Reset()
	})

	rc.Reset()
	rc.OverrideContextDir(config.DirContext)

	ctxDir := filepath.Join(dir, config.DirContext)
	if err := os.MkdirAll(ctxDir, 0750); err != nil {
		t.Fatal(err)
	}

	// Create an encrypted file but no key
	if err := os.WriteFile(
		filepath.Join(ctxDir, config.FileScratchpadEnc),
		[]byte("encrypted data here but dummy"),
		0600,
	); err != nil {
		t.Fatal(err)
	}

	_, err := runCmd(newPadCmd())
	if err == nil {
		t.Fatal("expected error when no key exists")
	}
	if !strings.Contains(err.Error(), "no key") {
		t.Errorf("error = %q, want 'no key' message", err.Error())
	}
}

func TestDecryptionFailure_WrongKey(t *testing.T) {
	setupEncrypted(t)

	// Add an entry
	if _, err := runCmd(newPadCmd("add", "secret")); err != nil {
		t.Fatal(err)
	}

	// Replace the key with a different one
	newKey, _ := crypto.GenerateKey()
	keyFile := filepath.Join(config.DirContext, config.FileScratchpadKey)
	if err := crypto.SaveKey(keyFile, newKey); err != nil {
		t.Fatal(err)
	}

	_, err := runCmd(newPadCmd())
	if err == nil {
		t.Fatal("expected decryption error with wrong key")
	}
	if !strings.Contains(err.Error(), "wrong key") {
		t.Errorf("error = %q, want 'wrong key' message", err.Error())
	}
}

func TestPlaintext_ListFormat(t *testing.T) {
	setupPlaintext(t)

	for _, e := range []string{"alpha", "beta", "gamma"} {
		if _, err := runCmd(newPadCmd("add", e)); err != nil {
			t.Fatal(err)
		}
	}

	out, err := runCmd(newPadCmd())
	if err != nil {
		t.Fatal(err)
	}

	// Check 2-space indent, 1-based numbering
	if !strings.Contains(out, "  1. alpha") {
		t.Errorf("output missing '  1. alpha': %q", out)
	}
	if !strings.Contains(out, "  2. beta") {
		t.Errorf("output missing '  2. beta': %q", out)
	}
	if !strings.Contains(out, "  3. gamma") {
		t.Errorf("output missing '  3. gamma': %q", out)
	}
}

func TestParseEntries_EmptyInput(t *testing.T) {
	entries := parseEntries(nil)
	if entries != nil {
		t.Errorf("parseEntries(nil) = %v, want nil", entries)
	}

	entries = parseEntries([]byte{})
	if entries != nil {
		t.Errorf("parseEntries(empty) = %v, want nil", entries)
	}
}

func TestParseEntries_SkipsEmpty(t *testing.T) {
	entries := parseEntries([]byte("a\n\nb\n"))
	if len(entries) != 2 {
		t.Fatalf("len = %d, want 2", len(entries))
	}
	if entries[0] != "a" || entries[1] != "b" {
		t.Errorf("entries = %v, want [a b]", entries)
	}
}

func TestFormatEntries_Empty(t *testing.T) {
	data := formatEntries(nil)
	if data != nil {
		t.Errorf("formatEntries(nil) = %v, want nil", data)
	}
}

func TestFormatEntries_TrailingNewline(t *testing.T) {
	data := formatEntries([]string{"a", "b"})
	if string(data) != "a\nb\n" {
		t.Errorf("formatEntries = %q, want %q", string(data), "a\nb\n")
	}
}
