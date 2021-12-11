package server

import "net/http"

const defaultDebugPort = 5004

type debugConfig struct {
	port        int
	handler     http.Handler
	pprofDebug  bool
	expvarDebug bool
}

type debugConfigOptionFunc func(*debugConfig)

func debugConfigWithDefaults() *debugConfig {
	cfg := &debugConfig{
		port: defaultDebugPort,
	}
	return cfg
}

func WithDebugPort(port int) debugConfigOptionFunc {
	return func(dc *debugConfig) {
		dc.port = port
	}
}

func WithDebugHandler(h http.Handler) debugConfigOptionFunc {
	return func(dc *debugConfig) {
		dc.handler = h
	}
}

func WithPprofDebug() debugConfigOptionFunc {
	return func(dc *debugConfig) {
		dc.pprofDebug = true
	}
}

func WithExpvarDebug() debugConfigOptionFunc {
	return func(dc *debugConfig) {
		dc.expvarDebug = true
	}
}
