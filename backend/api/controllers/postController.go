package controllers

import (
	"Server/database"
	"Server/models"
	"context"
	"math"
	"slices"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreatePost
// @Summary create a new post
// @Description create new post
// @Tags Posts
// @Accept json
// @Produce json
// @Param post body models.CreateOrUpdatePost true "post create details"
// @Success 201 {object} models.PostModel
// @Failure 400 {object} map[string]interface{}
// @security BearerAuth
// @Router /posts [post]
func CreatePost(c *fiber.Ctx) error {
	var UserSchema = database.DB.Collection("users")
	var PostSchema = database.DB.Collection("posts")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var body models.CreateOrUpdatePost
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	// start set data
	var post models.PostModel
	post.Creator = c.Locals("userID").(string)
	post.Likes = make([]string, 0)
	post.Comments = make([]string, 0)
	post.CreatedAt = time.Now()
	post.Title = body.Title
	post.Message = body.Message
	post.SelectedFile = body.SelectedFile

	var user models.UserModel
	objId, _ := primitive.ObjectIDFromHex(c.Locals("userID").(string))
	err := UserSchema.FindOne(ctx, bson.M{"_id": objId}).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"details": err.Error(),
		})
	}

	post.Name = user.Name
	// set data end
	// create post
	result, err := PostSchema.InsertOne(ctx, &post)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	} else {
		var createdPost models.PostModel
		query := bson.M{"_id": result.InsertedID}

		PostSchema.FindOne(ctx, query).Decode(&createdPost)
		return c.Status(fiber.StatusCreated).JSON(createdPost)
	}
}

// GetPost
// @Summary get post data
// @Description get post details
// @Tags Post
// @Accept json
// @Produce json
// @Param id path string true "Post ID"
// @Success 200 {object} models.PostModel
// @Failure 400 {object} map[string]interface{}
// @Router /posts/{id} [get]
func GetPost(c *fiber.Ctx) error {
	var PostSchema = database.DB.Collection("posts")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "post id is required",
		})
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	var getPost models.PostModel
	query := bson.M{"_id": objID}

	err = PostSchema.FindOne(ctx, query).Decode(&getPost)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "post not found",
			"error":   err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"post": getPost,
	})
}

// UpdatePost
// @Summary update post
// @Description update post
// @Tags Posts
// @Accept json
// @Produce json
// @Param id path string true "Post ID"
// @Param post body models.CreateOrUpdatePost true "update post details"
// @Success 200 {object} models.PostModel
// @Failure 400 {object} map[string]interface{}
// @security BearerAuth
// @Router /posts/{id} [patch]
func UpdatePost(c *fiber.Ctx) error {
	var PostSchema = database.DB.Collection("posts")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var newData models.CreateOrUpdatePost
	if err := c.BodyParser(&newData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	// authorization
	var authPost models.PostModel
	primID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := PostSchema.FindOne(ctx, bson.M{"_id": primID}).Decode(&authPost); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Post not found",
		})
	}

	if authPost.Creator != c.Locals("userID").(string) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "You are not authorized to update this post",
		})
	}

	// update post
	update := bson.M{
		"title":        newData.Title,
		"message":      newData.Message,
		"selectedFile": newData.SelectedFile,
	}
	_, err = PostSchema.UpdateOne(ctx, bson.M{"_id": authPost.ID}, bson.M{"$set": update})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	authPost.Title = newData.Title
	authPost.Message = newData.Message
	authPost.SelectedFile = newData.SelectedFile

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": authPost,
	})
}

// GetAllPosts
// @Summary get all post
// @Description get all post
// @Tags Posts
// @Accept json
// @Produce json
// @Param page query int false "page number"
// @Param id query string true "user id"
// @Success 200 {object} models.PostModel
// @Failure 400 {object} map[string]interface{}
// @security BearerAuth
// @Router /posts [get]
func GetAllPosts(c *fiber.Ctx) error {
	var PostSchema = database.DB.Collection("posts")
	var UserSchema = database.DB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var user models.UserModel
	var posts []models.PostModel

	userid := c.Query("id")
	page, _ := strconv.Atoi(c.Query("page", "1"))

	// get user following list ides and add our user id to it
	MainUserid, _ := primitive.ObjectIDFromHex(userid)
	UserSchema.FindOne(ctx, bson.M{"_id": MainUserid}).Decode(&user)

	user.Following = append(user.Following, userid)

	var LIMIT = 2

	findOptions := options.Find()
	filter := bson.M{"creator": bson.M{"$in": user.Following}}

	// get total num of posts
	total, err := PostSchema.CountDocuments(ctx, filter)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "No posts",
		})
	}

	findOptions.SetSkip((int64(page) - 1) * int64(LIMIT))
	findOptions.SetLimit(int64(LIMIT))
	findOptions.SetSort(bson.D{{Key: "_id", Value: -1}})

	// start get post posts
	cursor, err := PostSchema.Find(ctx, filter, findOptions)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "No posts",
		})
	}

	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var post models.PostModel
		cursor.Decode(&post)
		posts = append(posts, post)
	}

	if posts == nil {
		posts = make([]models.PostModel, 0)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":          posts,
		"currentPage":   page,
		"numberOfPages": math.Ceil(float64(total) / float64(LIMIT)),
	})
}

// GetPostsUsersBySearch Post
// @Summary get post by search query
// @Description get post and users matching the search query
// @Tags Posts
// @Accept json
// @Produce json
// @Param searchQuery query string true "Search query"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @security BearerAuth
// @Router /posts/search [get]
func GetPostsUsersBySearch(c *fiber.Ctx) error {
	var PostSchema = database.DB.Collection("posts")
	var UserSchema = database.DB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var users []models.UserModel
	var posts []models.PostModel

	//
	filterPost := bson.M{}
	filterUser := bson.M{}

	//
	findOptionsPost := options.Find()
	findOptionsUser := options.Find()

	if search := c.Query("searchQuery"); search != "" {
		// post
		filterPost = bson.M{
			"$or": []bson.M{
				{
					"title": bson.M{
						"$regex": primitive.Regex{
							Pattern: search,
							Options: "i",
						},
					},
				},
				{
					"message": bson.M{
						"$regex": primitive.Regex{
							Pattern: search,
							Options: "i",
						},
					},
				},
			},
		}
		// user
		filterUser = bson.M{
			"$or": []bson.M{
				{
					"name": bson.M{
						"$regex": primitive.Regex{
							Pattern: search,
							Options: "i",
						},
					},
				},
				{
					"email": bson.M{
						"$regex": primitive.Regex{
							Pattern: search,
							Options: "i",
						},
					},
				},
			},
		}
	}

	// end
	cursorPosts, _ := PostSchema.Find(ctx, filterPost, findOptionsPost)
	defer cursorPosts.Close(ctx)

	cursorUsers, _ := UserSchema.Find(ctx, filterUser, findOptionsUser)
	defer cursorUsers.Close(ctx)

	for cursorUsers.Next(ctx) {
		var user models.UserModel
		cursorUsers.Decode(&user)
		users = append(users, user)
	}

	for cursorPosts.Next(ctx) {
		var post models.PostModel
		cursorPosts.Decode(&post)
		posts = append(posts, post)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"users": users,
		"posts": posts,
	})

}

// CommentPost
// @Summary comment post
// @Description comment post
// @Tags Posts
// @Accept json
// @Produce json
// @Param id path string true "Post ID"
// @Param post body models.CommentPost true "comment value"
// @Success 200 {object} models.PostModel
// @Failure 400 {object} map[string]interface{}
// @security BearerAuth
// @Router /posts/{id}/commentPost [post]
func CommentPost(c *fiber.Ctx) error {
	var UserSchema = database.DB.Collection("users")
	var PostSchema = database.DB.Collection("posts")
	var NotificationSchema = database.DB.Collection("notifications")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var b models.CommentPost
	if err := c.BodyParser(&b); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	var post models.PostModel
	postid, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"details": err.Error(),
		})
	}

	err = PostSchema.FindOne(ctx, bson.M{"_id": postid}).Decode(&post)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"details": err.Error(),
		})
	}
	//
	newComment := bson.M{"comments": append(post.Comments, b.Value)}
	_, err = PostSchema.UpdateOne(ctx, bson.M{"_id": postid}, bson.M{"$set": newComment})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"details": err.Error(),
		})
	}

	err = PostSchema.FindOne(ctx, bson.M{"_id": postid}).Decode(&post)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"details": err.Error(),
		})
	}
	// Create notififcation start
	userID := c.Locals("userID").(string)
	objId, _ := primitive.ObjectIDFromHex(userID)
	var user models.UserModel

	// get user data
	userResult := UserSchema.FindOne(ctx, bson.M{"_id": objId})
	if userResult.Err() != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "user not found",
		})
	}
	userResult.Decode(&user)
	// Create Notification
	notification := models.Notification{
		MainUID:   post.Creator,
		TargetID:  postid.Hex(),
		Details:   user.Name + "Commented on you Post",
		User:      models.User{Name: user.Name, Avata: user.ImageUrl},
		CreatedAt: time.Now(),
	}
	_, err = NotificationSchema.InsertOne(ctx, notification)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Faild to create notification",
			"error":   err.Error(),
		})
	}
	// end
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": post,
	})

}

// LikePost
// @Summary like or unlike a post by it's id
// @Description like post
// @Tags Posts
// @Accept json
// @Produce json
// @Param id path string true "Post ID"
// @Success 200 {object} models.PostModel
// @Failure 400 {object} map[string]interface{}
// @security BearerAuth
// @Router /posts/{id}/likePost [patch]
func LikePost(c *fiber.Ctx) error {
	var UserSchema = database.DB.Collection("users")
	var PostSchema = database.DB.Collection("posts")
	var NotificationSchema = database.DB.Collection("notifications")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var post models.PostModel
	postid, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"details": err.Error(),
		})
	}

	err = PostSchema.FindOne(ctx, bson.M{"_id": postid}).Decode(&post)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"details": err.Error(),
		})
	}
	// after getting post
	userID, ok := c.Locals("userID").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"details": "You are not authorized",
		})
	}

	//check
	if slices.Contains(post.Likes, userID) {
		i := slices.Index(post.Likes, userID)
		post.Likes = slices.Delete(post.Likes, i, i+1)
	} else {
		post.Likes = append(post.Likes, userID)
		// TODO Statr create notification
		objId, _ := primitive.ObjectIDFromHex(userID)
		var user models.UserModel

		// get user data
		userResult := UserSchema.FindOne(ctx, bson.M{"_id": objId})
		if userResult.Err() != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "user not found",
			})
		}
		userResult.Decode(&user)
		// Create Notification
		notification := models.Notification{
			MainUID:   post.Creator,
			TargetID:  post.ID.Hex(),
			Details:   user.Name + "Liked your Post",
			User:      models.User{Name: user.Name, Avata: user.ImageUrl},
			CreatedAt: time.Now(),
		}
		_, err = NotificationSchema.InsertOne(ctx, notification)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Faild to create notification",
				"error":   err.Error(),
			})
		}
		// end
		// End create nofyication
	}

	// update post
	updateLikelist := bson.M{"likes": post.Likes}
	_, err = PostSchema.UpdateOne(ctx, bson.M{"_id": postid}, bson.M{"$set": updateLikelist})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"details": err.Error(),
		})
	}

	err = PostSchema.FindOne(ctx, bson.M{"_id": postid}).Decode(&post)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"details": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"post": post,
	})

}

// DeletePost
// @Summary delete post by id
// @Description delete post by post id need to provided auth token for post creator
// @Tags Posts
// @Accept json
// @Produce json
// @Param id path string true "Post ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @security BearerAuth
// @Router /posts/{id} [delete]
func DeletePost(c *fiber.Ctx) error {
	var PostSchema = database.DB.Collection("posts")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// authorization
	var authPost models.PostModel
	primID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := PostSchema.FindOne(ctx, bson.M{"_id": primID}).Decode(&authPost); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Post not found",
		})
	}

	if authPost.Creator != c.Locals("userID").(string) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "You are not authorized to delete this post",
		})
	}

	//
	result, err := PostSchema.DeleteOne(ctx, bson.M{"_id": primID})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"details": err.Error(),
		})
	}

	if result.DeletedCount == 1 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Post deleted successfully!",
		})
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Can't deleted post",
		})
	}
}
