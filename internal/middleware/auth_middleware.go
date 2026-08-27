package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/Ebiladou/wisp/internal/config"
	"github.com/Ebiladou/wisp/internal/repositories"
	"github.com/Ebiladou/wisp/internal/utils"
)

func AuthMiddleware(
	authRepository repositories.AuthRepository,
	applicationConfig *config.Config,
) gin.HandlerFunc {

	return func(context *gin.Context) {

		accessToken, err := context.Cookie("access_token")

		if err != nil {
			context.JSON(http.StatusUnauthorized, gin.H{
				"error": "authentication required",
			})
			context.Abort()
			return
		}

		claims, err := utils.ValidateJWT(
			accessToken,
			applicationConfig,
		)

		if err != nil {
			context.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired token",
			})
			context.Abort()
			return
		}

		if claims.TokenType != utils.AccessTokenType {
			context.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid access token",
			})
			context.Abort()
			return
		}

		userID, err := uuid.Parse(claims.UserID)

		if err != nil {
			context.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authentication token",
			})
			context.Abort()
			return
		}

		user, err := authRepository.FindByID(userID)

		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to authenticate user",
			})
			context.Abort()
			return
		}

		if user == nil {
			context.JSON(http.StatusUnauthorized, gin.H{
				"error": "user not found",
			})
			context.Abort()
			return
		}

		context.Set("user", user)

		context.Next()
	}
}
