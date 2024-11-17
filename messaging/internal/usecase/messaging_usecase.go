package usecase

import (
	"context"
	"errors"
	"job_portal/messaging/internal/domain/entity"
	"job_portal/messaging/internal/domain/repository"
	"time"
)

type MessagingUseCase interface {
	SendMessage(ctx context.Context, message *entity.Message) error
	GetMessages(ctx context.Context, conversationID string, isRead *bool) ([]entity.Message, error)
	DeleteMessages(ctx context.Context, conversationID string) error
	UpdateMessage(ctx context.Context, messageID string, content string) error
	UpdateMessageReadStatus(ctx context.Context, messageID string, isRead bool) error
	GetConversationBetweenUsers(ctx context.Context, user1ID, user2ID string) ([]entity.Conversation, error)
}

type messagingUseCase struct {
	messageRepo repository.MessageRepository
}

// NewMessagingUseCase creates a new instance of MessagingUseCase.
func NewMessagingUseCase(
	messageRepo repository.MessageRepository,

) MessagingUseCase {
	return &messagingUseCase{
		messageRepo: messageRepo,
	}
}

func (uc *messagingUseCase) SendMessage(ctx context.Context, message *entity.Message) error {
	// Set message timestamps
	message.CreatedAt = time.Now()
	message.UpdatedAt = time.Now()

	// Save the message
	if err := uc.messageRepo.SaveMessage(ctx, message); err != nil {
		return err
	}

	// Update the conversation's last message
	if err := uc.messageRepo.UpdateConversationLastMessage(ctx, message.ConversationID.String(), message.ID.String(), message.Content); err != nil {
		return err
	}

	return nil
}

func (uc *messagingUseCase) GetMessages(ctx context.Context, conversationID string, isRead *bool) ([]entity.Message, error) {
	// Retrieve messages from the repository
	messages, err := uc.messageRepo.GetMessages(ctx, conversationID, isRead)
	if err != nil {
		return nil, err
	}

	// Business logic: Additional filtering or processing can be added here if needed
	return messages, nil
}

func (uc *messagingUseCase) DeleteMessages(ctx context.Context, conversationID string) error {
	// Perform a soft delete on messages in the conversation
	if err := uc.messageRepo.DeleteMessages(ctx, conversationID); err != nil {
		return err
	}

	return nil
}

func (uc *messagingUseCase) UpdateMessage(ctx context.Context, messageID string, content string) error {
	// Update the message content
	if err := uc.messageRepo.UpdateMessage(ctx, messageID, content); err != nil {
		return err
	}

	return nil
}

func (uc *messagingUseCase) UpdateMessageReadStatus(ctx context.Context, messageID string, isRead bool) error {
	// Update the read status of a specific message
	if err := uc.messageRepo.UpdateMessageReadStatus(ctx, messageID, isRead); err != nil {
		return err
	}

	return nil
}

func (uc *messagingUseCase) GetConversationBetweenUsers(ctx context.Context, user1ID, user2ID string) ([]entity.Conversation, error) {
	// Retrieve conversations between two users
	conversations, err := uc.messageRepo.GetConversationBetweenUsers(ctx, user1ID, user2ID)
	if err != nil {
		return nil, err
	}

	// Ensure at least one conversation exists
	if len(conversations) == 0 {
		return nil, errors.New("no conversations found between the users")
	}

	return conversations, nil
}
