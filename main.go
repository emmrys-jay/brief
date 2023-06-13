package main

import (
	pgdb "brief/pkg/repository/storage/postgres"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"brief/utility"

	"brief/internal/config"
	"brief/pkg/router"

	"github.com/go-playground/validator/v10"
	// "brief/pkg/repository/storage/redis"
	// mdb "brief/pkg/repository/storage/mongo"
)

func init() {
	config.Setup()
	pgdb.ConnectToDB()
	// redis.SetupRedis() uncomment when you need redis
	// mdb.ConnectToDB() uncomment or remove if using mongodb
}

func main() {
	//Load config
	logger := utility.NewLogger()
	getConfig := config.GetConfig()
	validatorRef := validator.New()
	e := router.Setup(validatorRef, logger)

	// The HTTP Server
	server := &http.Server{
		Addr:    fmt.Sprintf("127.0.0.1:%s", getConfig.ServerPort),
		Handler: e,
	}

	// Server run context
	serverCtx, serverCancel := context.WithCancel(context.Background())

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, shutdownCancel := context.WithTimeout(serverCtx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		shutdownCancel()
		serverCancel()
	}()

	// Run the server
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
}