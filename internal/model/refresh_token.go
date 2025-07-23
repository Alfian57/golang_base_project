package model

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	TokenHash string    `json:"token_hash" db:"token_hash"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	ExpiresAt int64     `json:"expires_at" db:"expires_at"`
}
