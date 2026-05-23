package controllers

import (
	"Server/database"
	"Server/models"
	"context"
	"slices"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetUserBy ID
// @Summary Get User By ID
// @Description GetUser Details By ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} models.UserModel
// @Failure 400 {object} map[string]interface{}
// @Router /user/getUser/{id} [get]
func GetUserByID(c *fiber.Ctx) error {
	var UserSchema = database.DB.Collection("users")
	var PostSchema = database.DB.Collection("posts")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var user models.UserModel
	var posts []models.PostModel
	objId, _ := primitive.ObjectIDFromHex(c.Params("id"))
	strID := c.Params("id")
	// TODO GET and Return user posts
	findOptions := options.Find()
	postResult, err := PostSchema.Find(ctx, bson.M{"creator": strID}, findOptions)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err,
		})
	}

	defer postResult.Close(ctx)
	for postResult.Next(ctx) {
		var singlePost models.PostModel
		postResult.Decode(&singlePost)
		posts = append(posts, singlePost)
	}

	// get user data
	userResult := UserSchema.FindOne(ctx, bson.M{"_id": objId})

	if userResult.Err() != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "user not found",
		})
	}

	userResult.Decode(&user)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user":  user,
		"posts": posts,
	})
}

// UpdateUser
// @Summary update user data
// @Description update user details
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body models.UpdateUser true "details"
// @Success 200 {object} models.UserModel
// @Failure 400 {object} map[string]interface{}
// @security BearerAuth
// @Router /user/Update/{id} [patch]
func UpdateUser(c *fiber.Ctx) error {
	var UserSchema = database.DB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	//
	extUid := c.Locals("userID").(string)

	if extUid != c.Params("id") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "You are not authorized to update this profile",
		})
	}
	userid, _ := primitive.ObjectIDFromHex(c.Params("id"))
	var user models.UpdateUser
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	update := bson.M{"name": user.Name, "imageUrl": user.ImageUrl, "bio": user.Bio}
	result, err := UserSchema.UpdateOne(ctx, bson.M{"_id": userid}, bson.M{"$set": update})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "cannont update the user data",
			"details": err.Error(),
		})
	}
	//
	var updateUser models.UserModel
	if result.MatchedCount == 1 {
		err := UserSchema.FindOne(ctx, bson.M{"_id": userid}).Decode(&updateUser)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"details": err.Error(),
			})
		}
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": updateUser})
}

// Following Users
// @Summary Follow/UnFoloow User
// @Description follow or unfollow a user
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @security BearerAuth
// @Router /user/{id}/following [patch]
func FollowingUser(c *fiber.Ctx) error {
	var UserSchema = database.DB.Collection("users")
	var NotificationSchema = database.DB.Collection("notifications")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var FirstUser models.UserModel
	var SecondUser models.UserModel

	FirstUserID, _ := primitive.ObjectIDFromHex(c.Params("id"))
	SecondUserID, _ := primitive.ObjectIDFromHex(c.Locals("userID").(string))

	err := UserSchema.FindOne(ctx, bson.M{"_id": FirstUserID}).Decode(&FirstUser)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"details": err.Error(),
		})
	}

	err = UserSchema.FindOne(ctx, bson.M{"_id": SecondUserID}).Decode(&SecondUser)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"details": err.Error(),
		})
	}

	fuid := c.Params("id")
	suid := c.Locals("userID").(string)

	if slices.Contains(FirstUser.Followers, suid) {
		i := slices.Index(FirstUser.Followers, suid)
		FirstUser.Followers = slices.Delete(FirstUser.Followers, i, i+1)
		// remove form the following list on second user
		index := slices.Index(SecondUser.Following, fuid)
		SecondUser.Following = slices.Delete(SecondUser.Following, index, index+1)
	} else {
		FirstUser.Followers = append(FirstUser.Followers, suid)
		SecondUser.Following = append(SecondUser.Following, fuid)

		// Create Notification
		notification := models.Notification{
			MainUID:   FirstUser.ID.Hex(),
			TargetID:  SecondUser.ID.Hex(),
			Details:   SecondUser.Name + "Start Following You!",
			User:      models.User{Name: SecondUser.Name, Avata: SecondUser.ImageUrl},
			CreatedAt: time.Now(),
		}
		_, err := NotificationSchema.InsertOne(ctx, notification)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Faild to create notification",
				"error":   err.Error(),
			})
		}
	}

	updateFirst := bson.M{"followers": FirstUser.Followers}
	updateSecond := bson.M{"following": SecondUser.Following}

	_, err = UserSchema.UpdateOne(ctx, bson.M{"_id": FirstUserID}, bson.M{"$set": updateFirst})
	_, err = UserSchema.UpdateOne(ctx, bson.M{"_id": SecondUserID}, bson.M{"$set": updateSecond})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"details": err.Error(),
		})
	}

	err = UserSchema.FindOne(ctx, bson.M{"_id": FirstUserID}).Decode(&FirstUser)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"details": err.Error(),
		})
	}

	err = UserSchema.FindOne(ctx, bson.M{"_id": SecondUserID}).Decode(&SecondUser)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"SecondUser": SecondUser,
		"FirstUser":  FirstUser,
	})

}

// GetSugUser Users
// @Summary Get Suggersted users
// @Description get suggested userses based on the current user's following list
// @Tags Users
// @Accept json
// @Produce json
// @Param id query string true "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @security BearerAuth
// @Router /user/getSug [get]
func GetSugUser(c *fiber.Ctx) error {
	var UserSchema = database.DB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var MainUser models.UserModel
	var SugListId []string
	var AllSugUsers []models.UserModel

	MainUserID, err := primitive.ObjectIDFromHex(c.Query("id"))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"details": err.Error(),
		})
	}

	err = UserSchema.FindOne(ctx, bson.M{"_id": MainUserID}).Decode(&MainUser)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"details": err.Error(),
		})
	}

	// Get SugUsers id then put them in suglistid
	for _, FID := range MainUser.Following {
		var singleUser models.UserModel
		convFID, _ := primitive.ObjectIDFromHex(FID)
		err = UserSchema.FindOne(ctx, bson.M{"_id": convFID}).Decode(&singleUser)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"details": err.Error(),
			})
		}
		// following
		for _, id := range singleUser.Following {
			if !slices.Contains(SugListId, id) && id != c.Query("id") {
				SugListId = append(SugListId, id)
			}
		}

		// followers
		for _, id := range singleUser.Followers {
			if !slices.Contains(SugListId, id) && id != c.Query("id") {
				SugListId = append(SugListId, id)
			}
		}
	}

	// Gest Sug Users by id
	if len(SugListId) > 0 {
		var objectides []primitive.ObjectID
		for _, id := range SugListId {
			objid, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				continue
			}
			objectides = append(objectides, objid)
		}

		// fetch all users in one query using $in operator
		cursor, err := UserSchema.Find(ctx, bson.M{
			"_id": bson.M{"$in": objectides},
		})
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"details": err.Error(),
			})
		}

		defer cursor.Close(ctx)

		if err = cursor.All(ctx, &AllSugUsers); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"details": err.Error(),
			})
		}
	}

	if AllSugUsers == nil {
		AllSugUsers = make([]models.UserModel, 0)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"users": AllSugUsers})

}

// DeleteUser
// @Summary delete user data
// @Description delete user details
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @security BearerAuth
// @Router /user/delete/{id} [delete]
func DeleteUser(c *fiber.Ctx) error {
	var UserSchema = database.DB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	//
	extUid := c.Locals("userID").(string)

	if extUid != c.Params("id") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "You are not authorized to delete this user",
		})
	}

	userID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid user id",
		})
	}

	result, err := UserSchema.DeleteOne(ctx, bson.M{"_id": userID})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Faild to delete user",
			"error":   err.Error(),
		})
	}

	if result.DeletedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "user not found",
		})
	}
	// success
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "User Delete Successfully!",
	})
}
