package port

import (
	"context"
	"github.com/zaenalarifin12/grpc-course/protogen/go/bank"
	"google.golang.org/grpc"
)

type BankClientPort interface {
	GetCurrentBalance(ctx context.Context, in *bank.CurrentBalanceRequest, opts ...grpc.CallOption) (*bank.CurrentBalanceResponse, error)
	FetchExchangeRates(ctx context.Context, in *bank.ExchangeRateRequest, opts ...grpc.CallOption) (bank.BankService_FetchExchangeRatesClient, error)
	SummarizeTransaction(ctx context.Context, opts ...grpc.CallOption) (bank.BankService_SummarizeTransactionClient, error)
	TransferMultiple(ctx context.Context, opts ...grpc.CallOption) (bank.BankService_TransferMultipleClient, error)
	CreateAccount(ctx context.Context, in *bank.CreateAccountRequest, opts ...grpc.CallOption) (*bank.CreateAccountResponse, error)
}
