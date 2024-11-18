package messageHandler

import (
	"context"
	"fmt"
	"job_portal/messaging/internal/domain/entity"
	"job_portal/messaging/internal/usecase"

	pb "job_portal/messaging/interfaces/api/grpc"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type MessageHandler struct {
	messageUseCase usecase.MessagingUseCase
	pb.UnimplementedMessagingServiceServer
}

// NewMessageHandler creates a new instance of MessageHandler.
func NewMessageHandler(messageUseCase usecase.MessagingUseCase) *MessageHandler {
	return &MessageHandler{
		messageUseCase: messageUseCase,
	}
}

// SendMessage handles the SendMessage gRPC request.
func (h *MessageHandler) SendMessage(ctx context.Context, req *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	message := &entity.Message{
		ConversationID: req.GetConversationId(),
		SenderID:       req.GetSenderId(),
		ReceiverID:     req.GetReceiverId(),
		Content:        req.GetContent(),
		IsRead:         false,
		IsDeleted:      false,
	}

	err := h.messageUseCase.SendMessage(ctx, message)
	if err != nil {
		return nil, fmt.Errorf("failed to send message: %w", err)
	}

	return &pb.SendMessageResponse{Status: "Message sent successfully"}, nil
}

// GetMessages handles the GetMessages gRPC request.
func (h *MessageHandler) GetMessages(ctx context.Context, req *pb.GetMessagesRequest) (*pb.GetMessagesResponse, error) {
	messages, err := h.messageUseCase.GetMessages(ctx, req.GetConversationId(), &req.IncludeDeleted)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}

	var pbMessages []*pb.Message
	for _, msg := range messages {
		pbMessages = append(pbMessages, &pb.Message{
			Id:             msg.ID.Hex(),
			ConversationId: msg.ConversationID,
			SenderId:       msg.SenderID,
			ReceiverId:     msg.ReceiverID,
			Content:        msg.Content,
			IsRead:         msg.IsRead,
			CreatedAt:      timestamppb.New(msg.CreatedAt),
			UpdatedAt:      timestamppb.New(msg.UpdatedAt),
		})
	}

	return &pb.GetMessagesResponse{Messages: pbMessages}, nil
}

// DeleteMessages handles the DeleteMessages gRPC request.
func (h *MessageHandler) DeleteMessages(ctx context.Context, req *pb.DeleteMessagesRequest) (*pb.DeleteMessagesResponse, error) {
	err := h.messageUseCase.DeleteMessages(ctx, req.GetConversationId())
	if err != nil {
		return nil, fmt.Errorf("failed to delete messages: %w", err)
	}

	return &pb.DeleteMessagesResponse{Status: "Messages deleted successfully"}, nil
}

// UpdateMessage handles the UpdateMessage gRPC request.
func (h *MessageHandler) UpdateMessage(ctx context.Context, req *pb.UpdateMessageRequest) (*pb.UpdateMessageResponse, error) {
	err := h.messageUseCase.UpdateMessage(ctx, req.GetMessageId(), req.GetContent())
	if err != nil {
		return nil, fmt.Errorf("failed to update message: %w", err)
	}

	return &pb.UpdateMessageResponse{Status: "Message updated successfully"}, nil
}

// UpdateMessageReadStatus handles the UpdateMessageReadStatus gRPC request.
func (h *MessageHandler) UpdateMessageReadStatus(ctx context.Context, req *pb.UpdateMessageReadStatusRequest) (*pb.UpdateMessageReadStatusResponse, error) {
	err := h.messageUseCase.UpdateMessageReadStatus(ctx, req.GetMessageId(), req.GetIsRead())
	if err != nil {
		return nil, fmt.Errorf("failed to update read status: %w", err)
	}

	return &pb.UpdateMessageReadStatusResponse{Status: "Message read status updated successfully"}, nil
}

// GetConversationBetweenUsers handles the GetConversationBetweenUsers gRPC request.
func (h *MessageHandler) GetConversationBetweenUsers(ctx context.Context, req *pb.GetConversationBetweenUsersRequest) (*pb.GetConversationBetweenUsersResponse, error) {
	conversations, err := h.messageUseCase.GetConversationBetweenUsers(ctx, req.GetUser1Id(), req.GetUser2Id())
	if err != nil {
		return nil, fmt.Errorf("failed to get conversations between users: %w", err)
	}

	var pbConversations []*pb.Conversation
	for _, conv := range conversations {
		pbConversations = append(pbConversations, &pb.Conversation{
			Id:           conv.ID.Hex(),
			Participants: conv.Participants,
			LastMessage: &pb.LastMessage{
				Content:   conv.LastMessage.Content,
				SenderId:  conv.LastMessage.SenderID,
				Timestamp: timestamppb.New(conv.LastMessage.Timestamp),
			},
			CreatedAt: timestamppb.New(conv.CreatedAt),
			UpdatedAt: timestamppb.New(conv.UpdatedAt),
		})
	}

	return &pb.GetConversationBetweenUsersResponse{Conversations: pbConversations}, nil
}
