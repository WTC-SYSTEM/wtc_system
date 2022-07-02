package main

import (
	"errors"
	"fmt"
	redis "github.com/WTC-SYSTEM/wtc_system/api_gateway/internal/client/db"
	"github.com/WTC-SYSTEM/wtc_system/api_gateway/internal/client/user_service"
	"github.com/WTC-SYSTEM/wtc_system/api_gateway/internal/config"
	"github.com/WTC-SYSTEM/wtc_system/api_gateway/internal/handlers/auth"
	"github.com/WTC-SYSTEM/wtc_system/api_gateway/pkg/jwt"
	"github.com/WTC-SYSTEM/wtc_system/libs/logging"
	"github.com/WTC-SYSTEM/wtc_system/libs/utils"
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

	logger.Println("redis initializing")
	rc := redis.NewClient(cfg)

	logger.Println("helpers initializing")

	jwtHelper := jwt.NewHelper(rc, logger)

	logger.Println("create and register handlers")

	//metricHandler := metric.Handler{Logger: logger}
	//metricHandler.Register(router)

	userService := user_service.NewService(cfg.UserService.URL, "/v1/users", logger)
	authHandler := auth.Handler{JWTHelper: jwtHelper, UserService: userService, Logger: logger}
	authHandler.Register(router)

	logger.Println("start application")
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
