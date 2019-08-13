package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func MakePingHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
			"err_code": 0,
			"msg":      "pong",
		})
	}
}
