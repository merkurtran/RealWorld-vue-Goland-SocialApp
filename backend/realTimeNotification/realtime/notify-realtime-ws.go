package realtime

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/websocket/v2"
)

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

func StartWebSocketServer(ws map[string]*websocket.Conn, wsMu *sync.Mutex) {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowOriginsFunc: func(origin string) bool {
			return true
		},
	}))

	app.Get("/ws/:userId", websocket.New(func(c *websocket.Conn) {
		userId := c.Params("userId")
		wsMu.Lock()
		ws[userId] = c
		wsMu.Unlock()

		defer func() {
			fmt.Printf("user %s Disconnected\n", userId)
			wsMu.Lock()
			delete(ws, userId)
			wsMu.Unlock()

			c.Close()
		}()

		for {
			var notificationData Notification
			err := c.ReadJSON(&notificationData)
			if err != nil {
				fmt.Printf("user %s Disconnected\n", userId)
				break
			}
			c.WriteJSON(notificationData)
		}
	}))

	log.Fatal(app.Listen(":8088"))
}
