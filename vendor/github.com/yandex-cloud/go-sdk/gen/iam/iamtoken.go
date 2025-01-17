// Code generated by sdkgen. DO NOT EDIT.

//nolint
package iam

import (
	"context"

	"google.golang.org/grpc"

	iam "github.com/yandex-cloud/go-genproto/yandex/cloud/iam/v1"
)

//revive:disable

// IamTokenServiceClient is a iam.IamTokenServiceClient with
// lazy GRPC connection initialization.
type IamTokenServiceClient struct {
	getConn func(ctx context.Context) (*grpc.ClientConn, error)
}

var _ iam.IamTokenServiceClient = &IamTokenServiceClient{}

// Create implements iam.IamTokenServiceClient
func (c *IamTokenServiceClient) Create(ctx context.Context, in *iam.CreateIamTokenRequest, opts ...grpc.CallOption) (*iam.CreateIamTokenResponse, error) {
	conn, err := c.getConn(ctx)
	if err != nil {
		return nil, err
	}
	return iam.NewIamTokenServiceClient(conn).Create(ctx, in, opts...)
}
