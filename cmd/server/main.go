package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/inna-maikut/avito-shop/internal/api/auth"
	"github.com/inna-maikut/avito-shop/internal/api/info"
	"github.com/inna-maikut/avito-shop/internal/infrastructure/config"
	"github.com/inna-maikut/avito-shop/internal/infrastructure/jwt"
	"github.com/inna-maikut/avito-shop/internal/infrastructure/middleware"
	"github.com/inna-maikut/avito-shop/internal/infrastructure/pg"
	"github.com/inna-maikut/avito-shop/internal/repository"
	"github.com/inna-maikut/avito-shop/internal/usecases/authenticating"
	"github.com/inna-maikut/avito-shop/internal/usecases/info_collecting"
)

const (
	readHeaderTimeout = time.Second
)

func main() {
	cfg := config.Load()

	ctx := context.Background()

	db, cancel, err := pg.NewDB(ctx, cfg)
	if err != nil {
		log.Fatal("Unable to init database:", err)
	}
	defer cancel()

	_ = db

	tokenProvider, err := jwt.NewProviderFromEnv()
	if err != nil {
		panic(fmt.Errorf("create jwt provider: %w", err))
	}

	employeeRepo, err := repository.NewEmployeeRepository(db)
	if err != nil {
		panic(fmt.Errorf("create user repository: %w", err))
	}

	transactionRepo, err := repository.NewTransactionRepository(db)
	if err != nil {
		panic(fmt.Errorf("create transaction repository: %w", err))
	}

	inventoryRepo, err := repository.NewInventoryRepository(db)
	if err != nil {
		panic(fmt.Errorf("create inventory repository: %w", err))
	}

	authenticatingUseCase, err := authenticating.New(employeeRepo, tokenProvider)
	if err != nil {
		panic(fmt.Errorf("create authenticating use case: %w", err))
	}

	authHandler, err := auth.New(authenticatingUseCase)
	if err != nil {
		panic(fmt.Errorf("create auth handler: %w", err))
	}

	infoCollectingUseCase, err := info_collecting.New(employeeRepo, transactionRepo, inventoryRepo)
	if err != nil {
		panic(fmt.Errorf("create authenticating use case: %w", err))
	}

	infoHandler, err := info.New(infoCollectingUseCase)
	if err != nil {
		panic(fmt.Errorf("create info handler: %w", err))
	}

	noAuthMW, err := middleware.CreateNoAuthMiddleware()
	if err != nil {
		panic(fmt.Errorf("create no auth middleware: %w", err))
	}
	authMW, err := middleware.CreateAuthMiddleware(tokenProvider)
	if err != nil {
		panic(fmt.Errorf("create auth middleware: %w", err))
	}

	m := http.NewServeMux()

	authMux := http.NewServeMux()

	authMux.HandleFunc("GET /api/info", infoHandler.Handle)

	m.Handle("POST /api/auth", noAuthMW(http.HandlerFunc(authHandler.Handle)))
	m.Handle("/", authMW(authMux))

	s := &http.Server{
		Handler:           m,
		Addr:              "0.0.0.0:" + strconv.Itoa(cfg.ServerPort),
		ReadHeaderTimeout: readHeaderTimeout,
	}

	// And we serve HTTP until the world ends.
	err = s.ListenAndServe()
	if err != nil && !errors.Is(err, context.Canceled) {
		log.Default().Println("http server ListenAndServe:", err)
	}
}
