package handler

import (
	"github.com/Furkanturan8/goftr-template/internal/dto"
	"github.com/Furkanturan8/goftr-template/internal/model"
	"github.com/Furkanturan8/goftr-template/internal/service"
	"github.com/Furkanturan8/goftr-template/pkg/errorx"
	"github.com/Furkanturan8/goftr-template/pkg/response"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"time"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

var validate = validator.New()

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req dto.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return errorx.ErrInvalidRequest
	}

	// Validasyon
	if err := validate.Struct(req); err != nil {
		return errorx.ErrInvalidRequest
	}

	// Şifre uzunluğu kontrolü
	if len(req.Password) < 6 {
		return errorx.WithDetails(errorx.ErrInvalidRequest, "Password must be at least 6 characters")
	}

	user := req.ToDBModel(model.User{})

	err := h.authService.Register(c.Context(), user)
	if err != nil {
		return errorx.WithDetails(errorx.ErrInternal, err.Error())
	}

	resp := dto.RegisterResponse{
		Email: user.Email,
	}
	return response.Success(c, resp, "User registered successfully")
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return errorx.ErrValidation
	}

	// Validasyon
	if err := validate.Struct(req); err != nil {
		return errorx.ErrInvalidRequest
	}

	// Context'e client bilgilerini ekle
	ctx := c.Context()
	ctx.SetUserValue("user_agent", c.Get("User-Agent"))
	ctx.SetUserValue("client_ip", c.IP())

	token, err := h.authService.Login(ctx, req.Email, req.Password)
	if err != nil {
		return err
	}

	resp := dto.LoginResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresIn:    int(time.Until(token.ExpiresAt).Seconds()),
	}

	return response.Success(c, resp, "Login successful")
}

func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req dto.RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return errorx.ErrValidation
	}

	// Validasyon
	if err := validate.Struct(req); err != nil {
		return errorx.ErrInvalidRequest
	}

	token, err := h.authService.RefreshToken(c.Context(), req.RefreshToken)
	if err != nil {
		return errorx.ErrUnauthorized
	}

	resp := dto.LoginResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresIn:    int(time.Until(token.ExpiresAt).Seconds()),
	}

	return response.Success(c, resp)
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return errorx.ErrUnauthorized
	}

	// "Bearer " prefix'ini kaldır
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	if err := h.authService.Logout(c.Context(), token); err != nil {
		return errorx.ErrInternal
	}

	return response.Success(c, "Logged out successfully")
}

func (h *AuthHandler) ForgotPassword(c *fiber.Ctx) error {
	var req dto.ForgotPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return errorx.ErrInvalidRequest
	}

	// Validasyon
	if err := validate.Struct(req); err != nil {
		return errorx.ErrInvalidRequest
	}

	resetToken, err := h.authService.ForgotPassword(c.Context(), req.Email)
	if err != nil {
		return errorx.ErrInvalidRequest
	}

	// TODO: Send email with reset token
	// For development, return the token
	return response.Success(c, resetToken, "Password reset instructions have been sent to your email")
}

func (h *AuthHandler) ResetPassword(c *fiber.Ctx) error {
	var req dto.ResetPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return errorx.ErrInvalidRequest
	}

	// Validasyon
	if err := validate.Struct(req); err != nil {
		return errorx.ErrInvalidRequest
	}

	if err := h.authService.ResetPassword(c.Context(), req.Token, req.NewPassword); err != nil {
		return errorx.ErrInvalidRequest
	}

	return response.Success(c, "Password has been reset successfully")
}
