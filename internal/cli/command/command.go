package command

import (
	"context"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/umalmyha/authsrv/pkg/database/rdb"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Executor interface {
	Run() error
	Help()
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
		return nil, err
	}
	return db, nil
}

func newZapLogger() (*zap.SugaredLogger, error) {
	config := zap.NewProductionConfig()
	config.DisableCaller = true
	config.DisableStacktrace = true
	config.Encoding = "console"
	config.EncoderConfig.EncodeTime = func(t time.Time, pae zapcore.PrimitiveArrayEncoder) {}
	config.EncoderConfig.EncodeLevel = func(l zapcore.Level, pae zapcore.PrimitiveArrayEncoder) {}
	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return logger.Sugar(), nil
}
