package controllers

import (
	"Server/database"
	"Server/models"
	"context"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SendMessage
// @Summary send message to friend user
// @Description SendMessage from one user to another
// @Tags Chat
// @Accept json
// @Produce json
// @Param message body models.SendMessageModel true "user SendMessage details"
// @Success 201 {object} models.MessageModel
// @Failure 400 {object} map[string]interface{}
// @Router /chat/sendmessage [post]
func SendMessage(c *fiber.Ctx) error {
	var MessageSchema = database.DB.Collection("message")
	var UnReadedMsgSchema = database.DB.Collection("UnReadedmessage")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var body models.SendMessageModel
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	msg := models.MessageModel{
		Content:  body.Content,
		Sender:   body.Sender,
		Receiver: body.Receiver,
	}

	// save the mssage to db
	result, err := MessageSchema.InsertOne(ctx, &msg)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "faild to save msg",
			"details": err.Error(),
		})
	}

	// update or create the unreaded message count and is readed
	var unReadedMsg models.UnReadModel
	filtter := bson.M{"mainUserid": msg.Receiver, "otherUserid": msg.Sender}
	update := bson.M{"$inc": bson.M{"numOfUnreadedMessages": 1}, "$set": bson.M{"isReaded": false}}
	opts := options.FindOneAndUpdate().SetUpsert(true)
	err = UnReadedMsgSchema.FindOneAndUpdate(ctx, filtter, update, opts).Decode(&unReadedMsg)
	if err != nil && err != mongo.ErrNoDocuments {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Faild to update unreader message count",
			"details": err.Error(),
		})
	}
	// return the created message
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Message sent successfully",
		"result":  result.InsertedID,
	})
}

// GetMesgsByNums
// @Summary get message by pagenation
// @Description GetMesgsByNums two users by pagenation
// @Tags Chat
// @Accept json
// @Produce json
// @Param from query int true "Staring point page num"
// @Param firstuid query string true "first user id"
// @Param seconduid query string true "second user id"
// @Success 200 {object} []models.MessageModel
// @Failure 400 {object} map[string]interface{}
// @Router /chat/getmsgsbynums [get]
func GetMesgsByNums(c *fiber.Ctx) error {
	var MessageSchema = database.DB.Collection("message")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	from, err := strconv.Atoi(c.Query("from"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid value form from",
			"error":   err.Error(),
		})
	}
	firstuid := c.Query("firstuid")
	seconduid := c.Query("seconduid")

	// construct the filer
	senderFilter := bson.M{"sender": firstuid, "receiver": seconduid}
	receiverFilter := bson.M{"sender": seconduid, "receiver": firstuid}
	filter := bson.M{"$or": []bson.M{senderFilter, receiverFilter}}

	// pagenation options
	options := options.Find()
	options.SetSort(bson.D{{Key: "_id", Value: -1}})
	options.SetSkip(int64(from * 2))
	options.SetLimit(int64(2))

	// query the db
	cursor, err := MessageSchema.Find(ctx, filter, options)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Faild to retrivere messages",
			"error":   err.Error(),
		})
	}
	defer cursor.Close(ctx)

	// iterate over the cursor and build the res array
	var messages []models.MessageModel
	for cursor.Next(ctx) {
		var msg models.MessageModel
		err := cursor.Decode(&msg)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Faild to decode messages",
				"error":   err.Error(),
			})
		}
		messages = append(messages, msg)
	}

	// reciver the message array
	for i := 0; i < len(messages)/2; i++ {
		j := len(messages) - 1 - i
		messages[i], messages[j] = messages[j], messages[i]
	}
	if len(messages) == 0 {
		messages = []models.MessageModel{}
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"msgs":   "Message send successfully",
		"result": messages,
	})
}

// GetUserUnreadedMessage
// @Summary get unreaded message count & recodes for user
// @Description get unreaded message count & recodes for user
// @Tags Chat
// @Accept json
// @Produce json
// @Param userid query string true "user id"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /chat/get-user-unreadedmsg [get]
func GetUserUnreadedMsg(c *fiber.Ctx) error {
	var UnReadedMsgSchema = database.DB.Collection("UnReadedmessage")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	userid := c.Query("userid")
	if userid == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "user id query param is required",
		})
	}

	// filter
	filter := bson.M{"mainUserid": userid, "isReaded": false}

	// query the db
	cursor, err := UnReadedMsgSchema.Find(ctx, filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Faild to retrivere unreaded messages",
			"error":   err.Error(),
		})
	}
	defer cursor.Close(ctx)

	// iterate over the cursor and build the res array
	var urms []models.UnReadModel
	totalUnreadMessageCount := 0
	for cursor.Next(ctx) {
		var urm models.UnReadModel
		err := cursor.Decode(&urm)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Faild to decode unread messages",
				"error":   err.Error(),
			})
		}
		urms = append(urms, urm)
		totalUnreadMessageCount += urm.NumOfUnreadedMessages
	}

	if len(urms) == 0 {
		urms = []models.UnReadModel{}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"messages": urms,
		"total":    totalUnreadMessageCount,
	})
}

// MarkMsgAsReaded
// @Summary mark messages as read for user
// @Description mark messages as read for user update the recoded make is read true num 0
// @Tags Chat
// @Accept json
// @Produce json
// @Param mainuid query string true "main user id"
// @Param otheruid query string true "other user id"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /chat/mark-msg-asreaded [get]
func MarkMsgAsReaded(c *fiber.Ctx) error {
	var UnReadedMsgSchema = database.DB.Collection("UnReadedmessage")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	mainuid := c.Query("mainuid")
	otheruid := c.Query("otheruid")
	if mainuid == "" || otheruid == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "mainuid and otheruid query param is required",
		})
	}

	// filter
	filter := bson.M{"mainUserid": mainuid, "otherUserid": otheruid}
	update := bson.M{"$set": bson.M{"isReaded": true, "numOfUnreadedMessages": 0}}

	// update the document
	options := options.FindOneAndUpdate().SetReturnDocument(options.After)
	result := UnReadedMsgSchema.FindOneAndUpdate(ctx, filter, update, options)

	if result.Err() != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"messages": "Faild to mark message as readed",
			"error":    result.Err().Error(),
		})
	}

	// check
	var updateDoc bson.M
	if err := result.Decode(&updateDoc); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"messages": "Faild to decode update document",
			"error":    err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"isMarked": true,
	})
}
