package service

import (
	"context"
	"time"

	"github.com/Alfian57/belajar-golang/internal/constants"
	"github.com/Alfian57/belajar-golang/internal/dto"
	errs "github.com/Alfian57/belajar-golang/internal/errors"
	"github.com/Alfian57/belajar-golang/internal/logger"
	"github.com/Alfian57/belajar-golang/internal/model"
	"github.com/Alfian57/belajar-golang/internal/repository"
	"github.com/google/uuid"
)

type UserService struct {
	userRepository *repository.UserRepository
}

func NewUserService(r *repository.UserRepository) *UserService {
	return &UserService{
		userRepository: r,
	}
}

// GetAllUsers retrieves all users with optional filtering and pagination.
// It returns a paginated result containing user data.
func (s *UserService) GetAllUsers(ctx context.Context, query dto.GetUsersFilter) (dto.PaginatedResult[model.User], error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Validate pagination parameters
	users, err := s.userRepository.GetAll(ctx, query)
	if err != nil {
		logger.Log.Errorw("failed to retrieve users", "error", err)
		return dto.PaginatedResult[model.User]{}, errs.NewAppError(500, "failed to retrieve users", err)
	}

	// Count total users for pagination
	count, err := s.userRepository.CountAll(ctx, query)
	if err != nil {
		logger.Log.Errorw("failed to count users", "error", err)
		return dto.PaginatedResult[model.User]{}, errs.NewAppError(500, "failed to retrieve users", err)
	}

	// Create pagination response
	pagination := dto.NewPaginationResponse(query.Page, query.Limit, count)
	result := dto.PaginatedResult[model.User]{
		Data:       users,
		Pagination: pagination,
	}

	return result, nil
}

// CreateUser creates a new user with the provided request data.
// It checks for existing usernames and hashes the password before saving.
func (s *UserService) CreateUser(ctx context.Context, request dto.CreateUserRequest) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Check if email already exists
	_, err := s.userRepository.GetByEmail(ctx, request.Email)
	if err != nil && err != errs.ErrUserNotFound {
		logger.Log.Errorw("failed to check existing email", "email", request.Email, "error", err)
		return errs.NewAppError(500, "failed to validate email", err)
	}
	if err == nil {
		logger.Log.Infow("email already exists", "email", request.Email)
		fieldError := errs.NewFieldError("email", "email already exists")
		return errs.NewValidationError([]errs.FieldError{fieldError})
	}

	// Check if username already exists
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

	user := model.User{
		Email:    request.Email,
		Username: request.Username,
		Role:     constants.UserRoleMember, // Default role
	}

	// Password processing
	err = user.SetHashedPassword(request.Password)
	if err != nil {
		logger.Log.Errorw("failed to hash password", "error", err)
		return errs.NewAppError(500, "failed to process password", err)
	}

	// Set user email and role
	if err := s.userRepository.Create(ctx, &user); err != nil {
		logger.Log.Errorw("failed to create user", "username", request.Username, "error", err)
		return errs.NewAppError(500, "failed to create user", err)
	}

	logger.Log.Infow("user created successfully", "username", request.Username)
	return nil
}

// GetUserByID retrieves a user by their ID.
func (s *UserService) GetUserByID(ctx context.Context, id string) (model.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Check if the user exists
	user, err := s.userRepository.GetByID(ctx, id)
	if err != nil {
		if err == errs.ErrUserNotFound {
			return model.User{}, err
		}
		logger.Log.Errorw("failed to get user by ID", "id", id, "error", err)
		return model.User{}, errs.NewAppError(500, "failed to retrieve user", err)
	}
	return user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, request dto.UpdateUserRequest) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Validate user existence
	_, err := s.userRepository.GetByID(ctx, request.ID.String())
	if err != nil {
		if err == errs.ErrUserNotFound {
			return err
		}

		logger.Log.Errorw("failed to check user existence for update", "id", request.ID, "error", err)
		return errs.NewAppError(500, "failed to validate user", err)
	}

	// Check if email already exists
	existingUser, err := s.userRepository.GetByEmail(ctx, request.Email)
	if err != nil && err != errs.ErrUserNotFound {
		logger.Log.Errorw("failed to check email availability", "email", request.Email, "error", err)
		return errs.NewAppError(500, "failed to validate email", err)
	}
	if err == nil && existingUser.ID != request.ID {
		fieldError := errs.NewFieldError("email", "email already exists")
		return errs.NewValidationError([]errs.FieldError{fieldError})
	}

	// Check if username already exists
	existingUser, err = s.userRepository.GetByUsername(ctx, request.Username)
	if err != nil && err != errs.ErrUserNotFound {
		logger.Log.Errorw("failed to check username availability", "username", request.Username, "error", err)
		return errs.NewAppError(500, "failed to validate username", err)
	}
	if err == nil && existingUser.ID != request.ID {
		fieldError := errs.NewFieldError("username", "username already exists")
		return errs.NewValidationError([]errs.FieldError{fieldError})
	}

	// Prepare user data for update
	user := model.User{
		ID:       request.ID,
		Email:    request.Email,
		Username: request.Username,
	}

	// Update user password if provided
	if err := s.userRepository.Update(ctx, &user); err != nil {
		logger.Log.Errorw("failed to update user", "id", request.ID, "error", err)
		return errs.NewAppError(500, "failed to update user", err)
	}

	logger.Log.Infow("user updated successfully", "id", request.ID)
	return nil
}

// DeleteUser deletes a user by their ID.
// It returns an error if the user does not exist or if the deletion fails.
func (s *UserService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Check if the user exists
	if err := s.userRepository.Delete(ctx, id.String()); err != nil {
		if err == errs.ErrUserNotFound {
			return err
		}
		logger.Log.Errorw("failed to delete user", "id", id, "error", err)
		return errs.NewAppError(500, "failed to delete user", err)
	}

	logger.Log.Infow("user deleted successfully", "id", id)
	return nil
}
