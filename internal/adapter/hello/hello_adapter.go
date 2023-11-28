package hello

import (
	"context"
	"github.com/zaenalarifin12/grpc-course/protogen/go/hello"
	"github.com/zaenalarifin12/my-grpc-go-client/internal/port"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
)

type HelloAdapter struct {
	helloClient port.HelloClientPort
}

func NewHelloAdapter(conn *grpc.ClientConn) (*HelloAdapter, error) {
	client := hello.NewHelloServiceClient(conn)

	return &HelloAdapter{helloClient: client}, nil
}

func (a *HelloAdapter) SayHello(ctx context.Context, name string) (*hello.HelloResponse, error) {
	helloRequest := &hello.HelloRequest{Name: name}

	greet, err := a.helloClient.SayHello(ctx, helloRequest)

	if err != nil {
		log.Fatalf("Error on say hello :%v", err)
	}

	return greet, nil
}

func (a *HelloAdapter) SayManyHellos(ctx context.Context, name string) {
	helloRequest := &hello.HelloRequest{Name: name}

	greetStream, err := a.helloClient.SayManyHellos(ctx, helloRequest)

	if err != nil {
		log.Fatalf("Error on SayManyHellos: %v", err)
	}

	for {
		greet, err := greetStream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("Error on SayManyHellos: %v", err)
		}

		log.Println(greet.Greet)

	}
}

func (a *HelloAdapter) SayHelloToEveryone(ctx context.Context, names []string) {
	greetStream, err := a.helloClient.SayHelloToEveryone(ctx)

	if err != nil {
		log.Fatalf("Error on sayHelloToEveryone : ", err)
	}

	for _, name := range names {
		req := &hello.HelloRequest{
			Name: name,
		}

		greetStream.Send(req)
		time.Sleep(500 * time.Millisecond)
	}

	res, err := greetStream.CloseAndRecv()

	if err != nil {
		log.Fatalln("Error on SayHelloToEveryone : ", err)
	}

	log.Println(res.Greet)

}

func (a *HelloAdapter) SayHelloContinuous(ctx context.Context, names []string) {
	greetStream, err := a.helloClient.SayHelloContinuous(ctx)

	if err != nil {
		log.Fatalln("Error on SayHelloContinuous : ", err)
	}

	greetChan := make(chan struct{})

	go func() {
		for _, name := range names {
			req := &hello.HelloRequest{
				Name: name,
			}

			greetStream.Send(req)
		}

		greetStream.CloseSend()
	}()

	go func() {
		for {
			greet, err := greetStream.Recv()

			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatalln("Error on sayhellocontinuos : ,", err)
			}

			log.Println(greet.Greet)
		}

		close(greetChan)
	}()

	<-greetChan
}
