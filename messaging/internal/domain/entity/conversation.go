package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LastMessage struct {
	Content   string    `json:"content" bson:"content"`
	SenderID  string    `json:"senderId" bson:"senderId"`
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
}

type Conversation struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Participants []string           `json:"participants" bson:"participants"`
	LastMessage  *LastMessage       `json:"lastMessage" bson:"lastMessage"`
	CreatedAt    time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt    time.Time          `json:"updatedAt" bson:"updatedAt"`
}
