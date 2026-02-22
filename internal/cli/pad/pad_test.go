//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package pad

import (
	"bytes"
	"encoding/base64"
	"fmt"
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
// scratchpad_encrypt: false in .ctxrc.
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

	// Write .ctxrc with encryption disabled
	rcContent := "scratchpad_encrypt: false\n"
	if err := os.WriteFile(filepath.Join(dir, ".ctxrc"), []byte(rcContent), 0600); err != nil {
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

func TestValidateIndex(t *testing.T) {
	entries := []string{"a", "b", "c"}

	// Valid indices
	for _, n := range []int{1, 2, 3} {
		if err := validateIndex(n, entries); err != nil {
			t.Errorf("validateIndex(%d) should be valid: %v", n, err)
		}
	}

	// Invalid indices
	for _, n := range []int{0, -1, 4, 100} {
		if err := validateIndex(n, entries); err == nil {
			t.Errorf("validateIndex(%d) should be invalid", n)
		}
	}
}

func TestValidateIndex_EmptySlice(t *testing.T) {
	err := validateIndex(1, nil)
	if err == nil {
		t.Error("validateIndex on nil slice should fail")
	}
}

func TestErrEntryRange(t *testing.T) {
	msg := errEntryRange(5, 3)
	if !strings.Contains(msg, "5") || !strings.Contains(msg, "3") {
		t.Errorf("errEntryRange = %q, want indices 5 and 3 mentioned", msg)
	}
}

func TestCmd_HasSubcommands(t *testing.T) {
	cmd := Cmd()
	if cmd.Use != "pad" {
		t.Errorf("cmd.Use = %q, want 'pad'", cmd.Use)
	}

	names := make(map[string]bool)
	for _, sub := range cmd.Commands() {
		names[sub.Use] = true
	}
	for _, expected := range []string{"show N", "add TEXT", "rm N", "edit N [TEXT]", "mv N M", "resolve", "import FILE", "export [DIR]"} {
		if !names[expected] {
			t.Errorf("missing subcommand %q", expected)
		}
	}
}

func TestRm_InvalidIndex(t *testing.T) {
	setupEncrypted(t)

	if _, err := runCmd(newPadCmd("add", "solo")); err != nil {
		t.Fatal(err)
	}

	// Non-numeric argument
	_, err := runCmd(newPadCmd("rm", "abc"))
	if err == nil {
		t.Error("expected error for non-numeric rm argument")
	}
}

func TestMv_InvalidIndex(t *testing.T) {
	setupEncrypted(t)

	if _, err := runCmd(newPadCmd("add", "entry")); err != nil {
		t.Fatal(err)
	}

	// Non-numeric first argument
	_, err := runCmd(newPadCmd("mv", "abc", "1"))
	if err == nil {
		t.Error("expected error for non-numeric mv src argument")
	}

	// Non-numeric second argument
	_, err = runCmd(newPadCmd("mv", "1", "abc"))
	if err == nil {
		t.Error("expected error for non-numeric mv dst argument")
	}
}

func TestShow_InvalidIndex(t *testing.T) {
	setupEncrypted(t)

	_, err := runCmd(newPadCmd("show", "abc"))
	if err == nil {
		t.Error("expected error for non-numeric show argument")
	}
}

func TestEdit_InvalidIndex(t *testing.T) {
	setupEncrypted(t)

	_, err := runCmd(newPadCmd("edit", "abc", "text"))
	if err == nil {
		t.Error("expected error for non-numeric edit argument")
	}
}

func TestEnsureGitignore_NewFile(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(origDir) })

	err := ensureGitignore(".context", ".scratchpad.key")
	if err != nil {
		t.Fatalf("ensureGitignore error: %v", err)
	}

	data, err := os.ReadFile(".gitignore")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), filepath.Join(".context", ".scratchpad.key")) {
		t.Errorf(".gitignore = %q, want key entry", string(data))
	}
}

func TestEnsureGitignore_AlreadyPresent(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(origDir) })

	entry := filepath.Join(".context", ".scratchpad.key")
	if err := os.WriteFile(".gitignore", []byte(entry+"\n"), 0600); err != nil {
		t.Fatal(err)
	}

	err := ensureGitignore(".context", ".scratchpad.key")
	if err != nil {
		t.Fatalf("ensureGitignore error: %v", err)
	}

	data, _ := os.ReadFile(".gitignore")
	// Should not duplicate the entry
	count := strings.Count(string(data), entry)
	if count != 1 {
		t.Errorf("expected 1 occurrence of entry, got %d", count)
	}
}

func TestEnsureGitignore_AppendToExisting(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(origDir) })

	// Write file without trailing newline
	if err := os.WriteFile(".gitignore", []byte("node_modules"), 0600); err != nil {
		t.Fatal(err)
	}

	err := ensureGitignore(".context", ".scratchpad.key")
	if err != nil {
		t.Fatalf("ensureGitignore error: %v", err)
	}

	data, _ := os.ReadFile(".gitignore")
	if !strings.Contains(string(data), "node_modules\n") {
		t.Error("existing content should be preserved with newline")
	}
	if !strings.Contains(string(data), filepath.Join(".context", ".scratchpad.key")) {
		t.Error("new entry should be present")
	}
}

func TestScratchpadPath_Plaintext(t *testing.T) {
	setupPlaintext(t)

	path := scratchpadPath()
	if !strings.HasSuffix(path, config.FileScratchpadMd) {
		t.Errorf("scratchpadPath() = %q, want suffix %q", path, config.FileScratchpadMd)
	}
}

func TestScratchpadPath_Encrypted(t *testing.T) {
	setupEncrypted(t)

	path := scratchpadPath()
	if !strings.HasSuffix(path, config.FileScratchpadEnc) {
		t.Errorf("scratchpadPath() = %q, want suffix %q", path, config.FileScratchpadEnc)
	}
}

func TestKeyPath(t *testing.T) {
	setupEncrypted(t)

	path := keyPath()
	if !strings.HasSuffix(path, config.FileScratchpadKey) {
		t.Errorf("keyPath() = %q, want suffix %q", path, config.FileScratchpadKey)
	}
}

func TestEnsureKey_KeyAlreadyExists(t *testing.T) {
	setupEncrypted(t)

	// Key already exists from setup
	err := ensureKey()
	if err != nil {
		t.Fatalf("ensureKey should succeed when key already exists: %v", err)
	}
}

func TestEnsureKey_EncFileExistsNoKey(t *testing.T) {
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

	// Create enc file but no key
	encPath := filepath.Join(ctxDir, config.FileScratchpadEnc)
	if err := os.WriteFile(encPath, []byte("data"), 0600); err != nil {
		t.Fatal(err)
	}

	err := ensureKey()
	if err == nil {
		t.Fatal("expected error when enc file exists without key")
	}
	if !strings.Contains(err.Error(), "no key") {
		t.Errorf("error = %q, want 'no key' message", err.Error())
	}
}

func TestEnsureKey_GeneratesNewKey(t *testing.T) {
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

	// No key, no enc file -- should generate
	err := ensureKey()
	if err != nil {
		t.Fatalf("ensureKey error: %v", err)
	}

	kp := filepath.Join(ctxDir, config.FileScratchpadKey)
	if _, err := os.Stat(kp); err != nil {
		t.Error("key file should have been created")
	}
}

func TestWriteEntries_Plaintext(t *testing.T) {
	setupPlaintext(t)

	entries := []string{"one", "two"}
	if err := writeEntries(entries); err != nil {
		t.Fatalf("writeEntries error: %v", err)
	}

	path := scratchpadPath()
	data, err := os.ReadFile(path) //nolint:gosec // test temp path
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "one\ntwo\n" {
		t.Errorf("file = %q, want %q", string(data), "one\ntwo\n")
	}
}

func TestReadEntries_Plaintext(t *testing.T) {
	setupPlaintext(t)

	path := scratchpadPath()
	if err := os.WriteFile(path, []byte("alpha\nbeta\n"), 0600); err != nil {
		t.Fatal(err)
	}

	entries, err := readEntries()
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 || entries[0] != "alpha" || entries[1] != "beta" {
		t.Errorf("entries = %v, want [alpha beta]", entries)
	}
}

func TestReadEntries_NoFile(t *testing.T) {
	setupEncrypted(t)

	entries, err := readEntries()
	if err != nil {
		t.Fatalf("readEntries with no file should return nil, nil: %v", err)
	}
	if entries != nil {
		t.Errorf("entries = %v, want nil", entries)
	}
}

func TestResolve_PlaintextMode(t *testing.T) {
	setupPlaintext(t)

	_, err := runCmd(newPadCmd("resolve"))
	if err == nil {
		t.Fatal("expected error for resolve in plaintext mode")
	}
	if !strings.Contains(err.Error(), "only needed for encrypted") {
		t.Errorf("error = %q, want 'only needed for encrypted'", err.Error())
	}
}

func TestResolve_NoConflictFiles(t *testing.T) {
	setupEncrypted(t)

	_, err := runCmd(newPadCmd("resolve"))
	if err == nil {
		t.Fatal("expected error when no conflict files exist")
	}
	if !strings.Contains(err.Error(), "no conflict files found") {
		t.Errorf("error = %q, want 'no conflict files'", err.Error())
	}
}

func TestResolve_WithConflictFiles(t *testing.T) {
	setupEncrypted(t)

	// Load the key
	kp := filepath.Join(config.DirContext, config.FileScratchpadKey)
	key, err := crypto.LoadKey(kp)
	if err != nil {
		t.Fatal(err)
	}

	// Create encrypted "ours" file
	oursPlain := []byte("ours-entry\n")
	oursCipher, err := crypto.Encrypt(key, oursPlain)
	if err != nil {
		t.Fatal(err)
	}
	oursPath := filepath.Join(config.DirContext, config.FileScratchpadEnc+".ours")
	err = os.WriteFile(oursPath, oursCipher, 0600)
	if err != nil {
		t.Fatal(err)
	}

	// Create encrypted "theirs" file
	theirsPlain := []byte("theirs-entry\n")
	theirsCipher, err := crypto.Encrypt(key, theirsPlain)
	if err != nil {
		t.Fatal(err)
	}
	theirsPath := filepath.Join(config.DirContext, config.FileScratchpadEnc+".theirs")
	err = os.WriteFile(theirsPath, theirsCipher, 0600)
	if err != nil {
		t.Fatal(err)
	}

	out, err := runCmd(newPadCmd("resolve"))
	if err != nil {
		t.Fatalf("resolve error: %v", err)
	}
	if !strings.Contains(out, "OURS") {
		t.Error("output should contain OURS section")
	}
	if !strings.Contains(out, "THEIRS") {
		t.Error("output should contain THEIRS section")
	}
	if !strings.Contains(out, "ours-entry") {
		t.Error("output should contain ours-entry")
	}
	if !strings.Contains(out, "theirs-entry") {
		t.Error("output should contain theirs-entry")
	}
}

func TestResolve_OnlyOursFile(t *testing.T) {
	setupEncrypted(t)

	kp := filepath.Join(config.DirContext, config.FileScratchpadKey)
	key, err := crypto.LoadKey(kp)
	if err != nil {
		t.Fatal(err)
	}

	oursPlain := []byte("ours-only\n")
	oursCipher, err := crypto.Encrypt(key, oursPlain)
	if err != nil {
		t.Fatal(err)
	}
	oursPath := filepath.Join(config.DirContext, config.FileScratchpadEnc+".ours")
	err = os.WriteFile(oursPath, oursCipher, 0600)
	if err != nil {
		t.Fatal(err)
	}

	out, err := runCmd(newPadCmd("resolve"))
	if err != nil {
		t.Fatalf("resolve error: %v", err)
	}
	if !strings.Contains(out, "OURS") {
		t.Error("output should contain OURS section")
	}
	if strings.Contains(out, "THEIRS") {
		t.Error("output should NOT contain THEIRS section when only ours exists")
	}
}

func TestMv_SamePosition(t *testing.T) {
	setupEncrypted(t)

	for _, e := range []string{"A", "B", "C"} {
		if _, err := runCmd(newPadCmd("add", e)); err != nil {
			t.Fatal(err)
		}
	}

	// Move entry 2 to position 2 (noop)
	out, err := runCmd(newPadCmd("mv", "2", "2"))
	if err != nil {
		t.Fatalf("mv error: %v", err)
	}
	if !strings.Contains(out, "Moved entry 2 to 2.") {
		t.Errorf("output = %q", out)
	}
}

func TestList_PlaintextEmpty(t *testing.T) {
	setupPlaintext(t)

	out, err := runCmd(newPadCmd())
	if err != nil {
		t.Fatalf("list error: %v", err)
	}
	if !strings.Contains(out, msgEmpty) {
		t.Errorf("output = %q, want empty message", out)
	}
}

func TestAdd_MultiplePlaintext(t *testing.T) {
	setupPlaintext(t)

	for i, e := range []string{"first", "second", "third"} {
		out, err := runCmd(newPadCmd("add", e))
		if err != nil {
			t.Fatalf("add error: %v", err)
		}
		expected := strings.TrimSpace(out)
		_ = expected
		if !strings.Contains(out, "Added entry") {
			t.Errorf("add %d: output = %q, want 'Added entry'", i+1, out)
		}
	}

	out, err := runCmd(newPadCmd())
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "first") || !strings.Contains(out, "second") || !strings.Contains(out, "third") {
		t.Errorf("list output missing entries: %q", out)
	}
}

func TestEdit_AppendOutOfRange(t *testing.T) {
	setupEncrypted(t)

	_, err := runCmd(newPadCmd("edit", "1", "--append", "suffix"))
	if err == nil {
		t.Fatal("expected error for append on empty scratchpad")
	}
}

func TestEdit_PrependOutOfRange(t *testing.T) {
	setupEncrypted(t)

	_, err := runCmd(newPadCmd("edit", "1", "--prepend", "prefix"))
	if err == nil {
		t.Fatal("expected error for prepend on empty scratchpad")
	}
}

func TestDecryptFile_BadData(t *testing.T) {
	key, _ := crypto.GenerateKey()
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.enc")
	if err := os.WriteFile(path, []byte("not-encrypted"), 0600); err != nil {
		t.Fatal(err)
	}

	_, err := decryptFile(key, path)
	if err == nil {
		t.Fatal("expected decryption error for bad data")
	}
	if !strings.Contains(err.Error(), "wrong key") {
		t.Errorf("error = %q, want 'wrong key'", err.Error())
	}
}

func TestDecryptFile_MissingFile(t *testing.T) {
	key, _ := crypto.GenerateKey()

	_, err := decryptFile(key, "/nonexistent/path")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestDecryptFile_ValidData(t *testing.T) {
	key, _ := crypto.GenerateKey()
	dir := t.TempDir()
	path := filepath.Join(dir, "good.enc")

	plaintext := []byte("entry1\nentry2\n")
	ciphertext, err := crypto.Encrypt(key, plaintext)
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(path, ciphertext, 0600)
	if err != nil {
		t.Fatal(err)
	}

	entries, err := decryptFile(key, path)
	if err != nil {
		t.Fatalf("decryptFile error: %v", err)
	}
	if len(entries) != 2 || entries[0] != "entry1" || entries[1] != "entry2" {
		t.Errorf("entries = %v, want [entry1 entry2]", entries)
	}
}

func TestRm_Plaintext(t *testing.T) {
	setupPlaintext(t)

	for _, e := range []string{"one", "two"} {
		if _, err := runCmd(newPadCmd("add", e)); err != nil {
			t.Fatal(err)
		}
	}

	out, err := runCmd(newPadCmd("rm", "1"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "Removed entry 1.") {
		t.Errorf("output = %q", out)
	}

	out, err = runCmd(newPadCmd())
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out, "one") {
		t.Error("entry 'one' should be removed")
	}
	if !strings.Contains(out, "two") {
		t.Error("entry 'two' should remain")
	}
}

func TestEdit_PlaintextReplace(t *testing.T) {
	setupPlaintext(t)

	if _, err := runCmd(newPadCmd("add", "original")); err != nil {
		t.Fatal(err)
	}

	out, err := runCmd(newPadCmd("edit", "1", "replaced"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "Updated entry 1.") {
		t.Errorf("output = %q", out)
	}

	out, err = runCmd(newPadCmd("show", "1"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "replaced") {
		t.Error("entry should be replaced")
	}
}

func TestMv_Plaintext(t *testing.T) {
	setupPlaintext(t)

	for _, e := range []string{"A", "B", "C"} {
		if _, err := runCmd(newPadCmd("add", e)); err != nil {
			t.Fatal(err)
		}
	}

	out, err := runCmd(newPadCmd("mv", "1", "3"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "Moved entry 1 to 3.") {
		t.Errorf("output = %q", out)
	}

	out, err = runCmd(newPadCmd())
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if !strings.Contains(lines[0], "B") {
		t.Errorf("line 1 = %q, want B", lines[0])
	}
	if !strings.Contains(lines[1], "C") {
		t.Errorf("line 2 = %q, want C", lines[1])
	}
	if !strings.Contains(lines[2], "A") {
		t.Errorf("line 3 = %q, want A", lines[2])
	}
}

// --- Blob helper tests ---

func TestIsBlob(t *testing.T) {
	if !isBlob("my plan:::SGVsbG8=") {
		t.Error("expected isBlob to return true for blob entry")
	}
	if isBlob("just a plain entry") {
		t.Error("expected isBlob to return false for plain entry")
	}
}

func TestSplitBlob_Valid(t *testing.T) {
	data := []byte("hello world")
	encoded := base64.StdEncoding.EncodeToString(data)
	entry := "my label" + BlobSep + encoded

	label, decoded, ok := splitBlob(entry)
	if !ok {
		t.Fatal("splitBlob returned ok=false for valid blob")
	}
	if label != "my label" {
		t.Errorf("label = %q, want %q", label, "my label")
	}
	if string(decoded) != "hello world" {
		t.Errorf("data = %q, want %q", string(decoded), "hello world")
	}
}

func TestSplitBlob_NonBlob(t *testing.T) {
	_, _, ok := splitBlob("just a plain entry")
	if ok {
		t.Error("splitBlob should return ok=false for non-blob entry")
	}
}

func TestSplitBlob_MalformedBase64(t *testing.T) {
	_, _, ok := splitBlob("label:::not-valid-base64!!!")
	if ok {
		t.Error("splitBlob should return ok=false for malformed base64")
	}
}

func TestMakeBlob_Roundtrip(t *testing.T) {
	original := []byte("secret file content\nwith newlines\n")
	entry := makeBlob("my file", original)

	label, data, ok := splitBlob(entry)
	if !ok {
		t.Fatal("splitBlob failed on makeBlob output")
	}
	if label != "my file" {
		t.Errorf("label = %q, want %q", label, "my file")
	}
	if string(data) != string(original) {
		t.Errorf("data = %q, want %q", string(data), string(original))
	}
}

func TestDisplayEntry_Blob(t *testing.T) {
	entry := makeBlob("my plan", []byte("content"))
	display := displayEntry(entry)
	if display != "my plan [BLOB]" {
		t.Errorf("displayEntry = %q, want %q", display, "my plan [BLOB]")
	}
}

func TestDisplayEntry_Plain(t *testing.T) {
	entry := "just a note"
	display := displayEntry(entry)
	if display != entry {
		t.Errorf("displayEntry = %q, want %q", display, entry)
	}
}

// --- Blob add tests ---

func TestAdd_BlobEncrypted(t *testing.T) {
	dir := setupEncrypted(t)

	// Create a test file.
	testFile := filepath.Join(dir, "test-blob.md")
	content := "secret plan content\n"
	if err := os.WriteFile(testFile, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	out, err := runCmd(newPadCmd("add", "--file", testFile, "my plan"))
	if err != nil {
		t.Fatalf("add --file error: %v", err)
	}
	if !strings.Contains(out, "Added entry 1.") {
		t.Errorf("output = %q, want 'Added entry 1.'", out)
	}

	// Verify listing shows [BLOB].
	out, err = runCmd(newPadCmd())
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "my plan [BLOB]") {
		t.Errorf("list output = %q, want 'my plan [BLOB]'", out)
	}
}

func TestAdd_BlobTooLarge(t *testing.T) {
	dir := setupEncrypted(t)

	testFile := filepath.Join(dir, "big.bin")
	data := make([]byte, MaxBlobSize+1)
	if err := os.WriteFile(testFile, data, 0600); err != nil {
		t.Fatal(err)
	}

	_, err := runCmd(newPadCmd("add", "--file", testFile, "big blob"))
	if err == nil {
		t.Fatal("expected error for file exceeding MaxBlobSize")
	}
	if !strings.Contains(err.Error(), "file too large") {
		t.Errorf("error = %q, want 'file too large'", err.Error())
	}
}

func TestAdd_BlobFileNotFound(t *testing.T) {
	setupEncrypted(t)

	_, err := runCmd(newPadCmd("add", "--file", "/nonexistent/file.md", "missing"))
	if err == nil {
		t.Fatal("expected error for missing file")
	}
	if !strings.Contains(err.Error(), "read file") {
		t.Errorf("error = %q, want 'read file'", err.Error())
	}
}

// --- Blob list tests ---

func TestList_BlobDisplay(t *testing.T) {
	dir := setupEncrypted(t)

	// Add a plain entry.
	if _, err := runCmd(newPadCmd("add", "plain note")); err != nil {
		t.Fatal(err)
	}

	// Add a blob entry.
	testFile := filepath.Join(dir, "blob.txt")
	if err := os.WriteFile(testFile, []byte("file content"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := runCmd(newPadCmd("add", "--file", testFile, "my blob")); err != nil {
		t.Fatal(err)
	}

	out, err := runCmd(newPadCmd())
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "1. plain note") {
		t.Errorf("list missing plain entry: %q", out)
	}
	if !strings.Contains(out, "2. my blob [BLOB]") {
		t.Errorf("list missing blob entry: %q", out)
	}
}

// --- Blob show tests ---

func TestShow_BlobAutoDecodes(t *testing.T) {
	dir := setupEncrypted(t)

	content := "decoded file content\n"
	testFile := filepath.Join(dir, "blob.txt")
	if err := os.WriteFile(testFile, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := runCmd(newPadCmd("add", "--file", testFile, "my blob")); err != nil {
		t.Fatal(err)
	}

	out, err := runCmd(newPadCmd("show", "1"))
	if err != nil {
		t.Fatalf("show error: %v", err)
	}
	if out != content {
		t.Errorf("show output = %q, want %q", out, content)
	}
}

func TestShow_BlobOutFlag(t *testing.T) {
	dir := setupEncrypted(t)

	content := "file to recover\n"
	testFile := filepath.Join(dir, "blob.txt")
	if err := os.WriteFile(testFile, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := runCmd(newPadCmd("add", "--file", testFile, "my blob")); err != nil {
		t.Fatal(err)
	}

	outFile := filepath.Join(dir, "recovered.txt")
	out, err := runCmd(newPadCmd("show", "1", "--out", outFile))
	if err != nil {
		t.Fatalf("show --out error: %v", err)
	}
	if !strings.Contains(out, "Wrote") {
		t.Errorf("output = %q, want 'Wrote' confirmation", out)
	}

	recovered, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatal(err)
	}
	if string(recovered) != content {
		t.Errorf("recovered = %q, want %q", string(recovered), content)
	}
}

func TestShow_OutFlagOnPlainEntry(t *testing.T) {
	dir := setupEncrypted(t)

	if _, err := runCmd(newPadCmd("add", "plain note")); err != nil {
		t.Fatal(err)
	}

	outFile := filepath.Join(dir, "out.txt")
	_, err := runCmd(newPadCmd("show", "1", "--out", outFile))
	if err == nil {
		t.Fatal("expected error for --out on plain entry")
	}
	if !strings.Contains(err.Error(), "blob") {
		t.Errorf("error = %q, want mention of 'blob'", err.Error())
	}
}

// --- Blob edit tests ---

func TestEdit_BlobReplaceFile(t *testing.T) {
	dir := setupEncrypted(t)

	// Add a blob entry.
	v1 := filepath.Join(dir, "v1.txt")
	if err := os.WriteFile(v1, []byte("version 1"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := runCmd(newPadCmd("add", "--file", v1, "my blob")); err != nil {
		t.Fatal(err)
	}

	// Replace the file content.
	v2 := filepath.Join(dir, "v2.txt")
	if err := os.WriteFile(v2, []byte("version 2"), 0600); err != nil {
		t.Fatal(err)
	}
	out, err := runCmd(newPadCmd("edit", "1", "--file", v2))
	if err != nil {
		t.Fatalf("edit --file error: %v", err)
	}
	if !strings.Contains(out, "Updated entry 1.") {
		t.Errorf("output = %q", out)
	}

	// Verify content changed but label preserved.
	out, err = runCmd(newPadCmd("show", "1"))
	if err != nil {
		t.Fatal(err)
	}
	if out != "version 2" {
		t.Errorf("show = %q, want %q", out, "version 2")
	}

	listOut, err := runCmd(newPadCmd())
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(listOut, "my blob [BLOB]") {
		t.Errorf("list = %q, want label preserved", listOut)
	}
}

func TestEdit_BlobReplaceLabel(t *testing.T) {
	dir := setupEncrypted(t)

	v1 := filepath.Join(dir, "v1.txt")
	if err := os.WriteFile(v1, []byte("content"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := runCmd(newPadCmd("add", "--file", v1, "old label")); err != nil {
		t.Fatal(err)
	}

	out, err := runCmd(newPadCmd("edit", "1", "--label", "new label"))
	if err != nil {
		t.Fatalf("edit --label error: %v", err)
	}
	if !strings.Contains(out, "Updated entry 1.") {
		t.Errorf("output = %q", out)
	}

	// Verify label changed.
	listOut, err := runCmd(newPadCmd())
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(listOut, "new label [BLOB]") {
		t.Errorf("list = %q, want 'new label [BLOB]'", listOut)
	}

	// Verify content preserved.
	showOut, err := runCmd(newPadCmd("show", "1"))
	if err != nil {
		t.Fatal(err)
	}
	if showOut != "content" {
		t.Errorf("show = %q, want %q", showOut, "content")
	}
}

func TestEdit_BlobReplaceBoth(t *testing.T) {
	dir := setupEncrypted(t)

	v1 := filepath.Join(dir, "v1.txt")
	if err := os.WriteFile(v1, []byte("old content"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := runCmd(newPadCmd("add", "--file", v1, "old label")); err != nil {
		t.Fatal(err)
	}

	v2 := filepath.Join(dir, "v2.txt")
	if err := os.WriteFile(v2, []byte("new content"), 0600); err != nil {
		t.Fatal(err)
	}

	out, err := runCmd(newPadCmd("edit", "1", "--file", v2, "--label", "new label"))
	if err != nil {
		t.Fatalf("edit --file --label error: %v", err)
	}
	if !strings.Contains(out, "Updated entry 1.") {
		t.Errorf("output = %q", out)
	}

	listOut, err := runCmd(newPadCmd())
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(listOut, "new label [BLOB]") {
		t.Errorf("list = %q, want 'new label [BLOB]'", listOut)
	}

	showOut, err := runCmd(newPadCmd("show", "1"))
	if err != nil {
		t.Fatal(err)
	}
	if showOut != "new content" {
		t.Errorf("show = %q, want %q", showOut, "new content")
	}
}

func TestEdit_AppendOnBlobErrors(t *testing.T) {
	dir := setupEncrypted(t)

	testFile := filepath.Join(dir, "blob.txt")
	if err := os.WriteFile(testFile, []byte("content"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := runCmd(newPadCmd("add", "--file", testFile, "my blob")); err != nil {
		t.Fatal(err)
	}

	_, err := runCmd(newPadCmd("edit", "1", "--append", "suffix"))
	if err == nil {
		t.Fatal("expected error for --append on blob entry")
	}
	if !strings.Contains(err.Error(), "cannot append to a blob entry") {
		t.Errorf("error = %q, want 'cannot append to a blob entry'", err.Error())
	}
}

func TestEdit_PrependOnBlobErrors(t *testing.T) {
	dir := setupEncrypted(t)

	testFile := filepath.Join(dir, "blob.txt")
	if err := os.WriteFile(testFile, []byte("content"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := runCmd(newPadCmd("add", "--file", testFile, "my blob")); err != nil {
		t.Fatal(err)
	}

	_, err := runCmd(newPadCmd("edit", "1", "--prepend", "prefix"))
	if err == nil {
		t.Fatal("expected error for --prepend on blob entry")
	}
	if !strings.Contains(err.Error(), "cannot prepend to a blob entry") {
		t.Errorf("error = %q, want 'cannot prepend to a blob entry'", err.Error())
	}
}

func TestEdit_LabelOnNonBlobErrors(t *testing.T) {
	setupEncrypted(t)

	if _, err := runCmd(newPadCmd("add", "plain note")); err != nil {
		t.Fatal(err)
	}

	_, err := runCmd(newPadCmd("edit", "1", "--label", "new label"))
	if err == nil {
		t.Fatal("expected error for --label on non-blob entry")
	}
	if !strings.Contains(err.Error(), "not a blob entry") {
		t.Errorf("error = %q, want 'not a blob entry'", err.Error())
	}
}

func TestEdit_FileAndPositionalMutuallyExclusive(t *testing.T) {
	dir := setupEncrypted(t)

	testFile := filepath.Join(dir, "blob.txt")
	if err := os.WriteFile(testFile, []byte("content"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := runCmd(newPadCmd("add", "--file", testFile, "my blob")); err != nil {
		t.Fatal(err)
	}

	_, err := runCmd(newPadCmd("edit", "1", "replacement", "--file", testFile))
	if err == nil {
		t.Fatal("expected error for --file + positional text")
	}
	if !strings.Contains(err.Error(), "mutually exclusive") {
		t.Errorf("error = %q, want 'mutually exclusive'", err.Error())
	}
}

// --- Import tests ---

func TestImport_FromFile(t *testing.T) {
	dir := setupPlaintext(t)

	importFile := filepath.Join(dir, "notes.txt")
	if err := os.WriteFile(importFile, []byte("alpha\nbeta\ngamma\n"), 0600); err != nil {
		t.Fatal(err)
	}

	out, err := runCmd(newPadCmd("import", importFile))
	if err != nil {
		t.Fatalf("import error: %v", err)
	}
	if !strings.Contains(out, "Imported 3 entries.") {
		t.Errorf("output = %q, want 'Imported 3 entries.'", out)
	}

	// Verify entries
	out, err = runCmd(newPadCmd())
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range []string{"alpha", "beta", "gamma"} {
		if !strings.Contains(out, e) {
			t.Errorf("list missing entry %q", e)
		}
	}
}

func TestImport_SkipsEmpty(t *testing.T) {
	dir := setupPlaintext(t)

	importFile := filepath.Join(dir, "notes.txt")
	if err := os.WriteFile(importFile, []byte("alpha\n\n\nbeta\n\n"), 0600); err != nil {
		t.Fatal(err)
	}

	out, err := runCmd(newPadCmd("import", importFile))
	if err != nil {
		t.Fatalf("import error: %v", err)
	}
	if !strings.Contains(out, "Imported 2 entries.") {
		t.Errorf("output = %q, want 'Imported 2 entries.'", out)
	}
}

func TestImport_EmptyFile(t *testing.T) {
	dir := setupPlaintext(t)

	importFile := filepath.Join(dir, "empty.txt")
	if err := os.WriteFile(importFile, []byte(""), 0600); err != nil {
		t.Fatal(err)
	}

	out, err := runCmd(newPadCmd("import", importFile))
	if err != nil {
		t.Fatalf("import error: %v", err)
	}
	if !strings.Contains(out, "No entries to import.") {
		t.Errorf("output = %q, want 'No entries to import.'", out)
	}
}

func TestImport_AppendsToExisting(t *testing.T) {
	dir := setupPlaintext(t)

	// Add 2 entries first
	for _, e := range []string{"existing1", "existing2"} {
		if _, err := runCmd(newPadCmd("add", e)); err != nil {
			t.Fatal(err)
		}
	}

	importFile := filepath.Join(dir, "notes.txt")
	if err := os.WriteFile(importFile, []byte("new1\nnew2\nnew3\n"), 0600); err != nil {
		t.Fatal(err)
	}

	out, err := runCmd(newPadCmd("import", importFile))
	if err != nil {
		t.Fatalf("import error: %v", err)
	}
	if !strings.Contains(out, "Imported 3 entries.") {
		t.Errorf("output = %q, want 'Imported 3 entries.'", out)
	}

	// Verify all 5 entries exist
	out, err = runCmd(newPadCmd())
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range []string{"existing1", "existing2", "new1", "new2", "new3"} {
		if !strings.Contains(out, e) {
			t.Errorf("list missing entry %q", e)
		}
	}
}

func TestImport_Stdin(t *testing.T) {
	setupPlaintext(t)

	// Create a pipe to simulate stdin
	pr, pw, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	// Write data to the pipe
	go func() {
		_, _ = pw.WriteString("from stdin\nanother line\n")
		if err := pw.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "warning: close pipe writer: %v\n", err)
		}
	}()

	// Temporarily replace stdin
	origStdin := os.Stdin
	os.Stdin = pr
	t.Cleanup(func() { os.Stdin = origStdin })

	out, runErr := runCmd(newPadCmd("import", "-"))
	if runErr != nil {
		t.Fatalf("import stdin error: %v", runErr)
	}
	if !strings.Contains(out, "Imported 2 entries.") {
		t.Errorf("output = %q, want 'Imported 2 entries.'", out)
	}
}

func TestImport_FileNotFound(t *testing.T) {
	setupPlaintext(t)

	_, err := runCmd(newPadCmd("import", "/nonexistent/file.txt"))
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestImport_Encrypted(t *testing.T) {
	dir := setupEncrypted(t)

	importFile := filepath.Join(dir, "notes.txt")
	if err := os.WriteFile(importFile, []byte("secret1\nsecret2\n"), 0600); err != nil {
		t.Fatal(err)
	}

	out, err := runCmd(newPadCmd("import", importFile))
	if err != nil {
		t.Fatalf("import error: %v", err)
	}
	if !strings.Contains(out, "Imported 2 entries.") {
		t.Errorf("output = %q, want 'Imported 2 entries.'", out)
	}

	// Verify entries
	out, err = runCmd(newPadCmd())
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "secret1") || !strings.Contains(out, "secret2") {
		t.Errorf("list missing entries: %q", out)
	}
}

func TestImport_WhitespaceOnly(t *testing.T) {
	dir := setupPlaintext(t)

	importFile := filepath.Join(dir, "blanks.txt")
	if err := os.WriteFile(importFile, []byte("   \n\t\n  \t  \n"), 0600); err != nil {
		t.Fatal(err)
	}

	out, err := runCmd(newPadCmd("import", importFile))
	if err != nil {
		t.Fatalf("import error: %v", err)
	}
	if !strings.Contains(out, "No entries to import.") {
		t.Errorf("output = %q, want 'No entries to import.'", out)
	}
}

// --- Export tests ---

func TestExport_Basic(t *testing.T) {
	dir := setupPlaintext(t)

	// Add a plain entry and two blobs
	if _, err := runCmd(newPadCmd("add", "plain note")); err != nil {
		t.Fatal(err)
	}
	f1 := filepath.Join(dir, "file1.txt")
	if err := os.WriteFile(f1, []byte("content one"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := runCmd(newPadCmd("add", "--file", f1, "blob1.txt")); err != nil {
		t.Fatal(err)
	}
	f2 := filepath.Join(dir, "file2.md")
	if err := os.WriteFile(f2, []byte("content two"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := runCmd(newPadCmd("add", "--file", f2, "blob2.md")); err != nil {
		t.Fatal(err)
	}

	exportDir := filepath.Join(dir, "export")
	out, err := runCmd(newPadCmd("export", exportDir))
	if err != nil {
		t.Fatalf("export error: %v", err)
	}
	if !strings.Contains(out, "Exported 2 blobs.") {
		t.Errorf("output = %q, want 'Exported 2 blobs.'", out)
	}

	// Verify files
	data1, err := os.ReadFile(filepath.Join(exportDir, "blob1.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if string(data1) != "content one" {
		t.Errorf("blob1.txt = %q, want %q", string(data1), "content one")
	}

	data2, err := os.ReadFile(filepath.Join(exportDir, "blob2.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(data2) != "content two" {
		t.Errorf("blob2.md = %q, want %q", string(data2), "content two")
	}
}

func TestExport_EmptyPad(t *testing.T) {
	setupPlaintext(t)

	out, err := runCmd(newPadCmd("export"))
	if err != nil {
		t.Fatalf("export error: %v", err)
	}
	if !strings.Contains(out, "No blob entries to export.") {
		t.Errorf("output = %q, want 'No blob entries to export.'", out)
	}
}

func TestExport_NoBlobsOnly(t *testing.T) {
	setupPlaintext(t)

	if _, err := runCmd(newPadCmd("add", "plain one")); err != nil {
		t.Fatal(err)
	}
	if _, err := runCmd(newPadCmd("add", "plain two")); err != nil {
		t.Fatal(err)
	}

	out, err := runCmd(newPadCmd("export"))
	if err != nil {
		t.Fatalf("export error: %v", err)
	}
	if !strings.Contains(out, "No blob entries to export.") {
		t.Errorf("output = %q, want 'No blob entries to export.'", out)
	}
}

func TestExport_CollisionTimestamp(t *testing.T) {
	dir := setupPlaintext(t)

	// Add a blob
	f := filepath.Join(dir, "file.txt")
	if err := os.WriteFile(f, []byte("blob data"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := runCmd(newPadCmd("add", "--file", f, "existing.txt")); err != nil {
		t.Fatal(err)
	}

	exportDir := filepath.Join(dir, "export")
	if err := os.MkdirAll(exportDir, 0o750); err != nil {
		t.Fatal(err)
	}

	// Create a file at the expected path to cause collision
	if err := os.WriteFile(filepath.Join(exportDir, "existing.txt"), []byte("old"), 0600); err != nil {
		t.Fatal(err)
	}

	out, err := runCmd(newPadCmd("export", exportDir))
	if err != nil {
		t.Fatalf("export error: %v", err)
	}
	if !strings.Contains(out, "! existing.txt exists, writing as") {
		t.Errorf("output = %q, want collision warning", out)
	}
	if !strings.Contains(out, "Exported 1 blobs.") {
		t.Errorf("output = %q, want 'Exported 1 blobs.'", out)
	}

	// Verify old file is untouched
	oldData, _ := os.ReadFile(filepath.Join(exportDir, "existing.txt"))
	if string(oldData) != "old" {
		t.Errorf("existing file should not be overwritten, got %q", string(oldData))
	}
}

func TestExport_Force(t *testing.T) {
	dir := setupPlaintext(t)

	f := filepath.Join(dir, "file.txt")
	if err := os.WriteFile(f, []byte("new data"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := runCmd(newPadCmd("add", "--file", f, "target.txt")); err != nil {
		t.Fatal(err)
	}

	exportDir := filepath.Join(dir, "export")
	if err := os.MkdirAll(exportDir, 0o750); err != nil {
		t.Fatal(err)
	}

	// Create existing file
	if err := os.WriteFile(filepath.Join(exportDir, "target.txt"), []byte("old data"), 0600); err != nil {
		t.Fatal(err)
	}

	out, err := runCmd(newPadCmd("export", "--force", exportDir))
	if err != nil {
		t.Fatalf("export error: %v", err)
	}
	if !strings.Contains(out, "+ target.txt") {
		t.Errorf("output = %q, want '+ target.txt'", out)
	}

	// Verify file was overwritten
	data, _ := os.ReadFile(filepath.Join(exportDir, "target.txt"))
	if string(data) != "new data" {
		t.Errorf("target.txt = %q, want %q", string(data), "new data")
	}
}

func TestExport_DryRun(t *testing.T) {
	dir := setupPlaintext(t)

	f := filepath.Join(dir, "file.txt")
	if err := os.WriteFile(f, []byte("content"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := runCmd(newPadCmd("add", "--file", f, "test.txt")); err != nil {
		t.Fatal(err)
	}

	exportDir := filepath.Join(dir, "export")

	out, err := runCmd(newPadCmd("export", "--dry-run", exportDir))
	if err != nil {
		t.Fatalf("export error: %v", err)
	}
	if !strings.Contains(out, "Would export 1 blobs.") {
		t.Errorf("output = %q, want 'Would export 1 blobs.'", out)
	}

	// Verify directory was NOT created
	if _, err := os.Stat(exportDir); !os.IsNotExist(err) {
		t.Error("export directory should not be created in dry-run mode")
	}
}

func TestExport_DirCreated(t *testing.T) {
	dir := setupPlaintext(t)

	f := filepath.Join(dir, "file.txt")
	if err := os.WriteFile(f, []byte("data"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := runCmd(newPadCmd("add", "--file", f, "blob.txt")); err != nil {
		t.Fatal(err)
	}

	exportDir := filepath.Join(dir, "nested", "export", "dir")
	out, err := runCmd(newPadCmd("export", exportDir))
	if err != nil {
		t.Fatalf("export error: %v", err)
	}
	if !strings.Contains(out, "Exported 1 blobs.") {
		t.Errorf("output = %q, want 'Exported 1 blobs.'", out)
	}

	// Verify the directory was created
	if _, err := os.Stat(exportDir); err != nil {
		t.Errorf("export dir should exist: %v", err)
	}
}

func TestExport_Encrypted(t *testing.T) {
	dir := setupEncrypted(t)

	f := filepath.Join(dir, "secret.txt")
	if err := os.WriteFile(f, []byte("secret content"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := runCmd(newPadCmd("add", "--file", f, "secret.txt")); err != nil {
		t.Fatal(err)
	}

	exportDir := filepath.Join(dir, "export")
	out, err := runCmd(newPadCmd("export", exportDir))
	if err != nil {
		t.Fatalf("export error: %v", err)
	}
	if !strings.Contains(out, "Exported 1 blobs.") {
		t.Errorf("output = %q, want 'Exported 1 blobs.'", out)
	}

	data, err := os.ReadFile(filepath.Join(exportDir, "secret.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "secret content" {
		t.Errorf("exported = %q, want %q", string(data), "secret content")
	}
}

func TestExport_FilePermissions(t *testing.T) {
	dir := setupPlaintext(t)

	f := filepath.Join(dir, "file.txt")
	if err := os.WriteFile(f, []byte("data"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := runCmd(newPadCmd("add", "--file", f, "blob.txt")); err != nil {
		t.Fatal(err)
	}

	exportDir := filepath.Join(dir, "export")
	if _, err := runCmd(newPadCmd("export", exportDir)); err != nil {
		t.Fatal(err)
	}

	info, err := os.Stat(filepath.Join(exportDir, "blob.txt"))
	if err != nil {
		t.Fatal(err)
	}
	perm := info.Mode().Perm()
	if perm != 0o600 {
		t.Errorf("file perm = %o, want 600", perm)
	}
}

// --- Merge tests ---

// writePlaintextPad writes a plaintext scratchpad file at the given path.
func writePlaintextPad(t *testing.T, path string, entries []string) {
	t.Helper()
	content := strings.Join(entries, "\n") + "\n"
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
}

// writeEncryptedPad writes an encrypted scratchpad file at the given path
// using the provided key.
func writeEncryptedPad(t *testing.T, path string, key []byte, entries []string) {
	t.Helper()
	content := strings.Join(entries, "\n") + "\n"
	ciphertext, encErr := crypto.Encrypt(key, []byte(content))
	if encErr != nil {
		t.Fatal(encErr)
	}
	if err := os.WriteFile(path, ciphertext, 0600); err != nil {
		t.Fatal(err)
	}
}

func TestMerge_Basic(t *testing.T) {
	dir := setupPlaintext(t)

	// Add existing entries to current pad.
	for _, e := range []string{"existing1", "existing2"} {
		if _, err := runCmd(newPadCmd("add", e)); err != nil {
			t.Fatal(err)
		}
	}

	// Create a plaintext file with 3 entries (1 duplicate).
	mergeFile := filepath.Join(dir, "other.md")
	writePlaintextPad(t, mergeFile, []string{"existing1", "new1", "new2"})

	out, err := runCmd(newPadCmd("merge", mergeFile))
	if err != nil {
		t.Fatalf("merge error: %v", err)
	}
	if !strings.Contains(out, "Merged 2 new entries (1 duplicate skipped).") {
		t.Errorf("output = %q, want merge summary", out)
	}

	// Verify all 4 unique entries exist.
	listOut, err := runCmd(newPadCmd())
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range []string{"existing1", "existing2", "new1", "new2"} {
		if !strings.Contains(listOut, e) {
			t.Errorf("list missing entry %q", e)
		}
	}
}

func TestMerge_AllDuplicates(t *testing.T) {
	dir := setupPlaintext(t)

	for _, e := range []string{"alpha", "beta"} {
		if _, err := runCmd(newPadCmd("add", e)); err != nil {
			t.Fatal(err)
		}
	}

	mergeFile := filepath.Join(dir, "dupes.md")
	writePlaintextPad(t, mergeFile, []string{"alpha", "beta"})

	out, err := runCmd(newPadCmd("merge", mergeFile))
	if err != nil {
		t.Fatalf("merge error: %v", err)
	}
	if !strings.Contains(out, "No new entries to merge (2 duplicates skipped).") {
		t.Errorf("output = %q, want no-new summary", out)
	}
}

func TestMerge_EmptyFile(t *testing.T) {
	dir := setupPlaintext(t)

	mergeFile := filepath.Join(dir, "empty.md")
	writePlaintextPad(t, mergeFile, []string{})

	out, err := runCmd(newPadCmd("merge", mergeFile))
	if err != nil {
		t.Fatalf("merge error: %v", err)
	}
	if !strings.Contains(out, "No entries to merge.") {
		t.Errorf("output = %q, want empty summary", out)
	}
}

func TestMerge_MultipleFiles(t *testing.T) {
	dir := setupPlaintext(t)

	if _, err := runCmd(newPadCmd("add", "existing")); err != nil {
		t.Fatal(err)
	}

	fileA := filepath.Join(dir, "a.md")
	writePlaintextPad(t, fileA, []string{"from-a", "shared"})

	fileB := filepath.Join(dir, "b.md")
	writePlaintextPad(t, fileB, []string{"from-b", "shared"})

	out, err := runCmd(newPadCmd("merge", fileA, fileB))
	if err != nil {
		t.Fatalf("merge error: %v", err)
	}
	// "shared" appears in both files; second occurrence is a dupe.
	if !strings.Contains(out, "Merged 3 new entries (1 duplicate skipped).") {
		t.Errorf("output = %q, want multi-file summary", out)
	}

	listOut, err := runCmd(newPadCmd())
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range []string{"existing", "from-a", "from-b", "shared"} {
		if !strings.Contains(listOut, e) {
			t.Errorf("list missing entry %q", e)
		}
	}
}

func TestMerge_EncryptedInput(t *testing.T) {
	dir := setupEncrypted(t)

	if _, err := runCmd(newPadCmd("add", "current")); err != nil {
		t.Fatal(err)
	}

	// Create encrypted file using the same project key.
	key, loadErr := crypto.LoadKey(
		filepath.Join(dir, config.DirContext, config.FileScratchpadKey),
	)
	if loadErr != nil {
		t.Fatal(loadErr)
	}

	encFile := filepath.Join(dir, "other.enc")
	writeEncryptedPad(t, encFile, key, []string{"encrypted-entry"})

	out, err := runCmd(newPadCmd("merge", encFile))
	if err != nil {
		t.Fatalf("merge error: %v", err)
	}
	if !strings.Contains(out, "Merged 1 new entry") {
		t.Errorf("output = %q, want encrypted merge", out)
	}

	listOut, err := runCmd(newPadCmd())
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(listOut, "encrypted-entry") {
		t.Errorf("list missing encrypted-entry: %q", listOut)
	}
}

func TestMerge_PlaintextFallback(t *testing.T) {
	dir := setupPlaintext(t)

	// Plaintext file will fail decryption (no key) and fall back.
	mergeFile := filepath.Join(dir, "notes.md")
	writePlaintextPad(t, mergeFile, []string{"fallback-entry"})

	out, err := runCmd(newPadCmd("merge", mergeFile))
	if err != nil {
		t.Fatalf("merge error: %v", err)
	}
	if !strings.Contains(out, "Merged 1 new entry") {
		t.Errorf("output = %q, want fallback merge", out)
	}
}

func TestMerge_MixedEncPlain(t *testing.T) {
	dir := setupEncrypted(t)

	key, loadErr := crypto.LoadKey(
		filepath.Join(dir, config.DirContext, config.FileScratchpadKey),
	)
	if loadErr != nil {
		t.Fatal(loadErr)
	}

	encFile := filepath.Join(dir, "enc.enc")
	writeEncryptedPad(t, encFile, key, []string{"from-enc"})

	plainFile := filepath.Join(dir, "plain.md")
	writePlaintextPad(t, plainFile, []string{"from-plain"})

	out, err := runCmd(newPadCmd("merge", encFile, plainFile))
	if err != nil {
		t.Fatalf("merge error: %v", err)
	}
	if !strings.Contains(out, "Merged 2 new entries") {
		t.Errorf("output = %q, want mixed merge", out)
	}
}

func TestMerge_DryRun(t *testing.T) {
	dir := setupPlaintext(t)

	if _, err := runCmd(newPadCmd("add", "existing")); err != nil {
		t.Fatal(err)
	}

	mergeFile := filepath.Join(dir, "notes.md")
	writePlaintextPad(t, mergeFile, []string{"existing", "new-entry"})

	out, err := runCmd(newPadCmd("merge", "--dry-run", mergeFile))
	if err != nil {
		t.Fatalf("merge error: %v", err)
	}
	if !strings.Contains(out, "Would merge 1 new entry") {
		t.Errorf("output = %q, want dry-run summary", out)
	}

	// Verify pad was NOT modified.
	listOut, err := runCmd(newPadCmd())
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(listOut, "new-entry") {
		t.Error("dry-run should not write entries")
	}
}

func TestMerge_CustomKey(t *testing.T) {
	dir := setupPlaintext(t)

	// Generate a foreign key.
	foreignKey, genErr := crypto.GenerateKey()
	if genErr != nil {
		t.Fatal(genErr)
	}
	foreignKeyFile := filepath.Join(dir, "foreign.key")
	if err := crypto.SaveKey(foreignKeyFile, foreignKey); err != nil {
		t.Fatal(err)
	}

	// Create encrypted file with the foreign key.
	encFile := filepath.Join(dir, "foreign.enc")
	writeEncryptedPad(t, encFile, foreignKey, []string{"foreign-secret"})

	out, err := runCmd(newPadCmd("merge", "--key", foreignKeyFile, encFile))
	if err != nil {
		t.Fatalf("merge error: %v", err)
	}
	if !strings.Contains(out, "Merged 1 new entry") {
		t.Errorf("output = %q, want custom key merge", out)
	}

	listOut, err := runCmd(newPadCmd())
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(listOut, "foreign-secret") {
		t.Errorf("list missing foreign-secret: %q", listOut)
	}
}

func TestMerge_BlobEntries(t *testing.T) {
	dir := setupPlaintext(t)

	blobEntry := makeBlob("test.txt", []byte("hello world"))
	if _, err := runCmd(newPadCmd("add", "text-entry")); err != nil {
		t.Fatal(err)
	}

	// Create file with same blob + a new blob.
	newBlob := makeBlob("new.txt", []byte("new content"))
	mergeFile := filepath.Join(dir, "blobs.md")
	writePlaintextPad(t, mergeFile, []string{blobEntry, newBlob})

	out, err := runCmd(newPadCmd("merge", mergeFile))
	if err != nil {
		t.Fatalf("merge error: %v", err)
	}
	if !strings.Contains(out, "Merged 2 new entries") {
		t.Errorf("output = %q, want blob merge", out)
	}
	if !strings.Contains(out, "new.txt [BLOB]") {
		t.Errorf("output missing blob display: %q", out)
	}
}

func TestMerge_BlobConflict(t *testing.T) {
	dir := setupPlaintext(t)

	// Add a blob with label "config.json".
	blob1 := makeBlob("config.json", []byte(`{"v":1}`))
	mergeFile1 := filepath.Join(dir, "first.md")
	writePlaintextPad(t, mergeFile1, []string{blob1})
	if _, err := runCmd(newPadCmd("merge", mergeFile1)); err != nil {
		t.Fatal(err)
	}

	// Merge a different blob with the same label.
	blob2 := makeBlob("config.json", []byte(`{"v":2}`))
	mergeFile2 := filepath.Join(dir, "second.md")
	writePlaintextPad(t, mergeFile2, []string{blob2})

	out, err := runCmd(newPadCmd("merge", mergeFile2))
	if err != nil {
		t.Fatalf("merge error: %v", err)
	}
	if !strings.Contains(out, "different content across sources") {
		t.Errorf("output missing conflict warning: %q", out)
	}
	if !strings.Contains(out, "Merged 1 new entry") {
		t.Errorf("conflicting blob should still be added: %q", out)
	}
}

func TestMerge_BinaryWarning(t *testing.T) {
	dir := setupPlaintext(t)

	// Write raw binary data (not valid UTF-8).
	binFile := filepath.Join(dir, "garbage.bin")
	if err := os.WriteFile(binFile, []byte{0xff, 0xfe, 0x00, 0x01, 0x80, 0x90}, 0600); err != nil {
		t.Fatal(err)
	}

	out, err := runCmd(newPadCmd("merge", binFile))
	if err != nil {
		t.Fatalf("merge error: %v", err)
	}
	if !strings.Contains(out, "appears to contain binary data") {
		t.Errorf("output missing binary warning: %q", out)
	}
}

func TestMerge_FileNotFound(t *testing.T) {
	setupPlaintext(t)

	_, err := runCmd(newPadCmd("merge", "/nonexistent/file.md"))
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestMerge_EmptyPadMerge(t *testing.T) {
	dir := setupPlaintext(t)

	// Current pad is empty; merge entries into it.
	mergeFile := filepath.Join(dir, "fresh.md")
	writePlaintextPad(t, mergeFile, []string{"alpha", "beta", "gamma"})

	out, err := runCmd(newPadCmd("merge", mergeFile))
	if err != nil {
		t.Fatalf("merge error: %v", err)
	}
	if !strings.Contains(out, "Merged 3 new entries (0 duplicates skipped).") {
		t.Errorf("output = %q, want empty pad merge", out)
	}
}

func TestMerge_PlaintextMode(t *testing.T) {
	dir := setupPlaintext(t)

	if _, err := runCmd(newPadCmd("add", "plaintext-existing")); err != nil {
		t.Fatal(err)
	}

	mergeFile := filepath.Join(dir, "notes.md")
	writePlaintextPad(t, mergeFile, []string{"plaintext-new"})

	out, err := runCmd(newPadCmd("merge", mergeFile))
	if err != nil {
		t.Fatalf("merge error: %v", err)
	}
	if !strings.Contains(out, "Merged 1 new entry") {
		t.Errorf("output = %q, want plaintext merge", out)
	}

	// Verify the scratchpad.md file is plaintext.
	padPath := filepath.Join(dir, config.DirContext, config.FileScratchpadMd)
	data, readErr := os.ReadFile(padPath)
	if readErr != nil {
		t.Fatal(readErr)
	}
	content := string(data)
	if !strings.Contains(content, "plaintext-existing") ||
		!strings.Contains(content, "plaintext-new") {
		t.Errorf("pad content = %q, missing entries", content)
	}
}

func TestMerge_PreservesOrder(t *testing.T) {
	dir := setupPlaintext(t)

	// Add entries in specific order.
	for _, e := range []string{"first", "second", "third"} {
		if _, err := runCmd(newPadCmd("add", e)); err != nil {
			t.Fatal(err)
		}
	}

	mergeFile := filepath.Join(dir, "new.md")
	writePlaintextPad(t, mergeFile, []string{"fourth", "fifth"})

	if _, err := runCmd(newPadCmd("merge", mergeFile)); err != nil {
		t.Fatal(err)
	}

	// Read the raw pad and verify order.
	padPath := filepath.Join(dir, config.DirContext, config.FileScratchpadMd)
	data, readErr := os.ReadFile(padPath)
	if readErr != nil {
		t.Fatal(readErr)
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	expected := []string{"first", "second", "third", "fourth", "fifth"}
	if len(lines) != len(expected) {
		t.Fatalf("got %d lines, want %d: %v", len(lines), len(expected), lines)
	}
	for i, want := range expected {
		if lines[i] != want {
			t.Errorf("line %d = %q, want %q", i, lines[i], want)
		}
	}
}

func TestMerge_CrossFileDedup(t *testing.T) {
	dir := setupPlaintext(t)

	// Merge two files where entries overlap across files AND with current pad.
	if _, err := runCmd(newPadCmd("add", "base")); err != nil {
		t.Fatal(err)
	}

	fileA := filepath.Join(dir, "a.md")
	writePlaintextPad(t, fileA, []string{"base", "unique-a", "shared-ab"})

	fileB := filepath.Join(dir, "b.md")
	writePlaintextPad(t, fileB, []string{"shared-ab", "unique-b"})

	out, err := runCmd(newPadCmd("merge", fileA, fileB))
	if err != nil {
		t.Fatalf("merge error: %v", err)
	}
	// base: dupe (in pad), shared-ab from A: new, shared-ab from B: dupe
	// unique-a: new, unique-b: new
	if !strings.Contains(out, "Merged 3 new entries (2 duplicates skipped).") {
		t.Errorf("output = %q, want cross-file dedup summary", out)
	}
}

func TestMerge_EncryptedWithBlobDedup(t *testing.T) {
	dir := setupEncrypted(t)

	// Add a blob to the current pad.
	blob := makeBlob("readme.md", []byte("# README"))
	f := filepath.Join(dir, "tmp-readme.md")
	if err := os.WriteFile(f, []byte("# README"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := runCmd(newPadCmd("add", "--file", f, "readme.md")); err != nil {
		t.Fatal(err)
	}

	// Get the project key.
	key, loadErr := crypto.LoadKey(
		filepath.Join(dir, config.DirContext, config.FileScratchpadKey),
	)
	if loadErr != nil {
		t.Fatal(loadErr)
	}

	// Create encrypted file with the same blob.
	encFile := filepath.Join(dir, "merge.enc")
	writeEncryptedPad(t, encFile, key, []string{blob, "new-text"})

	out, err := runCmd(newPadCmd("merge", encFile))
	if err != nil {
		t.Fatalf("merge error: %v", err)
	}
	// blob is duplicate, "new-text" is new.
	if !strings.Contains(out, "Merged 1 new entry (1 duplicate skipped).") {
		t.Errorf("output = %q, want encrypted blob dedup", out)
	}
}

func TestPluralize(t *testing.T) {
	tests := []struct {
		word  string
		count int
		want  string
	}{
		{"entry", 1, "entry"},
		{"entry", 0, "entries"},
		{"entry", 2, "entries"},
		{"duplicate", 1, "duplicate"},
		{"duplicate", 3, "duplicates"},
	}

	for _, tt := range tests {
		got := pluralize(tt.word, tt.count)
		if got != tt.want {
			t.Errorf(
				"pluralize(%q, %d) = %q, want %q",
				tt.word, tt.count, got, tt.want,
			)
		}
	}
}

// Verify unused import doesn't cause issues.
var _ = base64.StdEncoding
