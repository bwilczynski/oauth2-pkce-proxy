package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	log := log.New(os.Stdout, "", log.LstdFlags)
	cfg := &Config{}
	cfg.ReadFromEnv(log)

	mux := http.NewServeMux()
	registerRoutes(mux, cfg, log)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.ListenPort),
		Handler:      mux,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	log.Printf("Server started at port %d.", cfg.ListenPort)

	sigChannel := make(chan os.Signal)
	signal.Notify(sigChannel, os.Interrupt)
	signal.Notify(sigChannel, os.Kill)

	osSignal := <-sigChannel
	log.Printf("Received %v signal, performing graceful shutdown.\n", osSignal)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	server.Shutdown(ctx)
}

func registerRoutes(mux *http.ServeMux, cfg *Config, log *log.Logger) {
	mux.Handle("/authorize", PrometheusMiddleware(NewAuthorizeHandler(log, cfg.Provider, "/code")))
	mux.Handle("/access_token", PrometheusMiddleware(NewAccessTokenHandler(log)))
	mux.Handle("/code", PrometheusMiddleware(NewAuthorizeCodeHandler(log)))
	mux.Handle("/metrics", promhttp.Handler())
}
