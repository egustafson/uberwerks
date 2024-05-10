package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Masterminds/log-go"
	"github.com/gin-gonic/gin"

	"github.com/egustafson/uberwerks/jsondb-go/server/config"
)

func Run(ctx context.Context, config *config.Config) (err error) {
	router := gin.Default()
	v0 := router.Group("/jsondb/v0")
	jdbAPI, err := InitJdbAPI(config, v0)
	if err != nil {
		return err
	}

	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Port),
		Handler: router,
	}
	log.Infof("listening on port %d", config.Port)
	idleConnsClosed := make(chan struct{})
	go func() {
		<-ctx.Done() // block until context canceled
		const timeoutTime = 5 * time.Second
		timeoutCtx, cancel := context.WithTimeout(context.Background(), timeoutTime)
		defer cancel()
		if err := s.Shutdown(timeoutCtx); err != nil {
			log.Warnf("http server Shutdown: %v", err)
			s.Close()
		}
		close(idleConnsClosed)
	}()
	if err = s.ListenAndServe(); err != http.ErrServerClosed {
		log.Errorf("http server ListenAndServe: %v", err)
	} else {
		err = nil // don't report the server closing, it's normal
	}

	<-idleConnsClosed
	jdbAPI.Close()
	log.Info("http server shutdown complete")
	return
}
