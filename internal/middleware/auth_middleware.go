package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/Ebiladou/wisp/internal/config"
	"github.com/Ebiladou/wisp/internal/repositories"
	"github.com/Ebiladou/wisp/internal/utils"
)

const AuthenticatedUserKey = "authenticated_user"

type AuthenticationMiddleware struct {
	applicationConfig *config.Config
	userRepository    repositories.UserRepository
}

func NewAuthenticationMiddleware(
	applicationConfig *config.Config,
	userRepository repositories.UserRepository,
) *AuthenticationMiddleware {

	return &AuthenticationMiddleware{
		applicationConfig: applicationConfig,
		userRepository:    userRepository,
	}
}

func (middleware *AuthenticationMiddleware) Authenticate() gin.HandlerFunc {
	return func(context *gin.Context) {
		accessToken, err := context.Cookie("access_token")

		if err != nil {
			context.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authentication required.",
			})
			context.Abort()
			return
		}

		claims, err := utils.ValidateJWT(
			accessToken,
			middleware.applicationConfig,
		)

		if err != nil {
			context.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired access token.",
			})
			context.Abort()
			return
		}

		if claims.TokenType != utils.AccessTokenType {
			context.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid access token.",
			})
			context.Abort()
			return
		}

		userID, err := uuid.Parse(claims.UserID)

		if err != nil {
			context.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid user ID.",
			})
			context.Abort()
			return
		}

		user, err := middleware.userRepository.FindByID(userID)

		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to retrieve user.",
			})
			context.Abort()
			return
		}

		if user == nil {
			context.JSON(http.StatusUnauthorized, gin.H{
				"error": "User not found.",
			})
			context.Abort()
			return
		}

		context.Set(
			AuthenticatedUserKey,
			user,
		)

		context.Next()
	}
}
