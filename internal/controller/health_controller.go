package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HealthController struct {
	db *gorm.DB
}

func NewHealthController(db *gorm.DB) *HealthController {
	return &HealthController{db: db}
}

func (c *HealthController) CheckAction(ctx *gin.Context) {

	sqlDB, err := c.db.DB()

	if err != nil {
		ctx.JSON(http.StatusServiceUnavailable, []gin.H{{"name": "database", "result": false, "message": err.Error()}})
		return
	}

	if err := sqlDB.PingContext(ctx.Request.Context()); err != nil {
		ctx.JSON(http.StatusServiceUnavailable, []gin.H{{"name": "database", "result": false, "message": err.Error()}})
		return
	}

	ctx.JSON(http.StatusOK, []gin.H{{"name": "database", "result": true, "message": "ok"}})
}

func (c *HealthController) PingAction(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, []gin.H{{"name": "status", "result": true, "message": "up"}})
}
