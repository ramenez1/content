package main

import (
	"Content/content_flow/internal/process"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
	"os/signal"
)

func main() {
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

	process.ExecContentFlow(mysqlDB)

	//监听退出信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	//等待退出信号
	<-quit
	log.Println("服务退出....")
}
