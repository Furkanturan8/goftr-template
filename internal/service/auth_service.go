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
	authRepo *repository.AuthRepository
}

func NewAuthService(userRepo *repository.UserRepository, authRepo *repository.AuthRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		authRepo: authRepo,
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

	// Access token oluştur
	accessToken, err := jwt.Generate(user)
	if err != nil {
		return nil, response.ErrInternal
	}

	// Refresh token oluştur (örnek olarak rastgele bir string)
	refreshToken := "refresh_" + accessToken // Gerçek uygulamada güvenli bir yöntem kullanılmalı

	// Token kaydını oluştur
	token := &model.Token{
		UserID:       user.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(time.Hour * 24), // 24 saat
	}

	if err := s.authRepo.CreateToken(ctx, token); err != nil {
		return nil, response.ErrInternal
	}

	// Oturum kaydı oluştur
	session := &model.Session{
		UserID:       user.ID,
		RefreshToken: refreshToken,
		UserAgent:    "web",                               // Gerçek uygulamada request'ten alınmalı
		ClientIP:     "0.0.0.0",                           // Gerçek uygulamada request'ten alınmalı
		ExpiresAt:    time.Now().Add(time.Hour * 24 * 30), // 30 gün
	}

	if err := s.authRepo.CreateSession(ctx, session); err != nil {
		return nil, response.ErrInternal
	}

	// Son giriş zamanını güncelle
	if err := s.userRepo.UpdateLastLogin(ctx, user.ID); err != nil {
		// Log error but don't fail the login
		// logger.Error("Failed to update last login: %v", err)
	}

	return &dto.TokenResponse{
		AccessToken:  accessToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(time.Hour * 24), // 24 saat
		RefreshToken: refreshToken,
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
	newAccessToken, err := jwt.Generate(user)
	if err != nil {
		return nil, response.ErrInternal
	}

	// Yeni refresh token oluştur
	newRefreshToken := "refresh_" + newAccessToken // Gerçek uygulamada güvenli bir yöntem kullanılmalı

	// Eski token'ı geçersiz kıl
	oldToken, err := s.authRepo.GetTokenByAccessToken(ctx, token)
	if err == nil {
		_ = s.authRepo.RevokeToken(ctx, oldToken.ID)
	}

	// Yeni token kaydı oluştur
	newToken := &model.Token{
		UserID:       user.ID,
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		ExpiresAt:    time.Now().Add(time.Hour * 24), // 24 saat
	}

	if err := s.authRepo.CreateToken(ctx, newToken); err != nil {
		return nil, response.ErrInternal
	}

	return &dto.TokenResponse{
		AccessToken:  newAccessToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(time.Hour * 24), // 24 saat
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, token string) error {
	// Token'ı blacklist'e ekle
	blacklist := &model.TokenBlacklist{
		Token:     token,
		ExpiresAt: time.Now().Add(time.Hour * 24), // 24 saat boyunca blacklist'te tut
	}

	if err := s.authRepo.AddToBlacklist(ctx, blacklist); err != nil {
		return response.ErrInternal
	}

	// Token kaydını bul ve geçersiz kıl
	tokenRecord, err := s.authRepo.GetTokenByAccessToken(ctx, token)
	if err == nil {
		_ = s.authRepo.RevokeToken(ctx, tokenRecord.ID)
	}

	return nil
}
