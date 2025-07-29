package app

import (
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"rest-api-marketplace/internal/config"
	"rest-api-marketplace/internal/repository"
	"rest-api-marketplace/internal/service"
	v1 "rest-api-marketplace/internal/transport/http/v1"
	"rest-api-marketplace/pkg/auth"
	postgres "rest-api-marketplace/pkg/client/postgresdb"
	"rest-api-marketplace/pkg/hash"
	"time"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func Run() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading .env file: %v", err)
		return
	}

	log.Println("config initializing")
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("cannot load config %v", err)
	}

	log := setupLogger(cfg.Env)
	log = log.With(slog.String("env", cfg.Env))

	log.Info("Initializing server", slog.String("address", cfg.Server.Host+":"+cfg.Server.Port))
	log.Debug("logger debug mode enabled")

	log.Info("Initializing dependencies")
	initCtx, cancelInit := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelInit()

	db, err := postgres.NewClient(initCtx, cfg.DB)
	if err != nil {
		log.Error("failed to connect database", slog.String("error", err.Error()))
		os.Exit(1)
	}

	tokenManager, err := auth.NewManager(cfg.Auth.SigningKey)
	if err != nil {
		log.Error("failed to init token manager", slog.String("error", err.Error()))
	}
	passwordHasher := hash.NewBcryptHasher(bcrypt.DefaultCost)
	v := validator.New()

	repos := repository.NewRepositories(db)

	services := service.NewServices(service.Deps{
		Repos:           repos,
		Hasher:          passwordHasher,
		TokenManager:    tokenManager,
		AccessTokenTTL:  cfg.Auth.AccessTokenTTL,
		RefreshTokenTTL: cfg.Auth.RefreshTokenTTL,
	})

	handler := v1.NewHandler(services, tokenManager)

	e := echo.New()
	e.Validator = &CustomValidator{validator: v}

	handler.Init(e.Group("/api"))

	listener, err := net.Listen("tcp", cfg.Server.Host+":"+cfg.Server.Port)
	if err != nil {
		log.Error("failed to bind to address", slog.String("error", err.Error()))
		os.Exit(1)
	}

	server := &http.Server{
		Handler:           e,
		WriteTimeout:      15 * time.Second,
		ReadHeaderTimeout: 15 * time.Second,
	}

	err = server.Serve(listener)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error("failed to start server", slog.String("error", err.Error()))
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	// TODO: сделать функцию, которая будет игнорировать сообщения, отпр-ые в логгер для тестов
	return log
}
