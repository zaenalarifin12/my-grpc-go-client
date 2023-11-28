package resiliency

import (
	"context"
	resl "github.com/zaenalarifin12/grpc-course/protogen/go/resiliency"
	"github.com/zaenalarifin12/my-grpc-go-client/internal/port"
	"google.golang.org/grpc"
	"io"
	"log"
)

type ResiliencyAdapter struct {
	resiliencyClient             port.ResiliencyClientPort
	resiliencyWithMetadataClient port.ResiliencyWithMetadataClientPort
}

func NewResiliencyAdapter(conn *grpc.ClientConn) (*ResiliencyAdapter, error) {
	client := resl.NewResiliencyServiceClient(conn)
	clientWithMetadata := resl.NewResiliencyWithMetadataServiceClient(conn)

	return &ResiliencyAdapter{
		resiliencyClient:             client,
		resiliencyWithMetadataClient: clientWithMetadata,
	}, nil
}

func (a *ResiliencyAdapter) UnaryResiliency(ctx context.Context, minDelaySecond int32, maxDelaySecond int32, statusCode []uint32) (*resl.ResiliencyResponse, error) {
	resiliencyRequest := &resl.ResiliencyRequest{
		MinDelaySecond: minDelaySecond,
		MaxDelaySecond: maxDelaySecond,
		StatusCode:     statusCode,
	}

	res, err := a.resiliencyClient.UnaryResiliency(ctx, resiliencyRequest)

	if err != nil {
		log.Println("Error on UnaryResiliency :", err)
		return nil, err
	}

	return res, nil
}

func (a *ResiliencyAdapter) StreamingResiliency(ctx context.Context, minDelaySecond int32, maxDelaySecond int32, statusCode []uint32) {
	resiliencyRequest := &resl.ResiliencyRequest{
		MinDelaySecond: minDelaySecond,
		MaxDelaySecond: maxDelaySecond,
		StatusCode:     statusCode,
	}

	reslStream, err := a.resiliencyClient.StreamingResiliency(ctx, resiliencyRequest)

	if err != nil {
		log.Fatalln("Error on StreamingResiliency : ", err)
	}

	for {
		res, err := reslStream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalln("Error on StreamingResiliency : ", err)
		}

		log.Println(res.DummyString)
	}
}

func (a *ResiliencyAdapter) ClientStreamingResiliency(ctx context.Context, minDelaySecond int32, maxDelaySecond int32, statusCode []uint32, count int) {
	reslStream, err := a.resiliencyClient.ClientStreamingResiliency(ctx)

	if err != nil {
		log.Fatalln("Error on ClientStreamingResiliency : ", err)
	}

	for i := 0; i < count; i++ {
		resiliencyRequest := &resl.ResiliencyRequest{
			MinDelaySecond: minDelaySecond,
			MaxDelaySecond: maxDelaySecond,
			StatusCode:     statusCode,
		}

		reslStream.Send(resiliencyRequest)
	}

	res, err := reslStream.CloseAndRecv()

	if err != nil {
		log.Fatalln("Error on ClientStreamingResiliency : ", err)
	}

	log.Println(res.DummyString)
}

func (a *ResiliencyAdapter) BiDirectionalResiliency(ctx context.Context, minDelaySecond int32, maxDelaySecond int32, statusCode []uint32, count int) {

	reslStream, err := a.resiliencyClient.BiDirectionalResiliency(ctx)

	if err != nil {
		log.Fatalln("Error on BiDirectionalResiliency : ", err)
	}

	reslChan := make(chan struct{})

	go func() {

		for i := 0; i < count; i++ {
			resiliencyRequest := &resl.ResiliencyRequest{
				MinDelaySecond: minDelaySecond,
				MaxDelaySecond: maxDelaySecond,
				StatusCode:     statusCode,
			}

			reslStream.Send(resiliencyRequest)
		}

		reslStream.CloseSend()
	}()

	go func() {
		for {
			res, err := reslStream.Recv()

			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatalln("Error on BiDirectionalResiliency : ", err)
			}

			log.Println(res.DummyString)

		}

		close(reslChan)
	}()

	<-reslChan
}
