package service

import (
	"context"
	"fmt"
	"goftr-v1/internal/dto"
	"goftr-v1/internal/model"
	"goftr-v1/internal/repository"
	"goftr-v1/pkg/jwt"
	"goftr-v1/pkg/response"
	"time"
)

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

func (s *AuthService) Register(ctx context.Context, req *dto.RegisterRequest) error {
	// Email kontrolü
	exists, err := s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return fmt.Errorf("email kontrol hatası: %v", err)
	}
	if exists {
		return fmt.Errorf("bu email adresi zaten kullanımda: %s", req.Email)
	}

	// Yeni kullanıcı oluştur
	user := &model.User{
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      model.UserRole,
		Status:    model.StatusActive,
	}

	// Şifreyi hashle
	if err := user.SetPassword(req.Password); err != nil {
		return fmt.Errorf("şifre hashleme hatası: %v", err)
	}

	// Kullanıcıyı kaydet
	if err := s.userRepo.Create(ctx, user); err != nil {
		return fmt.Errorf("kullanıcı oluşturma hatası: %v", err)
	}

	return nil
}

func (s *AuthService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.TokenResponse, error) {
	// Kullanıcıyı bul
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, response.ErrUnauthorized
	}

	// Şifreyi kontrol et
	if !user.CheckPassword(req.Password) {
		return nil, response.ErrUnauthorized
	}

	// Token oluştur
	token, err := jwt.Generate(user)
	if err != nil {
		return nil, response.ErrInternal
	}

	// Son giriş zamanını güncelle
	if err := s.userRepo.UpdateLastLogin(ctx, user.ID); err != nil {
		// Log error but don't fail the login
		// logger.Error("Failed to update last login: %v", err)
	}

	return &dto.TokenResponse{
		AccessToken:  token,
		TokenType:    "Bearer",
		ExpiresIn:    int64(time.Hour * 24), // 24 saat
		RefreshToken: "",                    // Refresh token implementation
	}, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, token string) (*dto.TokenResponse, error) {
	// Token'ı doğrula
	claims, err := jwt.Validate(token)
	if err != nil {
		return nil, response.ErrUnauthorized
	}

	// Kullanıcıyı kontrol et
	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, response.ErrNotFound
	}

	if user.Status != model.StatusActive {
		return nil, response.ErrForbidden
	}

	// Yeni token oluştur
	newToken, err := jwt.Generate(user)
	if err != nil {
		return nil, response.ErrInternal
	}

	return &dto.TokenResponse{
		AccessToken:  newToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(time.Hour * 24), // 24 saat
		RefreshToken: "",                    // Refresh token implementation
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, token string) error {
	// Token'ı blacklist'e ekle veya Redis'te invalidate et
	// Bu örnekte sadece başarılı dönüyoruz
	return nil
}
