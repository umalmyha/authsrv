package server

import (
	"log"
	"net/http"
	"time"
)

const defaultPort = 4004
const defaultReadTimeout = 5 * time.Second
const defaultWriteTimeout = 10 * time.Second
const defaultIdleTimeout = 30 * time.Second
const defaultShutdownTimeout = 60 * time.Second

type config struct {
	port            int
	readTimeout     time.Duration
	writeTimeout    time.Duration
	idleTimeout     time.Duration
	shutdownTimeout time.Duration
	maxHeaderBytes  int
	handler         http.Handler
	logger          *log.Logger
	debug           *debugConfig
}

type configOptionFunc func(*config)

func NewConfig(opts ...configOptionFunc) *config {
	cfg := &config{
		port:            defaultPort,
		readTimeout:     defaultReadTimeout,
		writeTimeout:    defaultWriteTimeout,
		idleTimeout:     defaultIdleTimeout,
		shutdownTimeout: defaultShutdownTimeout,
	}

	for _, optFn := range opts {
		optFn(cfg)
	}

	return cfg
}

func WithPort(port int) configOptionFunc {
	return func(sc *config) {
		sc.port = port
	}
}

func WithReadTimeout(rt time.Duration) configOptionFunc {
	return func(sc *config) {
		sc.readTimeout = rt
	}
}

func WithWriteTimeout(wt time.Duration) configOptionFunc {
	return func(sc *config) {
		sc.writeTimeout = wt
	}
}

func WithIdleTimeout(it time.Duration) configOptionFunc {
	return func(sc *config) {
		sc.idleTimeout = it
	}
}

func WithHandler(h http.Handler) configOptionFunc {
	return func(sc *config) {
		sc.handler = h
	}
}

func WithLogger(logger *log.Logger) configOptionFunc {
	return func(sc *config) {
		sc.logger = logger
	}
}

func WithMaxHeaderBytes(hb int) configOptionFunc {
	return func(sc *config) {
		sc.maxHeaderBytes = hb
	}
}

func WithShutdownTimeout(st time.Duration) configOptionFunc {
	return func(sc *config) {
		sc.shutdownTimeout = st
	}
}

func WithDebugConfig(opts ...debugConfigOptionFunc) configOptionFunc {
	return func(sc *config) {
		cfg := debugConfigWithDefaults()

		for _, optFn := range opts {
			optFn(cfg)
		}

		sc.debug = cfg
	}
}
