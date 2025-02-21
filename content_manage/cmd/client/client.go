package main

import (
	"content_manage/api/operate"
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

func main() {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("127.0.0.1:9000"),
		grpc.WithMiddleware(
			recovery.Recovery(),
		),
	)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := operate.NewAppClient(conn)
	//reply, err := client.CreateContent(context.Background(), &operate.CreateContentReq{
	//	Content: &operate.Content{
	//		Title:       "title",
	//		VideoUrl:    "http",
	//		Author:      "lucky",
	//		Description: "test",
	//	},
	//})
	//reply, err := client.UpdateContent(context.Background(), &operate.UpdateContentReq{
	//	Content: &operate.Content{
	//		Id:          13,
	//		Title:       "title",
	//		VideoUrl:    "http",
	//		Author:      "lucky",
	//		Description: "test update",
	//	},
	//})

	//reply, err := client.DeleteContent(context.Background(), &operate.DeleteContentReq{
	//	Id: 3,
	//})
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println(reply)

	reply, err := client.FindContent(context.Background(), &operate.FindContentReq{
		Id:       0,
		Author:   "",
		Title:    "",
		Page:     1,
		PageSize: 5,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("[grpc] find : ", reply)
}
