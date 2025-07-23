package constants

import "time"

// Database Query Constants
const (
	DefaultPageSize = 10
	MaxPageSize     = 100
	DefaultPage     = 1
)

// Timeout Constants
const (
	DefaultContextTimeout = 5 * time.Second
	DatabaseTimeout       = 10 * time.Second
)

// Validation Constants
const (
	MinUsernameLength = 3
	MaxUsernameLength = 100
	MinPasswordLength = 8
	MinTitleLength    = 3
	MaxTitleLength    = 100
	MaxSearchLength   = 255
)

// Default Values
const (
	DefaultOrderType = "ASC"
	DefaultOrderBy   = "created_at"
)
