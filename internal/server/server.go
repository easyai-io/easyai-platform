package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/easyai-io/easyai-platform/internal/server/config"
	_ "github.com/easyai-io/easyai-platform/internal/server/swagger" // swagger docs
	"github.com/easyai-io/easyai-platform/pkg/logger"
)

type options struct {
	ConfigFile         string
	WWWDir             string
	Version            string
	GracePeriodSeconds int
}

// Option server opts
type Option func(*options)

// SetConfigFile set config
func SetConfigFile(s string) Option {
	return func(o *options) {
		o.ConfigFile = s
	}
}

// SetWWWDir set html dir
func SetWWWDir(s string) Option {
	return func(o *options) {
		o.WWWDir = s
	}
}

// SetVersion set version
func SetVersion(s string) Option {
	return func(o *options) {
		o.Version = s
	}
}

// Init init server
func Init(ctx context.Context, opts ...Option) (func(), error) {
	// TODO init server
	return func() {}, nil
}

// InitHTTPServer init http server
func InitHTTPServer(ctx context.Context, handler http.Handler) func() {
	cfg := config.Global().HTTP
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 90 * time.Second, // case we will use a pprof with 60s
		IdleTimeout:  30 * time.Second,
	}

	go func() {
		logger.WithContext(ctx).Printf("HTTP server is running at %s.", addr)

		var err error
		if cfg.CertFile != "" && cfg.KeyFile != "" {
			srv.TLSConfig = &tls.Config{MinVersion: tls.VersionTLS12}
			err = srv.ListenAndServeTLS(cfg.CertFile, cfg.KeyFile)
		} else {
			err = srv.ListenAndServe()
		}
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}

	}()

	return func() {
		ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(cfg.ShutdownTimeout))
		defer cancel()

		srv.SetKeepAlivesEnabled(false)
		if err := srv.Shutdown(ctx); err != nil {
			logger.WithContext(ctx).Errorf(err.Error())
		}
	}
}

// Run app run
func Run(ctx context.Context, opts ...Option) error {
	state := 1
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	cleanFunc, err := Init(ctx, opts...)
	if err != nil {
		return err
	}

EXIT:
	for {
		sig := <-sc
		logger.WithContext(ctx).Infof("Receive signal[%s]", sig.String())
		switch sig {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			state = 0
			break EXIT
		case syscall.SIGHUP:
		default:
			break EXIT
		}
	}

	cleanFunc()
	os.Exit(state)
	return nil
}
