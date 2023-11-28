package main

import (
	"context"
	"fmt"
	"github.com/sony/gobreaker"
	resl_proto "github.com/zaenalarifin12/grpc-course/protogen/go/resiliency"
	"github.com/zaenalarifin12/my-grpc-go-client/internal/adapter/bank"
	"github.com/zaenalarifin12/my-grpc-go-client/internal/adapter/hello"
	"github.com/zaenalarifin12/my-grpc-go-client/internal/adapter/resiliency"
	dbank "github.com/zaenalarifin12/my-grpc-go-client/internal/application/domain/bank"
	"github.com/zaenalarifin12/my-grpc-go-client/internal/interceptor"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"math/rand"
	"time"
)

var cbreaker *gobreaker.CircuitBreaker

func init() {
	mybreaker := gobreaker.Settings{
		Name: "course-circuit-breaker",
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)

			log.Printf("Circuit breaker failure is %v, request is %v, means failure ratio %v\n", counts.TotalFailures, counts.Requests, failureRatio)

			return counts.Requests >= 3 && failureRatio >= 0.6
		},
		MaxRequests: 3,
		Timeout:     4 * time.Second,
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			log.Printf("Circuit breaker %v changed state, from %v to %v\n\n", name, from, to)
		},
		IsSuccessful: nil,
	}

	cbreaker = gobreaker.NewCircuitBreaker(mybreaker)
}

func main() {
	log.SetFlags(0)
	log.SetOutput(log.Writer())

	var opts []grpc.DialOption

	cred, err := credentials.NewClientTLSFromFile("ssl/ca.crt", "")

	if err != nil {
		log.Fatalln("Can't create client credentials : ", err)
	}

	opts = append(opts, grpc.WithTransportCredentials(cred))
	//opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	//opts = append(opts,
	//	grpc.WithUnaryInterceptor(
	//		grpc_retry.UnaryClientInterceptor(
	//			grpc_retry.WithCodes(codes.Unknown, codes.Internal),
	//			grpc_retry.WithMax(4),
	//			grpc_retry.WithBackoff(grpc_retry.BackoffExponential(2*time.Second)),
	//		),
	//	),
	//)
	//
	//opts = append(opts,
	//	grpc.WithStreamInterceptor(
	//		grpc_retry.StreamClientInterceptor(
	//			grpc_retry.WithCodes(codes.Unknown, codes.Internal),
	//			grpc_retry.WithMax(4),
	//			grpc_retry.WithBackoff(grpc_retry.BackoffLinear(3*time.Second)),
	//		),
	//	),
	//)
	opts = append(opts,
		grpc.WithChainUnaryInterceptor(
			interceptor.LogUnaryClientInterceptor(),
			interceptor.BasicUnaryClientInterceptor(),
			interceptor.TimeoutUnaryClientInterceptor(5*time.Second),
		),
	)

	opts = append(opts,
		grpc.WithChainStreamInterceptor(
			interceptor.LogStreamClientInterceptor(),
			interceptor.BasicClientStreamInterceptor(),
			interceptor.TimeoutClientInterceptor(15*time.Second),
		),
	)

	conn, err := grpc.Dial("localhost:9000", opts...)

	if err != nil {
		log.Fatalln("can't connect to gRPC Server :", err)
	}

	defer conn.Close()
	//
	helloAdapter, err := hello.NewHelloAdapter(conn)

	if err != nil {
		log.Fatalln("can't create hello adapter :", err)
	}

	//bankAdapter, err := bank.NewBankAdapter(conn)
	//
	//if err != nil {
	//	log.Fatalln("can't create bank adapter :", err)
	//}

	//resiliencyAdapter, err := resiliency.NewResiliencyAdapter(conn)
	//
	//if err != nil {
	//	log.Fatalln("can't create resiliencyAdapter :", err)
	//}

	runSayHello(helloAdapter, "zainal")
	//runSayManyHellos(helloAdapter, "arifin")
	//runSayHelloToEveryone(helloAdapter, []string{
	//	"tamvam1", "tamvam2", "tamvam3", "tamvam4", "tamvam5",
	//	"tamvam6", "tamvam7", "tamvam8", "tamvam9", "tamvam10",
	//})
	//runSayHelloContinuous(helloAdapter, []string{
	//	"tamvam1", "tamvam2", "tamvam3", "tamvam4", "tamvam5",
	//	"tamvam6", "tamvam7", "tamvam8", "tamvam9", "tamvam10",
	//})

	// run bank services
	//runGetCurrentBalance(bankAdapter, "7835697001")
	//runFetchExchangeRates(bankAdapter, "USD", "IDR")
	//runSummarizeTransaction(bankAdapter, "7835697001", 9)
	//runTransferMultiple(bankAdapter, "7835697001", "7835697002", 9)

	//	run resiliency services
	//runUnaryResiliencyWithTimeout(resiliencyAdapter, 2, 8, []uint32{dresl.OK}, 5*time.Second)
	//runServerStreamingResiliencyWithTimeout(resiliencyAdapter, 0, 3, []uint32{dresl.OK}, 15*time.Second)
	//runClientStreamingResiliencyWithTimeout(resiliencyAdapter, 0, 3, []uint32{dresl.OK}, 10, 15*time.Second)
	//runBiDirectionalResiliencyWithTimeout(resiliencyAdapter, 0, 3, []uint32{dresl.OK}, 10, 60*time.Second)

	//runUnaryResiliency(resiliencyAdapter, 0, 3, []uint32{dresl.UNKNOWN, dresl.OK})
	//runServerStreamingResiliency(resiliencyAdapter, 0, 3, []uint32{dresl.UNKNOWN, dresl.OK})
	//runServerStreamingResiliency(resiliencyAdapter, 0, 3, []uint32{dresl.UNKNOWN, dresl.OK})

	//

	//for i := 0; i < 300; i++ {
	//	runUnaryResiliencyWithCircuitBreaker(resiliencyAdapter, 0, 0, []uint32{dresl.UNKNOWN, dresl.OK})
	//	time.Sleep(time.Second)
	//}

	//runUnaryResiliencyWithMetadata(resiliencyAdapter, 6, 8, []uint32{dresl.OK})
	//runServerStreamingResiliencyWithMetadata(resiliencyAdapter, 1, 3, []uint32{dresl.OK})
	//runClientStreamingResiliencyWithMetadata(resiliencyAdapter, 0, 1, []uint32{dresl.OK}, 10)
	//runBiDirectionalResiliencyWithMetadata(resiliencyAdapter, 0, 1, []uint32{dresl.OK}, 10)

}

func runSayHello(adapter *hello.HelloAdapter, name string) {
	greet, err := adapter.SayHello(context.Background(), name)

	if err != nil {
		log.Fatalln("Can not call SayHello :", err)
	}

	log.Println(greet.Greet)
}

func runSayManyHellos(adapter *hello.HelloAdapter, name string) {
	adapter.SayManyHellos(context.Background(), name)
}

func runSayHelloToEveryone(adapter *hello.HelloAdapter, names []string) {
	adapter.SayHelloToEveryone(context.Background(), names)
}

func runSayHelloContinuous(adapter *hello.HelloAdapter, names []string) {
	adapter.SayHelloContinuous(context.Background(), names)
}

func runGetCurrentBalance(adapter *bank.BankAdapter, acct string) {
	bal, err := adapter.GetCurrentBalance(context.Background(), acct)

	if err != nil {
		log.Fatalln("Failed to call GetCurrentBalance")
	}

	log.Println(bal)
}

func runFetchExchangeRates(adapter *bank.BankAdapter, fromCur string, toCur string) {
	adapter.FetchExchangeRates(context.Background(), fromCur, toCur)
}

func runSummarizeTransaction(adapter *bank.BankAdapter, acct string, numDummyTransaction int) {
	var tx []dbank.Transaction

	for i := 0; i < numDummyTransaction; i++ {
		ttype := dbank.TransactionTypeIn

		if i%3 == 0 {
			ttype = dbank.TransactionTypeOut
		}

		t := dbank.Transaction{
			Amount:          float64(rand.Intn(500) + 10),
			TransactionType: ttype,
			Notes:           fmt.Sprintf("Dummy transaction %v", i),
		}

		tx = append(tx, t)
	}

	adapter.SummarizeTransactions(context.Background(), acct, tx)
}

func runTransferMultiple(adapter *bank.BankAdapter, fromAcct string, toAcct string, numDummyTransactions int) {
	var trf []dbank.TransferTransaction

	for i := 1; i <= numDummyTransactions; i++ {
		tr := dbank.TransferTransaction{
			FromAccountNumber: fromAcct,
			ToAccountNumber:   toAcct,
			Currency:          "USD",
			Amount:            float64(rand.Intn(200) + 5),
		}

		trf = append(trf, tr)
	}

	adapter.TransferMultiple(context.Background(), trf)
}

func runUnaryResiliencyWithTimeout(adapter *resiliency.ResiliencyAdapter, minDelaySecond int32, maxDelaySecond int32, statusCode []uint32, timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	defer cancel()

	res, err := adapter.UnaryResiliency(ctx, minDelaySecond, maxDelaySecond, statusCode)

	if err != nil {
		log.Fatalln("Failed to call runUnaryResiliency : ", err)
	}

	log.Println(res.DummyString)
}

func runServerStreamingResiliencyWithTimeout(adapter *resiliency.ResiliencyAdapter,
	minDelaySecond int32, maxDelaySecond int32, statusCodes []uint32, timeout time.Duration) {
	ctx, _ := context.WithTimeout(context.Background(), timeout)

	adapter.StreamingResiliency(ctx, minDelaySecond, maxDelaySecond, statusCodes)
}

func runClientStreamingResiliencyWithTimeout(adapter *resiliency.ResiliencyAdapter,
	minDelaySecond int32, maxDelaySecond int32, statusCodes []uint32,
	count int, timeout time.Duration) {
	ctx, _ := context.WithTimeout(context.Background(), timeout)

	adapter.ClientStreamingResiliency(ctx, minDelaySecond, maxDelaySecond, statusCodes, count)
}

func runBiDirectionalResiliencyWithTimeout(adapter *resiliency.ResiliencyAdapter,
	minDelaySecond int32, maxDelaySecond int32, statusCodes []uint32,
	count int, timeout time.Duration) {
	ctx, _ := context.WithTimeout(context.Background(), timeout)

	adapter.BiDirectionalResiliency(ctx, minDelaySecond, maxDelaySecond, statusCodes, count)
}

func runUnaryResiliency(adapter *resiliency.ResiliencyAdapter, minDelaySecond int32,
	maxDelaySecond int32, statusCodes []uint32) {
	res, err := adapter.UnaryResiliency(context.Background(), minDelaySecond, maxDelaySecond, statusCodes)

	if err != nil {
		log.Fatalln("Failed to call UnaryResiliency :", err)
	}

	log.Println(res.DummyString)
}

func runServerStreamingResiliency(adapter *resiliency.ResiliencyAdapter,
	minDelaySecond int32, maxDelaySecond int32, statusCodes []uint32) {
	adapter.StreamingResiliency(context.Background(), minDelaySecond, maxDelaySecond, statusCodes)
}

func runClientStreamingResiliency(adapter *resiliency.ResiliencyAdapter,
	minDelaySecond int32, maxDelaySecond int32, statusCodes []uint32,
	count int) {
	adapter.ClientStreamingResiliency(context.Background(), minDelaySecond,
		maxDelaySecond, statusCodes, count)
}

func runBiDirectionalResiliency(adapter *resiliency.ResiliencyAdapter,
	minDelaySecond int32, maxDelaySecond int32, statusCodes []uint32,
	count int) {
	adapter.BiDirectionalResiliency(context.Background(), minDelaySecond,
		maxDelaySecond, statusCodes, count)
}

func runUnaryResiliencyWithCircuitBreaker(adapter *resiliency.ResiliencyAdapter, minDelaySecond int32,
	maxDelaySecond int32, statusCodes []uint32) {
	cbreakerRes, cbreakerErr := cbreaker.Execute(func() (interface{}, error) {
		return adapter.UnaryResiliency(context.Background(), minDelaySecond, maxDelaySecond, statusCodes)
	})

	if cbreakerErr != nil {
		log.Println("Failed to call runUnaryResiliencyWithCircuitBreaker: ", cbreakerErr)
	} else {
		log.Println(cbreakerRes.(*resl_proto.ResiliencyResponse).DummyString)
	}
}

// ---
func runUnaryResiliencyWithMetadata(adapter *resiliency.ResiliencyAdapter, minDelaySecond int32,
	maxDelaySecond int32, statusCodes []uint32) {
	res, err := adapter.UnaryResiliencyWithMetadata(context.Background(), minDelaySecond, maxDelaySecond, statusCodes)

	if err != nil {
		log.Fatalln("Failed to call UnaryResiliency :", err)
	}

	log.Println(res.DummyString)
}

func runServerStreamingResiliencyWithMetadata(adapter *resiliency.ResiliencyAdapter,
	minDelaySecond int32, maxDelaySecond int32, statusCodes []uint32) {
	adapter.StreamingResiliencyWithMetadata(context.Background(), minDelaySecond, maxDelaySecond, statusCodes)
}

func runClientStreamingResiliencyWithMetadata(adapter *resiliency.ResiliencyAdapter,
	minDelaySecond int32, maxDelaySecond int32, statusCodes []uint32,
	count int) {
	adapter.ClientStreamingResiliencyWithMetadata(context.Background(), minDelaySecond,
		maxDelaySecond, statusCodes, count)
}

func runBiDirectionalResiliencyWithMetadata(adapter *resiliency.ResiliencyAdapter,
	minDelaySecond int32, maxDelaySecond int32, statusCodes []uint32,
	count int) {
	adapter.BiDirectionalResiliencyWithMetadata(context.Background(), minDelaySecond,
		maxDelaySecond, statusCodes, count)
}
