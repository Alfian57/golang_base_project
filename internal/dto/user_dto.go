package dto

import "github.com/google/uuid"

type CreateUserRequest struct {
	Username             string `form:"username" binding:"required,min=3,max=100"`
	Password             string `form:"password" binding:"required,min=8"`
	PasswordConfirmation string `form:"password_confirmation" binding:"required,eqfield=Password"`
}

type UpdateUserRequest struct {
	ID       uuid.UUID `form:"id"`
	Username string    `form:"username" binding:"required,min=3"`
}

type GetUsersFilter struct {
	PaginationRequest
	Search    string `json:"search" form:"search" binding:"omitempty,max=255"`
	OrderBy   string `json:"order_by" form:"order_by" binding:"omitempty,oneof=username created_at"`
	OrderType string `json:"order_type" form:"order_type" binding:"omitempty,oneof=ASC DESC asc desc"`
}
