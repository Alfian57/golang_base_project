# Belajar Golang - REST API

A clean, modular REST API built with Go following the Go Standard Layout and best practices.

## Features

- **Clean Architecture**: Organized with Repository, Service, and Handler layers
- **Dependency Injection**: Using Google Wire for dependency injection
- **Database**: MySQL with connection pooling
- **Authentication**: JWT-based authentication
- **Validation**: Request validation with custom error messages
- **Logging**: Structured logging with Zap
- **Graceful Shutdown**: Proper server shutdown handling
- **Migrations**: Database migrations support
- **CORS**: Cross-origin resource sharing support

## Project Structure

```
├── cmd/
│   └── api/
│       └── main.go           # Application entry point
├── internal/
│   ├── config/               # Configuration management
│   ├── constants/            # Application constants
│   ├── database/             # Database connection
│   ├── di/                   # Dependency injection (Wire)
│   ├── dto/                  # Data Transfer Objects
│   ├── errors/               # Custom error types
│   ├── handler/              # HTTP handlers
│   ├── logger/               # Logging setup
│   ├── middleware/           # HTTP middleware
│   ├── model/                # Data models
│   ├── repository/           # Database operations
│   ├── response/             # HTTP response utilities
│   ├── router/               # Route definitions
│   ├── service/              # Business logic
│   ├── utils/                # Utility functions
│   └── validation/           # Validation setup
├── migrations/               # Database migrations
├── build/                    # Build artifacts
├── tmp/                      # Temporary files (air)
├── go.mod                    # Go module definition
├── go.sum                    # Go module checksums
└── Makefile                  # Build automation
```

## Getting Started

### Prerequisites

- Go 1.19+
- MySQL 8.0+
- Air (for hot reloading during development)

### Installation

1. Clone the repository:

```bash
git clone https://github.com/Alfian57/belajar-golang.git
cd belajar-golang
```

2. Install dependencies:

```bash
go mod tidy
```

3. Set up your environment variables:

```bash
cp .env.example .env
# Edit .env with your database configuration
```

4. Run database migrations:

```bash
make migrate-up
```

### Running the Application

#### Development (with hot reload):

```bash
make dev
```

#### Production:

```bash
make build
make start
```

## API Endpoints

### Authentication

- `POST /api/v1/register` - Register a new user
- `POST /api/v1/login` - Login user
- `POST /api/v1/refresh` - Refresh access token
- `POST /api/v1/logout` - Logout user

### Users

- `GET /api/v1/users` - Get all users
- `POST /api/v1/users` - Create a new user
- `GET /api/v1/users/:id` - Get user by ID
- `PUT /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id` - Delete user

### Todos

- `GET /api/v1/todos` - Get all todos
- `POST /api/v1/todos` - Create a new todo
- `GET /api/v1/todos/:id` - Get todo by ID
- `PUT /api/v1/todos/:id` - Update todo
- `DELETE /api/v1/todos/:id` - Delete todo

## Development

### Database Migrations

Create a new migration:

```bash
make migrate-create
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
```

### Code Generation

Generate Wire dependencies:

```bash
wire gen ./internal/di
```

## Best Practices Implemented

- **Go Standard Layout**: Following the community standard project layout
- **Clean Architecture**: Separation of concerns with clear boundaries
- **Interface-based Design**: Using interfaces for better testability
- **Error Handling**: Proper error handling with custom error types
- **Validation**: Input validation with meaningful error messages
- **Logging**: Structured logging throughout the application
- **Configuration**: Environment-based configuration management
- **Database**: Connection pooling and proper transaction handling
- **Security**: JWT authentication and CORS configuration
- **Graceful Shutdown**: Proper server shutdown handling

## Environment Variables

```env
# Server Configuration
APP_URL=localhost:8000
TRUSTED_PROXIES=

# Database Configuration
DB_HOST=localhost
DB_PORT=3306
DB_USERNAME=your_username
DB_PASSWORD=your_password
DB_NAME=your_database

# CORS Configuration
CORS_ALLOW_ORIGINS=*
CORS_ALLOW_METHODS=GET,POST,PUT,DELETE,OPTIONS
CORS_ALLOW_CREDENTIALS=true
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
