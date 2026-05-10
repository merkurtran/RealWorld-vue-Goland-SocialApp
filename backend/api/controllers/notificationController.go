package controllers

import (
	"Server/database"
	"Server/models"
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MarkNotAsReaded
// @Summary Mark notification asreaded for a user
// @Description MarkNotAsReaded
// @Tags Notification
// @Accept json
// @Produce json
// @Param id query string true "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /notification/mark-notification-asreaded [get]
func MarkNotAsReaded(c *fiber.Ctx) error {

	// parse query paramter
	id := c.Query("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "id in query is required",
		})
	}
	// construct the filter and update
	filter := bson.M{"mainuid": bson.M{"$regex": id, "$options": "i"}}
	update := bson.M{"$set": bson.M{"isreaded": true}}

	// update
	var NotificationSchema = database.DB.Collection("notifications")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := NotificationSchema.UpdateMany(ctx, filter, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Faild to mark notification as read",
			"error":   err.Error(),
		})
	}

	// retreive the update notification
	options := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})
	cursor, err := NotificationSchema.Find(ctx, filter, options)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Faild to retrieve the update notifications",
			"error":   err.Error(),
		})
	}

	defer cursor.Close(ctx)
	var notifications []models.Notification
	if err := cursor.All(ctx, &notifications); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Faild to decoded notifications",
			"error":   err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"notifications": notifications,
	})
}

// GetUserNotification Post
// @Summary get user notification
// @Description get user notification
// @Tags Notification
// @Accept json
// @Produce json
// @Param userid path string true "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /notification/{userid} [get]
func GetUserNotification(c *fiber.Ctx) error {
	// parse query paramter
	id := c.Params("userid")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "userid in params is required",
		})
	}

	var NotificationSchema = database.DB.Collection("notifications")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{"mainuid": bson.M{"$regex": id, "$options": "i"}}
	options := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cursor, err := NotificationSchema.Find(ctx, filter, options)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Faild to retrieve the update notifications",
			"error":   err.Error(),
		})
	}

	defer cursor.Close(ctx)

	var notifications []models.Notification
	if err := cursor.All(ctx, &notifications); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Faild to decoded notifications",
			"error":   err.Error(),
		})
	}

	if notifications == nil {
		notifications = []models.Notification{}
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"notifications": notifications,
	})
}
