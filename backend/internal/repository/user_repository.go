package repository

import (
	"context"
	"fmt"
	"goftr-v1/backend/internal/model"
	"goftr-v1/backend/pkg/cache"

	"time"

	"github.com/uptrace/bun"
)

const (
	userCacheKeyPrefix = "user:"
	userCacheDuration  = 24 * time.Hour
)

type IUserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id int64) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id int64) error
	UpdateLastLogin(ctx context.Context, id int64) error
	List(ctx context.Context) ([]model.User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}

type UserRepository struct {
	db *bun.DB
}

func NewUserRepository(db *bun.DB) IUserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	_, err := r.db.NewInsert().Model(user).Exec(ctx)
	if err != nil {
		return fmt.Errorf("veritabanı insert hatası: %v", err)
	}

	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*model.User, error) {
	cacheKey := fmt.Sprintf("%s%d", userCacheKeyPrefix, id)

	// Önce cache'den kontrol et
	var user model.User
	err := cache.Get(ctx, cacheKey, &user)
	if err == nil {
		fmt.Printf("Kullanıcı (ID: %d) cache'den alındı\n", id)
		return &user, nil
	}

	// Cache'de yoksa veritabanından al
	user = model.User{}
	err = r.db.NewSelect().Model(&user).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Kullanıcı (ID: %d) veritabanından alındı\n", id)

	// Cache'e kaydet
	if err = cache.Set(ctx, cacheKey, &user, userCacheDuration); err != nil {
		// Cache hatası loglansın ama işlemi engellemeyecek
		return &user, nil
	}

	return &user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.db.NewSelect().Model(&user).Where("email = ?", email).Scan(ctx)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
	user.UpdatedAt = time.Now()
	// Sadece değişen alanları güncelle
	_, err := r.db.NewUpdate().
		Model(user).
		WherePK().
		Column("email", "first_name", "last_name", "password_hash", "role", "status", "updated_at").
		Exec(ctx)
	if err != nil {
		return err
	}

	// Cache'den sil
	cacheKey := fmt.Sprintf("%s%d", userCacheKeyPrefix, user.ID)
	if err = cache.Delete(ctx, cacheKey); err != nil {
		// Cache hatası loglansın ama işlemi engellemeyecek
		return nil
	}

	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.NewDelete().Model((*model.User)(nil)).Where("id = ?", id).Exec(ctx)
	if err != nil {
		return err
	}

	// Cache'den sil
	cacheKey := fmt.Sprintf("%s%d", userCacheKeyPrefix, id)
	if err = cache.Delete(ctx, cacheKey); err != nil {
		// Cache hatası loglansın ama işlemi engellemeyecek
		return nil
	}

	return nil
}

func (r *UserRepository) UpdateLastLogin(ctx context.Context, id int64) error {
	user := &model.User{ID: id}
	_, err := r.db.NewUpdate().
		Model(user).
		Column("last_login").
		WherePK().
		Exec(ctx)

	return err
}

func (r *UserRepository) List(ctx context.Context) ([]model.User, error) {
	cacheKey := fmt.Sprintf("%s%d", userCacheKeyPrefix)

	// Önce cache'den kontrol et
	var users []model.User
	err := cache.Get(ctx, cacheKey, &users)
	if err == nil {
		fmt.Printf("Kullanıcılar cache'den alındı\n")
		return users, nil
	}

	err = r.db.NewSelect().Model(&users).Scan(ctx)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Kullanıcılar veritabanından alındı\n")

	// Cache'e kaydet
	if err = cache.Set(ctx, cacheKey, &users, userCacheDuration); err != nil {
		// Cache hatası loglansın ama işlemi engellemeyecek
		return users, nil
	}
	return users, nil
}

func (r *UserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	exists, err := r.db.NewSelect().
		Model((*model.User)(nil)).
		Where("email = ?", email).
		Exists(ctx)

	if err != nil {
		return false, err
	}

	return exists, nil
}
