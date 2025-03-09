package handler

import (
	"goftr-v1/internal/dto"
	"goftr-v1/internal/service"
	"goftr-v1/pkg/errorx"
	"goftr-v1/pkg/response"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{service: s}
}

func (h *UserHandler) List(c *fiber.Ctx) error {
	resp, err := h.service.List(c.Context())
	if err != nil {
		return errorx.WithDetails(errorx.ErrInternal, err.Error())
	}

	return response.Success(c, resp)
}

func (h *UserHandler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return errorx.ErrInvalidRequest
	}

	resp, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		return errorx.WithDetails(errorx.ErrNotFound, "Kullanıcı bulunamadı")
	}
	return response.Success(c, resp)
}

func (h *UserHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return errorx.ErrInvalidRequest
	}

	var req dto.UpdateUserRequest
	if err = c.BodyParser(&req); err != nil {
		return errorx.WithDetails(errorx.ErrInvalidRequest, "Geçersiz giriş formatı")
	}

	if err = h.service.Update(c.Context(), id, &req); err != nil {
		return errorx.WithDetails(errorx.ErrInternal, err.Error())
	}

	return response.Success(c, nil, "Kullanıcı başarıyla güncellendi")
}

func (h *UserHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return errorx.ErrInvalidRequest
	}

	if err = h.service.Delete(c.Context(), id); err != nil {
		return errorx.WithDetails(errorx.ErrInternal, err.Error())
	}
	return response.Success(c, nil, "Kullanıcı başarıyla silindi")
}

func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int64)
	resp, err := h.service.GetByID(c.Context(), userID)
	if err != nil {
		return errorx.WithDetails(errorx.ErrNotFound, "Kullanıcı bulunamadı")
	}
	return response.Success(c, resp)
}

func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int64)

	var req dto.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return errorx.WithDetails(errorx.ErrInvalidRequest, "Geçersiz giriş formatı")
	}

	if err := h.service.Update(c.Context(), userID, &req); err != nil {
		return errorx.WithDetails(errorx.ErrInternal, err.Error())
	}

	return response.Success(c, nil, "Profil başarıyla güncellendi")
}
