// Code generated by sdkgen. DO NOT EDIT.

//nolint
package compute

import (
	"context"

	"google.golang.org/grpc"

	compute "github.com/yandex-cloud/go-genproto/yandex/cloud/compute/v1"
	"github.com/yandex-cloud/go-genproto/yandex/cloud/operation"
)

//revive:disable

// ImageServiceClient is a compute.ImageServiceClient with
// lazy GRPC connection initialization.
type ImageServiceClient struct {
	getConn func(ctx context.Context) (*grpc.ClientConn, error)
}

var _ compute.ImageServiceClient = &ImageServiceClient{}

// Create implements compute.ImageServiceClient
func (c *ImageServiceClient) Create(ctx context.Context, in *compute.CreateImageRequest, opts ...grpc.CallOption) (*operation.Operation, error) {
	conn, err := c.getConn(ctx)
	if err != nil {
		return nil, err
	}
	return compute.NewImageServiceClient(conn).Create(ctx, in, opts...)
}

// Delete implements compute.ImageServiceClient
func (c *ImageServiceClient) Delete(ctx context.Context, in *compute.DeleteImageRequest, opts ...grpc.CallOption) (*operation.Operation, error) {
	conn, err := c.getConn(ctx)
	if err != nil {
		return nil, err
	}
	return compute.NewImageServiceClient(conn).Delete(ctx, in, opts...)
}

// Get implements compute.ImageServiceClient
func (c *ImageServiceClient) Get(ctx context.Context, in *compute.GetImageRequest, opts ...grpc.CallOption) (*compute.Image, error) {
	conn, err := c.getConn(ctx)
	if err != nil {
		return nil, err
	}
	return compute.NewImageServiceClient(conn).Get(ctx, in, opts...)
}

// GetLatestByFamily implements compute.ImageServiceClient
func (c *ImageServiceClient) GetLatestByFamily(ctx context.Context, in *compute.GetImageLatestByFamilyRequest, opts ...grpc.CallOption) (*compute.Image, error) {
	conn, err := c.getConn(ctx)
	if err != nil {
		return nil, err
	}
	return compute.NewImageServiceClient(conn).GetLatestByFamily(ctx, in, opts...)
}

// List implements compute.ImageServiceClient
func (c *ImageServiceClient) List(ctx context.Context, in *compute.ListImagesRequest, opts ...grpc.CallOption) (*compute.ListImagesResponse, error) {
	conn, err := c.getConn(ctx)
	if err != nil {
		return nil, err
	}
	return compute.NewImageServiceClient(conn).List(ctx, in, opts...)
}

// ListOperations implements compute.ImageServiceClient
func (c *ImageServiceClient) ListOperations(ctx context.Context, in *compute.ListImageOperationsRequest, opts ...grpc.CallOption) (*compute.ListImageOperationsResponse, error) {
	conn, err := c.getConn(ctx)
	if err != nil {
		return nil, err
	}
	return compute.NewImageServiceClient(conn).ListOperations(ctx, in, opts...)
}

// Update implements compute.ImageServiceClient
func (c *ImageServiceClient) Update(ctx context.Context, in *compute.UpdateImageRequest, opts ...grpc.CallOption) (*operation.Operation, error) {
	conn, err := c.getConn(ctx)
	if err != nil {
		return nil, err
	}
	return compute.NewImageServiceClient(conn).Update(ctx, in, opts...)
}
