package repository

import (
	"context"
	"errors"

	"github.com/Alfian57/belajar-golang/internal/database"
	errs "github.com/Alfian57/belajar-golang/internal/errors"
	"github.com/Alfian57/belajar-golang/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RefreshTokenRepository struct {
	db *gorm.DB
}

func NewRefreshTokenRepository() *RefreshTokenRepository {
	return &RefreshTokenRepository{db: database.DB}
}

func (r *RefreshTokenRepository) Create(ctx context.Context, refreshToken *model.RefreshToken) error {
	refreshToken.ID = uuid.New()

	err := r.db.WithContext(ctx).Create(refreshToken).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *RefreshTokenRepository) GetByTokenHash(ctx context.Context, token string) (model.RefreshToken, error) {
	var refreshToken model.RefreshToken

	err := r.db.WithContext(ctx).First(&refreshToken, "token_hash = ?", token).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return refreshToken, errs.ErrTodoNotFound
		}
		return refreshToken, err
	}

	return refreshToken, nil
}

func (r *RefreshTokenRepository) DeleteByTokenHash(ctx context.Context, refreshToken string) error {
	result := r.db.WithContext(ctx).Delete(&model.RefreshToken{}, "token_hash = ?", refreshToken)

	if result.RowsAffected == 0 {
		return errs.ErrTodoNotFound
	}

	return result.Error
}
