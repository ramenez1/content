package service

import (
	"ContentSystem/internal/api/operate"
	"context"
	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/redis/go-redis/v9"
	goflow "github.com/s8sg/goflow/v1"
	clientv3 "go.etcd.io/etcd/client/v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type CmsApp struct {
	db               *gorm.DB
	rdb              *redis.Client
	flowService      *goflow.FlowService
	operateAppClient operate.AppClient
}

func NewCmsApp() *CmsApp {
	app := &CmsApp{}
	connDB(app)
	connRdb(app)
	connOperateAppClient(app)
	//app.flowService = flowService()
	//go func() {
	//	process.ExecContentFlow(app.db)
	//}()
	return app
}

func connDB(app *CmsApp) {
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
	//设置为debug模式
	mysqlDB = mysqlDB.Debug()
	app.db = mysqlDB
}

//func flowService() *goflow.FlowService {
//	fs := &goflow.FlowService{
//		RedisURL: "localhost:6379",
//	}
//	return fs
//}

func connRdb(app *CmsApp) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // use default Addr
		Password: "",               // no password set
		DB:       0,                // use default DB
	})
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	app.rdb = rdb
}

func connOperateAppClient(app *CmsApp) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"localhost:2379"},
	})
	if err != nil {
		panic(err)
	}
	dis := etcd.New(client)
	// ///后的为ETCD注册的名  即（etcdctl get --prefix ""）后的name 部分
	endpoint := "discovery:///content_manage"
	conn, err := grpc.DialInsecure(
		context.Background(),
		//grpc.WithEndpoint("127.0.0.1:9000"),
		grpc.WithMiddleware(
			recovery.Recovery(),
		),
		grpc.WithDiscovery(dis),
		grpc.WithEndpoint(endpoint),
	)
	if err != nil {
		panic(err)
	}
	appclient := operate.NewAppClient(conn)
	app.operateAppClient = appclient
}
