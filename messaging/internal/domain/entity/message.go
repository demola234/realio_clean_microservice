package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	ID             primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ConversationID string             `json:"conversationId" bson:"conversationId"`
	SenderID       string             `json:"senderId" bson:"senderId"`
	ReceiverID     string             `json:"receiverId" bson:"receiverId"`
	Content        string             `json:"content" bson:"content"`
	IsRead         bool               `json:"isRead" bson:"isRead"`
	CreatedAt      time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt      time.Time          `json:"updatedAt" bson:"updatedAt"`
	IsDeleted      bool               `json:"isDeleted" bson:"isDeleted"`
}
