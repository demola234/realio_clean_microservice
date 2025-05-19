package user_handler

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"time"

	pb "github.com/demola234/authentication/infrastructure/api/grpc"
	"github.com/demola234/authentication/internal/domain/entity"
	"github.com/demola234/authentication/internal/usecase"

	"google.golang.org/grpc/codes"
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

	session, err := h.userUsecase.GetSession(ctx, user.ID.String())
	if err != nil {
		return nil, status.Errorf(500, "failed to get session")
	}

	return &pb.LoginResponse{
		User: &pb.User{
			Email:      user.Email,
			FullName:   user.FullName,
			UserId:     user.ID.String(),
			Role:       user.Role,
			Phone:      user.Phone,
			IsVerified: session.OTPVerified,
			UpdatedAt:  timestamppb.New(user.UpdatedAt),
			CreatedAt:  timestamppb.New(user.CreatedAt),
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

	// Get User Info and Check if otp is valid
	user, err := h.userUsecase.GetUser(ctx, req.Email)
	if err != nil {
		return nil, status.Errorf(401, "invalid credentials %d", err)
	}

	token, err := h.userUsecase.GenerateToken(ctx, user.Email, user.ID.String())
	if err != nil {
		return nil, status.Errorf(500, "failed to generate token")
	}

	return &pb.VerifyUserResponse{
		Valid: valid,
		Session: &pb.Session{
			Token:     token,
			ExpiresAt: timestamppb.New(time.Now().Add(time.Hour * 24)),
		},
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

func (h *UserHandler) UploadImage(ctx context.Context, req *pb.UploadImageRequest) (*pb.UploadImageResponse, error) {
	user, err := h.userUsecase.GetUser(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid credentials: %v", err)
	}

	// DO NOT convert binary data to string - use bytes.NewReader directly
	reader := bytes.NewReader(req.Content)

	imageUrl, err := h.userUsecase.UppdateProfileImage(ctx, reader, user.ID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to upload image: %v", err)
	}

	return &pb.UploadImageResponse{
		Message:  "Image uploaded successfully",
		ImageUrl: imageUrl,
		UserId:   req.UserId,
	}, nil
}

// ForgotPassword initiates the password reset process with OTP
func (h *UserHandler) ForgotPassword(ctx context.Context, req *pb.ForgotPasswordRequest) (*pb.ForgotPasswordResponse, error) {
	// Validate email
	if req.Email == "" {
		return nil, status.Errorf(codes.InvalidArgument, "email is required")
	}

	// Call the usecase
	err := h.userUsecase.ForgetPassword(ctx, req.Email)
	if err != nil {
		// Log the error internally, but don't expose it to the client
		// This prevents email enumeration attacks

		// For development, you might want to log the actual error
		// fmt.Printf("Error in ForgotPassword: %v\n", err)

		// Still return a generic success message for security
	}

	// Return a generic success message - don't reveal if email exists
	return &pb.ForgotPasswordResponse{
		Message: "If your email exists in our system, you will receive a verification code shortly",
	}, nil
}

// VerifyResetPassword verifies the OTP for password reset
func (h *UserHandler) VerifyResetPassword(ctx context.Context, req *pb.VerifyResetPasswordRequest) (*pb.VerifyResetPasswordResponse, error) {
	// Validate inputs
	if req.Email == "" {
		return nil, status.Errorf(codes.InvalidArgument, "email is required")
	}

	if req.Otp == "" {
		return nil, status.Errorf(codes.InvalidArgument, "OTP is required")
	}

	// Call the usecase
	err := h.userUsecase.VerifyResetPassword(ctx, req.Email, req.Otp)
	if err != nil {
		// Return error with appropriate status code
		errMsg := strings.ToLower(err.Error())

		if strings.Contains(errMsg, "invalid otp format") {
			return &pb.VerifyResetPasswordResponse{
				Message: "Invalid verification code format",
				Valid:   false,
			}, nil
		}

		if strings.Contains(errMsg, "maximum otp attempts exceeded") {
			return &pb.VerifyResetPasswordResponse{
				Message: "Too many failed attempts, please request a new code",
				Valid:   false,
			}, nil
		}

		if strings.Contains(errMsg, "otp has expired") {
			return &pb.VerifyResetPasswordResponse{
				Message: "Verification code has expired, please request a new one",
				Valid:   false,
			}, nil
		}

		if strings.Contains(errMsg, "invalid otp") {
			return &pb.VerifyResetPasswordResponse{
				Message: "Invalid verification code",
				Valid:   false,
			}, nil
		}

		// For other errors, return a generic message
		return &pb.VerifyResetPasswordResponse{
			Message: "Verification failed",
			Valid:   false,
		}, nil
	}

	// Return success
	return &pb.VerifyResetPasswordResponse{
		Message: "Verification successful, you may now reset your password",
		Valid:   true,
	}, nil
}

// ResetPassword sets a new password after OTP verification
func (h *UserHandler) ResetPassword(ctx context.Context, req *pb.ResetPasswordRequest) (*pb.ResetPasswordResponse, error) {
	// Validate inputs
	if req.Email == "" {
		return nil, status.Errorf(codes.InvalidArgument, "email is required")
	}

	if req.NewPassword == "" {
		return nil, status.Errorf(codes.InvalidArgument, "new password is required")
	}

	// Call the usecase
	err := h.userUsecase.ResetPassword(ctx, req.Email, req.NewPassword)
	if err != nil {
		// Categorize errors for appropriate status codes
		errMsg := strings.ToLower(err.Error())

		if strings.Contains(errMsg, "otp not verified") {
			return nil, status.Errorf(codes.FailedPrecondition, "You must verify the code before resetting your password")
		}

		if strings.Contains(errMsg, "invalid password") {
			return nil, status.Errorf(codes.InvalidArgument, "Password does not meet requirements")
		}

		// For other errors, return a generic message
		return nil, status.Errorf(codes.Internal, "Failed to reset password")
	}

	// Return success message
	return &pb.ResetPasswordResponse{
		Message: "Your password has been successfully reset",
	}, nil
}

// Helper function to convert user entity to proto
func convertUserToProto(user *entity.User) *pb.User {
	return &pb.User{
		UserId:     user.ID.String(),
		Email:      user.Email,
		FullName:   user.FullName,
		Role:       user.Role,
		Phone:      user.Phone,
		IsVerified: user.EmailVerified,
		CreatedAt:  timestamppb.New(user.CreatedAt),
		UpdatedAt:  timestamppb.New(user.UpdatedAt),
	}
}

// ChangePassword handles changing the user's password
func (h *UserHandler) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {

	// Validate request
	if req.CurrentPassword == "" || req.NewPassword == "" {
		return nil, status.Errorf(codes.InvalidArgument, "current password and new password are required")
	}

	// Call usecase
	err := h.userUsecase.ChangePassword(ctx, req.CurrentPassword, req.NewPassword, req.UserId)
	if err != nil {
		// Check for specific error types
		if err.Error() == "current password is incorrect" {
			return nil, status.Errorf(codes.InvalidArgument, "current password is incorrect")
		}
		if err.Error() == "invalid new password" {
			return nil, status.Errorf(codes.InvalidArgument, "new password does not meet requirements")
		}
		return nil, status.Errorf(codes.Internal, "failed to change password: %v", err)
	}

	return &pb.ChangePasswordResponse{
		Message: "Password changed successfully",
	}, nil
}

// GetProfile handles getting the user's profile
func (h *UserHandler) GetProfile(ctx context.Context, req *pb.GetProfileRequest) (*pb.GetProfileResponse, error) {

	// Call usecase
	profile, err := h.userUsecase.GetUserProfile(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get profile: %v", err)
	}

	// Convert entity to proto
	user := convertUserToProto(profile.User)

	// Convert profile details
	profileDetails := &pb.ProfileDetails{
		Bio:      profile.Bio,
		Location: profile.Location,
		Website:  profile.Website,
		JoinedAt: timestamppb.New(profile.JoinedAt),
	}

	return &pb.GetProfileResponse{
		User:           user,
		ProfileDetails: profileDetails,
	}, nil
}

// UpdateProfile handles updating the user's profile
func (h *UserHandler) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.UpdateProfileResponse, error) {

	// Get current profile
	currentProfile, err := h.userUsecase.GetUserProfile(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get current profile: %v", err)
	}

	// Update profile fields if provided
	if req.FullName != "" {
		currentProfile.User.FullName = req.FullName
	}
	if req.Bio != "" {
		currentProfile.Bio = req.Bio
	}
	if req.Phone != "" {
		currentProfile.User.Phone = req.Phone
	}
	if req.Location != "" {
		currentProfile.Location = req.Location
	}
	if req.Website != "" {
		currentProfile.Website = req.Website
	}

	// Call usecase
	updatedProfile, err := h.userUsecase.UpdateUserProfile(ctx, currentProfile, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update profile: %v", err)
	}

	// Convert entity to proto
	user := convertUserToProto(updatedProfile.User)

	// Convert profile details
	profileDetails := &pb.ProfileDetails{
		Bio:      updatedProfile.Bio,
		Location: updatedProfile.Location,
		Website:  updatedProfile.Website,
		JoinedAt: timestamppb.New(updatedProfile.JoinedAt),
	}

	return &pb.UpdateProfileResponse{
		User:           user,
		ProfileDetails: profileDetails,
	}, nil
}

// GetSessions handles getting all active sessions for the user
func (h *UserHandler) GetSessions(ctx context.Context, req *pb.GetSessionsRequest) (*pb.GetSessionsResponse, error) {

	// Call usecase
	sessions, err := h.userUsecase.GetSessions(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get sessions: %v", err)
	}

	// Get current session ID (if available)
	var currentSessionID string
	if sessionID, ok := ctx.Value("session_id").(string); ok {
		currentSessionID = sessionID
	}

	// Convert entity to proto
	var sessionInfos []*pb.SessionInfo
	for _, session := range sessions {
		isCurrent := session.SessionID.String() == currentSessionID
		sessionInfos = append(sessionInfos, &pb.SessionInfo{
			SessionId:    session.SessionID.String(),
			DeviceInfo:   fmt.Sprintf("%s", session.DeviceInfo),
			IpAddress:    session.IpAddress,
			UserAgent:    session.UserAgent,
			LastActivity: timestamppb.New(session.LastActivity),
			IsCurrent:    isCurrent,
		})
	}

	return &pb.GetSessionsResponse{
		Sessions: sessionInfos,
	}, nil
}

// RevokeSession handles revoking a specific session
func (h *UserHandler) RevokeSession(ctx context.Context, req *pb.RevokeSessionRequest) (*pb.RevokeSessionResponse, error) {

	// Validate request
	if req.SessionId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "session ID is required")
	}

	// Call usecase
	err := h.userUsecase.RevokeSession(ctx, req.SessionId, req.UserId)
	if err != nil {
		// Check for specific errors
		if err.Error() == "unauthorized: session does not belong to the user" {
			return nil, status.Errorf(codes.PermissionDenied, "unauthorized: session does not belong to the user")
		}
		return nil, status.Errorf(codes.Internal, "failed to revoke session: %v", err)
	}

	return &pb.RevokeSessionResponse{
		Message: "Session revoked successfully",
	}, nil
}

// DeactivateAccount handles temporarily deactivating a user account
func (h *UserHandler) DeactivateAccount(ctx context.Context, req *pb.DeactivateAccountRequest) (*pb.DeactivateAccountResponse, error) {

	// Validate request
	if req.Password == "" {
		return nil, status.Errorf(codes.InvalidArgument, "password is required")
	}

	// Call usecase
	err := h.userUsecase.DeactivateAccount(ctx, req.Password, req.UserId)
	if err != nil {
		// Check for specific errors
		if err.Error() == "password is incorrect" {
			return nil, status.Errorf(codes.InvalidArgument, "password is incorrect")
		}
		return nil, status.Errorf(codes.Internal, "failed to deactivate account: %v", err)
	}

	return &pb.DeactivateAccountResponse{
		Message: "Account deactivated successfully",
	}, nil
}

// DeleteAccount handles permanently deleting a user account
func (h *UserHandler) DeleteAccount(ctx context.Context, req *pb.DeleteAccountRequest) (*pb.DeleteAccountResponse, error) {

	// Validate request
	if req.Password == "" {
		return nil, status.Errorf(codes.InvalidArgument, "password is required")
	}

	// Call usecase
	err := h.userUsecase.DeleteAccount(ctx, req.Password, req.UserId)
	if err != nil {
		// Check for specific errors
		if err.Error() == "password is incorrect" {
			return nil, status.Errorf(codes.InvalidArgument, "password is incorrect")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete account: %v", err)
	}

	return &pb.DeleteAccountResponse{
		Message: "Account deleted successfully",
	}, nil
}

// GetLoginHistory handles getting the login history for the user
func (h *UserHandler) GetLoginHistory(ctx context.Context, req *pb.GetLoginHistoryRequest) (*pb.GetLoginHistoryResponse, error) {

	// Set default limit if not provided
	limit := int(req.Limit)
	if limit <= 0 {
		limit = 10
	}

	// Call usecase
	history, err := h.userUsecase.GetLoginHistory(ctx, req.UserId, limit)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get login history: %v", err)
	}

	// Convert entity to proto
	var historyEntries []*pb.LoginHistoryEntry
	for _, entry := range history {
		historyEntries = append(historyEntries, &pb.LoginHistoryEntry{
			IpAddress: entry.IpAddress,
			UserAgent: entry.UserAgent,
			Location:  entry.Location,
		})
	}

	return &pb.GetLoginHistoryResponse{
		History: historyEntries,
	}, nil
}
