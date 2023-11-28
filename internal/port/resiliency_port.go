package port

import (
	"context"
	resl "github.com/zaenalarifin12/grpc-course/protogen/go/resiliency"
	"google.golang.org/grpc"
)

type ResiliencyClientPort interface {
	UnaryResiliency(ctx context.Context, in *resl.ResiliencyRequest, opts ...grpc.CallOption) (*resl.ResiliencyResponse, error)
	StreamingResiliency(ctx context.Context, in *resl.ResiliencyRequest, opts ...grpc.CallOption) (resl.ResiliencyService_StreamingResiliencyClient, error)
	ClientStreamingResiliency(ctx context.Context, opts ...grpc.CallOption) (resl.ResiliencyService_ClientStreamingResiliencyClient, error)
	BiDirectionalResiliency(ctx context.Context, opts ...grpc.CallOption) (resl.ResiliencyService_BiDirectionalResiliencyClient, error)
}

type ResiliencyWithMetadataClientPort interface {
	UnaryResiliencyWithMetadata(ctx context.Context, in *resl.ResiliencyRequest, opts ...grpc.CallOption) (*resl.ResiliencyResponse, error)
	StreamingResiliencyWithMetadata(ctx context.Context, in *resl.ResiliencyRequest, opts ...grpc.CallOption) (resl.ResiliencyWithMetadataService_StreamingResiliencyWithMetadataClient, error)
	ClientStreamingResiliencyWithMetadata(ctx context.Context, opts ...grpc.CallOption) (resl.ResiliencyWithMetadataService_ClientStreamingResiliencyWithMetadataClient, error)
	BiDirectionalResiliencyWithMetadata(ctx context.Context, opts ...grpc.CallOption) (resl.ResiliencyWithMetadataService_BiDirectionalResiliencyWithMetadataClient, error)
}
