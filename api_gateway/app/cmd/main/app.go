package main

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/hawkkiller/wtc_system/api_gateway/app/internal/config"
	"github.com/hawkkiller/wtc_system/api_gateway/app/pkg/logging"
	"github.com/hawkkiller/wtc_system/api_gateway/app/pkg/shutdown"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
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
	_, err := http.Get("http://wtc-user-service/api/v1/users")
	if err != nil {
		logger.Error(err)
	}

	//refreshTokenCache := freecache.NewCacheRepo(104857600) // 100MB
	//
	//logger.Println("helpers initializing")
	//jwtHelper := jwt.NewHelper(refreshTokenCache, logger)
	//
	//logger.Println("create and register handlers")
	//
	//metricHandler := metric.Handler{Logger: logger}
	//metricHandler.Register(router)
	//
	//userService := user_service.NewService(cfg.UserService.URL, "/users", logger)
	//authHandler := auth.Handler{JWTHelper: jwtHelper, UserService: userService, Logger: logger}
	//authHandler.Register(router)
	//
	//categoryService := category_service.NewService(cfg.CategoryService.URL, "/categories", logger)
	//categoriesHandler := categories.Handler{CategoryService: categoryService, Logger: logger}
	//categoriesHandler.Register(router)
	//
	//noteService := note_service.NewService(cfg.NoteService.URL, "/notes", logger)
	//notesHandler := notes.Handler{NoteService: noteService, Logger: logger}
	//notesHandler.Register(router)
	//
	//tagService := tag_service.NewService(cfg.TagService.URL, "/tags", logger)
	//tagsHandler := tags.Handler{TagService: tagService, Logger: logger}
	//tagsHandler.Register(router)

	logger.Println("start application")
	start(router, logger, cfg)
}

func start(router *mux.Router, logger logging.Logger, cfg *config.Config) {
	var server *http.Server
	var listener net.Listener

	if cfg.Listen.Type == "sock" {
		appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			logger.Fatal(err)
		}
		socketPath := path.Join(appDir, "app.sock")
		logger.Infof("socket path: %s", socketPath)

		logger.Info("create and listen unix socket")
		listener, err = net.Listen("unix", socketPath)
		if err != nil {
			logger.Fatal(err)
		}
	} else {
		logger.Infof("bind application to host: %s and port: %s", cfg.Listen.BindIP, cfg.Listen.Port)

		var err error

		listener, err = net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Listen.BindIP, cfg.Listen.Port))
		if err != nil {
			logger.Fatal(err)
		}
	}

	server = &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go shutdown.Graceful([]os.Signal{syscall.SIGABRT, syscall.SIGQUIT, syscall.SIGHUP, os.Interrupt, syscall.SIGTERM},
		server)

	logger.Println("application initialized and started")

	if err := server.Serve(listener); err != nil {
		switch {
		case errors.Is(err, http.ErrServerClosed):
			logger.Warn("server shutdown")
		default:
			logger.Fatal(err)
		}
	}
}
