package service

import (
	"context"
	"goftr-v1/backend/internal/dto"
	"goftr-v1/backend/internal/model"
	"goftr-v1/backend/internal/repository"
	"goftr-v1/backend/pkg/errorx"
	"goftr-v1/backend/pkg/jwt"

	"time"
)

type AuthService struct {
	authRepo repository.IAuthRepository
	userRepo repository.IUserRepository
}

func NewAuthService(a repository.IAuthRepository, u repository.IUserRepository) *AuthService {
	return &AuthService{
		authRepo: a,
		userRepo: u,
	}
}

func (s *AuthService) Register(ctx context.Context, req *dto.RegisterRequest) (*model.User, error) {
	// Email kontrolü
	exists, err := s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, errorx.ErrDatabaseOperation
	}
	if exists {
		return nil, errorx.WithDetails(errorx.ErrInvalidRequest, "Email already exists")
	}

	// Yeni kullanıcı oluştur
	user := &model.User{
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      model.UserRole,     // Varsayılan rol
		Status:    model.StatusActive, // Varsayılan durum
	}

	// Şifreyi hashle
	if err = user.SetPassword(req.Password); err != nil {
		return nil, errorx.ErrPasswordHash
	}

	// Kullanıcıyı kaydet
	if err = s.userRepo.Create(ctx, user); err != nil {
		return nil, errorx.ErrDatabaseOperation
	}

	return user, nil
}

func (s *AuthService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, jwt.ErrInvalidCredentials
	}

	if !user.CheckPassword(req.Password) {
		return nil, jwt.ErrInvalidCredentials
	}

	if user.Status != model.StatusActive {
		return nil, jwt.ErrAccountInactive
	}

	// Access token oluştur
	accessToken, err := jwt.Generate(user)
	if err != nil {
		return nil, jwt.ErrTokenGeneration
	}

	// Refresh token oluştur
	refreshToken, err := jwt.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, jwt.ErrTokenGeneration
	}

	// Token kaydını oluştur
	token := &model.Token{
		UserID:       user.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(time.Duration(24) * time.Hour), // 24 saat
	}

	if err = s.authRepo.SaveToken(ctx, token); err != nil {
		return nil, errorx.ErrDatabaseOperation
	}

	// Session oluştur
	session := &model.Session{
		UserID:       user.ID,
		RefreshToken: refreshToken,
		UserAgent:    ctx.Value("user_agent").(string),
		ClientIP:     ctx.Value("client_ip").(string),
		ExpiresAt:    time.Now().Add(time.Duration(168) * time.Hour), // 7 gün
	}

	if err = s.authRepo.CreateSession(ctx, session); err != nil {
		return nil, errorx.ErrDatabaseOperation
	}

	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    24 * 60 * 60, // 24 saat (saniye cinsinden)
	}, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*dto.LoginResponse, error) {
	// Refresh token'ı doğrula
	claims, err := jwt.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, jwt.ErrInvalidToken
	}

	// Session'ı kontrol et
	session, err := s.authRepo.GetSessionByRefreshToken(ctx, refreshToken)
	if err != nil || !session.IsValid() {
		return nil, jwt.ErrInvalidSession
	}

	// Kullanıcıyı getir
	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, errorx.ErrNotFound
	}

	if user.Status != model.StatusActive {
		return nil, jwt.ErrAccountInactive
	}

	// Yeni access token oluştur
	accessToken, err := jwt.Generate(user)
	if err != nil {
		return nil, jwt.ErrTokenGeneration
	}

	// Yeni refresh token oluştur
	newRefreshToken, err := jwt.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, jwt.ErrTokenGeneration
	}

	// Token kaydını güncelle
	token := &model.Token{
		UserID:       user.ID,
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresAt:    time.Now().Add(time.Duration(24) * time.Hour),
	}

	if err = s.authRepo.SaveToken(ctx, token); err != nil {
		return nil, errorx.ErrDatabaseOperation
	}

	// Session'ı güncelle
	session.RefreshToken = newRefreshToken
	session.ExpiresAt = time.Now().Add(time.Duration(168) * time.Hour)

	if err = s.authRepo.UpdateSession(ctx, session); err != nil {
		return nil, errorx.ErrDatabaseOperation
	}

	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    24 * 60 * 60,
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, token string) error {
	// Token'ı doğrula
	_, err := jwt.Validate(token)
	if err != nil {
		return jwt.ErrInvalidToken
	}

	// Session'ı bul ve sil
	session, err := s.authRepo.GetSessionByRefreshToken(ctx, token)
	if err == nil && session != nil {
		if err = s.authRepo.DeleteSession(ctx, session.ID); err != nil {
			return errorx.ErrDatabaseOperation
		}
	}

	// Token'ı blacklist'e ekle
	blacklist := &model.TokenBlacklist{
		Token:     token,
		ExpiresAt: time.Now().Add(time.Duration(24) * time.Hour),
	}

	if err = s.authRepo.AddToBlacklist(ctx, blacklist); err != nil {
		return errorx.ErrDatabaseOperation
	}

	return nil
}

func (s *AuthService) ForgotPassword(ctx context.Context, email string) (string, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", errorx.ErrNotFound
	}

	// Şifre sıfırlama token'ı oluştur
	resetToken, err := jwt.GeneratePasswordResetToken(user)
	if err != nil {
		return "", jwt.ErrTokenGeneration
	}

	return resetToken, nil
}

func (s *AuthService) ResetPassword(ctx context.Context, token, newPassword string) error {
	// Token'ı doğrula
	claims, err := jwt.ValidatePasswordResetToken(token)
	if err != nil {
		return jwt.ErrInvalidToken
	}

	// Kullanıcıyı bul
	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return errorx.ErrNotFound
	}

	// Şifreyi güncelle
	if err = user.SetPassword(newPassword); err != nil {
		return errorx.ErrInternal
	}

	if err = s.userRepo.Update(ctx, user); err != nil {
		return errorx.ErrDatabaseOperation
	}

	// Kullanıcının tüm oturumlarını sonlandır
	sessions, err := s.authRepo.GetSessionsByUserID(ctx, user.ID)
	if err == nil {
		for _, session := range sessions {
			err = s.authRepo.DeleteSession(ctx, session.ID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *AuthService) ValidateToken(ctx context.Context, token string) (*jwt.Claims, error) {
	// Token'ın blacklist'te olup olmadığını kontrol et
	isBlacklisted, err := s.authRepo.IsTokenBlacklisted(ctx, token)
	if err != nil {
		return nil, errorx.ErrDatabaseOperation
	}

	if isBlacklisted {
		return nil, jwt.ErrInvalidToken
	}

	// Token'ı doğrula
	claims, err := jwt.Validate(token)
	if err != nil {
		return nil, jwt.ErrInvalidToken
	}

	return claims, nil
}

// Cleanup işlemleri
func (s *AuthService) CleanupExpiredData(ctx context.Context) error {
	if err := s.authRepo.CleanupExpiredTokens(ctx); err != nil {
		return err
	}

	if err := s.authRepo.CleanupExpiredSessions(ctx); err != nil {
		return err
	}

	return nil
}
