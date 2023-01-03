package api

import (
	"errors"
	"github/dutt23/bank/token"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey = "authorization"
	authTypeBearer         = "Bearer"
	authPayloadKey         = "authotization_payload"
)

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader(authorizationHeaderKey)

		if len(authHeader) == 0 {
			err := errors.New("authorization header not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		fields := strings.Fields(authHeader)

		if len(fields) > 2 {
			err := errors.New("invalid auth format supplied")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		authType := fields[0]

		if authType != authTypeBearer {
			err := errors.New("auth type not supported by the server")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.Validate(accessToken)

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
		}

		ctx.Set(authPayloadKey, payload)
		ctx.Next()
	}
}
