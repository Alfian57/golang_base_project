package seeder

import (
	"context"
	"errors"

	errs "github.com/Alfian57/belajar-golang/internal/errors"
	"github.com/Alfian57/belajar-golang/internal/factory"
	"github.com/Alfian57/belajar-golang/internal/logger"
	"github.com/Alfian57/belajar-golang/internal/model"
	"github.com/Alfian57/belajar-golang/internal/repository"
	"github.com/google/uuid"
)

type UserSeeder struct {
	userRepository *repository.UserRepository
	useFactory     bool
	count          int
}

func NewUserSeeder(useFactory bool, count int) *UserSeeder {
	return &UserSeeder{
		userRepository: repository.NewUserRepository(),
		useFactory:     useFactory,
		count:          count,
	}
}

func (s *UserSeeder) Seed(ctx context.Context) error {
	logger.Log.Info("Starting user seeding...")

	if s.useFactory {
		return s.seedWithFactory(ctx)
	}
	return s.seedManual(ctx)
}

func (s *UserSeeder) seedManual(ctx context.Context) error {
	logger.Log.Info("Seeding users manually...")

	users := []model.User{
		{
			ID:       uuid.New(),
			Email:    "admin@example.com",
			Username: "admin",
			Role:     model.UserRoleAdmin,
		},
		{
			ID:       uuid.New(),
			Email:    "member@example.com",
			Username: "member",
			Role:     model.UserRoleMember,
		},
		{
			ID:       uuid.New(),
			Email:    "alfian@example.com",
			Username: "alfian57",
			Role:     model.UserRoleMember,
		},
	}

	// Set default password for all users
	for i := range users {
		if err := users[i].SetHashedPassword("password123"); err != nil {
			logger.Log.Errorw("Failed to set password for user", "username", users[i].Username, "error", err)
			return err
		}
	}

	return s.createUsers(ctx, users)
}

func (s *UserSeeder) seedWithFactory(ctx context.Context) error {
	logger.Log.Infof("Seeding %d users with factory...", s.count)

	var users []model.User

	// Create 1 admin user using AdminUserFactory
	adminFactory := factory.NewAdminFactory()
	adminUser := adminFactory.MustCreate().(*model.User)
	users = append(users, *adminUser)

	// Create remaining users using regular UserFactory
	memberFactory := factory.NewMemberFactory()
	for i := 1; i < s.count; i++ {
		member := memberFactory.MustCreate().(*model.User)
		users = append(users, *member)
	}

	return s.createUsers(ctx, users)
}

func (s *UserSeeder) createUsers(ctx context.Context, users []model.User) error {
	for _, user := range users {
		var err error

		// Check if user already exists by username
		_, err = s.userRepository.GetByUsername(ctx, user.Username)
		if err == nil {
			// User found, skip creation
			logger.Log.Infof("User with username %s already exists, skipping...", user.Username)
			continue
		} else if !errors.Is(err, errs.ErrUserNotFound) {
			// Other error occurred
			logger.Log.Errorw("Failed to check existing user", "username", user.Username, "error", err)
			return err
		}

		// Check if user already exists by email
		_, err = s.userRepository.GetByEmail(ctx, user.Email)
		if err == nil {
			// User found, skip creation
			logger.Log.Infof("User with email %s already exists, skipping...", user.Email)
			continue
		} else if !errors.Is(err, errs.ErrUserNotFound) {
			// Other error occurred
			logger.Log.Errorw("Failed to check existing user", "email", user.Email, "error", err)
			return err
		}

		// Create user (both username and email are available)
		if err := s.userRepository.Create(ctx, &user); err != nil {
			logger.Log.Errorw("Failed to create user", "user", user.Username, "error", err)
			return err
		}

		logger.Log.Infof("Successfully created user: %s (%s)", user.Username, user.Email)
	}

	return nil
}
