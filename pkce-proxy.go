package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/bwilczynski/oauth2-pkce-proxy/handlers"
)

const (
	ListenPort = 8080
)

func main() {
	log := log.New(os.Stdout, "", log.LstdFlags)

	mux := http.NewServeMux()
	mux.Handle("/authorize", handlers.NewAuthorize(log))
	mux.Handle("/access_token", handlers.NewAccessToken(log))

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", ListenPort),
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

	log.Printf("Server started at port %d.", ListenPort)

	sigChannel := make(chan os.Signal)
	signal.Notify(sigChannel, os.Interrupt)
	signal.Notify(sigChannel, os.Kill)

	osSignal := <-sigChannel
	log.Printf("Received %v signal, performing graceful shutdown.\n", osSignal)

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	server.Shutdown(ctx)
}
