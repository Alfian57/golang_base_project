package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Alfian57/belajar-golang/internal/database"
	"github.com/Alfian57/belajar-golang/internal/dto"
	errs "github.com/Alfian57/belajar-golang/internal/errors"
	"github.com/Alfian57/belajar-golang/internal/logger"
	"github.com/Alfian57/belajar-golang/internal/model"
	"github.com/Alfian57/belajar-golang/internal/utils/queryBuilder"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository() *UserRepository {
	return &UserRepository{db: database.DB}
}

func (r *UserRepository) GetAll(ctx context.Context, queryParam dto.GetUsersFilter) ([]model.User, error) {
	users := []model.User{}
	baseQuery := "SELECT id, email, username, password, role, created_at, updated_at FROM users"

	qb := queryBuilder.NewQueryBuilder(baseQuery)
	qb.Search("username", queryParam.Search).
		OrderBy(queryParam.OrderBy, queryParam.OrderType).
		Paginate(queryParam.PaginationRequest)

	query, args := qb.Build()
	err := r.db.SelectContext(ctx, &users, query, args...)
	return users, err
}

func (r *UserRepository) CountAll(ctx context.Context, queryParam dto.GetUsersFilter) (int64, error) {
	baseQuery := "SELECT COUNT(*) FROM users"

	qb := queryBuilder.NewQueryBuilder(baseQuery)
	qb.Search("username", queryParam.Search)

	query, args := qb.BuildCount(baseQuery)

	var count int64
	err := r.db.GetContext(ctx, &count, query, args...)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	user.ID = uuid.New()

	logger.Log.Info("Creating user", "id", user.ID, "email", user.Email, "username", user.Username)
	query := "INSERT INTO users(id, email, username, password, role) VALUES ($1, $2, $3, $4, $5)"
	_, err := r.db.ExecContext(ctx, query, user.ID.String(), user.Email, user.Username, user.Password, user.Role)
	logger.Log.Debug(err)
	return err
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (model.User, error) {
	user := model.User{}
	query := "SELECT id, email, username, password, role, created_at, updated_at FROM users WHERE id = $1"

	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, errs.ErrUserNotFound
		}
		return user, err
	}

	return user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (model.User, error) {
	user := model.User{}
	query := "SELECT id, email, username, password, role, created_at, updated_at FROM users WHERE email = $1"

	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, errs.ErrUserNotFound
		}
		return user, err
	}

	return user, nil
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (model.User, error) {
	user := model.User{}
	query := "SELECT id, email, username, password, created_at, updated_at FROM users WHERE username = $1"

	err := r.db.GetContext(ctx, &user, query, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, errs.ErrUserNotFound
		}
		return user, err
	}

	return user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
	query := "UPDATE users SET email = $1, username = $2 WHERE id = $3"
	_, err := r.db.ExecContext(ctx, query, user.Email, user.Username, user.ID.String())
	return err
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	query := "DELETE FROM users WHERE id = $1"
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errs.ErrUserNotFound
	}

	return nil
}
