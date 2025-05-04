package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/demola234/messaging/db/mongo"
	"github.com/demola234/messaging/internal/domain/entity"
	"github.com/demola234/messaging/internal/domain/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type messageRepository struct {
	collection mongo.Collection
}

func NewMessageRepository(db mongo.Database, collectionName string) repository.MessageRepository {
	return &messageRepository{
		collection: db.Collection(collectionName),
	}
}

func (m *messageRepository) SaveMessage(ctx context.Context, message *entity.Message) error {
	_, err := m.collection.InsertOne(ctx, message)
	return err
}

func (m *messageRepository) GetMessages(ctx context.Context, conversationID string, includeDeleted *bool) ([]entity.Message, error) {
	// objectID, err := primitive.ObjectIDFromHex(conversationID)
	// if err != nil {
	// 	return nil, fmt.Errorf("invalid conversation ID: %w", err)
	// }

	// Use the correct field name as per the BSON tag in the Message struct
	filter := bson.M{"conversationId": conversationID}
	// if includeDeleted != nil && !*includeDeleted {
	// 	filter["isDeleted"] = false
	// }

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

func (m *messageRepository) DeleteMessages(ctx context.Context, conversationID string) error {
	// Handle both ObjectID and string formats for conversationId
	filter := bson.M{}
	if conversationID != "" {
		filter["conversationId"] = conversationID
	}

	// objectID, err := primitive.ObjectIDFromHex(conversationID)
	// if err == nil {
	// 	// If conversationID can be converted to ObjectID, use it
	// 	filter["conversationId"] = objectID
	// } else {
	// 	// Otherwise, fallback to string format
	// 	filter["conversationId"] = conversationID
	// }

	// Log the query for debugging
	log.Printf("Deleting messages with filter: %+v", filter)

	// Perform delete operation
	result, err := m.collection.DeleteMany(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete messages: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("no messages found to delete for conversation ID: %s", conversationID)
	}

	return nil
}

func (m *messageRepository) UpdateMessage(ctx context.Context, messageID string, content string) error {
	id, err := primitive.ObjectIDFromHex(messageID)
	if err != nil {
		return fmt.Errorf("invalid message ID: %w", err)
	}

	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"content": content, "updatedAt": time.Now()}}

	result, err := m.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update message: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("no message found to update")
	}

	return nil
}

func (m *messageRepository) UpdateMessageReadStatus(ctx context.Context, messageID string, isRead bool) error {
	id, err := primitive.ObjectIDFromHex(messageID)
	if err != nil {
		return fmt.Errorf("invalid message ID: %w", err)
	}

	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"isRead": isRead, "updatedAt": time.Now()}}

	result, err := m.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update read status: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("no message found to update read status")
	}

	return nil
}
