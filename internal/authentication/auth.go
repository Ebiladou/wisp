package authentication

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Ebiladou/wisp/internal/middleware"
	"github.com/Ebiladou/wisp/internal/models"
)

func RequireUser() gin.HandlerFunc {

	return func(context *gin.Context) {

		_, exists := context.Get(
			middleware.AuthenticatedUserKey,
		)

		if !exists {
			context.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authenticated user is required.",
			})

			context.Abort()
			return
		}

		context.Next()
	}
}

func RequireActiveUser() gin.HandlerFunc {

	return func(context *gin.Context) {

		value, exists := context.Get(
			middleware.AuthenticatedUserKey,
		)

		if !exists {
			context.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authenticated user is required.",
			})

			context.Abort()
			return
		}

		user, ok := value.(*models.User)

		if !ok {
			context.JSON(http.StatusInternalServerError, gin.H{
				"error": "Invalid authenticated user.",
			})

			context.Abort()
			return
		}

		if !user.Active {
			context.JSON(http.StatusForbidden, gin.H{
				"error": "User account is not active.",
			})

			context.Abort()
			return
		}

		if user.DeletionRequested {
			context.JSON(http.StatusForbidden, gin.H{
				"error": "User account has been deleted.",
			})

			context.Abort()
			return
		}

		context.Next()
	}
}
