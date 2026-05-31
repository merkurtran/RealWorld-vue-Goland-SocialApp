package main

import (
	"Server/database"
	_ "Server/docs"
	pb "Server/protos"
	"Server/routes"
	"Server/servergrpc"
	"log"
	"net"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// @title Fiber Goland Mongo Grpc Websocet etc..
// @version 1.0
// @description This is Swagger docs for rest api goland fiber
// @host localhost:5005
// @BasePath /
// @schemes http
// @securityDefinitions.apiKey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and the token

func main() {
	// load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	database.Connect()
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOriginsFunc: func(origin string) bool {
			return true
		},
	}))

	// Setup Grpc Server
	lis, err := net.Listen("tcp", ":5001")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterRealtimeChatServiceServer(grpcServer, &servergrpc.Server{})
	reflection.Register(grpcServer)
	log.Println("grpc server running on port 5001")
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
	// end of setup grpc server

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello word")
	})

	//setup routes
	routes.SetupAuthRoutes(app)
	routes.SetupUserRoutes(app)
	routes.SetupPostRoutes(app)
	routes.SetupChatRoutes(app)
	routes.SetupNotificationRoutes(app)

	// Server swagger doctionation
	app.Get("/swagger/*", swagger.HandlerDefault)

	app.Listen(":5005")
}
