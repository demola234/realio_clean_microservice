package middleware

import (
	"errors"
	"fmt"

	token "github.com/demola234/api_gateway/infrastructure/middleware/token_maker"
	interfaces "github.com/demola234/authentication/infrastructure/error"

	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader     = "authorization"
	authorizationBearer     = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func authMiddleWare(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		authorizationHeader := ctx.GetHeader(authorizationHeader)
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header not found")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, err)
			return
		}

		stringSplit := strings.Fields(authorizationHeader)
		if len(stringSplit) < 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, interfaces.ErrorResponse(err, http.StatusUnauthorized))
			return
		}

		authType := strings.ToLower(stringSplit[0])
		if authType != authorizationBearer {
			err := fmt.Errorf("unsupported authorization type %s", authType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, interfaces.ErrorResponse(err, http.StatusUnauthorized))
			return
		}

		accessToken := stringSplit[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, interfaces.ErrorResponse(err, http.StatusUnauthorized))
			return
		}

		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}
