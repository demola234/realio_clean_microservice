package repository

// import (
// 	"context"
// 	"job_portal/messaging/db/mocks"
// 	"job_portal/messaging/internal/domain/entity"
// 	"testing"
// 	"time"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// 	"go.mongodb.org/mongo-driver/mongo"
// )

// func setupMockRepo() (*mocks.Database, *mocks.Collection, *messageRepository) {
// 	mockDatabase := &mocks.Database{}
// 	mockCollection := &mocks.Collection{}
// 	mockDatabase.On("Collection", "messages").Return(mockCollection)

// 	repo := &messageRepository{
// 		database:   mockDatabase,
// 		collection: mockCollection,
// 	}
// 	return mockDatabase, mockCollection, repo
// }

// func TestSaveMessage(t *testing.T) {
// 	_, mockCollection, repo := setupMockRepo()

// 	message := &entity.Message{
// 		ID:             primitive.NewObjectID(),
// 		ConversationID: primitive.NewObjectID().String(),
// 		SenderID:       primitive.NewObjectID().String(),
// 		Content:        "Hello, world!",
// 		CreatedAt:      time.Now(),
// 		UpdatedAt:      time.Now(),
// 		IsDeleted:      false,
// 	}

// 	mockCollection.On("InsertOne", mock.Anything, message).Return(&mongo.InsertOneResult{InsertedID: message.ID}, nil)

// 	err := repo.SaveMessage(context.Background(), message)
// 	assert.NoError(t, err)

// 	mockCollection.AssertExpectations(t)
// }

// // func TestGetMessages(t *testing.T) {
// // 	_, mockCollection, repo := setupMockRepo()

// // 	conversationID := primitive.NewObjectID().Hex()
// // 	messages := []entity.Message{
// // 		{
// // 			ID:             primitive.NewObjectID(),
// // 			ConversationID: primitive.NewObjectID().String(),
// // 			SenderID:       primitive.NewObjectID().String(),
// // 			Content:        "Hello!",
// // 			CreatedAt:      time.Now(),
// // 			UpdatedAt:      time.Now(),
// // 			IsDeleted:      false,
// // 		},
// // 	}

// // 	mockCursor := &mocks.Cursor{}
// // 	mockCollection.On("Find", mock.Anything, bson.M{"conversation_id": conversationID}).Return(mockCursor, nil)
// // 	mockCursor.On("All", mock.Anything, mock.AnythingOfType("*[]entity.Message")).Run(func(args mock.Arguments) {
// // 		arg := args.Get(1).(*[]entity.Message)
// // 		*arg = messages
// // 	}).Return(nil)
// // 	mockCursor.On("Close", mock.Anything).Return(nil)

// // 	result, err := repo.GetMessages(context.Background(), conversationID, nil)

// // 	assert.NoError(t, err)
// // 	assert.Equal(t, messages, result)

// // 	mockCollection.AssertExpectations(t)
// // 	mockCursor.AssertExpectations(t)
// // }

// func TestDeleteMessages(t *testing.T) {
// 	_, mockCollection, repo := setupMockRepo()

// 	conversationID := primitive.NewObjectID().Hex()
// 	filter := bson.M{"conversation_id": conversationID}

// 	// Mock DeleteOne to return a valid DeleteResult
// 	mockCollection.On("DeleteOne", mock.Anything, filter).
// 		Return(&mongo.DeleteResult{DeletedCount: 1}, nil)

// 	// Call the method under test
// 	err := repo.DeleteMessages(context.Background(), conversationID)

// 	// Assertions
// 	assert.NoError(t, err)
// 	mockCollection.AssertExpectations(t)
// }

// func TestUpdateMessage(t *testing.T) {
// 	_, mockCollection, repo := setupMockRepo()

// 	messageID := primitive.NewObjectID()
// 	content := "Updated Content"
// 	filter := bson.M{"_id": messageID}
// 	update := bson.M{"$set": bson.M{"content": content}}

// 	mockCollection.On("UpdateOne", mock.Anything, filter, update).Return(&mongo.UpdateResult{MatchedCount: 1, ModifiedCount: 1}, nil)

// 	err := repo.UpdateMessage(context.Background(), messageID.Hex(), content)

// 	assert.NoError(t, err)
// 	mockCollection.AssertExpectations(t)
// }

// func TestUpdateConversationLastMessage(t *testing.T) {
// 	_, mockCollection, repo := setupMockRepo()

// 	conversationID := primitive.NewObjectID()
// 	messageID := primitive.NewObjectID()
// 	content := "Updated message content"

// 	filter := bson.M{"_id": conversationID}
// 	update := bson.M{
// 		"$set": bson.M{
// 			"lastMessage": bson.M{
// 				"id":      messageID,
// 				"content": content,
// 			},
// 			"updatedAt": mock.MatchedBy(func(v interface{}) bool {
// 				_, ok := v.(time.Time)
// 				return ok
// 			}),
// 		},
// 	}

// 	mockCollection.On("UpdateOne", mock.Anything, filter, update).Return(&mongo.UpdateResult{MatchedCount: 1}, nil)

// 	err := repo.UpdateConversationLastMessage(context.Background(), conversationID.Hex(), messageID.Hex(), content)

// 	assert.NoError(t, err)
// 	mockCollection.AssertExpectations(t)
// }
