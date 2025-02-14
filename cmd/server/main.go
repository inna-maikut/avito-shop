package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/inna-maikut/avito-shop/internal/infrastructure/config"
	"github.com/inna-maikut/avito-shop/internal/infrastructure/pg"
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

	m := http.NewServeMux()

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
