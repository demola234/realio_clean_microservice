package error_response

import (
	"errors"

	"github.com/gin-gonic/gin"
)

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
