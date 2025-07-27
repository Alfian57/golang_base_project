package main

import (
	"flag"
	"fmt"

	"github.com/Alfian57/belajar-golang/internal/config"
	"github.com/Alfian57/belajar-golang/internal/database"
	"github.com/Alfian57/belajar-golang/internal/logger"
	"github.com/Alfian57/belajar-golang/internal/seeder"
)

func main() {
	// Command line flags
	var (
		useFactory = flag.Bool("factory", false, "Use factory to generate fake data")
		userCount  = flag.Int("users", 10, "Number of users to create (factory mode)")
	)
	flag.Parse()

	// Load environment variables
	config.LoadEnv()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}

	// Initialize logger
	logger.Init()

	// Initialize database connection
	database.Init(cfg.Database)

	// Create seeder configuration
	seederConfig := seeder.SeederConfig{
		UseFactory: *useFactory,
		UserCount:  *userCount,
	}

	// Initialize and run seeder
	databaseSeeder := seeder.NewDatabaseSeeder(seederConfig)

	if err := databaseSeeder.SeedWithTimeout(); err != nil {
		logger.Log.Fatalf("Database seeding failed: %v", err)
	}

	logger.Log.Info("Database seeding completed successfully!")
}
