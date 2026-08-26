package services

import (
	"errors"
	"log"
	"time"

	"github.com/Ebiladou/wisp/internal/dto"
	"github.com/Ebiladou/wisp/internal/models"
	"github.com/Ebiladou/wisp/internal/repositories"
	"github.com/Ebiladou/wisp/internal/utils"
	"github.com/google/uuid"
)

type AuthService interface {
	Create(request dto.CreateUser) (*dto.UserResponse, error)
	GetByID(id uuid.UUID) (*dto.UserResponse, error)
	ConfirmEmail(token string) error
	ResendConfirmation(email string) error
	ForgotPassword(email string) error
	ResetPassword(rawToken string, newPassword string) error
}

type DefaultAuthService struct {
	authRepository  repositories.AuthRepository
	tokenRepository repositories.TokenRepository
}

func NewAuthService(
	authRepository repositories.AuthRepository,
	tokenRepository repositories.TokenRepository,
) AuthService {
	return &DefaultAuthService{
		authRepository:  authRepository,
		tokenRepository: tokenRepository,
	}
}

func (service *DefaultAuthService) Create(request dto.CreateUser) (*dto.UserResponse, error) {

	existingUser, err := service.authRepository.FindByEmail(request.Email)

	if err != nil {
		return nil, err
	}

	if existingUser != nil {
		return nil, errors.New("email already registered")
	}

	existingUser, err = service.authRepository.FindByUsername(request.Username)

	if err != nil {
		return nil, err
	}

	if existingUser != nil {
		return nil, errors.New("username already taken")
	}

	hashedPassword, err := utils.HashPassword(request.Password)

	if err != nil {
		return nil, err
	}

	user := models.User{
		Name:        request.Name,
		Username:    request.Username,
		Email:       request.Email,
		Password:    hashedPassword,
		DateOfBirth: request.DateOfBirth,
		Active:      false,
	}

	err = service.authRepository.CreateUser(&user)

	if err != nil {
		return nil, err
	}

	rawToken, err := utils.GenerateToken()

	if err != nil {
		return nil, err
	}

	hashedToken := utils.HashToken(rawToken)

	token := models.Token{
		UserID:    user.ID,
		Token:     hashedToken,
		TokenType: models.TokenTypeEmailVerification,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	err = service.tokenRepository.Create(&token)

	if err != nil {
		return nil, err
	}

	// for testing, please. we'll move to an email service soon enough when my bag is up.
	log.Printf(
		"email verification token for user %s: %s",
		user.Email,
		rawToken,
	)

	response := &dto.UserResponse{
		ID:             user.ID.String(),
		Name:           user.Name,
		Username:       user.Username,
		Email:          user.Email,
		ProfilePicture: user.ProfilePicture,
		Active:         user.Active,
	}

	return response, nil
}

func (service *DefaultAuthService) GetByID(id uuid.UUID) (*dto.UserResponse, error) {

	user, err := service.authRepository.FindByID(id)

	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	response := &dto.UserResponse{
		ID:             user.ID.String(),
		Name:           user.Name,
		Username:       user.Username,
		Email:          user.Email,
		ProfilePicture: user.ProfilePicture,
		Active:         user.Active,
	}

	return response, nil
}

func (service *DefaultAuthService) ConfirmEmail(rawToken string) error {

	hashedToken := utils.HashToken(rawToken)

	token, err := service.tokenRepository.FindByToken(
		hashedToken,
		models.TokenTypeEmailVerification,
	)

	if err != nil {
		return err
	}

	if token == nil {
		return errors.New("invalid verification token")
	}

	if token.UsedAt != nil {
		return errors.New("verification token has already been used")
	}

	if token.ExpiresAt.Before(time.Now()) {
		return errors.New("verification token has expired")
	}

	user, err := service.authRepository.FindByID(token.UserID)

	if err != nil {
		return err
	}

	if user == nil {
		return errors.New("user not found")
	}

	if user.Active {
		return errors.New("user account is already verified")
	}

	err = service.authRepository.ConfirmUser(user)

	if err != nil {
		return err
	}

	err = service.tokenRepository.MarkAsUsed(token)

	if err != nil {
		return err
	}

	return nil
}

func (service *DefaultAuthService) ResendConfirmation(email string) error {

	user, err := service.authRepository.FindByEmail(email)

	if err != nil {
		return err
	}

	if user == nil {
		return errors.New("user not found")
	}

	if user.Active {
		return errors.New("user account is already verified")
	}

	rawToken, err := utils.GenerateToken()

	if err != nil {
		return err
	}

	hashedToken := utils.HashToken(rawToken)

	token := models.Token{
		UserID:    user.ID,
		Token:     hashedToken,
		TokenType: models.TokenTypeEmailVerification,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	err = service.tokenRepository.Create(&token)

	if err != nil {
		return err
	}

	log.Printf(
		"email verification token for user %s: %s",
		user.Email,
		rawToken,
	)

	return nil
}

func (service *DefaultAuthService) ForgotPassword(email string) error {

	user, err := service.authRepository.FindByEmail(email)

	if err != nil {
		return err
	}

	if user == nil {
		return nil
	}

	rawToken, err := utils.GenerateToken()

	if err != nil {
		return err
	}

	hashedToken := utils.HashToken(rawToken)

	token := models.Token{
		UserID:    user.ID,
		Token:     hashedToken,
		TokenType: models.TokenTypePasswordReset,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	err = service.tokenRepository.Create(&token)

	if err != nil {
		return err
	}

	log.Printf(
		"password reset token for user %s: %s",
		user.Email,
		rawToken,
	)

	return nil
}

func (service *DefaultAuthService) ResetPassword(rawToken string, newPassword string) error {

	hashedToken := utils.HashToken(rawToken)

	token, err := service.tokenRepository.FindByToken(
		hashedToken,
		models.TokenTypePasswordReset,
	)

	if err != nil {
		return err
	}

	if token == nil {
		return errors.New("invalid password reset token")
	}

	if token.UsedAt != nil {
		return errors.New("password reset token has already been used")
	}

	if token.ExpiresAt.Before(time.Now()) {
		return errors.New("password reset token has expired")
	}

	user, err := service.authRepository.FindByID(token.UserID)

	if err != nil {
		return err
	}

	if user == nil {
		return nil
	}

	hashedPassword, err := utils.HashPassword(newPassword)

	if err != nil {
		return err
	}

	user.Password = hashedPassword

	err = service.authRepository.UpdatePassword(user)

	if err != nil {
		return err
	}

	err = service.tokenRepository.MarkAsUsed(token)

	if err != nil {
		return err
	}

	return nil
}
