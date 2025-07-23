package model

import (
	"time"

	"github.com/Alfian57/belajar-golang/internal/utils/hash"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	Password  string    `json:"-" db:"password"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
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
