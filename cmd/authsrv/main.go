package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/go-chi/chi/v5"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/umalmyha/authsrv/internal/handler"
	"github.com/umalmyha/authsrv/internal/infrastruct"
	"github.com/umalmyha/authsrv/internal/service"
	redisdb "github.com/umalmyha/authsrv/pkg/database/redis"
	"github.com/umalmyha/authsrv/pkg/server"
	"github.com/umalmyha/authsrv/pkg/web"
	"go.uber.org/zap"
)

func main() {
	logger, err := infrastruct.NewZapProductionLogger("authentication server")
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
	err := infrastruct.LoadEnv()
	if err != nil {
		return errors.New(fmt.Sprintf("Error while loading environment variables: %s", err.Error()))
	}

	// init db
	db, err := infrastruct.ConnectToDb()
	if err != nil {
		return err
	}
	defer db.Close()

	redisOpts, err := infrastruct.RedisOptions()
	if err != nil {
		return err
	}

	rdb, err := redisdb.Connect(redisOpts)
	if err != nil {
		return err
	}

	// start server
	return startServer(db, rdb, logger)
}

func startServer(db *sqlx.DB, rdb *redis.Client, logger *zap.SugaredLogger) error {
	handler, err := handlerV1(db, rdb, logger)
	if err != nil {
		return err
	}

	srvCfg := server.NewConfig(
		server.WithLogger(zap.NewStdLog(logger.Desugar())),
		server.WithHandler(handler),
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

func handlerV1(db *sqlx.DB, rdb *redis.Client, logger *zap.SugaredLogger) (*chi.Mux, error) {
	r := chi.NewRouter()

	jwtCfg, err := infrastruct.JwtConfig()
	if err != nil {
		return nil, err
	}

	rfrCfg, err := infrastruct.RefreshTokenConfig()
	if err != nil {
		return nil, err
	}

	passCfg, err := infrastruct.PasswordConfig()
	if err != nil {
		return nil, err
	}

	authService := service.NewAuthService(db, rdb, jwtCfg, rfrCfg, passCfg)
	authHandler := handler.NewAuthHandler(authService, rfrCfg)

	r.Route("/api", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/signup", web.WithDefaultErrorHandler(authHandler.Signup))
			r.Post("/signin", web.WithDefaultErrorHandler(authHandler.Signin))
			r.Post("/logout", web.WithDefaultErrorHandler(authHandler.Logout))
			r.Post("/refresh", web.WithDefaultErrorHandler(authHandler.RefreshSession))
		})
	})

	return r, nil
}

func debugHandlerV1(db *sqlx.DB, logger *zap.SugaredLogger) *chi.Mux {
	// TODO: Add additional routes
	r := chi.NewRouter()

	dbgHandler := handler.NewDebugHandler()
	r.Get("/healthcheck", dbgHandler.Healthcheck)

	return r
}
