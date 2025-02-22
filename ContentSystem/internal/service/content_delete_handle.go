package service

import (
	"ContentSystem/internal/api/operate"
	//"ContentSystem/internal/dao"
	//"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ContentDeleteReq struct {
	ID int64 `json:"id" binding:"required"`
}

type ContentDeleteResp struct {
	Message string `json:"message"`
}

func (c *CmsApp) ContentDelete(ctx *gin.Context) {
	var req ContentDeleteReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	rsp, err := c.operateAppClient.DeleteContent(ctx, &operate.DeleteContentReq{Id: req.ID})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
		"data": rsp,
	})

	////判断内容存在
	//contentDao := dao.NewContentDao(c.db)
	//isExist, err := contentDao.IsExist(req.ID)
	//if err != nil {
	//	ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	//	return
	//}
	//if !isExist {
	//	ctx.JSON(http.StatusBadRequest, gin.H{"error": "Content not found"})
	//	return
	//}
	////删除内容
	//if err := contentDao.Delete(req.ID); err != nil {
	//	ctx.JSON(http.StatusInternalServerError, gin.H{"delete error": err.Error()})
	//}
	//
	//ctx.JSON(http.StatusOK, gin.H{
	//	"code": 0,
	//	"msg":  "ok",
	//	"data": &ContentDeleteResp{
	//		Message: fmt.Sprintf("ok"),
	//	},
	//})
}
