package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessageModel struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Content  string             `json:"content" bson:"content"`
	Sender   string             `json:"sender" bson:"sender"`
	Receiver string             `json:"receiver" bson:"receiver"`
}

type SendMessageModel struct {
	Content  string `json:"content" bson:"content" validate:"required,min=5"`
	Sender   string `json:"sender" bson:"sender" validate:"required"`
	Receiver string `json:"receiver" bson:"receiver" validate:"required"`
}
