package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/goccy/go-json"

	"affiliates-backoffice-backend/internal/batch"
	"affiliates-backoffice-backend/internal/database"
	"affiliates-backoffice-backend/internal/health"
	"affiliates-backoffice-backend/internal/middleware"
	"affiliates-backoffice-backend/internal/transaction"
	"affiliates-backoffice-backend/pkg/config"
	"affiliates-backoffice-backend/pkg/context"
	"affiliates-backoffice-backend/pkg/log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

var (
	shutdowns []func() error
	logger    log.LoggerI
)

func main() {
	config.Init()
	logger = log.New()

	defer func() {
		if err := logger.Sync(); err != nil {
			logger.Error("logger sync() error: ", err)
		}
	}()

	timeoutDuration, err := time.ParseDuration(config.Vars.ApiTimeout)
	if err != nil {
		logger.Fatal("error API_TIMEOUT parseDuration")
	}

	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	app.Use(requestid.New(requestid.Config{
		Header:     fiber.HeaderXRequestID,
		ContextKey: fmt.Sprint(context.KeyRequestID),
	}))
	app.Use(recover.New())
	app.Use(middleware.Cors())
	app.Use(middleware.Timeout(timeoutDuration))

	shutdown := make(chan struct{})
	go gracefulShutdown(app, shutdown)

	database.Connect(logger)

	health.RegisterHandlers(app)
	app.Use(middleware.ValidateHeaders())

	batchRepo := batch.NewRepository(logger, database.DB)
	batchSrv := batch.NewService(logger, batchRepo)
	batch.RegisterHandlers(app, logger, batchSrv)

	tranRepo := transaction.NewRepository(logger, database.DB)
	tranSrv := transaction.NewService(logger, tranRepo)
	transaction.RegisterHandlers(app, logger, tranSrv)

	logger.Info("server starting at port " + config.Vars.ApiPort)
	if err := app.Listen(":" + config.Vars.ApiPort); err != http.ErrServerClosed {
		logger.Fatal("server error: " + err.Error())
	}

	<-shutdown
}

func gracefulShutdown(server *fiber.App, shutdown chan struct{}) {
	var (
		sigint = make(chan os.Signal, 1)
	)

	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
	<-sigint

	logger.Info("shutting down server gracefully")

	if err := server.Shutdown(); err != nil {
		logger.Fatal("shutdown error")
	}

	for i := range shutdowns {
		err := shutdowns[i]()
		logger.Error("graceful shutdown error: ", err)
	}

	close(shutdown)
}
