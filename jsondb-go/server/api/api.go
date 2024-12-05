package api

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/egustafson/uberwerks/jsondb-go/server/config"
)

func Run(ctx context.Context, config *config.Config) (err error) {
	router := gin.Default()

	h0 := router.Group("/health")
	healthAPI, err := initHealth(config, h0)
	if err != nil {
		return err
	}

	v0 := router.Group("/jsondb/v0")
	jdbAPI, err := InitJdbAPI(config, v0)
	if err != nil {
		return err
	}

	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Port),
		Handler: router,
	}
	slog.Info("listening", slog.Int("port", config.Port))
	idleConnsClosed := make(chan struct{})
	go func() {
		<-ctx.Done() // block until context canceled
		const timeoutTime = 5 * time.Second
		timeoutCtx, cancel := context.WithTimeout(context.Background(), timeoutTime)
		defer cancel()
		if err := s.Shutdown(timeoutCtx); err != nil {
			slog.Warn("http server Shutdown", slog.Any("error", err))
			s.Close()
		}
		close(idleConnsClosed)
	}()
	if err = s.ListenAndServe(); err != http.ErrServerClosed {
		slog.Error("http server ListenAndServe", slog.Any("error", err))
	} else {
		err = nil // don't report the server closing, it's normal
	}

	<-idleConnsClosed
	jdbAPI.Close()
	healthAPI.Close()
	slog.Info("http server shutdown complete")
	return
}
