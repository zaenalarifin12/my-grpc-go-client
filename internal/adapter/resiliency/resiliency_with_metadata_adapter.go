package resiliency

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	resl "github.com/zaenalarifin12/grpc-course/protogen/go/resiliency"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"io"
	"log"
	"runtime"
	"time"
)

func sampleRequestMetadata() metadata.MD {
	md := map[string]string{
		"grpc-client-time":    fmt.Sprintf(time.Now().Format("15:04:05")),
		"grpc-client-request": runtime.GOOS,
		"grpc-request-uuid":   uuid.New().String(),
	}

	return metadata.New(md)
}

func sampleResponseMetadata(md metadata.MD) {
	if md.Len() == 0 {
		log.Println("Response metadata not found")
	} else {
		log.Println("Response metadata : ")
		for k, v := range md {
			log.Printf(" %v : %v\n", k, v)
		}
	}
}

func (a *ResiliencyAdapter) UnaryResiliencyWithMetadata(ctx context.Context, minDelaySecond int32, maxDelaySecond int32, statusCode []uint32) (*resl.ResiliencyResponse, error) {
	ctx = metadata.NewOutgoingContext(ctx, sampleRequestMetadata())

	resiliencyRequest := &resl.ResiliencyRequest{
		MinDelaySecond: minDelaySecond,
		MaxDelaySecond: maxDelaySecond,
		StatusCode:     statusCode,
	}

	var responseMetadata metadata.MD
	res, err := a.resiliencyWithMetadataClient.UnaryResiliencyWithMetadata(ctx, resiliencyRequest, grpc.Header(&responseMetadata))

	if err != nil {
		log.Println("Error on UnaryResiliencyWithMetadata :", err)
		return nil, err
	}

	sampleResponseMetadata(responseMetadata)

	return res, nil
}

func (a *ResiliencyAdapter) StreamingResiliencyWithMetadata(ctx context.Context, minDelaySecond int32, maxDelaySecond int32, statusCode []uint32) {
	ctx = metadata.NewOutgoingContext(ctx, sampleRequestMetadata())

	resiliencyRequest := &resl.ResiliencyRequest{
		MinDelaySecond: minDelaySecond,
		MaxDelaySecond: maxDelaySecond,
		StatusCode:     statusCode,
	}

	reslStream, err := a.resiliencyWithMetadataClient.StreamingResiliencyWithMetadata(ctx, resiliencyRequest)

	if err != nil {
		log.Fatalln("Error on StreamingResiliencyWithMetadata : ", err)
	}

	if responseMetadata, err := reslStream.Header(); err != nil {
		sampleResponseMetadata(responseMetadata)
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

func (a *ResiliencyAdapter) ClientStreamingResiliencyWithMetadata(ctx context.Context, minDelaySecond int32, maxDelaySecond int32, statusCode []uint32, count int) {
	reslStream, err := a.resiliencyWithMetadataClient.ClientStreamingResiliencyWithMetadata(ctx)

	if err != nil {
		log.Fatalln("Error on ClientStreamingResiliencyWithMetadata : ", err)
	}

	for i := 0; i < count; i++ {
		ctx = metadata.NewOutgoingContext(ctx, sampleRequestMetadata())
		resiliencyRequest := &resl.ResiliencyRequest{
			MinDelaySecond: minDelaySecond,
			MaxDelaySecond: maxDelaySecond,
			StatusCode:     statusCode,
		}

		reslStream.Send(resiliencyRequest)
	}

	res, err := reslStream.CloseAndRecv()

	if err != nil {
		log.Fatalln("Error on ClientStreamingResiliencyWithMetadata : ", err)
	}

	if responseMetadata, err := reslStream.Header(); err == nil {
		sampleResponseMetadata(responseMetadata)
	}

	log.Println(res.DummyString)
}

func (a *ResiliencyAdapter) BiDirectionalResiliencyWithMetadata(ctx context.Context, minDelaySecond int32, maxDelaySecond int32, statusCode []uint32, count int) {

	reslStream, err := a.resiliencyWithMetadataClient.BiDirectionalResiliencyWithMetadata(ctx)

	if err != nil {
		log.Fatalln("Error on BiDirectionalResiliencyWithMetadata : ", err)
	}

	if responseMetadata, err := reslStream.Header(); err == nil {
		sampleResponseMetadata(responseMetadata)
	}

	reslChan := make(chan struct{})

	go func() {

		for i := 0; i < count; i++ {
			ctx = metadata.NewOutgoingContext(ctx, sampleRequestMetadata())

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
				log.Fatalln("Error on BiDirectionalResiliencyWithMetadata : ", err)
			}

			log.Println(res.DummyString)

		}

		close(reslChan)
	}()

	<-reslChan
}
