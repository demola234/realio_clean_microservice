package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	ID             primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ConversationID primitive.ObjectID `json:"conversationId" bson:"conversationId"`
	SenderID       primitive.ObjectID `json:"senderId" bson:"senderId"`
	Content        string             `json:"content" bson:"content"`
	IsRead         bool               `json:"isRead" bson:"isRead"`
	CreatedAt      time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt      time.Time          `json:"updatedAt" bson:"updatedAt"`
	IsDeleted      bool               `json:"isDeleted" bson:"isDeleted"`
}
