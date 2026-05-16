package clients

import (
	"context"
	"errors"

	userpb "github.com/Moldirkab/ap2_final_user_service_generated"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserClient struct {
	client userpb.UserServiceClient
}

func NewUserClient(address string) (*UserClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &UserClient{
		client: userpb.NewUserServiceClient(conn),
	}, nil
}

func (c *UserClient) ValidateToken(token string) (string, error) {
	res, err := c.client.ValidateToken(context.Background(), &userpb.TokenRequest{
		Token: token,
	})
	if err != nil {
		return "", err
	}

	if !res.Valid {
		return "", errors.New("invalid token")
	}

	return res.UserId, nil
}
