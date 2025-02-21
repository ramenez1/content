package model

import "time"

type Account struct {
	ID         int64     `gorm:"column:id;primary_key"`
	UserID     string    `gorm:"column:user_id"`
	Password   string    `gorm:"column:password"`
	Nickname   string    `gorm:"column:nickname"`
	Creattime  time.Time `gorm:"column:created_at"`
	Updatetime time.Time `gorm:"column:updated_at"`
}

// TableName 自定义表名，重写方法
func (a Account) TableName() string {
	table := "cms_account.account"
	return table
}
