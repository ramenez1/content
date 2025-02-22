package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type HelloReq struct {
	Name string `json:"name" bingding:"required" message:"name is required"`
}

type HelloResp struct {
	Message string `json:"message" bingding:"required"`
}

func (c *CmsApp) Hello(ctx *gin.Context) {
	var req HelloReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok ",
		"data": &HelloResp{
			Message: fmt.Sprintf("hello %s", req.Name),
		},
	})
}
