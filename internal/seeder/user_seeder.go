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
	useFactory bool
	count      int
}

func NewUserSeeder(useFactory bool, count int) *UserSeeder {
	return &UserSeeder{
		useFactory: useFactory,
		count:      count,
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

	return s.createUsers(ctx, users)
}

func (s *UserSeeder) seedWithFactory(ctx context.Context) error {
	logger.Log.Infof("Seeding %d users with factory...", s.count)

	var users []model.User

	// Create 1 admin user using AdminUserFactory
	adminUser := factory.AdminUserFactory.MustCreate().(*model.User)
	users = append(users, *adminUser)

	// Create remaining users using regular UserFactory
	for i := 1; i < s.count; i++ {
		user := factory.UserFactory.MustCreate().(*model.User)
		users = append(users, *user)
	}

	return s.createUsers(ctx, users)
}

func (s *UserSeeder) createUsers(ctx context.Context, users []model.User) error {
	userRepisitory := repository.NewUserRepository()

	for i := range users {
		var err error

		// Check if user already exists by username
		_, err = userRepisitory.GetByUsername(ctx, users[i].Username)
		if err != nil {
			if !errors.Is(err, errs.ErrUserNotFound) {
				logger.Log.Errorw("Failed to check existing user", "username", users[i].Username, "error", err)
				return err
			} else {
				logger.Log.Infof("User %s already exists, skipping...", users[i].Username)
				continue
			}
		}

		// Check if user already exists by email
		_, err = userRepisitory.GetByEmail(ctx, users[i].Email)
		if err != nil {
			if !errors.Is(err, errs.ErrUserNotFound) {
				logger.Log.Errorw("Failed to check existing user", "email", users[i].Email, "error", err)
				return err
			} else {
				logger.Log.Infof("User %s already exists, skipping...", users[i].Email)
				continue
			}
		}

		// Create user
		if err := userRepisitory.Create(ctx, &users[i]); err != nil {
			logger.Log.Errorw("Failed to create user", "user", users[i].Username, "error", err)
			return err
		}
	}

	return nil
}
