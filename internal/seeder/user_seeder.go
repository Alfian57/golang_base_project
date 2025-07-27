package seeder

import (
	"context"

	"github.com/Alfian57/belajar-golang/internal/constants"
	"github.com/Alfian57/belajar-golang/internal/database"
	"github.com/Alfian57/belajar-golang/internal/factory"
	"github.com/Alfian57/belajar-golang/internal/logger"
	"github.com/Alfian57/belajar-golang/internal/model"
	"github.com/google/uuid"

	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
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
	for i := range users {
		// Check if user already exists
		var existingUser model.User
		err := database.DB.WithContext(ctx).Session(&gorm.Session{Logger: gormLogger.Default.LogMode(gormLogger.Silent)}).
			Where("email = ? OR username = ?", users[i].Email, users[i].Username).First(&existingUser).Error
		if err == nil {
			logger.Log.Infof("User %s already exists, skipping...", users[i].Username)
			continue
		}

		if users[i].Password == "" {
			err = users[i].SetHashedPassword(constants.DefaultPassword)
			if err != nil {
				logger.Log.Errorw("Failed to hash password", "user", users[i].Username, "error", err)
				return err
			}
		}

		// Create user
		if err := database.DB.WithContext(ctx).Create(&users[i]).Error; err != nil {
			logger.Log.Errorw("Failed to create user", "user", users[i].Username, "error", err)
			return err
		}
	}

	return nil
}
