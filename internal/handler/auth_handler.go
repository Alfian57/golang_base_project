package handler

import (
	"net/http"

	"github.com/Alfian57/belajar-golang/internal/dto"
	"github.com/Alfian57/belajar-golang/internal/response"
	"github.com/Alfian57/belajar-golang/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(s *service.AuthService) *AuthHandler {
	return &AuthHandler{
		service: s,
	}
}

func (h *AuthHandler) Login(ctx *gin.Context) {
	var request dto.LoginRequest
	if err := ctx.ShouldBind(&request); err != nil {
		response.WriteErrorResponse(ctx, err)
		return
	}

	credentials, err := h.service.Login(ctx, request)
	if err != nil {
		response.WriteErrorResponse(ctx, err)
		return
	}

	accessTokenMaxAge := 15 * 60        // 15 minutes in seconds
	refreshTokenMaxAge := 7 * 24 * 3600 // 7 days in seconds

	ctx.SetCookie("access_token", credentials.AccessToken, accessTokenMaxAge, "/", "", false, true)
	ctx.SetCookie("refresh_token", credentials.RefreshToken, refreshTokenMaxAge, "/", "", false, true)

	response.WriteMessageResponse(ctx, http.StatusOK, "user successfully logged in")
}

func (h *AuthHandler) Register(ctx *gin.Context) {
	var request dto.RegisterRequest
	if err := ctx.ShouldBind(&request); err != nil {
		response.WriteErrorResponse(ctx, err)
		return
	}

	if err := h.service.Register(ctx, request); err != nil {
		response.WriteErrorResponse(ctx, err)
		return
	}

	response.WriteMessageResponse(ctx, http.StatusCreated, "user successfully registered")
}

func (h *AuthHandler) Refresh(ctx *gin.Context) {
	refreshToken, err := ctx.Cookie("refresh_token")
	if err != nil {
		response.WriteErrorResponse(ctx, err)
		return
	}

	credentials, err := h.service.Refresh(ctx, refreshToken)
	if err != nil {
		response.WriteErrorResponse(ctx, err)
		return
	}

	accessTokenMaxAge := 15 * 60        // 15 minutes in seconds
	refreshTokenMaxAge := 7 * 24 * 3600 // 7 days in seconds

	ctx.SetCookie("access_token", credentials.AccessToken, accessTokenMaxAge, "/", "", true, true)
	ctx.SetCookie("refresh_token", credentials.RefreshToken, refreshTokenMaxAge, "/", "", true, true)

	response.WriteMessageResponse(ctx, http.StatusOK, "token successfully refreshed")
}

func (h *AuthHandler) Logout(ctx *gin.Context) {
	refreshToken, err := ctx.Cookie("refresh_token")
	if err != nil {
		response.WriteErrorResponse(ctx, err)
		return
	}

	h.service.Logout(ctx, refreshToken)

	ctx.SetCookie("access_token", "", -1, "/", "", false, true)
	ctx.SetCookie("refresh_token", "", -1, "/", "", false, true)

	response.WriteMessageResponse(ctx, http.StatusOK, "user successfully logged out")
}
