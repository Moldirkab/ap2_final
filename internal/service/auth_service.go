package service

import (
	"context"
	"strconv"
	"user-service/internal/auth"
	"user-service/internal/database"
	"user-service/internal/email"

	pb "github.com/Moldirkab/ap2_final_user_service_generated"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	pb.UnimplementedUserServiceServer
}

func (s *AuthService) Register(
	ctx context.Context,
	req *pb.RegisterRequest,
) (*pb.AuthResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, err
	}

	var userID int
	err = database.DB.QueryRow(
		`INSERT INTO users(email, password_hash) VALUES($1, $2) RETURNING id`,
		req.Email,
		string(hashedPassword),
	).Scan(&userID)
	if err != nil {
		return nil, err
	}

	token, err := auth.GenerateToken(
		strconv.Itoa(userID),
		"customer",
	)
	if err != nil {
		return nil, err
	}
	
	go email.SendWelcomeEmail(req.Email)
	return &pb.AuthResponse{Token: token}, nil
}

func (s *AuthService) Login(
	ctx context.Context,
	req *pb.LoginRequest,
) (*pb.AuthResponse, error) {
	var id int
	var email string
	var hashedPassword string
	var role string

	err := database.DB.QueryRow(
		`SELECT id, email, password_hash, role FROM users WHERE email = $1`,
		req.Email,
	).Scan(&id, &email, &hashedPassword, &role)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(req.Password),
	)
	if err != nil {
		return nil, err
	}

	token, err := auth.GenerateToken(
		strconv.Itoa(id),
		role,
	)
	if err != nil {
		return nil, err
	}

	return &pb.AuthResponse{Token: token}, nil
}

func (s *AuthService) ValidateToken(
	ctx context.Context,
	req *pb.TokenRequest,
) (*pb.ValidateResponse, error) {
	claims, err := auth.VerifyToken(req.Token)
	if err != nil {
		return &pb.ValidateResponse{
			Valid: false,
		}, nil
	}

	userID := claims["user_id"].(string)
	role := claims["role"].(string)

	return &pb.ValidateResponse{
		UserId: userID,
		Role:   role,
		Valid:  true,
	}, nil
}
