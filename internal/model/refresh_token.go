package model

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	TokenHash string    `json:"token_hash" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	ExpiresAt int64     `json:"expires_at" gorm:"not null"`
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}
