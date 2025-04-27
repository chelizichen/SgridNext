package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"sgridnext.com/src/domain/config"
)

func LoadRoutes(engine *gin.Engine) {
	engine.GET("/greet", Greet)
}

func Greet(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "ok",
		"data":    config.Conf["Hello"],
	})
}
