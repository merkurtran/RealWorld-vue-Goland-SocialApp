package main

import (
	"log"
	"realTimeNotification/realtime"
	"realTimeNotification/servergrpc"
	"sync"

	"github.com/gofiber/websocket/v2"
)

func main() {

	wsMu := sync.Mutex{}
	ws := make(map[string]*websocket.Conn)

	// call grpc server
	if err := servergrpc.StartGrpcServer(ws, &wsMu); err != nil {
		log.Fatalf("failed to start grpc server: %v", err)
	}

	go realtime.StartWebSocketServer(ws, &wsMu)
	// block main goroutine to keep program running
	select {}
}
