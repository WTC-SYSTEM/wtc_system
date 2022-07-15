package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/WTC-SYSTEM/logging"
	"github.com/WTC-SYSTEM/utils"
	"github.com/WTC-SYSTEM/wtc_system/user_service/internal/config"
	"github.com/WTC-SYSTEM/wtc_system/user_service/internal/user"
	"github.com/WTC-SYSTEM/wtc_system/user_service/internal/user/db"
	"github.com/WTC-SYSTEM/wtc_system/user_service/pkg/client/postgresql"
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

	postgresqlClient, err := postgresql.NewClient(context.Background(), cfg.Storage)
	if err != nil {
		logger.Println("failed to connect to postgresql")
		logger.Error(err)
	}

	userStorage := db.NewStorage(postgresqlClient, logger)

	userService, err := user.NewService(userStorage, logger)

	if err != nil {
		logger.Fatal(err)
	}

	usersHandler := user.Handler{
		Logger:      logger,
		UserService: userService,
		Validator:   validator.New(),
	}
	usersHandler.Register(router)

	logger.Println("Start user_service")

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
