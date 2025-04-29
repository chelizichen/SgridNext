package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"sgridnext.com/src/domain/config"
)

type PlatformCredentials struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

func Login(ctx *gin.Context) {
	var req struct {
		Account  string `json:"account"`
		Password string `json:"password"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}

	account := config.Conf.Get("plantform", "account")
	password := config.Conf.Get("plantform", "password")
	if req.Account == account && req.Password == password {
		ctx.JSON(http.StatusOK, gin.H{"success": true, "msg": "登录成功"})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "账号或密码错误"})
	}
}
