package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/inna-maikut/avito-shop/internal/infrastructure/config"
)

const (
	readHeaderTimeout = time.Second
)

func main() {
	cfg := config.Load()

	m := http.NewServeMux()

	s := &http.Server{
		Handler:           m,
		Addr:              "0.0.0.0:" + strconv.Itoa(cfg.ServerPort),
		ReadHeaderTimeout: readHeaderTimeout,
	}

	// And we serve HTTP until the world ends.
	log.Fatal(s.ListenAndServe())
}
