package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/WTC-SYSTEM/wtc_system/libs/logging"
	"github.com/WTC-SYSTEM/wtc_system/libs/utils"
	"github.com/WTC-SYSTEM/wtc_system/recipe_service/internal/config"
	"github.com/WTC-SYSTEM/wtc_system/recipe_service/internal/recipe"
	_ "github.com/WTC-SYSTEM/wtc_system/recipe_service/internal/recipe"
	"github.com/WTC-SYSTEM/wtc_system/recipe_service/internal/recipe/db"
	"github.com/WTC-SYSTEM/wtc_system/recipe_service/pkg/client/aws"
	"github.com/WTC-SYSTEM/wtc_system/recipe_service/pkg/client/postgresql"
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
	logger.Println("s3 initializing")

	awsCfg, err := aws.NewS3(cfg.AwsCfg)

	if err != nil {
		logger.Fatal(err)
	}
	logger.Println("s3 initialized")

	recipeStorage := db.NewStorage(postgresqlClient, logger, awsCfg)

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
