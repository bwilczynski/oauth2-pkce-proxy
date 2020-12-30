package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	h "github.com/bwilczynski/oauth2-pkce-proxy/handlers"
	m "github.com/bwilczynski/oauth2-pkce-proxy/models"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"

	"github.com/gorilla/mux"
)

func main() {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	logger := zerolog.New(output).With().Timestamp().Logger()

	cfg := &m.Config{}
	if err := cfg.ReadFromEnv(); err != nil {
		log.Fatal(err)
	}

	mux := mux.NewRouter()
	registerRoutes(mux, cfg, &logger)

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

func registerRoutes(mux *mux.Router, cfg *m.Config, logger *zerolog.Logger) {
	store := h.NewInMemoryChallengeStore()

	mux.Handle("/authorize", h.NewAuthorizeHandler(logger, cfg.Provider, "/code")).Methods("GET")
	mux.Handle("/token", h.NewAccessTokenHandler(logger, store, cfg.Provider)).Methods("POST")
	mux.Handle("/code", h.NewAuthorizeCodeHandler(logger, store)).Methods("GET")
	mux.Handle("/metrics", promhttp.Handler())

	mux.Use(h.LoggingMiddleware(logger))
	mux.Use(h.PrometheusMiddleware)
}
