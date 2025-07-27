package model

import (
	"time"

	"github.com/Alfian57/belajar-golang/internal/utils/hash"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	Email     string    `json:"email" gorm:"uniqueIndex;not null"`
	Username  string    `json:"username" gorm:"uniqueIndex;not null"`
	Password  string    `json:"-" gorm:"not null"`
	Role      string    `json:"role" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) SetHashedPassword(password string) error {
	hashedPass, err := hash.HashPassword(password)
	if err != nil {
		return err
	}

	u.Password = hashedPass

	return nil
}

func (u *User) CheckHashedPassword(password string) error {
	err := hash.CheckPasswordHash(password, u.Password)
	return err
}
