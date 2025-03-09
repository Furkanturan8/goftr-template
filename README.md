# Go Clean Architecture + Fiber + PostgreSQL + Redis Proje Şablonu

Bu proje, Go dilinde Clean Architecture prensiplerine uygun olarak geliştirilmiş, Fiber web framework'ü, PostgreSQL veritabanı ve Redis önbellek sistemi kullanan bir web uygulaması şablonudur.

## 1. Proje Yapısı

```
├── cmd/
│   └── api/
│       └── main.go           # Ana uygulama giriş noktası
├── config/
│   ├── config.go            # Yapılandırma yapıları
│   └── config.yaml          # Yapılandırma dosyası
├── internal/
│   ├── model/              # Veritabanı modelleri
│   ├── repository/         # Veritabanı işlemleri
│   ├── service/           # İş mantığı katmanı
│   ├── handler/           # HTTP işleyicileri
│   ├── dto/              # Veri transfer nesneleri
│   ├── middleware/       # HTTP ara yazılımları
│   └── router/          # Router yapılandırmaları
├── pkg/
│   ├── cache/           # Redis cache işlemleri
│   ├── errorx/         # Hata yönetimi
│   ├── jwt/            # JWT işlemleri
│   ├── logger/         # Loglama işlemleri
│   ├── query/          # Query işlemleri
│   └── response/       # Response işlemleri
├── migrations/         # Veritabanı migrasyon dosyaları
├── api/               # API dökümantasyonu
├── tests/            # Test dosyaları
├── logs/             # Log dosyaları
├── Dockerfile        # Docker yapılandırması
├── docker-compose.yml # Docker servisleri
└── go.mod            # Go modül tanımlamaları
```

## 2. Başlangıç

### 2.1. Gereksinimler

- Go 1.21 veya üzeri
- Docker ve Docker Compose
- PostgreSQL 15
- Redis 7

### 2.2. Projeyi Çalıştırma

```bash
# Projeyi klonlama
git clone [repo-url]
cd [proje-dizini]

# Bağımlılıkları yükleme
go mod download

# Docker servisleri başlatma
docker-compose up -d

# Uygulamayı başlatma
go run cmd/api/main.go
```

## 3. Yapılandırma

### 3.1. Temel Yapılandırma (config/config.yaml)

```yaml
app:
  name: "goftr-v1"
  port: 3000
  env: "development"

database:
  host: "localhost"
  port: 5432
  user: "postgres"
  password: "postgres"
  name: "goftr"

redis:
  host: "localhost"
  port: 6379
  password: ""
  db: 0
  pool_size: 10
```

### 3.2. Docker Compose Yapılandırması

```yaml
services:
  app:
    build: .
    ports:
      - "3000:3000"
    depends_on:
      - postgres
      - redis

  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: goftr

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
```

## 4. Temel Kullanım

### 4.1. Model Oluşturma

```go
// internal/model/user.go
type User struct {
    ID        int64     `json:"id" bun:",pk,autoincrement"`
    Email     string    `json:"email" bun:",unique,notnull"`
    Password  string    `json:"-" bun:",notnull"`
    FirstName string    `json:"first_name"`
    LastName  string    `json:"last_name"`
    CreatedAt time.Time `json:"created_at"`
}
```

### 4.2. Repository Katmanı

```go
// internal/repository/user_repository.go
type UserRepository struct {
    db *bun.DB
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*model.User, error) {
    cacheKey := fmt.Sprintf("user:%d", id)

    // Önce cache'den kontrol
    var user model.User
    if err := cache.Get(ctx, cacheKey, &user); err == nil {
        return &user, nil
    }

    // DB'den al ve cache'e kaydet
    // ...
}
```

### 4.3. Service Katmanı

```go
// internal/service/user_service.go
type UserService struct {
    repo *repository.UserRepository
}

func (s *UserService) CreateUser(ctx context.Context, user *model.User) error {
    // İş mantığı implementasyonu
    return s.repo.Create(ctx, user)
}
```

### 4.4. Handler Katmanı

```go
// internal/handler/user_handler.go
type UserHandler struct {
    service *service.UserService
}

func (h *UserHandler) RegisterRoutes(app *fiber.App) {
    app.Post("/users", h.Create)
    app.Get("/users/:id", h.GetByID)
}
```

## 5. Redis Cache Kullanımı

### 5.1. Cache İşlemleri

```go
// Veri kaydetme
cache.Set(ctx, "key", value, 24*time.Hour)

// Veri okuma
var result Type
cache.Get(ctx, "key", &result)

// Veri silme
cache.Delete(ctx, "key")

// Pattern ile silme
cache.DeleteMany(ctx, "user:*")
```

## 6. API Endpoint Örnekleri

### 6.1. Kullanıcı İşlemleri

```bash
# Kullanıcı oluşturma
POST /api/users
{
    "email": "user@example.com",
    "password": "password123",
    "first_name": "John",
    "last_name": "Doe"
}

# Kullanıcı getirme
GET /api/users/:id

# Kullanıcı güncelleme
PUT /api/users/:id

# Kullanıcı silme
DELETE /api/users/:id
```

## 7. Veritabanı İşlemleri

### 7.1. Migration Oluşturma

```sql
-- migrations/001_create_users.sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## 8. Güvenlik ve Performans

### 8.1. Redis Cache Stratejileri

- Sık erişilen veriler için cache kullanımı
- Cache süresinin doğru belirlenmesi
- Cache invalidation stratejileri

### 8.2. Güvenlik Önlemleri

- JWT token kullanımı
- Password hashing
- Rate limiting
- CORS yapılandırması

## 9. Deployment

### 9.1. Production Ortamı

```bash
# Production build
docker-compose -f docker-compose.prod.yml up -d
```

### 9.2. Monitoring (Henüz eklenmedi)

- Prometheus metrics
- Grafana dashboard
- Log aggregation

## 10. Katkıda Bulunma

1. Fork'layın
2. Feature branch oluşturun (`git checkout -b feature/amazing-feature`)
3. Değişikliklerinizi commit edin (`git commit -m 'feat: add amazing feature'`)
4. Branch'inizi push edin (`git push origin feature/amazing-feature`)
5. Pull Request oluşturun
