package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

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
	return "account"
}

func main() {
	db := connDB()
	var accounts []Account
	//默认读表为accounts(由上面的自动加s生成)
	if err := db.Find(&accounts).Error; err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(accounts)

}

func connDB() *gorm.DB {
	mysqlDB, err := gorm.Open(mysql.Open("root:root@tcp(127.0.0.1:3306)/cms_account?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	//连接池
	db, err := mysqlDB.DB()
	if err != nil {
		panic(err)
	}
	//设置最大连接数
	db.SetMaxOpenConns(4)
	//设置最大空闲连接数
	db.SetMaxIdleConns(2)
	//设置为debug模式
	mysqlDB = mysqlDB.Debug()
	return mysqlDB
}
