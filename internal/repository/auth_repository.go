package repository

import (
	"context"
	"goftr-v1/internal/model"
	"time"

	"github.com/uptrace/bun"
)

type AuthRepository struct {
	db *bun.DB
}

func NewAuthRepository(db *bun.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

// CreateToken yeni bir token kaydı oluşturur
func (r *AuthRepository) CreateToken(ctx context.Context, token *model.Token) error {
	_, err := r.db.NewInsert().Model(token).Exec(ctx)
	return err
}

// GetTokenByAccessToken access token ile token kaydını bulur
func (r *AuthRepository) GetTokenByAccessToken(ctx context.Context, accessToken string) (*model.Token, error) {
	token := new(model.Token)
	err := r.db.NewSelect().
		Model(token).
		Where("access_token = ?", accessToken).
		Scan(ctx)
	return token, err
}

// GetTokenByRefreshToken refresh token ile token kaydını bulur
func (r *AuthRepository) GetTokenByRefreshToken(ctx context.Context, refreshToken string) (*model.Token, error) {
	token := new(model.Token)
	err := r.db.NewSelect().
		Model(token).
		Where("refresh_token = ?", refreshToken).
		Scan(ctx)
	return token, err
}

// RevokeToken bir token'ı geçersiz kılar
func (r *AuthRepository) RevokeToken(ctx context.Context, tokenID int64) error {
	_, err := r.db.NewUpdate().
		Model((*model.Token)(nil)).
		Set("revoked_at = ?", time.Now()).
		Where("id = ?", tokenID).
		Exec(ctx)
	return err
}

// CreateSession yeni bir oturum kaydı oluşturur
func (r *AuthRepository) CreateSession(ctx context.Context, session *model.Session) error {
	_, err := r.db.NewInsert().Model(session).Exec(ctx)
	return err
}

// GetSessionByRefreshToken refresh token ile oturum kaydını bulur
func (r *AuthRepository) GetSessionByRefreshToken(ctx context.Context, refreshToken string) (*model.Session, error) {
	session := new(model.Session)
	err := r.db.NewSelect().
		Model(session).
		Where("refresh_token = ?", refreshToken).
		Scan(ctx)
	return session, err
}

// BlockSession bir oturumu bloklar
func (r *AuthRepository) BlockSession(ctx context.Context, sessionID int64) error {
	_, err := r.db.NewUpdate().
		Model((*model.Session)(nil)).
		Set("is_blocked = ?", true).
		Where("id = ?", sessionID).
		Exec(ctx)
	return err
}

// AddToBlacklist bir token'ı kara listeye ekler
func (r *AuthRepository) AddToBlacklist(ctx context.Context, blacklist *model.TokenBlacklist) error {
	_, err := r.db.NewInsert().Model(blacklist).Exec(ctx)
	return err
}

// IsTokenBlacklisted bir token'ın kara listede olup olmadığını kontrol eder
func (r *AuthRepository) IsTokenBlacklisted(ctx context.Context, token string) (bool, error) {
	exists, err := r.db.NewSelect().
		Model((*model.TokenBlacklist)(nil)).
		Where("token = ? AND expires_at > ?", token, time.Now()).
		Exists(ctx)
	return exists, err
}

// CleanupExpiredTokens süresi dolmuş token kayıtlarını temizler
func (r *AuthRepository) CleanupExpiredTokens(ctx context.Context) error {
	_, err := r.db.NewDelete().
		Model((*model.TokenBlacklist)(nil)).
		Where("expires_at < ?", time.Now()).
		Exec(ctx)
	return err
}
