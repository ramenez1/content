package api

import (
	"ContentSystem/internal/service"
	"github.com/gin-gonic/gin"
)

const (
	rootPath   = "/api/"
	noAuthPath = "/out/api/" //不需要鉴权的接口
)

func CmsRouters(r *gin.Engine) {
	cmsApp := service.NewCmsApp()
	session := NewSessionAuth() //鉴权
	root := r.Group(rootPath).Use(session.Auth)
	{
		//  /api/cms/ping
		root.GET("/cms/hello", cmsApp.Hello)
		//  /api/cms/content/create
		root.POST("/cms/content/create", cmsApp.ContentCreate)
		//  /api/cms/content/update
		root.POST("/cms/content/update", cmsApp.ContentUpdate)
		//  /api/cms/content/delete
		root.POST("/cms/content/delete", cmsApp.ContentDelete)
		//  /api/cms/content/find
		root.POST("/cms/content/find", cmsApp.ContentFind)
	}

	noAuth := r.Group(noAuthPath)
	{
		// /out/api/cms/register
		noAuth.POST("/cms/register", cmsApp.Register)
		// /out/api/cms/login
		noAuth.POST("/cms/login", cmsApp.Login)
	}
}
