package service

import (
	"ContentSystem/internal/api/operate"
	//"ContentSystem/internal/dao"
	//"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	//"time"
)

//	type Content struct {
//		ID    int    `json:"id"`                       // 内容标题
//		Title string `json:"title" binding:"required"` // 内容标题
//		//ContentID      string        `gorm:"column:content_id"`      // 内容标题
//		Description    string        `json:"description"`                  // 内容描述
//		Author         string        `json:"author" binding:"required"`    // 作者
//		VideoURL       string        `json:"video_url" binding:"required"` // 视频播放URL
//		Thumbnail      string        `json:"thumbnail"`                    // 封面图URL
//		Category       string        `json:"category"`                     // 内容分类
//		Duration       time.Duration `json:"duration"`                     // 内容时长
//		Resolution     string        `json:"resolution"`                   // 分辨率 如720p、1080p
//		FileSize       int64         `json:"fileSize"`                     // 文件大小
//		Format         string        `json:"format"`                       // 文件格式 如MP4、AVI
//		Quality        int           `json:"quality"`                      // 视频质量 1-高清 2-标清
//		ApprovalStatus int           `json:"approval_status"`              // 审核状态 1-审核中 2-审核通过 3-审核不通过
//	}
type ContentFindReq struct {
	ID       int64  `json:"id"` // 内容标题
	Author   string `json:"author"`
	Title    string `json:"title"`
	Page     int32  `json:"page"`      //页
	PageSize int32  `json:"page_size"` //每页数量
}

//type ContentFindResp struct {
//	Message  string    `json:"message"`
//	Contents []Content `json:"contents"`
//	Total    int64     `json:"total"`
//}

func (c *CmsApp) ContentFind(ctx *gin.Context) {
	var req ContentFindReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	rsp, err := c.operateAppClient.FindContent(ctx, &operate.FindContentReq{
		Id:       req.ID,
		Author:   req.Author,
		Title:    req.Title,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"operate error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
		"data": rsp,
	})
}

//contentDao := dao.NewContentDao(c.db)
//contentList, total, err := contentDao.Find(dao.FindParams{
//	ID:       req.ID,
//	Author:   req.Author,
//	Title:    req.Title,
//	Page:     req.Page,
//	PageSize: req.PageSize,
//})
//if err != nil {
//	ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//}
//contents := make([]Content, 0, len(contentList))
//for _, content := range contentList {
//	contents = append(contents, Content{
//		ID:             content.ID,
//		Title:          content.Title,
//		Description:    content.Description,
//		Author:         content.Author,
//		VideoURL:       content.VideoURL,
//		Thumbnail:      content.Thumbnail,
//		Category:       content.Category,
//		Duration:       content.Duration,
//		Resolution:     content.Resolution,
//		FileSize:       content.FileSize,
//		Format:         content.Format,
//		Quality:        content.Quality,
//		ApprovalStatus: content.ApprovalStatus,
//	})
//}
//	ctx.JSON(http.StatusOK, gin.H{
//		"code": 0,
//		"msg":  "ok",
//		"data": &ContentFindResp{
//			Message:  fmt.Sprintf("ok"),
//			Contents: contents,
//			Total:    total,
//		},
//	})
