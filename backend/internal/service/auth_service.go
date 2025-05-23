package service

import (
	"context"
	"github.com/Furkanturan8/goftr-template/internal/model"
	"github.com/Furkanturan8/goftr-template/internal/repository"
	"github.com/Furkanturan8/goftr-template/pkg/errorx"
	"github.com/Furkanturan8/goftr-template/pkg/jwt"
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

func (s *AuthService) Register(ctx context.Context, user model.User) error {
	// Email kontrolü
	exists, err := s.userRepo.ExistsByEmail(ctx, user.Email)
	if err != nil {
		return errorx.WrapErr(errorx.ErrInternal, err)
	}
	if exists {
		return errorx.WrapMsg(errorx.ErrDuplicate, "Bu e-posta adresi zaten kullanımda")
	}

	// Kullanıcıyı kaydet
	if err = s.userRepo.Create(ctx, &user); err != nil {
		return errorx.WrapErr(errorx.ErrInternal, err)
	}

	return nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*model.Token, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, errorx.WrapMsg(errorx.ErrNotFound, "Bu e-posta adresi ile kayıtlı kullanıcı bulunamadı")
	}

	if !user.CheckPassword(password) {
		return nil, errorx.WrapMsg(errorx.ErrInvalidCredentials, "Girdiğiniz şifre yanlış")
	}

	if user.Status != model.StatusActive {
		return nil, errorx.WrapMsg(errorx.ErrForbidden, "Hesabınız aktif değil. Lütfen yönetici ile iletişime geçin")
	}

	// Access token oluştur
	accessToken, err := jwt.Generate(user)
	if err != nil {
		return nil, errorx.WrapErr(errorx.ErrInternal, err)
	}

	// Refresh token oluştur
	refreshToken, err := jwt.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, errorx.WrapErr(errorx.ErrInternal, err)
	}

	// Token kaydını oluştur
	token := &model.Token{
		UserID:       user.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(time.Duration(24) * time.Hour), // 24 saat
	}

	if err = s.authRepo.SaveToken(ctx, token); err != nil {
		return nil, errorx.WrapErr(errorx.ErrInternal, err)
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
		return nil, errorx.WrapErr(errorx.ErrInternal, err)
	}

	return token, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*model.Token, error) {
	// Refresh token'ı doğrula
	claims, err := jwt.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, errorx.WrapMsg(errorx.ErrUnauthorized, "Geçersiz veya süresi dolmuş refresh token")
	}

	// Session'ı kontrol et
	session, err := s.authRepo.GetSessionByRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, errorx.WrapMsg(errorx.ErrNotFound, "Oturum bulunamadı")
	}
	if !session.IsValid() {
		return nil, errorx.WrapMsg(errorx.ErrUnauthorized, "Oturum geçersiz veya süresi dolmuş")
	}

	// Kullanıcıyı getir
	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, errorx.WrapMsg(errorx.ErrNotFound, "Kullanıcı bulunamadı")
	}

	if user.Status != model.StatusActive {
		return nil, errorx.WrapMsg(errorx.ErrForbidden, "Hesabınız aktif değil. Lütfen yönetici ile iletişime geçin")
	}

	// Yeni access token oluştur
	accessToken, err := jwt.Generate(user)
	if err != nil {
		return nil, errorx.WrapErr(errorx.ErrInternal, err)
	}

	// Yeni refresh token oluştur
	newRefreshToken, err := jwt.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, errorx.WrapErr(errorx.ErrInternal, err)
	}

	// Token kaydını güncelle
	token := &model.Token{
		UserID:       user.ID,
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresAt:    time.Now().Add(time.Duration(24) * time.Hour),
	}

	if err = s.authRepo.SaveToken(ctx, token); err != nil {
		return nil, errorx.WrapErr(errorx.ErrInternal, err)
	}

	// Session'ı güncelle
	session.RefreshToken = newRefreshToken
	session.ExpiresAt = time.Now().Add(time.Duration(168) * time.Hour)

	if err = s.authRepo.UpdateSession(ctx, session); err != nil {
		return nil, errorx.WrapErr(errorx.ErrInternal, err)
	}

	return token, nil
}

func (s *AuthService) Logout(ctx context.Context, token string) error {
	// Token'ı doğrula
	_, err := jwt.Validate(token)
	if err != nil {
		return errorx.WrapMsg(errorx.ErrUnauthorized, "Geçersiz veya süresi dolmuş token")
	}

	// Session'ı bul ve sil
	session, err := s.authRepo.GetSessionByRefreshToken(ctx, token)
	if err == nil && session != nil {
		if err = s.authRepo.DeleteSession(ctx, session.ID); err != nil {
			return errorx.WrapErr(errorx.ErrInternal, err)
		}
	}

	// Token'ı blacklist'e ekle
	blacklist := &model.TokenBlacklist{
		Token:     token,
		ExpiresAt: time.Now().Add(time.Duration(24) * time.Hour),
	}

	if err = s.authRepo.AddToBlacklist(ctx, blacklist); err != nil {
		return errorx.WrapErr(errorx.ErrInternal, err)
	}

	return nil
}

func (s *AuthService) ForgotPassword(ctx context.Context, email string) (string, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", errorx.WrapMsg(errorx.ErrNotFound, "Bu e-posta adresi ile kayıtlı kullanıcı bulunamadı")
	}

	// Şifre sıfırlama token'ı oluştur
	resetToken, err := jwt.GeneratePasswordResetToken(user)
	if err != nil {
		return "", errorx.WrapErr(errorx.ErrInternal, err)
	}

	return resetToken, nil
}

func (s *AuthService) ResetPassword(ctx context.Context, token, newPassword string) error {
	// Token'ı doğrula
	claims, err := jwt.ValidatePasswordResetToken(token)
	if err != nil {
		return errorx.WrapMsg(errorx.ErrUnauthorized, "Geçersiz veya süresi dolmuş şifre sıfırlama token'ı")
	}

	// Kullanıcıyı bul
	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return errorx.WrapMsg(errorx.ErrNotFound, "Kullanıcı bulunamadı")
	}

	// Şifreyi güncelle
	if err = user.SetPassword(newPassword); err != nil {
		return errorx.WrapErr(errorx.ErrInternal, err)
	}

	if err = s.userRepo.Update(ctx, user); err != nil {
		return errorx.WrapErr(errorx.ErrInternal, err)
	}

	// Kullanıcının tüm oturumlarını sonlandır
	sessions, err := s.authRepo.GetSessionsByUserID(ctx, user.ID)
	if err == nil {
		for _, session := range sessions {
			err = s.authRepo.DeleteSession(ctx, session.ID)
			if err != nil {
				return errorx.WrapErr(errorx.ErrInternal, err)
			}
		}
	}

	return nil
}

func (s *AuthService) ValidateToken(ctx context.Context, token string) (*jwt.Claims, error) {
	// Token'ın blacklist'te olup olmadığını kontrol et
	isBlacklisted, err := s.authRepo.IsTokenBlacklisted(ctx, token)
	if err != nil {
		return nil, errorx.WrapErr(errorx.ErrInternal, err)
	}

	if isBlacklisted {
		return nil, errorx.WrapMsg(errorx.ErrUnauthorized, "Bu token iptal edilmiş")
	}

	// Token'ı doğrula
	claims, err := jwt.Validate(token)
	if err != nil {
		return nil, errorx.WrapMsg(errorx.ErrUnauthorized, "Geçersiz veya süresi dolmuş token")
	}

	return claims, nil
}

// Cleanup işlemleri
func (s *AuthService) CleanupExpiredData(ctx context.Context) error {
	if err := s.authRepo.CleanupExpiredTokens(ctx); err != nil {
		return errorx.WrapErr(errorx.ErrInternal, err)
	}

	if err := s.authRepo.CleanupExpiredSessions(ctx); err != nil {
		return errorx.WrapErr(errorx.ErrInternal, err)
	}

	return nil
}
