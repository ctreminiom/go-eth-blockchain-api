package main

import (
	"context"
	"github.com/ctreminiom/go-eth-blockchain-api/services/block-chain-api/application"
	"github.com/ctreminiom/go-eth-blockchain-api/services/block-chain-api/presentation/web"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)

	ethClient, err := application.NewEthereumBlockchain("http://localhost:7545")
	if err != nil {
		return
	}

	handler := web.NewEthBlockchainHandlers(ethClient)

	router.Get("/api/v1/eth/transaction", handler.GetTransactionHandler)
	router.Get("/api/v1/eth/block/latest", handler.GetLatestBlockHandler)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Print("Server Started")

	<-done
	log.Print("Server Stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err = srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
}
