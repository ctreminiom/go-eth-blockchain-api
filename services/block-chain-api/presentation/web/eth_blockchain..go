package web

import (
	"context"
	"encoding/json"
	"github.com/ctreminiom/go-eth-blockchain-api/services/block-chain-api/service"
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func NewEthBlockchainHandlers(service service.BlockchainService) *EthBlockchainHandlers {
	return &EthBlockchainHandlers{service: service}
}

type EthBlockchainHandlers struct{ service service.BlockchainService }

func (e *EthBlockchainHandlers) GetTransactionHandler(w http.ResponseWriter, r *http.Request) {

	hash := chi.URLParam(r, "hash")

	if len(hash) == 0 {
		respondWithError(w, http.StatusBadRequest, "Malformed request!")
		return
	}

	transaction, err := e.service.GetTransactionByHash(context.Background(), common.HexToHash(hash))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Server Error")
		return
	}

	if transaction != nil {
		respondWithJSON(w, http.StatusOK, transaction)
		return
	}

	respondWithError(w, http.StatusNoContent, "transaction not found!")
}

func (e *EthBlockchainHandlers) GetLatestBlockHandler(w http.ResponseWriter, r *http.Request) {

	block, err := e.service.GetLatestBlock(context.Background())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, block)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJSON(w, code, map[string]string{"message": msg})
}
