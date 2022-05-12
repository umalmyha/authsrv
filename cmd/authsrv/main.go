package main

import (
	"log"

	"github.com/pkg/errors"

	"github.com/go-chi/chi/v5"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jmoiron/sqlx"
	valueobj "github.com/umalmyha/authsrv/internal/business/value-object"
	"github.com/umalmyha/authsrv/internal/infra"
	"github.com/umalmyha/authsrv/internal/infra/handler"
	"github.com/umalmyha/authsrv/internal/infra/service"
	redisdb "github.com/umalmyha/authsrv/pkg/database/redis"
	"github.com/umalmyha/authsrv/pkg/web"
	"github.com/umalmyha/authsrv/pkg/web/middleware"
	"github.com/umalmyha/authsrv/pkg/web/server"
	"go.uber.org/zap"
)

func main() {
	logger, err := infra.NewZapProductionLogger("authentication server")
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	if err := start(logger); err != nil {
		logger.Error(err)
	}
}

func start(logger *zap.SugaredLogger) error {
	// load environment variables
	err := infra.LoadEnv()
	if err != nil {
		return errors.Wrap(err, "error while loading environment variables")
	}

	// init db
	db, err := infra.ConnectToDb()
	if err != nil {
		return errors.Wrap(err, "failed to connect to db")
	}
	defer db.Close()

	redisOpts, err := infra.RedisOptions()
	if err != nil {
		return errors.Wrap(err, "failed to build redis options")
	}

	rdb, err := redisdb.Connect(redisOpts)
	if err != nil {
		return errors.Wrap(err, "failed to connect to redis")
	}

	// start server
	return startServer(db, rdb, logger)
}

func startServer(db *sqlx.DB, rdb *redis.Client, logger *zap.SugaredLogger) error {
	stdLoger := zap.NewStdLog(logger.Desugar())
	handler, err := handlerV1(db, rdb, stdLoger)
	if err != nil {
		return errors.Wrap(err, "failed to build handler")
	}

	srvCfg := server.NewConfig(
		server.WithLogger(stdLoger),
		server.WithHandler(handler),
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

func handlerV1(db *sqlx.DB, rdb *redis.Client, logger *log.Logger) (*chi.Mux, error) {
	r := chi.NewRouter()

	jwtCfg, err := infra.JwtConfig()
	if err != nil {
		return nil, errors.Wrap(err, "failed to JWT config")
	}

	rfrCfg, err := infra.RefreshTokenConfig()
	if err != nil {
		return nil, errors.Wrap(err, "failed to build refresh token config")
	}

	passCfg, err := infra.PasswordConfig()
	if err != nil {
		return nil, errors.Wrap(err, "failed to build password config")
	}

	// servcices and handlers
	authService := service.NewAuthService(db, rdb, jwtCfg, rfrCfg, passCfg)
	authHandler := handler.NewAuthHandler(authService, rfrCfg)

	scopeService := service.NewScopeService(db)
	scopeHandler := handler.NewScopeHandler(scopeService)

	roleService := service.NewRoleService(db)
	roleHandler := handler.NewRoleHandler(roleService)

	userService := service.NewUserService(db, rdb)
	userHandler := handler.NewUserHandler(userService)

	// middleware
	loggerMw := middleware.RequestLogger(logger)

	jwtValidator := func(rawToken string) (middleware.AuthClaimsProvider, error) {
		var claims valueobj.JwtClaims
		parser := jwt.NewParser(jwt.WithValidMethods([]string{jwtCfg.Algorithm()}))

		keyFunc := func(token *jwt.Token) (any, error) {
			if token.Method.Alg() != jwtCfg.Algorithm() {
				return nil, errors.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtCfg.PublicKey(), nil
		}

		if _, err := parser.ParseWithClaims(rawToken, &claims, keyFunc); err != nil {
			return nil, err
		}

		return claims, nil
	}

	jwtAuthMw := middleware.JwtAuthentication(jwtValidator)

	r.Route("/api", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/signup", web.HttpHandlerFunc(middleware.Wrap(authHandler.Signup, middleware.RequestId, loggerMw)))
			r.Post("/signin", web.HttpHandlerFunc(middleware.Wrap(authHandler.Signin, middleware.RequestId, loggerMw)))
			r.Post("/logout", web.HttpHandlerFunc(middleware.Wrap(authHandler.Logout, middleware.RequestId, loggerMw, jwtAuthMw)))
			r.Post("/refresh", web.HttpHandlerFunc(middleware.Wrap(authHandler.RefreshSession, middleware.RequestId, loggerMw)))
		})

		r.Route("/scopes", func(r chi.Router) {
			r.Post("/", web.HttpHandlerFunc(middleware.Wrap(scopeHandler.CreateScope, middleware.RequestId, loggerMw, jwtAuthMw)))
		})

		r.Route("/roles", func(r chi.Router) {
			r.Post("/", web.HttpHandlerFunc(middleware.Wrap(roleHandler.CreateRole, middleware.RequestId, loggerMw, jwtAuthMw)))
			r.Post("/assign", web.HttpHandlerFunc(middleware.Wrap(roleHandler.AssignScope, middleware.RequestId, loggerMw, jwtAuthMw)))
			r.Post("/unassign", web.HttpHandlerFunc(middleware.Wrap(roleHandler.UnassignScope, middleware.RequestId, loggerMw, jwtAuthMw)))
		})

		r.Route("/users", func(r chi.Router) {
			r.Post("/assign", web.HttpHandlerFunc(middleware.Wrap(userHandler.AssignRole, middleware.RequestId, loggerMw, jwtAuthMw)))
			r.Post("/unassign", web.HttpHandlerFunc(middleware.Wrap(userHandler.UnassignRole, middleware.RequestId, loggerMw, jwtAuthMw)))
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
