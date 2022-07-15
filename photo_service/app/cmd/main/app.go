package main

import (
	"errors"
	"fmt"
	"github.com/WTC-SYSTEM/logging"
	"github.com/WTC-SYSTEM/utils"
	"github.com/WTC-SYSTEM/wtc_system/photo_service/internal/config"
	"github.com/WTC-SYSTEM/wtc_system/photo_service/internal/photo"
	"github.com/WTC-SYSTEM/wtc_system/photo_service/pkg/client/aws"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"syscall"
	"time"
)

func main() {
	logging.Init()
	logger := logging.GetLogger()
	logger.Println("logger initialized")
	logger.Println("config initializing")
	cfg := config.GetConfig()

	logger.Println("router initializing")
	router := mux.NewRouter()

	logger.Println("s3 initializing")

	awsCfg, err := aws.NewS3(cfg.AwsCfg)

	if err != nil {
		logger.Fatal(err)
	}
	logger.Println("s3 initialized")

	storage := photo.NewStorage(logger, awsCfg)

	service := photo.NewService(storage, logger)

	handler := photo.NewHandler(
		logger,
		service,
		validator.New(),
	)
	handler.Register(router)
	logger.Println("Start photos service")

	start(router, logger, cfg)
}

func start(router http.Handler, logger logging.Logger, cfg *config.Config) {
	var server *http.Server

	server = &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Listen.Port),
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go utils.Graceful(
		[]os.Signal{
			syscall.SIGABRT,
			syscall.SIGQUIT,
			syscall.SIGHUP,
			syscall.SIGTERM,
			os.Interrupt,
		},
		server,
	)

	logger.Println("application initialized and started")

	if err := server.ListenAndServe(); err != nil {
		switch {
		case errors.Is(err, http.ErrServerClosed):
			logger.Warn("server shutdown")
		default:
			logger.Fatal(err)
		}
	}

}
