package api

import (
	"blogapi/pkg/token"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

// AuthMiddleware creates a gin middleware for authorization
func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		cookie, err := ctx.Cookie("gin_cookie")

		if err != nil {
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
				return
			}	
		}

		fmt.Printf("Cookie value: %s \n", cookie)

		// authorizationHeader := ctx.GetHeader(authorizationHeaderKey)

		// if len(authorizationHeader) == 0 {
		// 	err := errors.New("authorization header is not provided")
		// 	ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
		// 	return
		// }

		// fields := strings.Fields(authorizationHeader)
		// if len(fields) < 2 {
		// 	err := errors.New("invalid authorization header format")
		// 	ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
		// 	return
		// }

		// authorizationType := strings.ToLower(fields[0])
		// if authorizationType != authorizationTypeBearer {
		// 	err := fmt.Errorf("unsupported authorization type %s", authorizationType)
		// 	ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
		// 	return
		// }

		// accessToken := fields[1]
		// payload, err := tokenMaker.VerifyToken(accessToken)
		// if err != nil {
		// 	ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
		// 	return
		// }

		payload, err := tokenMaker.VerifyToken(cookie)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}
