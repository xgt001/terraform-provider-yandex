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

// DiskServiceClient is a compute.DiskServiceClient with
// lazy GRPC connection initialization.
type DiskServiceClient struct {
	getConn func(ctx context.Context) (*grpc.ClientConn, error)
}

var _ compute.DiskServiceClient = &DiskServiceClient{}

// Create implements compute.DiskServiceClient
func (c *DiskServiceClient) Create(ctx context.Context, in *compute.CreateDiskRequest, opts ...grpc.CallOption) (*operation.Operation, error) {
	conn, err := c.getConn(ctx)
	if err != nil {
		return nil, err
	}
	return compute.NewDiskServiceClient(conn).Create(ctx, in, opts...)
}

// Delete implements compute.DiskServiceClient
func (c *DiskServiceClient) Delete(ctx context.Context, in *compute.DeleteDiskRequest, opts ...grpc.CallOption) (*operation.Operation, error) {
	conn, err := c.getConn(ctx)
	if err != nil {
		return nil, err
	}
	return compute.NewDiskServiceClient(conn).Delete(ctx, in, opts...)
}

// Get implements compute.DiskServiceClient
func (c *DiskServiceClient) Get(ctx context.Context, in *compute.GetDiskRequest, opts ...grpc.CallOption) (*compute.Disk, error) {
	conn, err := c.getConn(ctx)
	if err != nil {
		return nil, err
	}
	return compute.NewDiskServiceClient(conn).Get(ctx, in, opts...)
}

// List implements compute.DiskServiceClient
func (c *DiskServiceClient) List(ctx context.Context, in *compute.ListDisksRequest, opts ...grpc.CallOption) (*compute.ListDisksResponse, error) {
	conn, err := c.getConn(ctx)
	if err != nil {
		return nil, err
	}
	return compute.NewDiskServiceClient(conn).List(ctx, in, opts...)
}

// ListOperations implements compute.DiskServiceClient
func (c *DiskServiceClient) ListOperations(ctx context.Context, in *compute.ListDiskOperationsRequest, opts ...grpc.CallOption) (*compute.ListDiskOperationsResponse, error) {
	conn, err := c.getConn(ctx)
	if err != nil {
		return nil, err
	}
	return compute.NewDiskServiceClient(conn).ListOperations(ctx, in, opts...)
}

// Update implements compute.DiskServiceClient
func (c *DiskServiceClient) Update(ctx context.Context, in *compute.UpdateDiskRequest, opts ...grpc.CallOption) (*operation.Operation, error) {
	conn, err := c.getConn(ctx)
	if err != nil {
		return nil, err
	}
	return compute.NewDiskServiceClient(conn).Update(ctx, in, opts...)
}
