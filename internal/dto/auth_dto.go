package dto

type LoginRequest struct {
	Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

type RegisterRequest struct {
	Email                string `json:"email" form:"email" binding:"required,email,min=3,max=100"`
	Username             string `json:"username" form:"username" binding:"required,min=3,max=100"`
	Password             string `json:"password" form:"password" binding:"required,min=8"`
	PasswordConfirmation string `json:"password_confirmation" form:"password_confirmation" binding:"required,eqfield=Password"`
}

type Credentials struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
