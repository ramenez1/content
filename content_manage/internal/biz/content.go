package biz

import (
	"context"
	"errors"
	"github.com/go-kratos/kratos/v2/log"
	"golang.org/x/sync/errgroup"
	"time"
)

// Content is a Content model.
type Content struct {
	ID             int           `json:"id"`                           // 内容ID
	Title          string        `json:"title" binding:"required"`     // 内容标题
	ContentID      string        `gorm:"column:content_id"`            // 内容标题
	Description    string        `json:"description"`                  // 内容描述
	Author         string        `json:"author" binding:"required"`    // 作者
	VideoURL       string        `json:"video_url" binding:"required"` // 视频播放URL
	Thumbnail      string        `json:"thumbnail"`                    // 封面图URL
	Category       string        `json:"category"`                     // 内容分类
	Duration       time.Duration `json:"duration"`                     // 内容时长
	Resolution     string        `json:"resolution"`                   // 分辨率 如720p、1080p
	FileSize       int64         `json:"fileSize"`                     // 文件大小
	Format         string        `json:"format"`                       // 文件格式 如MP4、AVI
	Quality        int           `json:"quality"`                      // 视频质量 1-高清 2-标清
	ApprovalStatus int           `json:"approval_status"`              // 审核状态 1-审核中 2-审核通过 3-审核不通过
	UpdatedAt      time.Time     `json:"updated_at"`                   // 内容更新时间
	CreatedAt      time.Time     `json:"created_at"`                   // 内容创建时间
}

type FindParams struct {
	ID       int    `json:"id"`
	Author   string `json:"author"`
	Title    string `json:"title"`
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
}

type ContextIndex struct {
	ID        int64  `json:"id"`         // 自增ID
	ContentID string `json:"content_id"` // 内容ID
}

// ContentRepo is a Content repo.
type ContentRepo interface {
	Create(ctx context.Context, c *Content) (int, error)
	Update(ctx context.Context, id int, c *Content) error
	IsExist(ctx context.Context, contentID int) (bool, error)
	Delete(ctx context.Context, id int) error
	Find(ctx context.Context, params *FindParams) ([]*Content, int, error)
	FindIndex(ctx context.Context, params *FindParams) ([]*ContextIndex, int, error)
	First(ctx context.Context, idx *ContextIndex) (*Content, error)
}

// ContentUsecase is a Content usecase.
type ContentUsecase struct {
	repo ContentRepo
	log  *log.Helper
}

// NewContentUsecase new a Content usecase.
func NewContentUsecase(repo ContentRepo, logger log.Logger) *ContentUsecase {
	return &ContentUsecase{repo: repo, log: log.NewHelper(logger)}
}

// CreateContent creates a Content, and returns the new Content.
func (uc *ContentUsecase) CreateContent(ctx context.Context, c *Content) (int, error) {
	uc.log.WithContext(ctx).Infof("CreateContent: %v+v", c)
	return uc.repo.Create(ctx, c)
}

// UpdateContent updates a Content, and returns the updated Content.
func (uc *ContentUsecase) UpdateContent(ctx context.Context, c *Content) error {
	uc.log.WithContext(ctx).Infof("UpdateContent: %v", c)
	return uc.repo.Update(ctx, c.ID, c)
}

// DeleteContent deletes a Content, and returns the deleted Content.
func (uc *ContentUsecase) DeleteContent(ctx context.Context, id int) error {
	ok, err := uc.repo.IsExist(ctx, id)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("content not exist")
	}
	return uc.repo.Delete(ctx, id)
}

// FindContent find Content.
func (uc *ContentUsecase) FindContent(ctx context.Context, params *FindParams) ([]*Content, int, error) {
	repo := uc.repo
	//contents, total, err := repo.Find(ctx, params)
	//if err != nil {
	//	return nil, 0, err
	//}
	indices, total, err := repo.FindIndex(ctx, params)
	if err != nil {
		return nil, 0, err
	}
	var eg errgroup.Group
	contents := make([]*Content, len(indices))
	for index, idx := range indices {
		tempIndex := index
		tempIdx := idx
		eg.Go(func() error {
			content, err := repo.First(ctx, tempIdx)
			if err != nil {
				return err
			}
			contents[tempIndex] = &Content{
				ID:             content.ID,
				Title:          content.Title,
				ContentID:      content.ContentID,
				Description:    content.Description,
				Author:         content.Author,
				VideoURL:       content.VideoURL,
				Thumbnail:      content.Thumbnail,
				Category:       content.Category,
				Duration:       content.Duration,
				Resolution:     content.Resolution,
				FileSize:       content.FileSize,
				Format:         content.Format,
				Quality:        content.Quality,
				ApprovalStatus: content.ApprovalStatus,
				UpdatedAt:      content.UpdatedAt,
				CreatedAt:      content.CreatedAt,
			}
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, 0, err
	}
	return contents, total, err
}
