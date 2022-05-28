package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/hawkkiller/wtc_system/recipe_service/internal/config"
	"github.com/hawkkiller/wtc_system/recipe_service/internal/recipe"
	_ "github.com/hawkkiller/wtc_system/recipe_service/internal/recipe"
	"github.com/hawkkiller/wtc_system/recipe_service/internal/recipe/db"
	"github.com/hawkkiller/wtc_system/recipe_service/pkg/client/postgresql"
	"github.com/hawkkiller/wtc_system/recipe_service/pkg/logging"
	"github.com/hawkkiller/wtc_system/recipe_service/pkg/shutdown"
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

	recipeStorage := db.NewStorage(postgresqlClient, logger)

	recipeService, err := recipe.NewService(recipeStorage, logger)

	if err != nil {
		logger.Fatal(err)
	}

	recipeHandler := recipe.Handler{
		Logger:        logger,
		RecipeService: recipeService,
		Validator:     validator.New(),
	}

	recipeHandler.Register(router)

	logger.Println("Start recipe service")

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
