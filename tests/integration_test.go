package integration_test

import (
	"context"
	"os"
	"testing"

	"user-service/internal/auth"
	"user-service/internal/database"
	"user-service/internal/service"

	pb "github.com/Moldirkab/ap2_final_user_service_generated"
	_ "github.com/lib/pq"
)

func setupTestDB(t *testing.T) {
	t.Helper()

	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_PASSWORD", "070906")
	os.Setenv("DB_NAME", "user_db")

	database.ConnectDB()

	_, err := database.DB.Exec("DELETE FROM users WHERE email LIKE 'test_%@test.com'")
	if err != nil {
		t.Fatalf("failed to clean test data: %v", err)
	}
}

func TestRegister_Success(t *testing.T) {
	setupTestDB(t)

	svc := &service.AuthService{}
	res, err := svc.Register(context.Background(), &pb.RegisterRequest{
		Email:    "test_register@test.com",
		Password: "password123",
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res.Token == "" {
		t.Error("expected token, got empty string")
	}

	claims, err := auth.VerifyToken(res.Token)
	if err != nil {
		t.Fatalf("expected valid token, got %v", err)
	}
	if claims["role"] != "customer" {
		t.Errorf("expected role 'customer', got %v", claims["role"])
	}
}

func TestRegister_DuplicateEmail(t *testing.T) {
	setupTestDB(t)

	svc := &service.AuthService{}

	_, err := svc.Register(context.Background(), &pb.RegisterRequest{
		Email:    "test_duplicate@test.com",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("first register failed: %v", err)
	}

	_, err = svc.Register(context.Background(), &pb.RegisterRequest{
		Email:    "test_duplicate@test.com",
		Password: "password456",
	})
	if err == nil {
		t.Fatal("expected error for duplicate email, got nil")
	}
}

func TestLogin_Success(t *testing.T) {
	setupTestDB(t)

	svc := &service.AuthService{}

	_, err := svc.Register(context.Background(), &pb.RegisterRequest{
		Email:    "test_login@test.com",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	res, err := svc.Login(context.Background(), &pb.LoginRequest{
		Email:    "test_login@test.com",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res.Token == "" {
		t.Error("expected token, got empty string")
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	setupTestDB(t)

	svc := &service.AuthService{}

	_, err := svc.Register(context.Background(), &pb.RegisterRequest{
		Email:    "test_wrongpass@test.com",
		Password: "correctpassword",
	})
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	_, err = svc.Login(context.Background(), &pb.LoginRequest{
		Email:    "test_wrongpass@test.com",
		Password: "wrongpassword",
	})
	if err == nil {
		t.Fatal("expected error for wrong password, got nil")
	}
}

func TestLogin_UserNotFound(t *testing.T) {
	setupTestDB(t)

	svc := &service.AuthService{}
	_, err := svc.Login(context.Background(), &pb.LoginRequest{
		Email:    "nonexistent@test.com",
		Password: "password123",
	})
	if err == nil {
		t.Fatal("expected error for nonexistent user, got nil")
	}
}

func TestValidateToken_Integration(t *testing.T) {
	setupTestDB(t)

	svc := &service.AuthService{}

	res, err := svc.Register(context.Background(), &pb.RegisterRequest{
		Email:    "test_validate@test.com",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	validateRes, err := svc.ValidateToken(context.Background(), &pb.TokenRequest{
		Token: res.Token,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !validateRes.Valid {
		t.Error("expected valid=true")
	}
	if validateRes.Role != "customer" {
		t.Errorf("expected role 'customer', got %s", validateRes.Role)
	}
}
