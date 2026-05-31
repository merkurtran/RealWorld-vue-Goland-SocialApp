package servergrpc

import (
	"Server/database"
	"Server/models"
	pb "Server/protos"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Server struct {
	pb.UnimplementedRealtimeChatServiceServer
}

func (s *Server) GetUserFollowingFollowers(ctx context.Context, req *pb.UserID) (*pb.UserIDsListResponse, error) {
	UserSchema := database.DB.Collection("users")
	userID := req.GetUserid()

	if userID == "" {
		return nil, fmt.Errorf("user id is required")
	}
	var user models.UserModel

	uid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %s (error: %v)", userID, err)
	}

	err = UserSchema.FindOne(ctx, bson.M{"_id": uid}).Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("user not found: %s", userID)
	}

	friendsMap := make(map[string]bool)
	for _, id := range user.Following {
		friendsMap[id] = true
	}
	for _, id := range user.Followers {
		friendsMap[id] = true
	}
	friends := make([]string, 0)
	for id := range friendsMap {
		friends = append(friends, id)
	}
	return &pb.UserIDsListResponse{
		UserIDsLists: []*pb.UserIDsList{
			{UseridsList: friends},
		},
	}, nil
}

func (s *Server) SendMessage(ctx context.Context, req *pb.MessageRequest) (*pb.MessageResponse, error) {
	if req.GetSender() == "" || req.GetReceiver() == "" || req.GetContent() == "" {
		return nil, fmt.Errorf("sender, receiver and content are required")
	}
	UserSchema := database.DB.Collection("users")

	// check s & r are exists
	senderID, _ := primitive.ObjectIDFromHex(req.GetSender())
	receiverID, _ := primitive.ObjectIDFromHex(req.GetReceiver())

	var sender, reciver models.UserModel
	err := UserSchema.FindOne(ctx, bson.M{"_id": senderID}).Decode(&sender)
	if err != nil {
		return nil, fmt.Errorf("sender not found")
	}
	err = UserSchema.FindOne(ctx, bson.M{"_id": receiverID}).Decode(&reciver)
	if err != nil {
		return nil, fmt.Errorf("receiver not found")
	}

	// save msg
	message := models.MessageModel{
		Content:  req.GetContent(),
		Sender:   req.GetSender(),
		Receiver: req.GetReceiver(),
	}
	_, err = database.DB.Collection("messages").InsertOne(ctx, message)
	if err != nil {
		return nil, fmt.Errorf("error saving message: %s", err.Error())
	}
	// update unreaded messages
	unreadedmessagesSchema := database.DB.Collection("UnReadedmessage")

	existingRecored := bson.M{}
	err = unreadedmessagesSchema.FindOneAndUpdate(
		ctx,
		bson.M{"mainUserid": req.GetReceiver(), "otherUserid": req.GetSender()},
		bson.M{"$inc": bson.M{"numOfUnreadedMessages": 1}, "$set": bson.M{"isReaded": false}},
	).Decode(&existingRecored)

	if err != nil {
		_, err = unreadedmessagesSchema.InsertOne(
			ctx,
			bson.M{"mainUserid": req.GetReceiver(), "otherUserid": req.GetSender(), "numOfUnreadedMessages": 1, "isReaded": false},
		)
		if err != nil {
			return nil, fmt.Errorf("error saving unreaded message: %s", err.Error())
		}
	}

	// res
	return &pb.MessageResponse{
		Message: req.GetContent(),
	}, nil
}
