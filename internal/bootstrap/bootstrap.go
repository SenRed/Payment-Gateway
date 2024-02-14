package bootstrap

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/payment-gateway/internal/config"
	"github.com/payment-gateway/internal/domain/api"
	"github.com/payment-gateway/internal/infrastructure/mock"
	"github.com/payment-gateway/internal/infrastructure/postgres"
	"github.com/payment-gateway/internal/infrastructure/processor"
	"github.com/payment-gateway/internal/ui"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"
)

// Bootstrap struct
type Bootstrap struct {
	Config config.AppConfig
	Router *gin.Engine
}

// Init function, bootstrap all the application configuration
func Init() Bootstrap {
	// Load app config
	c := config.NewAppConfig()
	zerolog.SetGlobalLevel(c.LoggerConfig.Level)
	if !c.LoggerConfig.JSON {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	log.Info().Interface("config", c).Msg("Configuration")

	// Load PostgresSQL as the default database provider
	paymentRepository, err := postgres.NewPostgresClient(c.PostgresConfig)
	if err != nil {
		log.Fatal().Err(err).Msg("Error when creating Postgres client")
	}

	// Mock HTTP client, used to simulate acquiring bank processor
	acquiringBankHTTPClientMock := mock.NewAcquiringBankHTTPClientMock()
	paymentProcessor := processor.NewBankProcessor(acquiringBankHTTPClientMock)

	// Load paymentService
	paymentService := api.NewPaymentService(paymentRepository, paymentProcessor)

	// Create the router
	r := ui.CreateRouter(paymentService)
	return Bootstrap{
		Config: c,
		Router: r,
	}
}

func (b Bootstrap) Run() {
	if b.Router != nil {
		hostURL := fmt.Sprintf("%s%s", ":", strconv.Itoa(b.Config.HTTPPort))
		server := &http.Server{
			Addr:    hostURL,
			Handler: b.Router,
		}

		go func() {
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatal().Err(err).Str("reason", "server listen and serve error").Msg("The server stopped listening")
			}
			log.Info().Msg("Server is started")
		}()

		// Graceful shutdown
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)
		<-quit

		log.Info().Msg("Shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Fatal().Err(err).Msg("Server shutdown failed")
		}

		log.Info().Msg("Server stopped gracefully")
	}
}
