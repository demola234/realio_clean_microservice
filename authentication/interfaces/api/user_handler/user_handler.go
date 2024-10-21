package user_handler

import (
	"context"
	pb "job_portal/authentication/interfaces/api/grpc"
	"job_portal/authentication/internal/usecase"

	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserHandler struct {
	userUsecase usecase.UserUsecase
	pb.UnimplementedAuthServiceServer
}

// NewUserHandler creates a new instance of UserHandler
func NewUserHandler(userUsecase usecase.UserUsecase) *UserHandler {
	return &UserHandler{userUsecase: userUsecase}
}

func (h *UserHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	user, err := h.userUsecase.RegisterUser(ctx, req.FullName, req.Password, req.Email)
	if err != nil {
		return nil, status.Errorf(400, err.Error())
	}

	return &pb.RegisterResponse{
		User: &pb.User{
			Email:             user.Email,
			FullName:          user.FullName,
			PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
			CreatedAt:         timestamppb.New(user.CreatedAt),
		},
		AccessToken:  user.AccessToken,
		RefreshToken: user.RefreshToken,
	}, nil
}

func (h *UserHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, err := h.userUsecase.LoginUser(ctx, req.Password, req.Email)
	if err != nil {
		return nil, status.Errorf(401, "invalid credentials %d", err)
	}

	token, err := h.userUsecase.GenerateToken(ctx, user.Email)
	if err != nil {
		return nil, status.Errorf(500, "failed to generate token")
	}

	refreshToken, err := h.userUsecase.GenerateToken(ctx, user.Email)
	if err != nil {
		return nil, status.Errorf(500, "failed to generate token")
	}

	return &pb.LoginResponse{
		User: &pb.User{
			Email:             user.Email,
			FullName:          user.FullName,
			PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
			CreatedAt:         timestamppb.New(user.CreatedAt),
		},
		AccessToken:  token,
		RefreshToken: refreshToken,
	}, nil
}

func (h *UserHandler) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	// Implementation for RefreshToken if needed
	return nil, status.Errorf(501, "not implemented")
}
