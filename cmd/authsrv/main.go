package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/umalmyha/authsrv/internal/handler"
	"github.com/umalmyha/authsrv/internal/service"
	"github.com/umalmyha/authsrv/pkg/database/rdb"
	"github.com/umalmyha/authsrv/pkg/server"
	"github.com/umalmyha/authsrv/pkg/web"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	logger, err := newZapLogger("authentication server")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer logger.Sync()

	if err := start(logger); err != nil {
		logger.Error(err.Error())
	}
}

func start(logger *zap.SugaredLogger) error {
	// load environment variables
	err := loadEnv()
	if err != nil {
		return errors.New(fmt.Sprintf("Error while loading environment variables: %s", err.Error()))
	}

	// init db
	db, err := connectToDb()
	if err != nil {
		return err
	}
	defer db.Close()

	// start server
	return startServer(db, logger)
}

func newZapLogger(service string) (*zap.SugaredLogger, error) {
	config := zap.NewProductionConfig()
	config.DisableStacktrace = true
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.CallerKey = "src"
	config.InitialFields = map[string]interface{}{
		"service": service,
	}

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return logger.Sugar(), nil
}

func loadEnv() error {
	if os.Getenv("APP_ENV") != "production" { // TODO: add normal handling later
		return godotenv.Load()
	}
	return nil
}

func connectToDb() (*sqlx.DB, error) {
	dbConfig := rdb.NewConfig(
		rdb.DatabasePostgres,
		rdb.WithUser(os.Getenv("AUTHSRV_DB_USERNAME")),
		rdb.WithPassword(os.Getenv("AUTHSRV_DB_PASSWORD")),
		rdb.WithDatabase(os.Getenv("AUTHSRV_DB_DBNAME")),
		rdb.WithHost(os.Getenv("AUTHSRV_DB_HOST")),
		rdb.WithParams(
			rdb.Param("sslmode", "disable"),
		),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := rdb.Connect(ctx, dbConfig)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("database connection error: %s", err.Error()))
	}

	return db, nil
}

func startServer(db *sqlx.DB, logger *zap.SugaredLogger) error {
	srvCfg := server.NewConfig(
		server.WithLogger(zap.NewStdLog(logger.Desugar())),
		server.WithHandler(handlerV1(db, logger)),
		server.WithDebugConfig(
			server.WithExpvarDebug(),
			server.WithPprofDebug(),
			server.WithDebugHandler(debugHandlerV1(db, logger)),
		),
	)
	srv := server.New(srvCfg)

	if err := srv.ListenAndServe(); err != nil {
		return errors.New(fmt.Sprintf("server startup error: %s", err.Error()))
	}

	return nil
}

func handlerV1(db *sqlx.DB, logger *zap.SugaredLogger) *chi.Mux {
	r := chi.NewRouter()

	userHandler := handler.NewUserHandler(service.NewUserService(db))

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			r.Post("/signup", web.WithDefaultErrorHandler(userHandler.Signup))
		})

		// r.Route("/roles", func(r chi.Router) {
		// 	r.Get("/", web.WithDefaultErrorHandler(scopeHandler.GetAll))
		// 	r.Post("/", web.WithDefaultErrorHandler(scopeHandler.Create))
		// })
	})

	return r
}

func debugHandlerV1(db *sqlx.DB, logger *zap.SugaredLogger) *chi.Mux {
	// TODO: Add additional routes
	r := chi.NewRouter()

	dbgHandler := handler.NewDebugHandler()
	r.Get("/healthcheck", dbgHandler.Healthcheck)

	return r
}
