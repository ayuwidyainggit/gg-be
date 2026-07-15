package server

import (
	"log"
	"sales/pkg/config"
	"sales/pkg/config/env"

	"os"
	"os/signal"

	"github.com/go-co-op/gocron/v2"
	"github.com/gofiber/fiber/v2"
)

// StartServerWithGracefulShutdown function for starting server with a graceful shutdown.
func FiberServerWithGracefulShutdown(a *fiber.App, scheduler gocron.Scheduler) {
	envCfg := env.NewCfgEnv()
	// Create channel for idle connections.
	idleConnsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt) // Catch OS signals.
		<-sigint

		// Received an interrupt signal, shutdown.
		if err := a.Shutdown(); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("Oops... Server is not shutting down! Reason: %v", err)
		}

		// shutdown scheduler
		if scheduler != nil {
			if err := scheduler.Shutdown(); err != nil {
				log.Printf("Oops... Scheduler is not shutting down! Reason: %v", err)
			}
		}

		close(idleConnsClosed)
	}()

	// Build Fiber connection URL.
	fiberConnURL, _ := config.ConnectionURL("fiber", envCfg)

	// Run server.
	if err := a.Listen(fiberConnURL); err != nil {
		log.Printf("Oops... Server is not running! Reason: %v", err)
	}

	<-idleConnsClosed
}

// StartServer func for starting a simple server.
func FiberServer(a *fiber.App, envCfg env.ConfigEnv) {
	// Build Fiber connection URL.
	fiberConnURL, _ := config.ConnectionURL("fiber", env.NewCfgEnv())

	// Run server.
	if err := a.Listen(fiberConnURL); err != nil {
		log.Printf("Oops... Server is not running! Reason: %v", err)
	}
}
