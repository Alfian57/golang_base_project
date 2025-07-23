package service

import (
	"context"
	"net/http"
	"time"

	"github.com/Alfian57/belajar-golang/internal/constants"
	"github.com/Alfian57/belajar-golang/internal/dto"
	errs "github.com/Alfian57/belajar-golang/internal/errors"
	"github.com/Alfian57/belajar-golang/internal/logger"
	"github.com/Alfian57/belajar-golang/internal/model"
	"github.com/Alfian57/belajar-golang/internal/repository"
	"github.com/Alfian57/belajar-golang/internal/utils/jwt"
)

type AuthService struct {
	userRepository         *repository.UserRepository
	refreshTokenRepository *repository.RefreshTokenRepository
}

func NewAuthService(userRepository *repository.UserRepository, refreshTokenRepository *repository.RefreshTokenRepository) *AuthService {
	return &AuthService{
		userRepository:         userRepository,
		refreshTokenRepository: refreshTokenRepository,
	}
}

// Login authenticates a user using username and password.
// It generates access and refresh tokens upon successful authentication.
func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (dto.Credentials, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	credentials := dto.Credentials{}

	// Get user by username
	user, err := s.userRepository.GetByUsername(ctx, req.Username)
	if err != nil {
		if err == errs.ErrUserNotFound {
			return credentials, errs.NewAppError(http.StatusUnauthorized, "username or password is incorrect", err)
		}
		return credentials, errs.NewAppError(http.StatusInternalServerError, "failed to get user", err)
	}

	// Check password
	if err := user.CheckHashedPassword(req.Password); err != nil {
		return credentials, errs.NewAppError(http.StatusUnauthorized, "username or password is incorrect", err)
	}

	// Create access token
	accessToken, err := jwt.CreateAccessToken(user)
	if err != nil {
		return credentials, errs.NewAppError(http.StatusInternalServerError, "failed to create access token", err)
	}

	// Create refresh token
	refreshToken, err := jwt.CreateRefreshToken(user)
	if err != nil {
		return credentials, errs.NewAppError(http.StatusInternalServerError, "failed to create refresh token", err)
	}

	// Save refresh token to repository
	rt := &model.RefreshToken{
		UserID:    user.ID,
		TokenHash: refreshToken,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(),
	}
	if err := s.refreshTokenRepository.Create(ctx, rt); err != nil {
		logger.Log.Info("error bwang", err)
		return credentials, errs.NewAppError(http.StatusInternalServerError, "failed to save refresh token", err)
	}

	credentials.AccessToken = accessToken
	credentials.RefreshToken = refreshToken

	return credentials, nil
}

// Register creates a new user with the provided registration details.
// It checks for unique email and username, hashes the password, and saves the user.
func (s *AuthService) Register(ctx context.Context, request dto.RegisterRequest) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Unique email validation
	_, err := s.userRepository.GetByEmail(ctx, request.Email)
	if err != nil && err != errs.ErrUserNotFound {
		logger.Log.Errorw("failed to check existing email", "email", request.Email, "error", err)
		return errs.NewAppError(500, "failed to validate email", err)
	}
	if err == nil {
		logger.Log.Infow("email already exists", "email", request.Username)
		fieldError := errs.NewFieldError("email", "email already exists")
		return errs.NewValidationError([]errs.FieldError{fieldError})
	}

	// Unique username validation
	_, err = s.userRepository.GetByUsername(ctx, request.Username)
	if err != nil && err != errs.ErrUserNotFound {
		logger.Log.Errorw("failed to check existing username", "username", request.Username, "error", err)
		return errs.NewAppError(500, "failed to validate username", err)
	}
	if err == nil {
		logger.Log.Infow("username already exists", "username", request.Username)
		fieldError := errs.NewFieldError("username", "username already exists")
		return errs.NewValidationError([]errs.FieldError{fieldError})
	}

	// Password processing
	user := model.User{
		Email:    request.Email,
		Username: request.Username,
		Role:     constants.UserRoleMember,
	}
	err = user.SetHashedPassword(request.Password)
	if err != nil {
		logger.Log.Errorw("failed to hash password", "error", err)
		return errs.NewAppError(500, "failed to process password", err)
	}

	// Create user
	if err := s.userRepository.Create(ctx, &user); err != nil {
		logger.Log.Errorw("failed to create user", "username", request.Username, "error", err)
		return errs.NewAppError(500, "failed to create user", err)
	}

	logger.Log.Infow("user registered successfully", "username", request.Username)
	return nil
}

// Refresh generates new access and refresh tokens using a valid refresh token.
// It deletes the old refresh token and saves the new one.
func (s *AuthService) Refresh(ctx context.Context, refreshTokenParam string) (dto.Credentials, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	credentials := dto.Credentials{}

	// Get refresh token from repository
	refreshToken, err := s.refreshTokenRepository.GetByTokenHash(ctx, refreshTokenParam)
	if err != nil {
		if err == errs.ErrRefreshTokenNotFound {
			return credentials, errs.NewAppError(http.StatusUnauthorized, "refresh token not valid", err)
		}
		return credentials, errs.NewAppError(http.StatusInternalServerError, "failed to get refresh token", err)
	}

	// Get user by ID from refresh token
	user, err := s.userRepository.GetByID(ctx, refreshToken.UserID.String())
	if err != nil {
		return credentials, errs.NewAppError(http.StatusInternalServerError, "failed to get user", err)
	}

	// Delete old refresh token
	err = s.refreshTokenRepository.DeleteByTokenHash(ctx, refreshTokenParam)
	if err != nil {
		return credentials, errs.NewAppError(http.StatusInternalServerError, "failed to delete refresh token", err)
	}

	// Create new access token
	newAccessToken, err := jwt.CreateAccessToken(user)
	if err != nil {
		return credentials, errs.NewAppError(http.StatusInternalServerError, "failed to create access token", err)
	}

	// Create new refresh token
	newRefreshToken, err := jwt.CreateRefreshToken(user)
	if err != nil {
		return credentials, errs.NewAppError(http.StatusInternalServerError, "failed to create refresh token", err)
	}

	// Save new refresh token to repository
	rt := &model.RefreshToken{
		UserID:    user.ID,
		TokenHash: newRefreshToken,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(),
	}
	if err := s.refreshTokenRepository.Create(ctx, rt); err != nil {
		return credentials, errs.NewAppError(http.StatusInternalServerError, "failed to save refresh token", err)
	}

	credentials.AccessToken = newAccessToken
	credentials.RefreshToken = newRefreshToken

	return credentials, nil
}

// Logout logs out a user by invalidating the provided refresh token.
func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Delete the refresh token from the repository
	err := s.refreshTokenRepository.DeleteByTokenHash(ctx, refreshToken)

	return err
}
