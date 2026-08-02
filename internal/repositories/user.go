package repositories

import (
	"github.com/Ebiladou/wisp/internal/models"
)

type UserRepository interface {
	Create(user *models.User) error
}
