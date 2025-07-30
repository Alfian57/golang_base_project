package repository

import (
	"context"
	"errors"

	"github.com/Alfian57/belajar-golang/internal/database"
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

// GetAllWithFilterPagination retrieves users with optional filters, ordering, and pagination
func (r *UserRepository) GetAllWithFilterPagination(ctx context.Context, search string, orderBy string, orderType string, limit int, offset int) ([]model.User, error) {
	var users []model.User

	query := r.db.WithContext(ctx)

	// Apply search filter
	if search != "" {
		query = query.Where("username LIKE ?", "%"+search+"%")
	}

	// Apply ordering
	if orderBy != "" && orderType != "" {
		query = query.Order(orderBy + " " + orderType)
	}

	// Apply limit and offset
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&users).Error
	return users, err
}

// GetAll retrieves all users without any filters
func (r *UserRepository) GetAll(ctx context.Context) ([]model.User, error) {
	var users []model.User
	err := r.db.WithContext(ctx).Find(&users).Error
	return users, err
}

// CountWithFilter returns the total number of users matching the search criteria
func (r *UserRepository) CountWithFilter(ctx context.Context, search string) (int64, error) {
	var count int64

	query := r.db.WithContext(ctx).Model(&model.User{})

	// Apply search filter
	if search != "" {
		query = query.Where("username LIKE ?", "%"+search+"%")
	}

	err := query.Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

// Count returns the total number of all users
func (r *UserRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.User{}).Count(&count).Error
	return count, err
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	user.ID = uuid.New()

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
