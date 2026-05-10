package routes

import (
	"Server/controllers"

	"github.com/gofiber/fiber/v2"
)

func SetupNotificationRoutes(app *fiber.App) {
	app.Get("/notification/mark-notification-asreaded", controllers.MarkNotAsReaded)
	app.Get("/notification/:userid", controllers.GetUserNotification)
}
