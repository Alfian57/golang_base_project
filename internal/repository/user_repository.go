package repository

import (
	"context"
	"errors"

	"github.com/Alfian57/belajar-golang/internal/database"
	"github.com/Alfian57/belajar-golang/internal/dto"
	errs "github.com/Alfian57/belajar-golang/internal/errors"
	"github.com/Alfian57/belajar-golang/internal/logger"
	"github.com/Alfian57/belajar-golang/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository() *UserRepository {
	return &UserRepository{db: database.DB}
}

func (r *UserRepository) GetAllWithPagination(ctx context.Context, queryParam dto.GetUsersFilter) ([]model.User, error) {
	var users []model.User

	query := r.db.WithContext(ctx)

	// Apply search filter
	if queryParam.Search != "" {
		query = query.Where("username LIKE ?", "%"+queryParam.Search+"%")
	}

	// Apply ordering
	orderBy := queryParam.OrderBy
	if orderBy == "" {
		orderBy = "created_at"
	}
	orderType := queryParam.OrderType
	if orderType != "ASC" && orderType != "DESC" {
		orderType = "ASC"
	}
	query = query.Order(orderBy + " " + orderType)

	// Apply pagination
	query = query.Limit(queryParam.PaginationRequest.Limit).Offset(queryParam.PaginationRequest.GetOffset())

	err := query.Find(&users).Error
	return users, err
}

func (r *UserRepository) CountAll(ctx context.Context, queryParam dto.GetUsersFilter) (int64, error) {
	var count int64

	query := r.db.WithContext(ctx).Model(&model.User{})

	// Apply search filter
	if queryParam.Search != "" {
		query = query.Where("username LIKE ?", "%"+queryParam.Search+"%")
	}

	err := query.Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	user.ID = uuid.New()

	logger.Log.Info("Creating user", "id", user.ID, "email", user.Email, "username", user.Username)
	err := r.db.WithContext(ctx).Create(user).Error
	logger.Log.Debug(err)
	return err
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (model.User, error) {
	var user model.User

	err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, errs.ErrUserNotFound
		}
		return user, err
	}

	return user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (model.User, error) {
	var user model.User

	err := r.db.WithContext(ctx).First(&user, "email = ?", email).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, errs.ErrUserNotFound
		}
		return user, err
	}

	return user, nil
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (model.User, error) {
	var user model.User

	err := r.db.WithContext(ctx).First(&user, "username = ?", username).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, errs.ErrUserNotFound
		}
		return user, err
	}

	return user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
	err := r.db.WithContext(ctx).Model(user).Select("email", "username").Updates(user).Error
	return err
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&model.User{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errs.ErrUserNotFound
	}

	return nil
}
