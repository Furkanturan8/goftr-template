package handler

import (
	"fmt"
	"github.com/Furkanturan8/goftr-template/internal/dto"
	"github.com/Furkanturan8/goftr-template/internal/model"
	"github.com/Furkanturan8/goftr-template/internal/service"
	"github.com/Furkanturan8/goftr-template/pkg/errorx"
	"github.com/Furkanturan8/goftr-template/pkg/response"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{service: s}
}

func (h *UserHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return errorx.WrapErr(errorx.ErrInvalidRequest, err)
	}

	user := req.ToDBModel(model.User{})
	if user.Password == "" { // when admin create a new user, password is empty. so we set default password
		// maybe we can use a link to send a mail to the user to set a password
		// todo: send email to user to set a password
		_ = user.SetPassword("goftr-template-default-password-1907") // default password
	}

	if err := h.service.Create(c.Context(), user); err != nil {
		return errorx.WrapErr(errorx.ErrInternal, err)
	}

	return response.Success(c, nil, "Kullanıcı başarıyla oluşturuldu")
}

func (h *UserHandler) List(c *fiber.Ctx) error {
	resp, err := h.service.List(c.Context())
	if err != nil {
		return errorx.WrapErr(errorx.ErrInternal, err)
	}

	users := make([]dto.UserResponse, len(resp))

	for i, user := range resp {
		users[i] = dto.UserResponse{}.ToResponseModel(user)
	}

	return response.Success(c, users)
}

func (h *UserHandler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return errorx.WrapErr(errorx.ErrInvalidRequest, err)
	}

	resp, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		return errorx.WrapMsg(errorx.ErrNotFound, "Kullanıcı Bulunamadı!")
	}

	user := dto.UserResponse{}.ToResponseModel(*resp)
	return response.Success(c, user)
}

func (h *UserHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return errorx.WrapErr(errorx.ErrInvalidRequest, err)
	}

	currentUser, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		return errorx.WrapMsg(errorx.ErrNotFound, "Kullanıcı Bulunamadı!")
	}

	var req dto.UpdateUserRequest
	if err = c.BodyParser(&req); err != nil {
		return errorx.WrapErr(errorx.ErrInvalidRequest, err)
	}

	user := req.ToDBModel(model.User{})
	user.ID = id
	// Eğer şifre değiştirilmek isteniyorsa
	if req.NewPassword != "" {
		// Eski şifre doğrulaması zorunlu
		if currentUser.CheckPassword(req.CurrentPassword) {
			// Yeni şifre hashlenerek ayarlanır
			_ = user.SetPassword(req.NewPassword)
			fmt.Println("şifre değiştirildi")
		} else {
			user.Password = currentUser.Password
			fmt.Println("hatalı şifre değiştirilmedi!")
		}
	} else if req.CurrentPassword == "" || req.NewPassword == "" {
		fmt.Println("şifre değiştirilmedi!")
		user.Password = currentUser.Password
	}

	if err = h.service.Update(c.Context(), id, user); err != nil {
		return errorx.WrapErr(errorx.ErrInternal, err)
	}

	return response.Success(c, nil, "Kullanıcı başarıyla güncellendi")
}

func (h *UserHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return errorx.WrapErr(errorx.ErrInvalidRequest, err)
	}

	if err = h.service.Delete(c.Context(), id); err != nil {
		return errorx.WrapErr(errorx.ErrInternal, err)
	}
	return response.Success(c, nil, "Kullanıcı başarıyla silindi")
}

func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int64)
	resp, err := h.service.GetByID(c.Context(), userID)
	if err != nil {
		return errorx.WrapErr(errorx.ErrNotFound, err)
	}

	user := dto.UserResponse{}.ToResponseModel(*resp)
	return response.Success(c, user)
}

func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int64)
	role := c.Locals("role").(model.Role)
	status := c.Locals("status").(model.Status)

	currentUser, err := h.service.GetByID(c.Context(), userID)
	if err != nil {
		return errorx.WrapErr(errorx.ErrNotFound, err)
	}

	var req dto.UpdateUserRequest
	if err = c.BodyParser(&req); err != nil {
		return errorx.WrapErr(errorx.ErrInvalidRequest, err)
	}

	user := req.ToDBModel(model.User{})
	user.ID = userID
	user.Role = role
	user.Status = status

	// Eğer şifre değiştirilmek isteniyorsa
	if req.NewPassword != "" {
		// Eski şifre doğrulaması zorunlu
		if currentUser.CheckPassword(req.CurrentPassword) {
			// Yeni şifre hashlenerek ayarlanır
			_ = user.SetPassword(req.NewPassword)
			fmt.Println("şifre değiştirildi")
		} else {
			user.Password = currentUser.Password
			fmt.Println("hatalı şifre değiştirilmedi!")
		}
	} else if req.CurrentPassword == "" || req.NewPassword == "" {
		fmt.Println("şifre değiştirilmedi!")
		user.Password = currentUser.Password
	}

	if err = h.service.Update(c.Context(), userID, user); err != nil {
		return errorx.WrapErr(errorx.ErrInternal, err)
	}

	return response.Success(c, nil, "Profil başarıyla güncellendi")
}
