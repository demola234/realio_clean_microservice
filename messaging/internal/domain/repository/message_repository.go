package repository

import (
	"context"

	"github.com/demola234/messaging/internal/domain/entity"
)

type MessageRepository interface {
	// Save a new message
	SaveMessage(ctx context.Context, message *entity.Message) error

	// Get all messages for a conversation between two users
	GetMessages(ctx context.Context, conversationID string, includeDeleted *bool) ([]entity.Message, error)

	// Delete messages for a specific conversation
	DeleteMessages(ctx context.Context, conversationID string) error

	// Update the content of an existing message
	UpdateMessage(ctx context.Context, messageID string, content string) error

	// Update the read status of a message
	UpdateMessageReadStatus(ctx context.Context, messageID string, isRead bool) error
}
