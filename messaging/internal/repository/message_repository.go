package repository

import (
	"context"
	"errors"
	"fmt"
	"job_portal/messaging/db/mongo"
	"job_portal/messaging/internal/domain/entity"
	"job_portal/messaging/internal/domain/repository"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type messageRepository struct {
	database   mongo.Database
	collection mongo.Collection
}

// NewMessageRepository creates a new instance of messageRepository.
func NewMessageRepository(db mongo.Database, collectionName string) repository.MessageRepository {
	return &messageRepository{
		database:   db,
		collection: db.Collection(collectionName),
	}
}

// GetMessages retrieves messages for a conversation with an optional filter for deleted messages.
func (m *messageRepository) GetMessages(ctx context.Context, conversationID string, includeDeleted *bool) ([]entity.Message, error) {
	filter := bson.M{"conversation_id": conversationID}
	if includeDeleted != nil && !*includeDeleted {
		filter["deleted"] = false
	}

	cursor, err := m.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch messages: %w", err)
	}
	defer cursor.Close(ctx)

	var messages []entity.Message
	if err := cursor.All(ctx, &messages); err != nil {
		return nil, fmt.Errorf("failed to decode messages: %w", err)
	}
	return messages, nil
}

// DeleteMessages deletes all messages for a specific conversation.
// DeleteMessages deletes all messages for a specific conversation.
func (m *messageRepository) DeleteMessages(ctx context.Context, conversationID string) error {
	filter := bson.M{"conversation_id": conversationID}
	result, err := m.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete messages: %w", err)
	}

	if result.DeletedCount == 0 {
		return errors.New("no messages found to delete")
	}
	return nil
}

// GetConversationBetweenUsers retrieves conversations between two users.
func (m *messageRepository) GetConversationBetweenUsers(ctx context.Context, user1ID, user2ID string) ([]entity.Conversation, error) {
	filter := bson.M{"participants": bson.M{"$all": []string{user1ID, user2ID}}}

	cursor, err := m.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch conversations: %w", err)
	}
	defer cursor.Close(ctx)

	var conversations []entity.Conversation
	if err := cursor.All(ctx, &conversations); err != nil {
		return nil, fmt.Errorf("failed to decode conversations: %w", err)
	}
	return conversations, nil
}

// SaveMessage saves a new message to the database.
func (m *messageRepository) SaveMessage(ctx context.Context, message *entity.Message) error {
	message.ID = primitive.NewObjectID()
	message.IsRead = false
	message.IsDeleted = false

	_, err := m.collection.InsertOne(ctx, message)
	if err != nil {
		return fmt.Errorf("failed to save message: %w", err)
	}
	return nil
}

// UpdateMessage updates the content of an existing message.
func (m *messageRepository) UpdateMessage(ctx context.Context, messageID, content string) error {
	id, err := primitive.ObjectIDFromHex(messageID)
	if err != nil {
		return fmt.Errorf("invalid message ID: %w", err)
	}

	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"content": content}}

	result, err := m.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update message: %w", err)
	}

	if result.MatchedCount == 0 {
		return errors.New("no message found to update")
	}
	return nil
}

// UpdateMessageReadStatus updates the read status of a message.
func (m *messageRepository) UpdateMessageReadStatus(ctx context.Context, messageID string, isRead bool) error {
	id, err := primitive.ObjectIDFromHex(messageID)
	if err != nil {
		return fmt.Errorf("invalid message ID: %w", err)
	}

	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"isRead": isRead}}

	result, err := m.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update read status: %w", err)
	}

	if result.MatchedCount == 0 {
		return errors.New("no message found to update read status")
	}
	return nil
}

// UpdateConversationLastMessage updates the last message details in a conversation.
func (m *messageRepository) UpdateConversationLastMessage(ctx context.Context, conversationID, messageID, content string) error {
	convID, err := primitive.ObjectIDFromHex(conversationID)
	if err != nil {
		return fmt.Errorf("invalid conversation ID: %w", err)
	}

	msgID, err := primitive.ObjectIDFromHex(messageID)
	if err != nil {
		return fmt.Errorf("invalid message ID: %w", err)
	}

	filter := bson.M{"_id": convID}
	update := bson.M{
		"$set": bson.M{
			"lastMessage": bson.M{
				"id":      msgID,
				"content": content,
			},
			"updatedAt": time.Now(),
		},
	}

	result, err := m.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update conversation: %w", err)
	}

	if result.MatchedCount == 0 {
		return errors.New("no conversation found to update")
	}
	return nil
}
