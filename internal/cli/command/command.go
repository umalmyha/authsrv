package command

import (
	"context"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/umalmyha/authsrv/pkg/database"
	"go.uber.org/zap"
)

type Executor interface {
	Run() error
	Help()
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
		return nil, err
	}
	return db, nil
}

func newZapLogger() (*zap.SugaredLogger, error) {
	config := zap.NewDevelopmentConfig()
	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return logger.Sugar(), nil
}
