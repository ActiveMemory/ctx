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
	for _, expected := range []string{"show N", "add TEXT", "rm N", "edit N [TEXT]", "mv N M", "resolve"} {
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
	data, err := os.ReadFile(path)
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
	if err := os.WriteFile(oursPath, oursCipher, 0600); err != nil {
		t.Fatal(err)
	}

	// Create encrypted "theirs" file
	theirsPlain := []byte("theirs-entry\n")
	theirsCipher, err := crypto.Encrypt(key, theirsPlain)
	if err != nil {
		t.Fatal(err)
	}
	theirsPath := filepath.Join(config.DirContext, config.FileScratchpadEnc+".theirs")
	if err := os.WriteFile(theirsPath, theirsCipher, 0600); err != nil {
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
	if err := os.WriteFile(oursPath, oursCipher, 0600); err != nil {
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
	if err := os.WriteFile(path, ciphertext, 0600); err != nil {
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
