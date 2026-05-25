package services

import (
	"testing"

	"noams/config"
	"noams/models"
	"noams/utils"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	if err := db.AutoMigrate(&models.User{}); err != nil {
		t.Fatalf("failed to migrate user table: %v", err)
	}
	return db
}

func setupAuthService(db *gorm.DB) *AuthService {
	config.AppConfig = &config.Config{}
	config.AppConfig.JWT.Secret = "test-jwt-secret"
	config.AppConfig.JWT.ExpireHour = 24
	utils.JWTSecret = []byte("test-jwt-secret")
	return NewAuthService(db)
}

func TestLogin_Success(t *testing.T) {
	db := setupTestDB(t)
	svc := setupAuthService(db)

	svc.Register(&RegisterRequest{Username: "testuser", Password: "password123"})

	resp, err := svc.Login(&LoginRequest{Username: "testuser", Password: "password123"})
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}
	if resp.Token == "" {
		t.Fatal("Login returned empty token")
	}
	if resp.User.Username != "testuser" {
		t.Fatalf("expected username testuser, got %s", resp.User.Username)
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	db := setupTestDB(t)
	svc := setupAuthService(db)
	svc.Register(&RegisterRequest{Username: "user", Password: "correctpw"})

	_, err := svc.Login(&LoginRequest{Username: "user", Password: "wrongpw"})
	if err == nil {
		t.Fatal("Login with wrong password should fail")
	}
}

func TestLogin_UserNotFound(t *testing.T) {
	db := setupTestDB(t)
	svc := setupAuthService(db)

	_, err := svc.Login(&LoginRequest{Username: "nonexistent", Password: "pw"})
	if err == nil {
		t.Fatal("Login with nonexistent user should fail")
	}
}

func TestRegister_DuplicateUsername(t *testing.T) {
	db := setupTestDB(t)
	svc := setupAuthService(db)
	svc.Register(&RegisterRequest{Username: "dupuser", Password: "password123"})

	err := svc.Register(&RegisterRequest{Username: "dupuser", Password: "otherpw"})
	if err == nil {
		t.Fatal("Register with duplicate username should fail")
	}
}

func TestRegister_Success(t *testing.T) {
	db := setupTestDB(t)
	svc := setupAuthService(db)

	err := svc.Register(&RegisterRequest{
		Username: "newuser",
		Password: "securePass1!",
		Email:    "user@example.com",
		Nickname: "New User",
	})
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}
}

func TestGetUserInfo(t *testing.T) {
	db := setupTestDB(t)
	svc := setupAuthService(db)
	svc.Register(&RegisterRequest{Username: "infouser", Password: "pass123"})

	resp, _ := svc.Login(&LoginRequest{Username: "infouser", Password: "pass123"})

	info, err := svc.GetUserInfo(resp.User.ID)
	if err != nil {
		t.Fatalf("GetUserInfo failed: %v", err)
	}
	if info.Username != "infouser" {
		t.Fatalf("expected username infouser, got %s", info.Username)
	}
}
