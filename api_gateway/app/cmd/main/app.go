package main

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/hawkkiller/wtc_system/api_gateway/internal/client/user_service"
	"github.com/hawkkiller/wtc_system/api_gateway/internal/config"
	"github.com/hawkkiller/wtc_system/api_gateway/internal/handlers/auth"
	"github.com/hawkkiller/wtc_system/api_gateway/pkg/cache/freecache"
	"github.com/hawkkiller/wtc_system/api_gateway/pkg/jwt"
	"github.com/hawkkiller/wtc_system/api_gateway/pkg/logging"
	"github.com/hawkkiller/wtc_system/api_gateway/pkg/shutdown"
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

	logger.Println("cache initializing")

	refreshTokenCache := freecache.NewCacheRepo(104857600) // 100MB

	logger.Println("helpers initializing")
	jwtHelper := jwt.NewHelper(refreshTokenCache, logger)

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

	go shutdown.Graceful(
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
