package server

import (
	"context"
	"expvar"
	"fmt"
	"log"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/errors"
)

type server struct {
	logger          *log.Logger
	httpServer      *http.Server
	debugServer     *http.Server
	shutdownTimeout time.Duration
}

func New(cfg *config) *server {
	httpServer := httpServer(cfg)

	srv := &server{
		httpServer:      httpServer,
		shutdownTimeout: cfg.shutdownTimeout,
	}

	if cfg.debug != nil {
		srv.debugServer = debugServer(cfg.debug)
	}

	if cfg.logger != nil {
		srv.logger = cfg.logger
	} else {
		srv.logger = log.New(os.Stdout, "server: ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
	}

	return srv
}

func (s *server) ListenAndServe() error {
	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, syscall.SIGINT, syscall.SIGTERM)

	errorsCh := make(chan error, 1)

	go func() {
		s.logger.Printf("starting server on %s", s.httpServer.Addr)
		errorsCh <- s.httpServer.ListenAndServe()
	}()

	if s.debugServer != nil {
		s.logger.Printf("starting debug server on %s", s.debugServer.Addr)
		go func() {
			if err := s.debugServer.ListenAndServe(); err != nil {
				s.logger.Print("failed to start debug server, main server might start normally: %w", err)
			}
		}()
	}

	select {
	case err := <-errorsCh:
		s.debugServer.Close()
		return errors.Wrap(err, "server runtime error")

	case sig := <-shutdownCh:
		s.logger.Printf("%s shutdown signal has been sent", sig)
		return s.handleShutdown()
	}
}

func (s *server) handleShutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	defer s.debugServer.Close()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.httpServer.Close()
		return errors.Errorf("failed to stop server gracefully: %w", err)
	}

	s.logger.Print("server has been stopped gracefully")

	return nil
}

func httpServer(cfg *config) *http.Server {
	return &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      cfg.handler,
		ReadTimeout:  cfg.readTimeout,
		WriteTimeout: cfg.writeTimeout,
		IdleTimeout:  cfg.idleTimeout,
		ErrorLog:     cfg.logger,
	}
}

func debugServer(cfg *debugConfig) *http.Server {
	mux := http.NewServeMux()

	if cfg.pprofDebug {
		mux.HandleFunc("/debug/pprof/", pprof.Index)
		mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	}

	if cfg.expvarDebug {
		mux.Handle("/debug/expvar", expvar.Handler())
	}

	if cfg.handler != nil {
		mux.Handle("/debug", cfg.handler)
	}

	return &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.port),
		Handler: mux,
	}
}
