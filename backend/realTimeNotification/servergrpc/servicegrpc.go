package servergrpc

import (
	"context"
	"fmt"
	"log"
	"net"
	pb "realTimeNotification/protos"
	"sync"
	"time"

	"github.com/gofiber/websocket/v2"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
)

type notificationServer struct {
	pb.UnimplementedNotificationServer
	wsMu *sync.Mutex
	ws   map[string]*websocket.Conn
}

type Notification struct {
	ID        string    `json:"id"`
	Details   string    `json:"details"`
	MainUID   string    `json:"mainuid"`
	TargetID  string    `json:"targetid"`
	IsReaded  bool      `json:"isreaded"`
	CreatedAt time.Time `json:"createdAt"`
	User      User      `json:"user"`
}

type User struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

func (s *notificationServer) SendGrpcNotification(ctx context.Context, req *pb.NotificationGrpcRequest) (*empty.Empty, error) {
	fmt.Printf("Sending notification to user %s : %s\n", req.Mainuid, req.Details)

	// send the notification to ws
	s.wsMu.Lock()
	defer s.wsMu.Unlock()
	if conn, ok := s.ws[req.Mainuid]; ok {
		notification := Notification{
			ID:        req.XId,
			MainUID:   req.Mainuid,
			TargetID:  req.Targetid,
			Details:   req.Details,
			IsReaded:  req.Isreaded,
			CreatedAt: time.Unix(req.CreatedAt.GetSeconds(), 0),
			User: User{
				Name:   req.User.Name,
				Avatar: req.User.Avatar,
			},
		}
		err := conn.WriteJSON(notification)
		if err != nil {
			return nil, fmt.Errorf("failed to write notification to ws: %v", err)
		}

	}
	return &empty.Empty{}, nil
}

func StartGrpcServer(ws map[string]*websocket.Conn, wsMu *sync.Mutex) error {
	lis, err := net.Listen("tcp", ":8090")
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	notificationService := &notificationServer{
		ws:   ws,
		wsMu: wsMu,
	}
	pb.RegisterNotificationServer(grpcServer, notificationService)
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
	return nil
}
