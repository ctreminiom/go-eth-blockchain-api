package application

import (
	"context"
	"crypto/ecdsa"
	internal "github.com/ctreminiom/go-eth-blockchain-api/internal/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
)

func NewEthereumBlockchain(host string) (*EthereumApplication, error) {

	client, err := ethclient.Dial(host)
	if err != nil {
		return nil, err
	}

	return &EthereumApplication{client: client}, nil
}

type EthereumApplication struct{ client *ethclient.Client }

func (e *EthereumApplication) GetLatestBlock(ctx context.Context) (*internal.Block, error) {

	header, err := e.client.HeaderByNumber(ctx, nil)
	if err != nil {
		return nil, err
	}

	blockNumber := big.NewInt(header.Number.Int64())

	chainBlock, err := e.client.BlockByNumber(ctx, blockNumber)
	if err != nil {
		return nil, err
	}

	var transactions []internal.Transaction
	for _, transaction := range chainBlock.Transactions() {

		transactions = append(transactions, internal.Transaction{
			Hash:     transaction.Hash().String(),
			Value:    transaction.Value().String(),
			Gas:      transaction.Gas(),
			GasPrice: transaction.GasPrice().Uint64(),
			Nonce:    transaction.Nonce(),
			To:       transaction.To().String(),
		})
	}

	// create the model response
	block := &internal.Block{
		BlockNumber:       chainBlock.Number().Int64(),
		Timestamp:         chainBlock.Time(),
		Difficulty:        chainBlock.Difficulty().Uint64(),
		Hash:              chainBlock.Hash().String(),
		TransactionsCount: len(chainBlock.Transactions()),
		Transactions:      transactions,
	}

	return block, nil
}

func (e *EthereumApplication) GetTransactionByHash(ctx context.Context, hash common.Hash) (*internal.Transaction, error) {

	transaction, pending, err := e.client.TransactionByHash(ctx, hash)
	if err != nil {
		return nil, err
	}

	return &internal.Transaction{
		Hash:     transaction.Hash().String(),
		Value:    transaction.Value().String(),
		Gas:      transaction.Gas(),
		GasPrice: transaction.GasPrice().Uint64(),
		To:       transaction.To().String(),
		Pending:  pending,
		Nonce:    transaction.Nonce(),
	}, nil
}

func (e *EthereumApplication) GetAddressBalance(ctx context.Context, addressee string) (string, error) {

	account := common.HexToAddress(addressee)
	balance, err := e.client.BalanceAt(ctx, account, nil)

	if err != nil {
		return "0", err
	}

	return balance.String(), nil
}

func (e *EthereumApplication) TransferEthereum(ctx context.Context, privateKey, addressee string, amount int64) (string, error) {

	privateKeyAsECDSA, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return "", err
	}

	publicKey := privateKeyAsECDSA.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", err
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// Now we can read the nonce that we should use for the account's transaction.
	nonce, err := e.client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		return "", err
	}

	value := big.NewInt(amount) // in wei (1 eth)
	gasLimit := uint64(21000)   // in units
	gasPrice, err := e.client.SuggestGasPrice(ctx)
	if err != nil {
		return "", err
	}

	// We figure out who we're sending the ETH to.
	toAddress := common.HexToAddress(addressee)
	var data []byte

	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)
	chainID, err := e.client.NetworkID(ctx)
	if err != nil {
		return "", err
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKeyAsECDSA)
	if err != nil {
		return "", err
	}

	if err = e.client.SendTransaction(ctx, signedTx); err != nil {
		return "", err
	}

	return signedTx.Hash().String(), nil
}
