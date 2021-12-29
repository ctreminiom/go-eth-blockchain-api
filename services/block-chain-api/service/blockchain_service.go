package service

import (
	"context"
	internal "github.com/ctreminiom/go-eth-blockchain-api/internal/models"
	"github.com/ethereum/go-ethereum/common"
)

type BlockchainService interface {
	GetLatestBlock(ctx context.Context) (*internal.Block, error)
	GetTransactionByHash(ctx context.Context, hash common.Hash) (*internal.Transaction, error)
	GetAddressBalance(ctx context.Context, addressee string) (string, error)
	TransferEthereum(ctx context.Context, privateKey, addressee string, amount int64) (string, error)
}
