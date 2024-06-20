package api

import (
	"net/http"

	"github.com/egustafson/uberwerks/jsondb-go/server/config"
	"github.com/gin-gonic/gin"
)

type HealthAPI struct{}

func initHealth(config *config.Config, rg *gin.RouterGroup) (*HealthAPI, error) {
	healthAPI := new(HealthAPI)
	rg.GET("/", healthAPI.Healthz)
	rg.GET("/healthz", healthAPI.Healthz)
	return healthAPI, nil
}

func (api *HealthAPI) Close() {
	// nothing to do
}

func (api *HealthAPI) Healthz(c *gin.Context) {
	c.JSON(http.StatusOK, "ok")
}
