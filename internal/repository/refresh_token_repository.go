package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Alfian57/belajar-golang/internal/database"
	errs "github.com/Alfian57/belajar-golang/internal/errors"
	"github.com/Alfian57/belajar-golang/internal/model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type RefreshTokenRepository struct {
	db *sqlx.DB
}

func NewRefreshTokenRepository() *RefreshTokenRepository {
	return &RefreshTokenRepository{db: database.DB}
}

func (r *RefreshTokenRepository) Create(ctx context.Context, refreshToken *model.RefreshToken) error {
	query := "INSERT INTO refresh_tokens(id, token_hash, user_id, expires_at) VALUES (?, ?, ?, ?)"

	_, err := r.db.ExecContext(ctx, query, uuid.New().String(), refreshToken.TokenHash, refreshToken.UserID, refreshToken.ExpiresAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *RefreshTokenRepository) GetByTokenHash(ctx context.Context, token string) (model.RefreshToken, error) {
	refreshToken := model.RefreshToken{}
	query := "SELECT id, token_hash, user_id, created_at, expires_at FROM refresh_tokens WHERE token_hash = ?"

	err := r.db.GetContext(ctx, &refreshToken, query, token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return refreshToken, errs.ErrTodoNotFound
		}
		return refreshToken, err
	}

	return refreshToken, nil
}

func (r *RefreshTokenRepository) DeleteByTokenHash(ctx context.Context, refreshToken string) error {
	query := "DELETE FROM refresh_tokens WHERE token_hash = ?"
	result, err := r.db.ExecContext(ctx, query, refreshToken)

	rowsAffected, err := result.RowsAffected()

	if rowsAffected == 0 {
		return errs.ErrTodoNotFound
	}

	return err
}
