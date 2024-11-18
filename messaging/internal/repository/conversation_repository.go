package repository

import (
	"context"
	"fmt"
	"job_portal/messaging/db/mongo"
	"job_portal/messaging/internal/domain/entity"
	"job_portal/messaging/internal/domain/repository"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type conversationRepository struct {
	collection mongo.Collection
}

func NewConversationRepository(db mongo.Database, collectionName string) repository.ConversationRepository {
	return &conversationRepository{
		collection: db.Collection(collectionName),
	}
}

func (c *conversationRepository) GetConversations(ctx context.Context, userID string) ([]entity.Conversation, error) {
	filter := bson.M{"participants": userID}

	cursor, err := c.collection.Find(ctx, filter)
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

func (c *conversationRepository) GetConversationBetweenUsers(ctx context.Context, user1ID, user2ID string) ([]entity.Conversation, error) {
	filter := bson.M{"participants": bson.M{"$all": []string{user1ID, user2ID}}}

	cursor, err := c.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch conversations: %w", err)
	}
	defer cursor.Close(ctx)

	var conversations []entity.Conversation
	if err := cursor.All(ctx, &conversations); err != nil {
		return nil, fmt.Errorf("failed to decode conversations: %w", err)
	}

	if len(conversations) == 0 {
		return []entity.Conversation{}, nil
	}

	return conversations, nil
}

func (c *conversationRepository) UpdateConversationLastMessage(ctx context.Context, conversationID string, messageID string, content string) error {
	objectID, err := primitive.ObjectIDFromHex(conversationID)
	if err != nil {
		return fmt.Errorf("invalid conversation ID: %w", err)
	}

	update := bson.M{
		"$set": bson.M{
			"lastMessage.messageID": messageID,
			"lastMessage.content":   content,
			"lastMessage.timestamp": time.Now(),
			"updatedAt":             time.Now(),
		},
	}

	result, err := c.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return fmt.Errorf("failed to update last message: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("no conversation found to update")
	}

	return nil
}

func (c *conversationRepository) CreateConversation(ctx context.Context, conversation *entity.Conversation) error {
	_, err := c.collection.InsertOne(ctx, conversation)
	if err != nil {
		return fmt.Errorf("failed to create conversation: %w", err)
	}
	return nil
}
