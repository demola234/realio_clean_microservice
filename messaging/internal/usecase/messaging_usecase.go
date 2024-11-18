package usecase

import (
	"context"
	"errors"
	"fmt"
	"job_portal/messaging/internal/domain/entity"
	"job_portal/messaging/internal/domain/repository"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessagingUseCase interface {
	SendMessage(ctx context.Context, message *entity.Message) (entity.Message, error)
	GetMessages(ctx context.Context, conversationID string, includeDeleted *bool) ([]entity.Message, error)
	DeleteMessages(ctx context.Context, conversationID string) error
	UpdateMessage(ctx context.Context, messageID string, content string) error
	UpdateMessageReadStatus(ctx context.Context, messageID string, isRead bool) error
	GetConversationBetweenUsers(ctx context.Context, user1ID, user2ID string) ([]entity.Conversation, error)
	GetAllConversations(ctx context.Context, userID string) ([]entity.Conversation, error)
}

type messagingUseCase struct {
	messageRepo      repository.MessageRepository
	conversationRepo repository.ConversationRepository
}

// NewMessagingUseCase creates a new instance of MessagingUseCase.
func NewMessagingUseCase(
	messageRepo repository.MessageRepository,
	conversationRepo repository.ConversationRepository,
) MessagingUseCase {
	return &messagingUseCase{
		messageRepo:      messageRepo,
		conversationRepo: conversationRepo,
	}
}

func (uc *messagingUseCase) SendMessage(ctx context.Context, message *entity.Message) (entity.Message, error) {
	// Initialize message details
	message.ID = primitive.NewObjectID()
	message.IsRead = false
	message.IsDeleted = false
	message.CreatedAt = time.Now().UTC()
	message.UpdatedAt = time.Now().UTC()

	// Check for an existing conversation
	conversations, err := uc.conversationRepo.GetConversationBetweenUsers(ctx, message.SenderID, message.ReceiverID)
	if err != nil {
		return entity.Message{}, fmt.Errorf("failed to retrieve conversations: %w", err)
	}

	var conversationID primitive.ObjectID
	if len(conversations) > 0 {
		// Use the existing conversation
		conversationID = conversations[0].ID
	} else {
		// Create a new conversation if none exists
		conversationID = primitive.NewObjectID()
		conversation := &entity.Conversation{
			ID:           conversationID,
			Participants: []string{message.SenderID, message.ReceiverID},
			LastMessage: &entity.LastMessage{
				SenderID:  message.SenderID,
				Content:   message.Content,
				Timestamp: message.CreatedAt,
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := uc.conversationRepo.CreateConversation(ctx, conversation); err != nil {
			return entity.Message{}, fmt.Errorf("failed to create conversation: %w", err)
		}
	}

	// Set the conversation ID for the message and save it
	message.ConversationID = conversationID.Hex()
	if err := uc.messageRepo.SaveMessage(ctx, message); err != nil {
		return entity.Message{}, fmt.Errorf("failed to save message: %w", err)
	}

	// Update the last message in the conversation
	if err := uc.conversationRepo.UpdateConversationLastMessage(
		ctx,
		conversationID.Hex(),
		message.ID.Hex(),
		message.Content,
	); err != nil {
		return entity.Message{}, fmt.Errorf("failed to update conversation last message: %w", err)
	}

	return *message, nil
}

func (uc *messagingUseCase) GetMessages(ctx context.Context, conversationID string, includeDeleted *bool) ([]entity.Message, error) {
	// Validate input
	if conversationID == "" {
		return nil, errors.New("conversation ID cannot be empty")
	}

	// Retrieve messages from the repository
	messages, err := uc.messageRepo.GetMessages(ctx, conversationID, includeDeleted)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages for conversation ID %s: %w", conversationID, err)
	}

	return messages, nil
}

func (uc *messagingUseCase) DeleteMessages(ctx context.Context, conversationID string) error {
	// Perform a delete on messages in the conversation
	if err := uc.messageRepo.DeleteMessages(ctx, conversationID); err != nil {
		return fmt.Errorf("failed to delete messages for conversation ID %s: %w", conversationID, err)
	}

	return nil
}

func (uc *messagingUseCase) UpdateMessage(ctx context.Context, messageID string, content string) error {
	// Update the message content
	if err := uc.messageRepo.UpdateMessage(ctx, messageID, content); err != nil {
		return fmt.Errorf("failed to update message with ID %s: %w", messageID, err)
	}

	return nil
}

func (uc *messagingUseCase) UpdateMessageReadStatus(ctx context.Context, messageID string, isRead bool) error {
	// Update the read status of a specific message
	if err := uc.messageRepo.UpdateMessageReadStatus(ctx, messageID, isRead); err != nil {
		return fmt.Errorf("failed to update read status for message ID %s: %w", messageID, err)
	}

	return nil
}

func (uc *messagingUseCase) GetConversationBetweenUsers(ctx context.Context, user1ID, user2ID string) ([]entity.Conversation, error) {
	// Retrieve conversations between two users
	conversations, err := uc.conversationRepo.GetConversationBetweenUsers(ctx, user1ID, user2ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversations between users %s and %s: %w", user1ID, user2ID, err)
	}

	if len(conversations) == 0 {
		return nil, errors.New("no conversations found between the users")
	}

	return conversations, nil
}

// GetAllConversations implements MessagingUseCase.
func (uc *messagingUseCase) GetAllConversations(ctx context.Context, userID string) ([]entity.Conversation, error) {
	// Retrieve all conversations for a user
	conversations, err := uc.conversationRepo.GetAllConversations(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get all conversations for user %s: %w", userID, err)
	}

	if len(conversations) == 0 {
		return nil, errors.New("no conversations found for the user")
	}

	return conversations, nil
}
