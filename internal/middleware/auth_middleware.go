package middleware

import (
	"goftr-v1/internal/model"
	"goftr-v1/pkg/jwt"
	"goftr-v1/pkg/response"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Authorization header kontrolü
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return response.Errorx(c, response.StatusUnauthorized, response.WithDetails(
				response.ErrUnauthorized,
				"Authorization header eksik",
			))
		}

		// Bearer token formatı kontrolü
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return response.Errorx(c, response.StatusUnauthorized, response.WithDetails(
				response.ErrUnauthorized,
				"Geçersiz Authorization header formatı. 'Bearer <token>' formatında olmalı",
			))
		}

		// Token doğrulama
		claims, err := jwt.Validate(tokenParts[1])
		if err != nil {
			return response.Errorx(c, response.StatusUnauthorized, response.WithDetails(
				response.ErrUnauthorized,
				"Geçersiz veya süresi dolmuş token",
			))
		}

		// Context'e kullanıcı bilgilerini ekle
		c.Locals("userID", claims.UserID)
		c.Locals("role", claims.Role)
		c.Locals("email", claims.Email)

		return c.Next()
	}
}

func AdminOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role")
		if role == nil {
			return response.Errorx(c, response.StatusUnauthorized, response.WithDetails(
				response.ErrUnauthorized,
				"Yetkilendirme bilgisi bulunamadı",
			))
		}

		if role.(model.Role) != model.AdminRole {
			return response.Errorx(c, response.StatusForbidden, response.WithDetails(
				response.ErrForbidden,
				"Bu işlem için admin yetkisi gerekli",
			))
		}

		return c.Next()
	}
}

// Belirli rollere sahip kullanıcılar için middleware
func HasRole(roles ...model.Role) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role")
		if role == nil {
			return response.Errorx(c, response.StatusUnauthorized, response.WithDetails(
				response.ErrUnauthorized,
				"Yetkilendirme bilgisi bulunamadı",
			))
		}

		userRole := role.(model.Role)
		for _, allowedRole := range roles {
			if userRole == allowedRole {
				return c.Next()
			}
		}

		return response.Errorx(c, response.StatusForbidden, response.WithDetails(
			response.ErrForbidden,
			"Bu işlem için yeterli yetkiniz yok",
		))
	}
}
