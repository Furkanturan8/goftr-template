package main

import (
	"context"
	"database/sql"
	"fmt"
	"goftr-v1/config"
	"goftr-v1/internal/handler"
	"goftr-v1/internal/repository"
	"goftr-v1/internal/router"
	"goftr-v1/internal/service"
	"goftr-v1/pkg/cache"
	"goftr-v1/pkg/jwt"
	"goftr-v1/pkg/logger"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func main() {
	// Yapılandırmayı yükle
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Printf("Config yükleme hatası: %v", err)
		os.Exit(1)
	}

	// Logger'ı başlat
	if err = logger.Init(cfg.App.LogDir); err != nil {
		log.Printf("Logger başlatma hatası: %v", err)
		os.Exit(1)
	}

	// Redis cache'i başlat
	if err = cache.InitDefaultCache(cfg.Redis.GetAddr(), cfg.Redis.Password, cfg.Redis.DB); err != nil {
		logger.Error("Redis cache başlatma hatası: %v", err)
		os.Exit(1)
	}

	// JWT yapılandırmasını başlat
	jwt.Init(&cfg.JWT)

	// Database bağlantısı
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(cfg.Database.GetDSN())))
	db := bun.NewDB(sqldb, pgdialect.New())
	defer db.Close()

	// Veritabanı bağlantısını kontrol et
	if err = db.Ping(); err != nil {
		logger.Error("Veritabanı bağlantı hatası: %v", err)
		os.Exit(1)
	}
	logger.Info("Veritabanı bağlantısı başarılı")

	// Repository'ler
	userRepo := repository.NewUserRepository(db)
	authRepo := repository.NewAuthRepository(db)

	// Service'ler
	authService := service.NewAuthService(userRepo, authRepo)
	userService := service.NewUserService(userRepo)

	// Handler'lar
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)

	// Router'ı oluştur ve yapılandır
	r := router.NewRouter(authHandler, userHandler)
	r.SetupRoutes()

	// Graceful shutdown için kanal oluştur
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// HTTP sunucusunu başlat
	serverShutdown := make(chan struct{})
	go func() {
		addr := fmt.Sprintf(":%d", cfg.App.Port)
		logger.Info("Sunucu %s portunda başlatılıyor...", addr)
		if err = r.GetApp().Listen(addr); err != nil {
			logger.Error("Sunucu hatası: %v", err)
		}
		close(serverShutdown)
	}()

	// Shutdown sinyalini bekle
	<-shutdown
	logger.Info("Graceful shutdown başlatılıyor...")

	// Shutdown timeout context'i oluştur
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.App.ShutdownTimeout)*time.Second)
	defer cancel()

	// Sunucuyu durdur
	if err = r.GetApp().ShutdownWithContext(ctx); err != nil {
		logger.Error("Sunucu kapatma hatası: %v", err)
	}

	// Veritabanı bağlantısını kapat
	if err = db.Close(); err != nil {
		logger.Error("Veritabanı bağlantısı kapatma hatası: %v", err)
	}

	logger.Info("Sunucu başarıyla kapatıldı")
}
