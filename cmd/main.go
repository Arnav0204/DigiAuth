package main

import (
	"context"

	"digiauth/main-app/db"
	"digiauth/main-app/interfaces/auth"
	"digiauth/main-app/interfaces/issuer"
	"digiauth/main-app/interfaces/receiver"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Server struct {
	name    string
	addr    string
	handler http.Handler
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("Application failed: %v", err)
	}
}

func run() error {
	if err := db.InitDB(); err != nil {
		return err
	}
	defer db.CloseDB()

	servers := []Server{
		{"Auth", ":1010", auth.RegisterRoutes()},
		{"Issuer", ":8080", issuer.RegisterRoutes()},
		{"Receiver", ":6060", receiver.RegisterRoutes()},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	for _, s := range servers {
		wg.Add(1)
		go func(s Server) {
			defer wg.Done()
			runServerWithRestart(ctx, s)
		}(s)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("Received interrupt signal. Shutting down servers...")

	cancel()
	wg.Wait()

	log.Println("All servers have been shut down")
	return nil
}

func runServerWithRestart(ctx context.Context, s Server) {
	for {
		if err := runServer(ctx, s); err != nil {
			log.Printf("%s server failed: %v. Restarting...", s.name, err)
			select {
			case <-ctx.Done():
				return
			case <-time.After(5 * time.Second):
				// Wait for 5 seconds before restarting
			}
		} else {
			return
		}
	}
}

func runServer(ctx context.Context, s Server) error {
	server := &http.Server{
		Addr:    s.addr,
		Handler: s.handler,
	}

	serverErr := make(chan error, 1)
	go func() {
		log.Printf("Starting %s server on %s", s.name, s.addr)
		serverErr <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErr:
		return err
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return server.Shutdown(shutdownCtx)
	}
}
