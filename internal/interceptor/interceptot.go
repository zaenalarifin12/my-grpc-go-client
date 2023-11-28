package interceptor

import (
	"context"
	"github.com/zaenalarifin12/grpc-course/protogen/go/hello"
	"github.com/zaenalarifin12/grpc-course/protogen/go/resiliency"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"time"
)

func LogUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		log.Println("[LOGGED BY CLIENT INTERCEPTOR]", req)

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func BasicUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		// modify request
		switch request := req.(type) {
		case *hello.HelloRequest:
			request.Name = "[MODIFIED BY CLIENT INTERCEPTOR - 1]" + request.Name
		}

		// add request metadata
		ctx = metadata.AppendToOutgoingContext(ctx, "my-request-metadata-key-1", "my-request-metadata-value-1")
		ctx = metadata.AppendToOutgoingContext(ctx, "my-request-metadata-key-2", "my-request-metadata-value-2")

		// invoke grpc method
		err := invoker(ctx, method, req, reply, cc, opts...)

		if err != nil {
			return err
		}

		// modify response
		switch response := reply.(type) {
		case *hello.HelloResponse:
			response.Greet = "[MODIFIED BY CLIENT INTERCEPTOR - 2]" + response.Greet
		case *resiliency.ResiliencyResponse:
			response.DummyString = "[MODIFIED BY CLIENT INTERCEPTOR - 3]" + response.DummyString
		}

		return nil
	}
}

func TimeoutUnaryClientInterceptor(timeout time.Duration) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		newCtx, _ := context.WithTimeout(ctx, timeout)

		return invoker(newCtx, method, req, reply, cc, opts...)
	}
}

func LogStreamClientInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		log.Println("[LOG CLIENT INTERCEPTOR]", method)

		return streamer(ctx, desc, cc, method, opts...)
	}
}

type InterpretedClientStream struct {
	grpc.ClientStream
}

func BasicClientStreamInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		ctx = metadata.AppendToOutgoingContext(ctx, "my-request-metadata-key-1", "my-request-metadata-value-1")
		ctx = metadata.AppendToOutgoingContext(ctx, "my-request-metadata-key-2", "my-request-metadata-value-2")

		clientStream, err := streamer(ctx, desc, cc, method, opts...)

		if err != nil {
			log.Printf("Failed to start %v streaming call to %v : %v\n", desc.StreamName, method, err)
			return nil, err
		}

		intercepterClientStream := &InterpretedClientStream{clientStream}

		return intercepterClientStream, nil
	}
}

func (s *InterpretedClientStream) SendMsg(msg interface{}) error {
	switch request := msg.(type) {
	case *hello.HelloRequest:
		request.Name = "[MODIFIED BY CLIENT INTERCEPTOR - 4] " + request.Name
	}

	return s.ClientStream.SendMsg(msg)
}

func (s *InterpretedClientStream) RecvMsg(msg interface{}) error {
	err := s.ClientStream.RecvMsg(msg)

	if err != nil {
		return err
	}

	switch response := msg.(type) {
	case *hello.HelloResponse:
		response.Greet = "[MODIFIED BY CLIENT INTERCEPTOR - 5]" + response.Greet
	case *resiliency.ResiliencyResponse:
		response.DummyString = "[MODIFIED BY CLIENT INTERCEPTOR - 6]" + response.DummyString
	}

	return nil
}

func TimeoutClientInterceptor(timeout time.Duration) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		newCtx, _ := context.WithTimeout(ctx, timeout)

		return streamer(newCtx, desc, cc, method, opts...)
	}
}
