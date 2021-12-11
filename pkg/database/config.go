package database

import "time"

type config struct {
	database        database
	user            string
	password        string
	host            string
	databaseName    string
	maxOpenConns    int
	maxIdleConns    int
	connMaxLifetime time.Duration
	connMaxIdleTime time.Duration
	params          map[string]string
}

type configParam struct {
	prop  string
	value string
}

type configOptionFunc func(*config)

func NewConfig(db database, opts ...configOptionFunc) *config {
	var cfg *config
	switch db {
	case DatabasePostgres:
		cfg = postgresDefaultConfig()
	default:
		cfg = &config{params: make(map[string]string)}
	}

	for _, optFn := range opts {
		optFn(cfg)
	}

	return cfg
}

func WithUser(user string) configOptionFunc {
	return func(c *config) {
		c.user = user
	}
}

func WithPassword(password string) configOptionFunc {
	return func(c *config) {
		c.password = password
	}
}

func WithHost(host string) configOptionFunc {
	return func(c *config) {
		c.host = host
	}
}

func WithDatabase(dbName string) configOptionFunc {
	return func(c *config) {
		c.databaseName = dbName
	}
}

func WithMaxOpenConns(maxOpenConns int) configOptionFunc {
	return func(c *config) {
		c.maxOpenConns = maxOpenConns
	}
}

func WithMaxIdleConns(maxIdleConns int) configOptionFunc {
	return func(c *config) {
		c.maxIdleConns = maxIdleConns
	}
}

func WithConnMaxLifetime(connMaxLifetime time.Duration) configOptionFunc {
	return func(c *config) {
		c.connMaxLifetime = connMaxLifetime
	}
}

func WithConnMaxIdleTime(connMaxIdleTime time.Duration) configOptionFunc {
	return func(c *config) {
		c.connMaxIdleTime = connMaxIdleTime
	}
}

func WithParams(params ...configParam) configOptionFunc {
	return func(c *config) {
		for _, param := range params {
			if param.prop != "" && param.value != "" {
				c.params[param.prop] = param.value
			}
		}
	}
}

func Param(prop string, value string) configParam {
	return configParam{
		prop:  prop,
		value: value,
	}
}
