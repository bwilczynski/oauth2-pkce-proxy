package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/bwilczynski/oauth2-pkce-proxy/handlers"
)

var (
	listenPort int
)

func init() {
	if lp := os.Getenv("PKCE_PROXY_PORT"); lp != "" {
		listenPort, _ = strconv.Atoi(lp)
	}
	if listenPort == 0 {
		listenPort = 8080
	}
}

func main() {
	log := log.New(os.Stdout, "", log.LstdFlags)

	mux := http.NewServeMux()
	handlers.RegisterRoutes(mux, log)

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

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	server.Shutdown(ctx)
}
