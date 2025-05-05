package repository

import (
	"context"

	"github.com/demola234/messaging/internal/domain/entity"
)

type ConversationRepository interface {
	// Get all conversations for a specific user
	GetConversations(ctx context.Context, userID string) ([]entity.Conversation, error)

	// Get a specific conversation between two users
	GetConversationBetweenUsers(ctx context.Context, user1ID, user2ID string) ([]entity.Conversation, error)

	// Update the last message details in a conversation
	UpdateConversationLastMessage(ctx context.Context, conversationID string, messageID string, content string) error

	// Create a new conversation
	CreateConversation(ctx context.Context, conversation *entity.Conversation) error

	// Get All conversations betweens different users
	GetAllConversations(ctx context.Context, userID string) ([]entity.Conversation, error)
}
