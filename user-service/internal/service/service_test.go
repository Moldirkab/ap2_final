package service

import (
	"context"
	"testing"
	"user-service/internal/auth"

	pb "github.com/Moldirkab/ap2_final_user_service_generated"
)

func TestValidateToken_ValidToken(t *testing.T) {
	token, err := auth.GenerateToken("1", "customer")
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	svc := &AuthServiceMock{}
	res, err := svc.ValidateToken(context.Background(), &pb.TokenRequest{Token: token})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !res.Valid {
		t.Error("expected valid=true")
	}
	if res.UserId != "1" {
		t.Errorf("expected user_id '1', got %s", res.UserId)
	}
	if res.Role != "customer" {
		t.Errorf("expected role 'customer', got %s", res.Role)
	}
}

func TestValidateToken_InvalidToken(t *testing.T) {
	svc := &AuthServiceMock{}
	res, err := svc.ValidateToken(context.Background(), &pb.TokenRequest{Token: "bad.token"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res.Valid {
		t.Error("expected valid=false for invalid token")
	}
}

func TestValidateToken_EmptyToken(t *testing.T) {
	svc := &AuthServiceMock{}
	res, err := svc.ValidateToken(context.Background(), &pb.TokenRequest{Token: ""})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res.Valid {
		t.Error("expected valid=false for empty token")
	}
}

// AuthServiceMock replicates ValidateToken logic without DB dependency
type AuthServiceMock struct{}

func (s *AuthServiceMock) ValidateToken(ctx context.Context, req *pb.TokenRequest) (*pb.ValidateResponse, error) {
	claims, err := auth.VerifyToken(req.Token)
	if err != nil {
		return &pb.ValidateResponse{Valid: false}, nil
	}

	userID, ok1 := claims["user_id"].(string)
	role, ok2 := claims["role"].(string)

	if !ok1 || !ok2 {
		return &pb.ValidateResponse{Valid: false}, nil
	}

	return &pb.ValidateResponse{
		UserId: userID,
		Role:   role,
		Valid:  true,
	}, nil
}
