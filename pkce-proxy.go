package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/bwilczynski/oauth2-pkce-proxy/handlers"
)

var (
	listenPort int
	provider   *handlers.OAuth2Provider
)

func init() {
	if lp := os.Getenv("PKCE_PROXY_PORT"); lp != "" {
		listenPort, _ = strconv.Atoi(lp)
	}
	if listenPort == 0 {
		listenPort = 8080
	}

	authURL, err := url.Parse(os.Getenv("PKCE_PROXY_OAUTH2_AUTH_URL"))
	if err != nil {
		log.Fatal(err)
	}

	tokenURL, err := url.Parse(os.Getenv("PKCE_PROXY_OAUTH2_TOKEN_URL"))
	if err != nil {
		log.Fatal(err)
	}

	provider = handlers.NewOAuth2Provider(authURL, tokenURL)
}

func main() {
	log := log.New(os.Stdout, "", log.LstdFlags)

	mux := http.NewServeMux()
	mux.Handle("/authorize", handlers.NewAuthorize(log, provider))
	mux.Handle("/access_token", handlers.NewAccessToken(log))

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", listenPort),
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

	log.Printf("Server started at port %d.", listenPort)

	sigChannel := make(chan os.Signal)
	signal.Notify(sigChannel, os.Interrupt)
	signal.Notify(sigChannel, os.Kill)

	osSignal := <-sigChannel
	log.Printf("Received %v signal, performing graceful shutdown.\n", osSignal)

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	server.Shutdown(ctx)
}
