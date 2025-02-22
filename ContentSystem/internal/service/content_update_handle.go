package service

import (
	"ContentSystem/internal/api/operate"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type ContentUpdateReq struct {
	ID    int64  `json:"id" binding:"required"`
	Title string `json:"title"` // 内容标题
	//ContentID      string        `gorm:"column:content_id"`      // 内容标题
	Description    string        `json:"description"`     // 内容描述
	Author         string        `json:"author"`          // 作者
	VideoURL       string        `json:"video_url"`       // 视频播放URL
	Thumbnail      string        `json:"thumbnail"`       // 封面图URL
	Category       string        `json:"category"`        // 内容分类
	Duration       time.Duration `json:"duration"`        // 内容时长
	Resolution     string        `json:"resolution"`      // 分辨率 如720p、1080p
	FileSize       int64         `json:"fileSize"`        // 文件大小
	Format         string        `json:"format"`          // 文件格式 如MP4、AVI
	Quality        int32         `json:"quality"`         // 视频质量 1-高清 2-标清
	ApprovalStatus int32         `json:"approval_status"` // 审核状态 1-审核中 2-审核通过 3-审核不通过
	UpdatedAt      time.Time     `json:"updated_at"`      // 内容更新时间
	CreatedAt      time.Time     `json:"created_at"`      // 内容创建时间
}

type ContentUpdateResp struct {
	Message string `json:"message"`
}

func (c *CmsApp) ContentUpdate(ctx *gin.Context) {
	var req ContentUpdateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rsp, err := c.operateAppClient.UpdateContent(ctx, &operate.UpdateContentReq{
		Content: &operate.Content{
			Id:             req.ID,
			Title:          req.Title,
			Description:    req.Description,
			Author:         req.Author,
			VideoUrl:       req.VideoURL,
			Thumbnail:      req.Thumbnail,
			Category:       req.Category,
			Duration:       req.Duration.Milliseconds(),
			Resolution:     req.Resolution,
			FileSize:       req.FileSize,
			Format:         req.Format,
			Quality:        req.Quality,
			ApprovalStatus: req.ApprovalStatus,
		},
	})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Content not found"})
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
	//err = contentDao.Update(req.ID, model.ContentDetail{
	//	Title:          req.Title,
	//	Description:    req.Description,
	//	Author:         req.Author,
	//	VideoURL:       req.VideoURL,
	//	Thumbnail:      req.Thumbnail,
	//	Category:       req.Category,
	//	Duration:       req.Duration,
	//	Resolution:     req.Resolution,
	//	FileSize:       req.FileSize,
	//	Format:         req.Format,
	//	Quality:        req.Quality,
	//	ApprovalStatus: req.ApprovalStatus,
	//})
	//if err != nil {
	//	ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	//}
	//ctx.JSON(http.StatusOK, gin.H{
	//	"code": 0,
	//	"msg":  "ok",
	//	"data": &ContentUpdateResp{
	//		Message: "Content updated successfully",
	//	},
	//})
}
