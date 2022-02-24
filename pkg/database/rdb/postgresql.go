package rdb

import (
	"context"
	"net/url"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

func postgresDefaultConfig() *config {
	return &config{
		database:     DatabasePostgres,
		user:         "postgres",
		password:     "postgres",
		host:         "localhost:5432",
		databaseName: "postgres",
		params: map[string]string{
			"sslmode": "disable",
		},
	}
}

func connectToPostgesql(ctx context.Context, cfg *config) (*sqlx.DB, error) {
	connStr := buildPostgresConnString(cfg)
	return sqlx.ConnectContext(ctx, "pgx", connStr)
}

func buildPostgresConnString(cfg *config) string {
	q := make(url.Values)
	for param, value := range cfg.params {
		if param != "" && value != "" {
			q.Set(param, value)
		}
	}

	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(cfg.user, cfg.password),
		Host:     cfg.host,
		Path:     cfg.databaseName,
		RawQuery: q.Encode(),
	}

	return u.String()
}
