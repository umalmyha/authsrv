package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/umalmyha/authsrv/internal/dbg"
	"github.com/umalmyha/authsrv/internal/user"
	userStore "github.com/umalmyha/authsrv/internal/user/store"
	"github.com/umalmyha/authsrv/pkg/database"
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
		return errors.Wrap(err, "loading environment variables")
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
	dbConfig := database.NewConfig(
		database.DatabasePostgres,
		database.WithUser(os.Getenv("AUTHSRV_DB_USERNAME")),
		database.WithPassword(os.Getenv("AUTHSRV_DB_PASSWORD")),
		database.WithDatabase(os.Getenv("AUTHSRV_DB_DBNAME")),
		database.WithHost(os.Getenv("AUTHSRV_DB_HOST")),
		database.WithParams(
			database.Param("sslmode", "disable"),
		),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := database.Connect(ctx, dbConfig)
	if err != nil {
		return nil, errors.Wrap(err, "database connection")
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
		return errors.Wrap(err, "server startup error")
	}

	return nil
}

func handlerV1(db *sqlx.DB, logger *zap.SugaredLogger) *chi.Mux {
	r := chi.NewRouter()

	userHandler := user.Handler(user.Service(logger, userStore.NewStore(db)))

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			r.Post("/signup", web.WithDefaultErrorHandler(userHandler.Signup))
			r.Get("/", web.WithDefaultErrorHandler(userHandler.GetAll))
			r.Get("/{id}", web.WithDefaultErrorHandler(userHandler.Get))
			r.Patch("/{id}", web.WithDefaultErrorHandler(userHandler.Update))
			r.Put("/{id}", web.WithDefaultErrorHandler(userHandler.Update))
			r.Delete("/{id}", web.WithDefaultErrorHandler(userHandler.Delete))
		})
	})

	return r
}

func debugHandlerV1(db *sqlx.DB, logger *zap.SugaredLogger) *chi.Mux {
	// TODO: Add additional routes
	r := chi.NewRouter()

	dbgHandler := dbg.Handler()
	r.Get("/healthcheck", dbgHandler.Healthcheck)

	return r
}
