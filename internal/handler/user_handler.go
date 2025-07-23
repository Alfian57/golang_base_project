package handler

import (
	"net/http"

	"github.com/Alfian57/belajar-golang/internal/dto"
	"github.com/Alfian57/belajar-golang/internal/response"
	"github.com/Alfian57/belajar-golang/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{
		service: s,
	}
}

func (h *UserHandler) GetAllUsers(ctx *gin.Context) {
	var query dto.GetUsersFilter
	if err := ctx.ShouldBindQuery(&query); err != nil {
		response.WriteErrorResponse(ctx, err)
		return
	}

	result, err := h.service.GetAllUsers(ctx, query)
	if err != nil {
		response.WriteErrorResponse(ctx, err)
		return
	}

	response.WritePaginatedResponse(ctx, http.StatusOK, result)
}

func (h *UserHandler) CreateUser(ctx *gin.Context) {
	var request dto.CreateUserRequest
	if err := ctx.ShouldBind(&request); err != nil {
		response.WriteErrorResponse(ctx, err)
		return
	}

	if err := h.service.CreateUser(ctx, request); err != nil {
		response.WriteErrorResponse(ctx, err)
		return
	}

	response.WriteMessageResponse(ctx, http.StatusCreated, "user successfully created")
}

func (h *UserHandler) GetUserByID(ctx *gin.Context) {
	idParam := ctx.Param("id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		response.WriteErrorResponse(ctx, err)
		return
	}

	user, err := h.service.GetUserByID(ctx, id.String())
	if err != nil {
		response.WriteErrorResponse(ctx, err)
		return
	}

	response.WriteDataResponse(ctx, http.StatusOK, user)
}

func (h *UserHandler) UpdateUser(ctx *gin.Context) {
	var request dto.UpdateUserRequest
	if err := ctx.ShouldBind(&request); err != nil {
		response.WriteErrorResponse(ctx, err)
		return
	}

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		response.WriteErrorResponse(ctx, err)
		return
	}
	request.ID = id

	if err := h.service.UpdateUser(ctx, request); err != nil {
		response.WriteErrorResponse(ctx, err)
		return
	}

	response.WriteMessageResponse(ctx, http.StatusOK, "user successfully updated")
}

func (h *UserHandler) DeleteUser(ctx *gin.Context) {
	idParam := ctx.Param("id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		response.WriteErrorResponse(ctx, err)
		return
	}

	if err := h.service.DeleteUser(ctx, id); err != nil {
		response.WriteErrorResponse(ctx, err)
		return
	}

	response.WriteMessageResponse(ctx, http.StatusOK, "user successfully deleted")
}
