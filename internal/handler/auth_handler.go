package handler

import (
	"goftr-v1/internal/dto"
	"goftr-v1/internal/service"
	"goftr-v1/pkg/errorx"
	"goftr-v1/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return errorx.ErrInvalidRequest
	}

	// Validasyon
	if req.Email == "" || req.Password == "" || req.FirstName == "" || req.LastName == "" {
		return errorx.ErrInvalidRequest
	}

	// Şifre uzunluğu kontrolü
	if len(req.Password) < 6 {
		return errorx.WithDetails(errorx.ErrInvalidRequest, "Password must be at least 6 characters")
	}

	user, err := h.authService.Register(c.Context(), &req)
	if err != nil {
		return errorx.WithDetails(errorx.ErrInternal, err.Error())
	}

	return response.Success(c, dto.RegisterResponse{ID: user.ID, Email: user.Email}, "User registered successfully")
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return errorx.ErrValidation
	}

	// Validasyon
	if req.Email == "" || req.Password == "" {
		return errorx.ErrValidation
	}

	// Context'e client bilgilerini ekle
	ctx := c.Context()
	ctx.SetUserValue("user_agent", c.Get("User-Agent"))
	ctx.SetUserValue("client_ip", c.IP())

	token, err := h.authService.Login(ctx, &req)
	if err != nil {
		return errorx.ErrInvalidRequest
	}

	return response.Success(c, dto.LoginResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	})
}

func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req dto.RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return errorx.ErrValidation
	}

	if req.RefreshToken == "" {
		return errorx.ErrInvalidRequest
	}

	token, err := h.authService.RefreshToken(c.Context(), req.RefreshToken)
	if err != nil {
		return errorx.ErrUnauthorized
	}

	return response.Success(c, dto.LoginResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	})
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

	if req.Email == "" {
		return errorx.WithDetails(errorx.ErrValidation, "Email is required")
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

	if req.Token == "" || req.NewPassword == "" {
		return errorx.WithDetails(errorx.ErrInvalidRequest, "Token and new password are required")
	}

	if err := h.authService.ResetPassword(c.Context(), req.Token, req.NewPassword); err != nil {
		return errorx.ErrInvalidRequest
	}

	return response.Success(c, "Password has been reset successfully")
}
