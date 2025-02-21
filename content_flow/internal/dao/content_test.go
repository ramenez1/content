package dao

import (
	"ContentSystem/internal/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func connDB() *gorm.DB {
	mysqlDB, err := gorm.Open(mysql.Open("root:root@tcp(127.0.0.1:3306)/?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
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
	return mysqlDB
}

func TestContentDao_Create(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		detail model.ContentDetail
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "内容插入",
			fields: fields{
				db: connDB(),
			},
			args: args{
				detail: model.ContentDetail{
					Title: "test1",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ContentDao{
				db: tt.fields.db,
			}
			if err := c.Create(tt.args.detail); (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
