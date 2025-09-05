package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"

	"rest-api-marketplace/internal/config"
	"rest-api-marketplace/internal/repository"
	"rest-api-marketplace/internal/service"
	v1 "rest-api-marketplace/internal/transport/http/v1"
	"rest-api-marketplace/pkg/auth"
	postgres "rest-api-marketplace/pkg/client/postgresdb"
	"rest-api-marketplace/pkg/hash"

	"github.com/go-playground/validator/v10"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
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
		slog.Error("error loading .env file", slog.Any("error", err))
		os.Exit(1)
	}

	slog.Info("config initializing")
	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("cannot load config", slog.Any("error", err))
		os.Exit(1)
	}

	log := setupLogger(cfg.Env)
	log = log.With(slog.String("env", cfg.Env))

	log.Info("initializing server", slog.String("address", cfg.Server.Host+":"+cfg.Server.Port))
	log.Debug("logger debug mode enabled")

	log.Info("initializing dependencies")
	initCtx, cancelInit := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelInit()

	runMigrations(cfg.DB, log)

	db, err := postgres.NewClient(initCtx, cfg.DB, log)
	if err != nil {
		log.Error("failed to connect database", slog.String("error", err.Error()))
		os.Exit(1)
	}

	tokenManager, err := auth.NewManager(cfg.Auth.SigningKey)
	if err != nil {
		log.Error("failed to init token manager", slog.String("error", err.Error()))
		os.Exit(1)
	}
	passwordHasher := hash.NewBcryptHasher(bcrypt.DefaultCost)
	v := validator.New()

	repos := repository.NewRepositories(db)

	services := service.NewServices(service.Deps{
		Logger:          log,
		Repos:           repos,
		Hasher:          passwordHasher,
		TokenManager:    tokenManager,
		AccessTokenTTL:  cfg.Auth.AccessTokenTTL,
		RefreshTokenTTL: cfg.Auth.RefreshTokenTTL,
	})

	handler := v1.NewHandler(services, tokenManager)

	e := echo.New()
	e.Validator = &CustomValidator{validator: v}
	e.HTTPErrorHandler = customErrorHandler(log)

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

func customErrorHandler(log *slog.Logger) echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		code := http.StatusInternalServerError
		message := "internal server error"

		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
			if msg, ok := he.Message.(string); ok {
				message = msg
			}
		}

		if code >= 500 {
			log.Error("internal server error", slog.String("error", err.Error()), slog.String("request_uri", c.Request().RequestURI))
		}

		if !c.Response().Committed {
			c.JSON(code, map[string]string{
				"error": message,
			})
		}
	}
}

func runMigrations(cfg config.PostgresConfig, log *slog.Logger) {
	migrationDNS := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
	)

	m, err := migrate.New("file://migrations", migrationDNS)
	if err != nil {
		log.Error("failed to create migrate instance", slog.Any("error", err))
		os.Exit(1)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Error("failed to apply migrations", slog.Any("error", err))
		os.Exit(1)
	}

	log.Info("migrations successfully applied")
}
