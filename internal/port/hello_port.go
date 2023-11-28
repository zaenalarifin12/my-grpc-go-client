package port

import (
	"context"
	"github.com/zaenalarifin12/grpc-course/protogen/go/hello"
	"google.golang.org/grpc"
)

type HelloClientPort interface {
	SayHello(ctx context.Context, in *hello.HelloRequest, opts ...grpc.CallOption) (*hello.HelloResponse, error)
	SayManyHellos(ctx context.Context, in *hello.HelloRequest, opts ...grpc.CallOption) (hello.HelloService_SayManyHellosClient, error)
	SayHelloToEveryone(ctx context.Context, opts ...grpc.CallOption) (hello.HelloService_SayHelloToEveryoneClient, error)
	SayHelloContinuous(ctx context.Context, opts ...grpc.CallOption) (hello.HelloService_SayHelloContinuousClient, error)
}
