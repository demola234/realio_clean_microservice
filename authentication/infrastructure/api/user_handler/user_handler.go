package user_handler

import (
	"context"
	"time"

	pb "github.com/demola234/authentication/infrastructure/api/grpc"
	"github.com/demola234/authentication/internal/usecase"

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
	user, _, err := h.userUsecase.RegisterUser(ctx, req.FullName, req.Password, req.Email, req.Role, req.Phone)
	if err != nil {
		return nil, status.Errorf(400, err.Error())
	}

	return &pb.RegisterResponse{
		User: &pb.User{
			Email:     user.Email,
			FullName:  user.FullName,
			UserId:    user.ID.String(),
			Role:      user.Role,
			Phone:     user.Phone,
			UpdatedAt: timestamppb.New(user.UpdatedAt),
			CreatedAt: timestamppb.New(user.CreatedAt),
		},
	}, nil
}

func (h *UserHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, err := h.userUsecase.LoginUser(ctx, req.Password, req.Email)
	if err != nil {
		return nil, status.Errorf(401, "invalid credentials %d", err)
	}

	token, err := h.userUsecase.GenerateToken(ctx, user.Email, user.ID.String())
	if err != nil {
		return nil, status.Errorf(500, "failed to generate token")
	}

	return &pb.LoginResponse{
		User: &pb.User{
			Email:     user.Email,
			FullName:  user.FullName,
			UserId:    user.ID.String(),
			Role:      user.Role,
			Phone:     user.Phone,
			UpdatedAt: timestamppb.New(user.UpdatedAt),
			CreatedAt: timestamppb.New(user.CreatedAt),
		},
		Session: &pb.Session{
			Token:     token,
			ExpiresAt: timestamppb.New(time.Now().Add(time.Hour * 24)),
		},
	}, nil
}

func (h *UserHandler) VerifyUser(ctx context.Context, req *pb.VerifyUserRequest) (*pb.VerifyUserResponse, error) {

	// Check if user is already verified
	valid, err := h.userUsecase.VerifyOtp(ctx, req.Email, req.Otp)
	if err != nil {
		return nil, status.Errorf(401, "invalid credentials %d", err)
	}

	if !valid {
		return &pb.VerifyUserResponse{
			Valid: false,
		}, status.Errorf(401, "invalid credentials %d", err)
	}

	return &pb.VerifyUserResponse{
		Valid: valid,
	}, nil

}

func (h *UserHandler) ResendOtp(ctx context.Context, req *pb.ResendOtpRequest) (*pb.ResendOtpResponse, error) {
	// Get User Info and Check if otp is valid
	_, err := h.userUsecase.GetUser(ctx, req.Email)
	if err != nil {
		return nil, status.Errorf(401, "invalid credentials %d", err)
	}

	// Generate OTP
	err = h.userUsecase.ResendOtp(ctx, req.Email)
	if err != nil {
		return nil, status.Errorf(401, "invalid credentials %d", err)
	}

	return &pb.ResendOtpResponse{
		Message: "OTP sent successfully",
	}, nil

}
func (h *UserHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user, err := h.userUsecase.GetUser(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(401, "invalid credentials %d", err)
	}

	return &pb.GetUserResponse{
		User: &pb.User{
			Email:     user.Email,
			FullName:  user.FullName,
			UserId:    user.ID.String(),
			Role:      user.Role,
			Phone:     user.Phone,
			UpdatedAt: timestamppb.New(user.UpdatedAt),
			CreatedAt: timestamppb.New(user.CreatedAt),
		},
	}, nil
}

func (h *UserHandler) LogOut(ctx context.Context, req *pb.LogOutRequest) (*pb.LogOutResponse, error) {
	err := h.userUsecase.LogOut(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(401, "invalid credentials %d", err)
	}

	return &pb.LogOutResponse{
		Message: "Logged out successfully",
	}, nil
}
