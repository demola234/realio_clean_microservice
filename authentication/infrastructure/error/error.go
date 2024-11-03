package error

import (
	"errors"

	"github.com/gin-gonic/gin"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func FieldViolation(field string, err error) *errdetails.BadRequest_FieldViolation {
	return &errdetails.BadRequest_FieldViolation{
		Field:       field,
		Description: err.Error(),
	}
}

func InvalidArgErr(violations []*errdetails.BadRequest_FieldViolation) error {
	badRequest := &errdetails.BadRequest{FieldViolations: violations}
	statusInvalid := status.New(codes.InvalidArgument, "invalid request")
	statusDetails, err := statusInvalid.WithDetails(badRequest)
	if err != nil {
		return statusInvalid.Err()
	}
	return statusDetails.Err()
}

func AuthorizationError(err error) error {
	return status.Errorf(codes.Unauthenticated, "unauthorized: %s", err)

}

func ErrorResponse(err error, status int) gin.H {
	return gin.H{
		"status":  status,
		"message": err.Error(),
	}
}

var ErrInvalidCredentials = errors.New("invalid-credentials")
var ErrUnauthorized = errors.New("unauthorized")
var ErrUserAlreadyExists = errors.New("user-already-exists")
var ErrUserAlreadyVerified = errors.New("user-already-verified")
var ErrOTPCreationFailed = errors.New("token-creation-failed")
var ErrInvalidOTP = errors.New("invalid-otp")
var ErrOTPUsed = errors.New("otp-used")
var ErrOTPExpired = errors.New("otp-expired")
var UnknownError = errors.New("unknown-error")
var ErrInvalidRequest = errors.New("invalid-request")
var ErrPasswordMismatch = errors.New("password-mismatch")
var ErrInvalidJWT = errors.New("invalid-jwt")
var ErrAccountInactive = errors.New("account-inactive")
var ErrUsernameAlreadyExists = errors.New("username-already-exists")
var ErrBadRequest = errors.New("bad-request")
