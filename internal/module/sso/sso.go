package sso

import (
	"context"
	"errors"
	"fmt"
	ssov1 "github.com/synthao/orders/gen/go/sso"
	"google.golang.org/grpc"
)

var ErrIsAuthorized = errors.New("failed to check if user is authorized")

func newSSOClient(grpcConn *grpc.ClientConn) *Client {
	return &Client{
		api: ssov1.NewServiceClient(grpcConn),
	}
}

func (c *Client) IsAuthorized(token string) (bool, error) {
	res, err := c.api.IsAuthorized(context.Background(), &ssov1.IsAuthorizedRequest{Token: token})
	if err != nil {
		return false, fmt.Errorf("%w, %w", ErrIsAuthorized, err)
	}

	return res.IsAuthorized, nil
}
