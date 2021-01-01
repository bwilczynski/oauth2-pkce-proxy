package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	h "github.com/bwilczynski/oauth2-pkce-proxy/handlers"
	m "github.com/bwilczynski/oauth2-pkce-proxy/models"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/gorilla/mux"
)

func main() {
	pflag.Int("port", 8080, "HTTP port to run server on")
	pflag.String("client-id", "", "OAuth Client ID")
	pflag.String("client-secret", "", "OAuth Client Secret")
	pflag.String("authorization-endpoint", "https://www.strava.com/oauth/authorize", "OAuth Authorization Endpoint")
	pflag.String("token-endpoint", "https://www.strava.com/oauth/token", "OAuth Token Endpoint")
	pflag.Bool("console", false, "Pretty logging on the console")
	pflag.Parse()

	viper.BindPFlags(pflag.CommandLine)
	viper.SetEnvPrefix("PKCE_PROXY")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	var output io.Writer
	if viper.GetBool("console") {
		output = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	} else {
		output = os.Stdout
	}
	logger := zerolog.New(output).With().Timestamp().Logger()

	var provider m.OAuth2Provider
	if err := viper.Unmarshal(&provider); err != nil {
		logger.Panic().Err(err).Msg("")
	}

	mux := mux.NewRouter()
	registerRoutes(mux, &provider, &logger)

	port := viper.GetInt("port")
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
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

	logger.Info().
		Int("port", port).
		Str("client-id", provider.ClientId).
		Str("authorization-endpoint", provider.AuthorizationEndpoint).
		Str("token-endpoint", provider.AuthorizationEndpoint).
		Msg("server started")

	sigChannel := make(chan os.Signal)
	signal.Notify(sigChannel, os.Interrupt)
	signal.Notify(sigChannel, os.Kill)

	osSignal := <-sigChannel
	log.Printf("Received %v signal, performing graceful shutdown.\n", osSignal)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	server.Shutdown(ctx)
}

func registerRoutes(mux *mux.Router, provider *m.OAuth2Provider, logger *zerolog.Logger) {
	store := h.NewInMemoryChallengeStore()

	mux.Handle("/authorize", h.NewAuthorizeHandler(logger, provider, "/code")).Methods("GET")
	mux.Handle("/token", h.NewAccessTokenHandler(logger, store, provider)).Methods("POST")
	mux.Handle("/code", h.NewAuthorizeCodeHandler(logger, store)).Methods("GET")
	mux.Handle("/metrics", promhttp.Handler())

	mux.Use(h.LoggingMiddleware(logger))
	mux.Use(h.PrometheusMiddleware)
}
