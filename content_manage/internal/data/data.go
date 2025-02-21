package data

import (
	"content_manage/internal/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewcontentRepo)

// Data .
type Data struct {
	// TODO wrapped database client
	db *gorm.DB
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resourses")
	}
	mysqlDB, err := gorm.Open(mysql.Open(c.GetDatabase().GetSource()))
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
	return &Data{
		db: mysqlDB,
	}, cleanup, nil
}
