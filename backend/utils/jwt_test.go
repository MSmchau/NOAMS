package utils

import (
	"testing"
	"time"
)

func init() {
	JWTSecret = []byte("test-secret-key-for-unit-tests")
}

func TestGenerateToken(t *testing.T) {
	token, err := GenerateToken(1, "admin", "admin", 24)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}
	if token == "" {
		t.Fatal("GenerateToken returned empty token")
	}
}

func TestParseToken(t *testing.T) {
	token, err := GenerateToken(42, "testuser", "operator", 24)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	claims, err := ParseToken(token)
	if err != nil {
		t.Fatalf("ParseToken failed: %v", err)
	}

	if claims.UserID != 42 {
		t.Fatalf("expected UserID 42, got %d", claims.UserID)
	}
	if claims.Username != "testuser" {
		t.Fatalf("expected Username testuser, got %s", claims.Username)
	}
	if claims.Role != "operator" {
		t.Fatalf("expected Role operator, got %s", claims.Role)
	}
}

func TestParseToken_Expired(t *testing.T) {
	token, err := GenerateToken(1, "user", "operator", -1) // negative = expired
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	_, err = ParseToken(token)
	if err == nil {
		t.Fatal("ParseToken should fail for expired token")
	}
}

func TestParseToken_Invalid(t *testing.T) {
	_, err := ParseToken("invalid-token-string")
	if err == nil {
		t.Fatal("ParseToken should fail for invalid token")
	}
}

func TestParseToken_WrongSecret(t *testing.T) {
	oldSecret := JWTSecret
	JWTSecret = []byte("original-secret")
	token, err := GenerateToken(1, "user", "admin", 24)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	JWTSecret = []byte("different-secret")
	_, err = ParseToken(token)
	if err == nil {
		t.Fatal("ParseToken should fail with wrong secret")
	}
	JWTSecret = oldSecret
}

func TestGenerateToken_Unique(t *testing.T) {
	t1, _ := GenerateToken(1, "user", "admin", 24)
	time.Sleep(2 * time.Second) // JWT iat has second precision
	t2, _ := GenerateToken(1, "user", "admin", 24)
	if t1 == t2 {
		t.Fatal("Successive tokens should be unique")
	}
}
