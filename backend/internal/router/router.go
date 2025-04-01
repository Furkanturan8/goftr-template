package router

import (
	"goftr-v1/backend/config"
	"goftr-v1/backend/internal/handler"
	"goftr-v1/backend/internal/middleware"
	"goftr-v1/backend/pkg/monitoring"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type Router struct {
	app         *fiber.App
	authHandler *handler.AuthHandler
	userHandler *handler.UserHandler
	// Diğer handler'lar buraya eklenecek
}

func NewRouter(authHandler *handler.AuthHandler, userHandler *handler.UserHandler) *Router {
	return &Router{
		app:         fiber.New(),
		authHandler: authHandler,
		userHandler: userHandler,
	}
}

func (r *Router) SetupRoutes(cfg *config.Config) {
	// Middleware'leri ekle
	r.app.Use(logger.New())
	r.app.Use(recover.New())
	r.app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:63342, http://localhost:3005",
		AllowMethods: "GET,POST,PUT,DELETE",
		AllowHeaders: "Content-Type, Authorization",
	}))

	// Rate limiting middleware'i ekle (30 sn de 10 istek olsun)
	r.app.Use(limiter.New(limiter.Config{
		Max:        10,               // Maksimum istek sayısı
		Expiration: 30 * time.Second, // Zaman aralığı
		KeyGenerator: func(c *fiber.Ctx) string {
			// /metrics endpoint'i için rate limiting'i devre dışı bırak
			if c.Path() == "/metrics" {
				return "metrics_no_limit"
			}
			// Her route'u ayrı ayrı sınırla (örneğin: "/users", "/users/:id", "/auth/login")
			return c.IP() + ":" + c.Path()
		},
	}))

	// Prometheus Middleware ekleyelim
	r.app.Use(monitoring.PrometheusMiddleware())

	// API versiyonu
	api := r.app.Group("/api")
	v1 := api.Group("/v1")

	// Prometheus'un topladığı metrikleri görüntülemek için /metrics endpoint'i
	if cfg.MonitoringConfig.Prometheus.Enabled {
		r.app.Get(cfg.MonitoringConfig.Prometheus.Endpoint, monitoring.MetricsHandler())
	}

	// Auth routes
	auth := v1.Group("/auth")
	auth.Post("/register", r.authHandler.Register)
	auth.Post("/login", r.authHandler.Login)
	auth.Post("/refresh", r.authHandler.RefreshToken)
	auth.Post("/forgot-password", r.authHandler.ForgotPassword)
	auth.Post("/reset-password", r.authHandler.ResetPassword)
	auth.Post("/logout", middleware.AuthMiddleware(), r.authHandler.Logout)

	// User routes - Base group
	users := v1.Group("/users")

	// Normal user routes (profil yönetimi)
	userProfile := users.Group("/me")
	userProfile.Use(middleware.AuthMiddleware()) // Sadece authentication gerekli
	userProfile.Get("/", r.userHandler.GetProfile)
	userProfile.Put("/", r.userHandler.UpdateProfile)

	// Admin only routes
	adminUsers := users.Group("/")
	adminUsers.Use(middleware.AuthMiddleware(), middleware.AdminOnly()) // Admin yetkisi gerekli
	adminUsers.Post("/", r.userHandler.Create)
	adminUsers.Get("/", r.userHandler.List)
	adminUsers.Get("/:id", r.userHandler.GetByID)
	adminUsers.Put("/:id", r.userHandler.Update)
	adminUsers.Delete("/:id", r.userHandler.Delete)

	// Diğer route grupları buraya eklenecek
}

func (r *Router) GetApp() *fiber.App {
	return r.app
}
