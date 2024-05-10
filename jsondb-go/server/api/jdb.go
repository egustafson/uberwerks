package api

/* GET  /      - return list of obj ids
 * POST /      - add new object, return object w/ _id
 * GET  /<id>  - return object or error
 * PUT  /<id>  - update object or create
 * DEL  /<id>  - delete object
 *
 */

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/egustafson/uberwerks/jsondb-go/jsondb"
	"github.com/egustafson/uberwerks/jsondb-go/server/config"
	"github.com/egustafson/uberwerks/jsondb-go/server/jdb"
)

type JdbAPI struct {
	db jdb.JDB
}

func InitJdbAPI(config *config.Config, rg *gin.RouterGroup) (*JdbAPI, error) {

	jDB, err := jdb.InitSqliteJDB(config.DSN)
	if err != nil {
		return nil, err
	}
	jdbAPI := &JdbAPI{db: jDB}

	rg.GET("/", jdbAPI.List)
	rg.POST("/", jdbAPI.Create)
	rg.GET("/:id", jdbAPI.Get)
	rg.PUT("/:id", jdbAPI.Update)
	rg.DELETE("/:id", jdbAPI.Delete)

	rg.GET("/ping", jdbAPI.Ping)

	return jdbAPI, nil
}

func (api *JdbAPI) Close() {
	//
	// TODO: impl close down into jDB
	//
}

func (api *JdbAPI) Ping(c *gin.Context) {
	err := api.db.Ping()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("db err: %v", err.Error()),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func (api *JdbAPI) List(c *gin.Context) {
	idList, err := api.db.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("db err: %v", err.Error()),
		})
		return
	}
	c.JSON(http.StatusOK, idList)
}

func (api *JdbAPI) Create(c *gin.Context) {
	var jo jsondb.JSONObj
	if err := c.BindJSON(&jo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("parse err: %v", err.Error()),
		})
		return
	}
	jo, err := api.db.Put(jo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("db err: %v", err.Error()),
		})
	}
	c.JSON(http.StatusOK, jo)
}

func (api *JdbAPI) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "id is malformed uuid",
		})
		return
	}
	obj, ok, err := api.db.Get(jsondb.JID(id.String()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("db error: %v", err.Error()),
		})
		return
	}
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "id not found",
		})
		return
	}
	c.JSON(http.StatusOK, obj)
}

func (api *JdbAPI) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "id is malformed uuid",
		})
		return
	}
	//
	// TODO implement
	//
	_ = id
}

func (api *JdbAPI) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "id is malformed uuid",
		})
		return
	}
	ok, err := api.db.Del(jsondb.JID(id.String()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("db error: %v", err.Error()),
		})
		return
	}
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "id not found",
		})
		return
	}
	c.Status(http.StatusNoContent)
}
