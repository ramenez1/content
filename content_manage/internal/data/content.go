package data

import (
	"content_manage/internal/biz"
	"context"
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
	"hash/fnv"
	"math/big"
	"time"
)

const contentNumTables = 4

type contentRepo struct {
	data *Data
	log  *log.Helper
}

// NewcontentRepo .
func NewcontentRepo(data *Data, logger log.Logger) biz.ContentRepo {
	return &contentRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

type ContentDetail struct {
	ID             int           `gorm:"column:id;primary_key"`  // 内容ID
	Title          string        `gorm:"column:title"`           // 内容标题
	ContentID      string        `gorm:"column:content_id"`      // 内容标题
	Description    string        `gorm:"column:description"`     // 内容描述
	Author         string        `gorm:"column:author"`          // 作者
	VideoURL       string        `gorm:"column:video_url"`       // 视频播放URL
	Thumbnail      string        `gorm:"column:thumbnail"`       // 封面图URL
	Category       string        `gorm:"column:category"`        // 内容分类
	Duration       time.Duration `gorm:"column:duration"`        // 内容时长
	Resolution     string        `gorm:"column:resolution"`      // 分辨率 如720p、1080p
	FileSize       int64         `gorm:"column:fileSize"`        // 文件大小
	Format         string        `gorm:"column:format"`          // 文件格式 如MP4、AVI
	Quality        int           `gorm:"column:quality"`         // 视频质量 1-高清 2-标清
	ApprovalStatus int           `gorm:"column:approval_status"` // 审核状态 1-审核中 2-审核通过 3-审核不通过
	UpdatedAt      time.Time     `gorm:"column:updated_at"`      // 内容更新时间
	CreatedAt      time.Time     `gorm:"column:created_at"`      // 内容创建时间
}

type IdxContentDetail struct {
	ID        int       `gorm:"column:id;primary_key"` // 内容ID
	Title     string    `gorm:"column:title"`          // 内容标题
	ContentID string    `gorm:"column:content_id"`     // 内容标题
	Author    string    `gorm:"column:author"`         // 作者
	UpdatedAt time.Time `gorm:"column:updated_at"`     // 内容更新时间
	CreatedAt time.Time `gorm:"column:created_at"`     // 内容创建时间
}

//func (c ContentDetail) TableName() string {
//	return "cms_content.t_content_details"
//}

func getContentDetailsTable(contentID string) string {
	tableIndex := getContentTableIndex(contentID)
	table := fmt.Sprintf("cms_content.t_content_details_%d", tableIndex)
	log.Infof("content_id = %s, table = %s", contentID, table)
	return table
}

func getContentTableIndex(uuid string) int {
	hash := fnv.New64()
	_, _ = hash.Write([]byte(uuid))
	hashValue := hash.Sum64()
	fmt.Println("hashValue = ", hashValue)

	bigNum := big.NewInt(int64(hashValue))
	mod := big.NewInt(contentNumTables)
	tableIndex := bigNum.Mod(bigNum, mod).Int64()
	return int(tableIndex)
}

func (c *contentRepo) Create(ctx context.Context, content *biz.Content) (int, error) {
	c.log.Infof("contentRepo Create content: %v", content)
	db := c.data.db
	idx := IdxContentDetail{
		Title:     content.Title,
		ContentID: content.ContentID,
		Author:    content.Author,
	}
	if err := db.Table("cms_content.t_idx_content_details").Create(&idx).Error; err != nil {
		return 0, err
	}
	detail := ContentDetail{
		ContentID:      content.ContentID,
		Title:          content.Title,
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
	}

	if err := db.Table(getContentDetailsTable(content.ContentID)).Create(&detail).Error; err != nil {
		c.log.Errorf("create content error: %v", err)
		return 0, err
	}
	return idx.ID, nil
}

func (c *contentRepo) Update(ctx context.Context, id int, content *biz.Content) error {
	c.log.Infof("contentRepo Update content: %v", content)
	detail := ContentDetail{
		//ID:             content.ID,
		Title:          content.Title,
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
	}
	var idx IdxContentDetail
	db := c.data.db
	if err := db.Table("cms_content.t_idx_content_details").Where("id = ?", id).First(&idx).Error; err != nil {
		c.log.Errorf("get content_idx  error: %v", err)
		return err
	}
	if err := db.Table(getContentDetailsTable(idx.ContentID)).Where("content_id = ?", idx.ContentID).Updates(&detail).Error; err != nil {
		c.log.Errorf("update content error: %v", err)
		return err
	}
	return nil
}

func (c *contentRepo) IsExist(ctx context.Context, id int) (bool, error) {
	db := c.data.db
	var detail IdxContentDetail
	err := db.Table("cms_content.t_idx_content_details").Where("id = ?", id).First(&detail).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		c.log.WithContext(ctx).Errorf("ContentDao isexist error: %v", err)
		return false, err
	}
	return true, nil
}

func (c *contentRepo) Delete(ctx context.Context, id int) error {
	db := c.data.db
	//查询索引
	var idx IdxContentDetail
	if err := db.Table("cms_content.t_idx_content_details").Where("id = ?", id).First(&idx).Error; err != nil {
		c.log.Errorf("get content_idx  error: %v", err)
		return err
	}
	//删除索引
	err := db.Table("cms_content.t_idx_content_details").Where("id = ?", id).Delete(&IdxContentDetail{}).Error
	if err != nil {
		return err
	}
	//删除详情
	err = db.Table(getContentDetailsTable(idx.ContentID)).Where("content_id = ?", idx.ContentID).Delete(&ContentDetail{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *contentRepo) Find(ctx context.Context, params *biz.FindParams) ([]*biz.Content, int, error) {
	//构造条件
	query := c.data.db.Model(&ContentDetail{})
	if params.ID != 0 {
		query = query.Where("id = ?", params.ID)
	}
	if params.Author != "" {
		query = query.Where("author = ?", params.Author)
	}
	if params.Title != "" {
		query = query.Where("title = ?", params.Title)
	}
	//总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var page, pageSize = 1, 10
	if params.Page > 0 {
		page = params.Page
	}
	if params.PageSize > 0 {
		pageSize = params.PageSize
	}
	offset := (page - 1) * pageSize
	var result []*ContentDetail
	if err := query.Offset(offset).Limit(pageSize).
		Find(&result).Error; err != nil {
		c.log.WithContext(ctx).Errorf("find content error: %v", err)
		return nil, 0, err
	}
	var contents []*biz.Content
	for _, r := range result {
		contents = append(contents, &biz.Content{
			ID:             r.ID,
			Title:          r.Title,
			Description:    r.Description,
			Author:         r.Author,
			VideoURL:       r.VideoURL,
			Thumbnail:      r.Thumbnail,
			Category:       r.Category,
			Duration:       r.Duration,
			Resolution:     r.Resolution,
			FileSize:       r.FileSize,
			Format:         r.Format,
			Quality:        r.Quality,
			ApprovalStatus: r.ApprovalStatus,
			UpdatedAt:      r.UpdatedAt,
			CreatedAt:      r.CreatedAt,
		})
	}
	return contents, int(total), nil
}

// FindIndex 查询索引
func (c *contentRepo) FindIndex(ctx context.Context, params *biz.FindParams) ([]*biz.ContextIndex, int, error) {
	query := c.data.db.Model(&IdxContentDetail{})
	if params.ID != 0 {
		query = query.Where("id = ?", params.ID)
	}
	if params.Author != "" {
		query = query.Where("author = ?", params.Author)
	}
	if params.Title != "" {
		query = query.Where("title = ?", params.Title)
	}

	var total int64
	if err := query.Table("cms_content.t_idx_content_details").Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var page, pageSize = 1, 10
	if params.Page > 0 {
		page = params.Page
	}
	if params.PageSize > 0 {
		pageSize = params.PageSize
	}
	offset := (page - 1) * pageSize
	var result []*IdxContentDetail
	if err := query.Offset(offset).Limit(pageSize).
		Find(&result).Error; err != nil {
		c.log.WithContext(ctx).Errorf("find content error: %v", err)
		return nil, 0, err
	}
	var contents []*biz.ContextIndex
	for _, r := range result {
		contents = append(contents, &biz.ContextIndex{
			ID:        int64(r.ID),
			ContentID: r.ContentID,
		})
	}
	return contents, int(total), nil
}

func (c *contentRepo) First(ctx context.Context, idx *biz.ContextIndex) (*biz.Content, error) {
	db := c.data.db
	var detail ContentDetail
	if err := db.Table(getContentDetailsTable(idx.ContentID)).
		Where("content_id = ?", idx.ContentID).First(&detail).Error; err != nil {
		c.log.WithContext(ctx).Errorf("content first error = %v", err)
		return nil, err
	}
	content := &biz.Content{
		ID:             detail.ID,
		Title:          detail.Title,
		Description:    detail.Description,
		Author:         detail.Author,
		VideoURL:       detail.VideoURL,
		Thumbnail:      detail.Thumbnail,
		Category:       detail.Category,
		Duration:       detail.Duration,
		Resolution:     detail.Resolution,
		FileSize:       detail.FileSize,
		Format:         detail.Format,
		Quality:        detail.Quality,
		ApprovalStatus: detail.ApprovalStatus,
		UpdatedAt:      detail.UpdatedAt,
		CreatedAt:      detail.CreatedAt,
	}
	return content, nil
}
