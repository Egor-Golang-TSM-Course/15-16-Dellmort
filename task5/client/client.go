package client

import (
	"context"
	"fmt"
	"lesson15_16/task5/api"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func StartClient() {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	conn, err := grpc.Dial(
		":8080",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	c := api.NewStorageClient(conn)
	resp, err := c.Save(ctx, &api.TaskRequest{
		Message: "hello world",
	})
	if err != nil {
		fmt.Println("Save err: ", err)
		return
	}
	fmt.Println("save:\nnew task id =", resp.GetId())

	getResp, err := c.Get(ctx, &api.GetRequest{
		Id: 55,
	})
	if err != nil {
		fmt.Println("Get err: ", err)
	}
	if getResp.GetId() > 0 {
		fmt.Printf("get: id = %d, message = %s\n", getResp.GetId(), getResp.GetMessage())
	}
}
