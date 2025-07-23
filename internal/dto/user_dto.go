package dto

import "github.com/google/uuid"

type CreateUserRequest struct {
	Email                string `json:"email" form:"email" binding:"required,min=3,max=100,email"`
	Username             string `json:"username" form:"username" binding:"required,min=3,max=100"`
	Password             string `json:"password" form:"password" binding:"required,min=8"`
	PasswordConfirmation string `json:"password_confirmation" form:"password_confirmation" binding:"required,eqfield=Password"`
}

type UpdateUserRequest struct {
	ID       uuid.UUID `json:"id" form:"id"`
	Email    string    `json:"email" form:"email" binding:"required,min=3,max=100,email"`
	Username string    `json:"username" form:"username" binding:"required,min=3"`
}

type GetUsersFilter struct {
	PaginationRequest
	Search    string `json:"search" form:"search" binding:"omitempty,max=255"`
	OrderBy   string `json:"order_by" form:"order_by" binding:"omitempty,oneof=username created_at"`
	OrderType string `json:"order_type" form:"order_type" binding:"omitempty,oneof=ASC DESC asc desc"`
}
