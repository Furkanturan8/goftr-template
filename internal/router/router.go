package router

import (
	"goftr-v1/internal/handler"
	"goftr-v1/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type Router struct {
	app         *fiber.App
	authHandler *handler.AuthHandler
	userHandler *handler.UserHandler
	// Diğer handler'lar buraya eklenecek
}

func NewRouter(
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	// Diğer handler'lar buraya eklenecek
) *Router {
	return &Router{
		app:         fiber.New(),
		authHandler: authHandler,
		userHandler: userHandler,
	}
}

func (r *Router) SetupRoutes() {
	// Middleware'leri ekle
	r.app.Use(logger.New())
	r.app.Use(recover.New())
	r.app.Use(cors.New())

	// API versiyonu
	api := r.app.Group("/api")
	v1 := api.Group("/v1")

	// Auth routes
	auth := v1.Group("/auth")
	auth.Post("/register", r.authHandler.Register)
	auth.Post("/login", r.authHandler.Login)
	auth.Post("/refresh", r.authHandler.RefreshToken)
	auth.Post("/logout", middleware.AuthMiddleware(), r.authHandler.Logout)

	// User routes
	users := v1.Group("/users")
	users.Use(middleware.AuthMiddleware()) // Tüm user route'ları için auth gerekli

	// Normal user routes (profil yönetimi)
	users.Get("/me", r.userHandler.GetProfile)
	users.Put("/me", r.userHandler.UpdateProfile)

	// Admin only routes
	users.Use(middleware.AdminOnly()) // Bundan sonraki tüm route'lar için admin yetkisi gerekli
	users.Get("/", r.userHandler.List)
	users.Get("/:id", r.userHandler.GetByID)
	users.Put("/:id", r.userHandler.Update)
	users.Delete("/:id", r.userHandler.Delete)

	// Diğer route grupları buraya eklenecek
}

func (r *Router) GetApp() *fiber.App {
	return r.app
}
