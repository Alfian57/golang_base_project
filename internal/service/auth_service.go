package service

import (
	"context"
	"net/http"
	"time"

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

func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (dto.Credentials, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	credentials := dto.Credentials{}

	user, err := s.userRepository.GetByUsername(ctx, req.Username)
	if err != nil {
		if err == errs.ErrUserNotFound {
			return credentials, errs.NewAppError(http.StatusUnauthorized, "username or password is incorrect", err)
		}
		return credentials, errs.NewAppError(http.StatusInternalServerError, "failed to get user", err)
	}

	if err := user.CheckHashedPassword(req.Password); err != nil {
		return credentials, errs.NewAppError(http.StatusUnauthorized, "username or password is incorrect", err)
	}

	accessToken, err := jwt.CreateAccessToken(user)
	if err != nil {
		return credentials, errs.NewAppError(http.StatusInternalServerError, "failed to create access token", err)
	}

	refreshToken, err := jwt.CreateRefreshToken(user)
	if err != nil {
		return credentials, errs.NewAppError(http.StatusInternalServerError, "failed to create refresh token", err)
	}

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

func (s *AuthService) Register(ctx context.Context, request dto.RegisterRequest) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := s.userRepository.GetByUsername(ctx, request.Username)
	if err != nil && err != errs.ErrUserNotFound {
		logger.Log.Errorw("failed to check existing username", "username", request.Username, "error", err)
		return errs.NewAppError(500, "failed to validate username", err)
	}

	if err == nil {
		logger.Log.Infow("username already exists", "username", request.Username)
		fieldError := errs.NewFieldError("username", "username already exists")
		return errs.NewValidationError([]errs.FieldError{fieldError})
	}

	user := model.User{
		Username: request.Username,
	}
	err = user.SetHashedPassword(request.Password)
	if err != nil {
		logger.Log.Errorw("failed to hash password", "error", err)
		return errs.NewAppError(500, "failed to process password", err)
	}

	if err := s.userRepository.Create(ctx, &user); err != nil {
		logger.Log.Errorw("failed to create user", "username", request.Username, "error", err)
		return errs.NewAppError(500, "failed to create user", err)
	}

	logger.Log.Infow("user registered successfully", "username", request.Username)
	return nil
}

func (s *AuthService) Refresh(ctx context.Context, refreshTokenParam string) (dto.Credentials, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	credentials := dto.Credentials{}

	refreshToken, err := s.refreshTokenRepository.GetByTokenHash(ctx, refreshTokenParam)
	if err != nil {
		if err == errs.ErrRefreshTokenNotFound {
			return credentials, errs.NewAppError(http.StatusUnauthorized, "refresh token not valid", err)
		}
		return credentials, errs.NewAppError(http.StatusInternalServerError, "failed to get refresh token", err)
	}

	user, err := s.userRepository.GetByID(ctx, refreshToken.UserID.String())
	if err != nil {
		return credentials, errs.NewAppError(http.StatusInternalServerError, "failed to get user", err)
	}

	err = s.refreshTokenRepository.DeleteByTokenHash(ctx, refreshTokenParam)
	if err != nil {
		return credentials, errs.NewAppError(http.StatusInternalServerError, "failed to delete refresh token", err)
	}

	newAccessToken, err := jwt.CreateAccessToken(user)
	if err != nil {
		return credentials, errs.NewAppError(http.StatusInternalServerError, "failed to create access token", err)
	}

	newRefreshToken, err := jwt.CreateRefreshToken(user)
	if err != nil {
		return credentials, errs.NewAppError(http.StatusInternalServerError, "failed to create refresh token", err)
	}

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

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := s.refreshTokenRepository.DeleteByTokenHash(ctx, refreshToken)

	return err
}
