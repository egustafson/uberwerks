package api

import (
	"context"
	"net/http"
	"time"

	"github.com/Masterminds/log-go"
	"github.com/gin-gonic/gin"
)

func Run(ctx context.Context) (err error) {
	router := gin.Default()
	v0 := router.Group("/v0")
	v0.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	s := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
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
	log.Info("http server shutdown complete")
	return
}
