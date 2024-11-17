package repository

import (
	"context"
	"job_portal/messaging/internal/domain/entity"
)

type MessageRepository interface {
	SaveMessage(c context.Context, message *entity.Message) error
	GetMessages(c context.Context, conversationID string, isRead *bool) ([]entity.Message, error)
	DeleteMessages(c context.Context, conversationID string) error
	UpdateMessage(c context.Context, messageID string, content string) error
	UpdateMessageReadStatus(c context.Context, messageID string, isRead bool) error
	GetConversationBetweenUsers(c context.Context, user1ID, user2ID string) ([]entity.Conversation, error)
	UpdateConversationLastMessage(c context.Context, conversationID string, messageID string, content string) error
}
