package servergrpc

import (
	"context"
	"log"
	"realTimeChat/protos"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func GetFollowingFollowersClient(id string) ([]*protos.UserIDsList, error) {
	// conn, err := grpc.NewClient(":5001", grpc.WithTransportCredentials(insecure.NewCredentials()))

	// DevOps docker Compose usage
	conn, err := grpc.NewClient("GolangApiServer:5001", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Printf("did not connect: %v", err)
		return nil, err
	}
	defer conn.Close()
	client := protos.NewRealtimeChatServiceClient(conn)

	// call grpc method
	ctx := context.Background()
	req := &protos.UserID{Userid: id}
	resp, err := client.GetUserFollowingFollowers(ctx, req)
	if err != nil {
		log.Printf("could not get user following followers: %v", err)
		return nil, err
	}
	return resp.GetUserIDsLists(), nil
}

func SendMessageClient(sender, receiver, content string) error {
	// conn, err := grpc.NewClient(":5001", grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.NewClient("GolangApiServer:5001", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Printf("did not connect: %v", err)
		return err
	}
	defer conn.Close()
	client := protos.NewRealtimeChatServiceClient(conn)
	ctx := context.Background()
	req := &protos.MessageRequest{
		Sender:   sender,
		Receiver: receiver,
		Content:  content,
	}
	_, err = client.SendMessage(ctx, req)
	if err != nil {
		log.Printf("could not send message: %v", err)
		return err
	}
	return nil
}
