package utils

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "testPassword123!"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}
	if hash == "" {
		t.Fatal("HashPassword returned empty hash")
	}
	if !CheckPassword(password, hash) {
		t.Fatal("CheckPassword should return true for correct password")
	}
	if CheckPassword("wrongPassword", hash) {
		t.Fatal("CheckPassword should return false for wrong password")
	}
}

func TestEncryptDecryptAES(t *testing.T) {
	key := "0123456789abcdef0123456789abcdef" // 32 bytes for AES-256
	plaintext := "sensitive-data-to-encrypt"

	ciphertext, err := EncryptAES(plaintext, key)
	if err != nil {
		t.Fatalf("EncryptAES failed: %v", err)
	}
	if ciphertext == "" {
		t.Fatal("EncryptAES returned empty ciphertext")
	}
	if ciphertext == plaintext {
		t.Fatal("EncryptAES should not return plaintext")
	}

	decrypted, err := DecryptAES(ciphertext, key)
	if err != nil {
		t.Fatalf("DecryptAES failed: %v", err)
	}
	if decrypted != plaintext {
		t.Fatalf("DecryptAES returned wrong data: got %s, want %s", decrypted, plaintext)
	}
}

func TestEncryptDecryptAES_WrongKey(t *testing.T) {
	key := "0123456789abcdef0123456789abcdef"
	wrongKey := "fedcba9876543210fedcba9876543210"
	plaintext := "test-data"

	ciphertext, err := EncryptAES(plaintext, key)
	if err != nil {
		t.Fatalf("EncryptAES failed: %v", err)
	}

	_, err = DecryptAES(ciphertext, wrongKey)
	if err == nil {
		t.Fatal("DecryptAES with wrong key should fail")
	}
}

func TestHashPassword_Empty(t *testing.T) {
	hash, err := HashPassword("")
	if err != nil {
		t.Fatalf("HashPassword with empty string failed: %v", err)
	}
	if !CheckPassword("", hash) {
		t.Fatal("CheckPassword should verify empty password hash")
	}
}

func TestEncryptDecryptAES_Empty(t *testing.T) {
	key := "0123456789abcdef0123456789abcdef"
	ciphertext, err := EncryptAES("", key)
	if err != nil {
		t.Fatalf("EncryptAES with empty text failed: %v", err)
	}
	decrypted, err := DecryptAES(ciphertext, key)
	if err != nil {
		t.Fatalf("DecryptAES failed: %v", err)
	}
	if decrypted != "" {
		t.Fatal("DecryptAES of empty should return empty")
	}
}
