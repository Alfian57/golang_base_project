# Go REST API Template

A clean, modular REST API template built with Go, following Go Standard Layout and best practices. This serves as a robust foundation for all future Go projects, featuring modern architecture patterns and production-ready configurations.

## Features

- **Clean Architecture**: Repository, Service, and Handler layers with clear separation of concerns
- **Dependency Injection**: Google Wire for compile-time dependency injection
- **Database**: PostgreSQL with GORM ORM and connection pooling
- **Authentication**: JWT-based authentication with access and refresh tokens
- **Validation**: Custom error messages with struct validation
- **Logging**: Zap structured logging for production-ready logging
- **Graceful Shutdown**: Proper server shutdown handling
- **Migrations**: Database migration support with golang-migrate
- **CORS**: Configurable cross-origin resource sharing
- **Seeding**: Database seeding with factory pattern support
- **Hot Reload**: Development server with Air for automatic reloading
- **Middleware**: Authentication, authorization, and error handling middleware

## Project Structure

```
├── cmd/
│   ├── api/
│   │   └── main.go           # API server entry point
│   └── seeder/
│       └── main.go           # Database seeder entry point
├── internal/
│   ├── config/               # Configuration management
│   ├── constants/            # Application constants
│   ├── database/             # Database connection setup
│   ├── di/                   # Dependency injection (Wire)
│   ├── dto/                  # Data Transfer Objects
│   ├── errors/               # Custom error definitions
│   ├── factory/              # Factory pattern for data generation
│   ├── handler/              # HTTP request handlers
│   ├── logger/               # Structured logging setup
│   ├── middleware/           # HTTP middleware
│   ├── model/                # Database models
│   ├── repository/           # Data access layer
│   ├── response/             # Response utilities
│   ├── router/               # Route definitions
│   ├── seeder/               # Database seeding logic
│   ├── service/              # Business logic layer
│   ├── utils/                # Utility functions
│   │   ├── auth/             # Authentication utilities
│   │   ├── hash/             # Password hashing
│   │   └── jwt/              # JWT token management
│   └── validation/           # Input validation
├── migrations/               # Database migration files
├── logs/                     # Application logs
├── build/                    # Build artifacts
├── tmp/                      # Temporary files (Air)
├── go.mod                    # Go module definition
├── go.sum                    # Go module checksums
├── Makefile                  # Build and development commands
└── .env.example             # Environment variables template
```

## Getting Started

### Prerequisites

- Go 1.24+ (latest stable version recommended)
- PostgreSQL 13+ (with database created)
- Air (optional, for hot reload during development)
- golang-migrate (for database migrations)

### Installation

1. Clone this template repository:

   ```bash
   git clone https://github.com/Alfian57/blinkr.git your-project-name
   cd your-project-name
   ```

2. Update module name in `go.mod`:

   ```bash
   go mod edit -module github.com/yourusername/your-project-name
   ```

3. Install dependencies:

   ```bash
   go mod tidy
   ```

4. Set up environment variables:

   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

5. Create PostgreSQL database and run migrations:

   ```bash
   # Create your database first, then run:
   make migrate-up
   ```

6. (Optional) Seed the database:

   ```bash
   make seed
   # or with factory data:
   make seed-factory
   ```

### Running the Application

#### Development (hot reload):

```bash
make dev
```

#### Production:

```bash
make build
make start
```

## API Endpoints

This template includes the following API endpoints. Customize them according to your project needs.

### Authentication

- `POST /api/v1/register` - Register new user
- `POST /api/v1/login` - User login
- `POST /api/v1/refresh` - Refresh access token
- `POST /api/v1/logout` - User logout

### Users (Protected Routes)

- `GET /api/v1/users` - List users (Admin only)
- `POST /api/v1/users` - Create user (Admin only)
- `GET /api/v1/users/:id` - Get user by ID
- `PUT /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id` - Delete user (Admin only)

## Development

### Available Make Commands

#### Application Commands

```bash
make dev          # Start development server with hot reload
make build        # Build the application
make start        # Start the built application
```

#### Database Migration Commands

```bash
make migrate-create    # Create new migration (interactive)
make migrate-up       # Apply all pending migrations
make migrate-down     # Rollback last migration
make migrate-force    # Force migration to specific version (interactive)
```

#### Database Seeding Commands

```bash
make seed             # Run basic seeder
make seed-factory     # Run factory seeder with sample data
make seed-factory-custom  # Run factory seeder with custom user count
make seed-build       # Build seeder binary
make seed-run         # Run built seeder
make seed-run-factory # Run built seeder with factory data
```

### Database Migrations

Create migration:

```bash
make migrate-create
# Enter migration name when prompted (use underscores)
```

Apply migrations:

```bash
make migrate-up
```

Rollback migrations:

```bash
make migrate-down
```

Force migration version:

```bash
make migrate-force
# Enter version number when prompted
```

### Code Generation

Generate Wire dependencies after adding new dependencies:

```bash
wire gen ./internal/di
```

### Database Seeding

The project includes a flexible seeding system:

1. **Basic Seeding**: Predefined data

   ```bash
   make seed
   ```

2. **Factory Seeding**: Generated fake data
   ```bash
   make seed-factory
   # or with custom count:
   make seed-factory-custom
   ```

## Architecture & Best Practices

This template follows industry best practices and proven patterns:

### Architecture Patterns

- **Clean Architecture**: Clear separation between layers
- **Repository Pattern**: Data access abstraction
- **Service Layer**: Business logic encapsulation
- **Dependency Injection**: Loose coupling with Google Wire

### Code Quality

- **Go Standard Layout**: Organized project structure
- **Interface-based design**: Easy testing and mocking
- **Custom error handling**: Comprehensive error management
- **Input validation**: Robust data validation
- **Structured logging**: Production-ready logging with Zap

### Security & Performance

- **JWT Authentication**: Secure token-based auth with refresh tokens
- **Password Hashing**: bcrypt for secure password storage
- **CORS Configuration**: Configurable cross-origin policies
- **Connection Pooling**: Efficient database connections
- **Graceful Shutdown**: Proper resource cleanup

### Development Experience

- **Hot Reload**: Fast development with Air
- **Database Migrations**: Version-controlled schema changes
- **Factory Pattern**: Easy test data generation
- **Environment Configuration**: Flexible config management

## Environment Variables

Create a `.env` file based on `.env.example`:

### Configuration Notes

- **JWT Secrets**: Use strong, random strings in production
- **GIN_MODE**: Set to `release` for production deployment
- **Database**: Ensure PostgreSQL is running and database exists
- **CORS**: Configure allowed origins according to your frontend setup

## Getting Started with New Projects

When using this template for a new project:

1. **Clone and Rename**:

   ```bash
   git clone https://github.com/Alfian57/blinkr.git your-new-project
   cd your-new-project
   ```

2. **Update Module Reference**:

   ```bash
   # Update go.mod
   go mod edit -module github.com/yourusername/your-new-project

   # Update import paths in all files
   find . -type f -name "*.go" -exec sed -i 's|github.com/Alfian57/belajar-golang|github.com/yourusername/your-new-project|g' {} +
   ```

3. **Customize Project**:

   - Update `README.md` with your project details
   - Modify API endpoints in handlers and routes
   - Add your business logic to services
   - Create new models and migrations as needed
   - Update environment variables

4. **Initialize Git**:
   ```bash
   rm -rf .git
   git init
   git add .
   git commit -m "Initial commit from blinkr template"
   ```

## Contributing

Contributions to improve this template are welcome!

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License. See the LICENSE file for details.

---

## About

**Blinkr** is a production-ready Go REST API template created by [Alfian57](https://github.com/Alfian57). It serves as a solid foundation for building scalable and maintainable Go applications with modern architecture patterns and best practices.

> **Template Usage**: This is a template repository. Use it as a starting point for your Go projects by following the "Getting Started with New Projects" section above.
