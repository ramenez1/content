package dao

import (
	"Content/content_flow/internal/model"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"hash/fnv"
	"log"
	"math/big"
)

type ContentDao struct {
	db *gorm.DB
}

const contentNumTables = 4

func NewContentDao(db *gorm.DB) *ContentDao {
	return &ContentDao{db: db}
}

func getContentDetailsTable(contentID string) string {
	tableIndex := getContentTableIndex(contentID)
	table := fmt.Sprintf("cms_content.t_content_details_%d", tableIndex)
	log.Printf("content_id = %s, table = %s", contentID, table)
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

func (c *ContentDao) First(contentID string) (model.ContentDetail, error) {
	var detail model.ContentDetail
	err := c.db.Table(getContentDetailsTable(contentID)).Where("content_id = ?", contentID).First(&detail).Error
	if err != nil {
		log.Printf("contentDao First error = %v", err)
		return detail, err
	}
	return detail, nil
}

func (c *ContentDao) Create(detail model.ContentDetail) (int, error) {
	if err := c.db.Create(&detail).Error; err != nil {
		log.Printf("contentDao Create error = %v", err)
		return 0, err
	}
	return detail.ID, nil
}

// 校验存在
func (c *ContentDao) IsExist(contentID int) (bool, error) {
	var detail model.ContentDetail
	err := c.db.Where("id = ?", contentID).First(&detail).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		log.Printf("content error isExist = %v", err)
		return false, err
	}
	return true, nil
}

// 更新
func (c *ContentDao) Update(id int, detail model.ContentDetail) error {
	if err := c.db.Where("id = ?", id).
		Updates(&detail).Error; err != nil {
		log.Printf("contentDao Update error = %v", err)
		return err
	}
	return nil
}

// 临时使用
func (c *ContentDao) UpdateByID(contentID string, column string, value interface{}) error {
	if err := c.db.Table(getContentDetailsTable(contentID)).Where("content_id = ?", contentID).Update(column, value).Error; err != nil {
		log.Printf("contentDao UpdateByID error = %v", err)
		return err
	}
	return nil
}

// 删除
func (c *ContentDao) Delete(id int) error {
	if err := c.db.Where("id = ?", id).Delete(&model.ContentDetail{}).Error; err != nil {
		log.Printf("contentDao Delete error = %v", err)
		return err
	}
	return nil
}

// 分页查询（id查询，author,title）
type FindParams struct {
	ID       int
	Author   string
	Title    string
	Page     int
	PageSize int
}

func (c *ContentDao) Find(params FindParams) ([]model.ContentDetail, int64, error) {
	//构建查询条件
	query := c.db.Model(&model.ContentDetail{})
	if params.ID != 0 {
		query = query.Where("id = ?", params.ID)
	}
	if params.Author != "" {
		query = query.Where("author = ?", params.Author)
	}
	if params.Title != "" {
		query = query.Where("title = ?", params.Title)
	}
	//计算总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		log.Printf("contentDao Find total error = %v", err)
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
	var data []model.ContentDetail
	if err := query.Offset(offset).
		Limit(pageSize).
		Find(&data).Error; err != nil {
		log.Printf("contentDao Find error = %v", err)
		return nil, 0, err
	}
	return data, total, nil

}
