//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package crypto

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestGenerateKey(t *testing.T) {
	key, err := GenerateKey()
	if err != nil {
		t.Fatalf("GenerateKey() error: %v", err)
	}
	if len(key) != KeySize {
		t.Errorf("key length = %d, want %d", len(key), KeySize)
	}

	// Two keys should be different
	key2, err := GenerateKey()
	if err != nil {
		t.Fatalf("GenerateKey() second call error: %v", err)
	}
	if bytes.Equal(key, key2) {
		t.Error("two generated keys should not be equal")
	}
}

func TestEncryptDecrypt_RoundTrip(t *testing.T) {
	key, err := GenerateKey()
	if err != nil {
		t.Fatalf("GenerateKey() error: %v", err)
	}

	plaintext := []byte("remember to check DNS config on staging")

	ciphertext, err := Encrypt(key, plaintext)
	if err != nil {
		t.Fatalf("Encrypt() error: %v", err)
	}

	// Ciphertext should be longer than plaintext (nonce + tag)
	if len(ciphertext) <= len(plaintext) {
		t.Errorf("ciphertext length %d should be greater than plaintext length %d",
			len(ciphertext), len(plaintext))
	}

	decrypted, err := Decrypt(key, ciphertext)
	if err != nil {
		t.Fatalf("Decrypt() error: %v", err)
	}

	if !bytes.Equal(decrypted, plaintext) {
		t.Errorf("decrypted = %q, want %q", decrypted, plaintext)
	}
}

func TestEncryptDecrypt_EmptyPlaintext(t *testing.T) {
	key, err := GenerateKey()
	if err != nil {
		t.Fatalf("GenerateKey() error: %v", err)
	}

	ciphertext, err := Encrypt(key, []byte{})
	if err != nil {
		t.Fatalf("Encrypt() error: %v", err)
	}

	decrypted, err := Decrypt(key, ciphertext)
	if err != nil {
		t.Fatalf("Decrypt() error: %v", err)
	}

	if len(decrypted) != 0 {
		t.Errorf("decrypted length = %d, want 0", len(decrypted))
	}
}

func TestDecrypt_WrongKey(t *testing.T) {
	key1, _ := GenerateKey()
	key2, _ := GenerateKey()

	plaintext := []byte("secret note")
	ciphertext, err := Encrypt(key1, plaintext)
	if err != nil {
		t.Fatalf("Encrypt() error: %v", err)
	}

	_, err = Decrypt(key2, ciphertext)
	if err == nil {
		t.Error("Decrypt() with wrong key should return error")
	}
}

func TestDecrypt_ShortCiphertext(t *testing.T) {
	key, _ := GenerateKey()

	_, err := Decrypt(key, []byte("short"))
	if err == nil {
		t.Error("Decrypt() with short ciphertext should return error")
	}
	if err.Error() != "ciphertext too short" {
		t.Errorf("error = %q, want %q", err.Error(), "ciphertext too short")
	}
}

func TestDecrypt_TamperedCiphertext(t *testing.T) {
	key, _ := GenerateKey()
	plaintext := []byte("important data")

	ciphertext, err := Encrypt(key, plaintext)
	if err != nil {
		t.Fatalf("Encrypt() error: %v", err)
	}

	// Tamper with the ciphertext (flip a byte after the nonce)
	ciphertext[NonceSize+1] ^= 0xFF

	_, err = Decrypt(key, ciphertext)
	if err == nil {
		t.Error("Decrypt() with tampered ciphertext should return error")
	}
}

func TestSaveKey_LoadKey_RoundTrip(t *testing.T) {
	key, err := GenerateKey()
	if err != nil {
		t.Fatalf("GenerateKey() error: %v", err)
	}

	path := filepath.Join(t.TempDir(), "test.key")

	if saveErr := SaveKey(path, key); saveErr != nil {
		t.Fatalf("SaveKey() error: %v", saveErr)
	}

	// Verify file permissions
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat() error: %v", err)
	}
	if perm := info.Mode().Perm(); perm != 0600 {
		t.Errorf("key file permissions = %o, want 0600", perm)
	}

	loaded, err := LoadKey(path)
	if err != nil {
		t.Fatalf("LoadKey() error: %v", err)
	}

	if !bytes.Equal(loaded, key) {
		t.Error("loaded key does not match saved key")
	}
}

func TestLoadKey_WrongSize(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bad.key")
	if err := os.WriteFile(path, []byte("too short"), 0600); err != nil {
		t.Fatalf("WriteFile() error: %v", err)
	}

	_, err := LoadKey(path)
	if err == nil {
		t.Error("LoadKey() with wrong-size file should return error")
	}
}

func TestLoadKey_NotFound(t *testing.T) {
	_, err := LoadKey(filepath.Join(t.TempDir(), "nonexistent.key"))
	if err == nil {
		t.Error("LoadKey() with missing file should return error")
	}
}
