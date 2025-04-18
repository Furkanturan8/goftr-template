# Go Clean Architecture + Fiber + PostgreSQL + Redis Proje Şablonu

Bu proje, Go dilinde Clean Architecture prensiplerine uygun olarak geliştirilmiş, Fiber web framework'ü, PostgreSQL veritabanı ve Redis önbellek sistemi kullanan bir web uygulaması şablonudur.

## 1. Proje Yapısı

```
├── cmd/
│   └── server/
│       └── main.go         # Ana uygulama giriş noktası
│   └── migrate/
│       └── main.go         # DB Migrate başlangıcı
├── config/
│    └── config.go          # Yapılandırma kodları
├── internal/
│   ├── model/              # Veritabanı modelleri
│   ├── repository/         # Veritabanı işlemleri
│   ├── service/            # İş mantığı katmanı
│   ├── handler/            # HTTP işleyicileri
│   ├── dto/                # Veri transfer nesneleri
│   ├── middleware/         # HTTP ara yazılımları
│   └── router/             # Router yapılandırmaları
├── pkg/
│   ├── cache/              # Redis cache işlemleri
│   ├── errorx/             # Hata yönetimi
│   ├── jwt/                # JWT işlemleri
│   ├── logger/             # Loglama işlemleri
│   ├── monitoring/         # Monitoring işlemleri  
│   ├── query/              # Query işlemleri
│   └── response/           # Response işlemleri
├── migrations/             # Veritabanı migrasyon dosyaları
├── mock_data/              # generative-mock.sh ile oluşturulmuş mock veriler 
├── tests/                  # Test dosyaları
├── logs/                   # Log dosyaları
├── Dockerfile              # Docker yapılandırması
├── docker-compose.yml      # Docker servisleri
├── generative-structure.sh # Otomatik dosya oluşturma scripti
├── generative-mock.sh      # Otomatik mock oluşturma scripti
├── .env                    # Ortam değişkenleri
└── go.mod                  # Go modül tanımlamaları
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
go mod init goftr-template/backend
go mod tidy

# Docker servisleri başlatma
docker-compose up -d

# Uygulamayı başlatma
go run cmd/api/main.go
```

## 3. Yapılandırma

### 3.1. Temel Yapılandırma (.env)

```dotenv
# App
APP_NAME=goftr-v1
APP_VERSION=1.0.0
APP_ENV=development
APP_PORT=3005
APP_SHUTDOWN_TIMEOUT=10
APP_LOG_DIR=./logs

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=goftr
DB_SSLMODE=disable

# Other environments
# Redis, JWT, Prometheus, Grafana etc
# Check .env.example for more details
```

### 3.2. Docker Compose Yapılandırması

```yaml
services:
 postgres:
 image: postgres:15-alpine
 container_name: goftr-postgres
 environment:
  - POSTGRES_USER=${DB_USER}
  - POSTGRES_PASSWORD=${DB_PASSWORD}
  - POSTGRES_DB=${DB_NAME}
 ports:
  - "${DB_PORT}:5432"
 volumes:
  - postgres_data:/var/lib/postgresql/data
 networks:
  - goftr-network
 healthcheck:
  test: ["CMD-SHELL", "pg_isready -U ${DB_USER} -d ${DB_NAME}"]
  interval: 5s
  timeout: 5s
  retries: 5
 restart: unless-stopped

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
type IUserRepository interface {
    GetByID(ctx context.Context, id int64) (*model.User, error) 
}

type UserRepository struct {
    db *bun.DB
}

func NewUserRepository(db *bun.DB) IUserRepository {
	return &UserRepository{db: db}
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
	userRepo repository.IUserRepository
}

func NewUserService(u repository.IUserRepository) *UserService {
	return &UserService{
		userRepo: u,
	}
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

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{service: s}
}
```

### 4.5. Otomatik Dosya Oluşturma (generative-structure.sh)

Yukarıda kendimiz oluşturduğumuz dosyaları/kodları (model, dto, handler, service, repository) otomatik olarak oluşturmak için `generative-structure.sh` (kendi yapımıza göre tasarlanmış shell script kodları mevcuttur) dosyasını kullanabilirsiniz.

```shell
  ./generate-structure.sh Product 'Name string' 'Description string' 'Price float64'
```

- -h veya --help komutu 
	<img width="1018" alt="1" src="https://github.com/user-attachments/assets/a98a3250-023d-4763-b8b1-8b970917934f" />

- Product örneği için komutu çalıştırma
	<img width="1018" alt="2" src="https://github.com/user-attachments/assets/031a0fd0-f971-4719-b15a-505c2e04b718" />

- Daha önceden var olan Product modeli için tekrar deneme yapalım
	<img width="1018" alt="3" src="https://github.com/user-attachments/assets/aae9f0c8-6df3-4a7d-a783-bc8b0a4484a8" />

### 4.6. cmd/api/main.go Katmanı

```go

func main(){
	// Yapılandırmayı yükle
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Printf("Config yükleme hatası: %v", err)
		os.Exit(1)
	}

        // Bazı fonxların initialize'ı yapılıyor
        // Logger'ı başlat, Redis cache'i başlat, JWT yapılandırmasını başlat,  Database bağlantısı

  	// Router'ı oluştur ve yapılandır
	r := router.NewRouter(db, cfg)
	r.SetupRoutes()

        // Graceful shutdown ve sunucu açılması/kapatılması kodları
}

```

### 4.7. Router Katmanı

```go
type Router struct {
	app *fiber.App
	db  *bun.DB
	cfg *config.Config
}

func NewRouter(db *bun.DB, cfg *config.Config) *Router {
	return &Router{
		app: fiber.New(),
		db:  db,
		cfg: cfg,
	}
}

func (r *Router) SetupRoutes() {
	// Middleware'leri ekle
	r.app.Use(logger.New())
	r.app.Use(recover.New())
	r.app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:63342,http://localhost:3005,http://localhost:5173",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Content-Type, Authorization",
	}))

	// API versiyonu
	api := r.app.Group("/api")
	v1 := api.Group("/v1")

	// Repository'ler
	userRepo := repository.NewUserRepository(r.db)
	authRepo := repository.NewAuthRepository(r.db)
	
	// Service'ler
	authService := service.NewAuthService(authRepo, userRepo)
	userService := service.NewUserService(userRepo)
	
	// Handler'lar
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	
	// Auth routes
	auth := v1.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/refresh", authHandler.RefreshToken)
	auth.Post("/forgot-password", authHandler.ForgotPassword)
	auth.Post("/reset-password", authHandler.ResetPassword)
	auth.Post("/logout", middleware.AuthMiddleware(), authHandler.Logout)

	// User routes - Base group
	users := v1.Group("/users")

	// Normal user routes (profil yönetimi)
	userProfile := users.Group("/me")
	userProfile.Use(middleware.AuthMiddleware()) // Sadece authentication gerekli
	userProfile.Get("/", userHandler.GetProfile)
	userProfile.Put("/", userHandler.UpdateProfile)
	
	// Admin only routes
	adminUsers := users.Group("/")
	adminUsers.Use(middleware.AuthMiddleware(), middleware.AdminOnly()) // Admin yetkisi gerekli
	adminUsers.Post("/", userHandler.Create)
	adminUsers.Get("/", userHandler.List)
	adminUsers.Get("/:id", userHandler.GetByID)
	adminUsers.Put("/:id", userHandler.Update)
	adminUsers.Delete("/:id", userHandler.Delete)

	// Diğer route grupları buraya eklenecek
}

func (r *Router) GetApp() *fiber.App {
	return r.app
}
```

### 4.8. Otomatik Mock Dataları Oluşturma (generative-mock.sh)

Model ismi ve sayısını girerek otomatik mock json dataları oluşturur. 
generative-mock.sh dosyası sayesinde frontend geliştirme ve test süreçlerinde hızlı ve pratik kullanım imkânı sunar.

- Kullanımı: 

 	<img width="535" alt="Ekran Resmi 2025-04-05 00 16 20" src="https://github.com/user-attachments/assets/7f579a59-9d79-4890-b877-c414ebb237bf" />

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

### 6.2. Api istekleri test/api_test.html 

Api isteklerini api_test.html sayfasında deneyebilirsiniz

## 7. Veritabanı İşlemleri

### 7.1. Tablo oluşturma SQL Kodları

```sql
-- migrations/sql/000001_create_users.sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```
### 7.2. Migration Komutları Uygulanması

```bash
  # Tüm migration'ları uygula
go run cmd/migrate/main.go -action up

# Son 2 migration'ı geri al
go run cmd/migrate/main.go -action down -step 2

# Migration listesini görüntüle
go run cmd/migrate/main.go -action status
```

## 8. Güvenlik ve Performans

### 8.1. Redis Cache Stratejileri

- Sık erişilen veriler için cache kullanımı
- Cache süresinin doğru belirlenmesi
- Cache invalidation stratejileri

### 8.2. Güvenlik Önlemleri

Güvenlik, bir uygulamanın en önemli bileşenlerinden biridir. Bu bölümde uygulamada alınan temel güvenlik önlemleri listelenmiştir.

#### 1. JWT Token Kullanımı
Kullanıcı oturum yönetimi ve kimlik doğrulama işlemleri **JSON Web Token (JWT)** ile sağlanmaktadır.
- Kullanıcı giriş yaptığında kendisine şifrelenmiş bir JWT token verilir.
- Token, her istekte doğrulanarak kullanıcının yetkilendirilmesi sağlanır.
- Token süresi dolduğunda kullanıcı tekrar giriş yapmalıdır.

#### 2. Şifreleme (Password Hashing)
Kullanıcı şifreleri **bcrypt** algoritması kullanılarak hashlenir ve güvenli bir şekilde veritabanında saklanır.

Şifreleme adımları:
1. Kullanıcı kayıt olurken şifre **bcrypt** ile hashlenir.
2. Giriş yaparken girilen şifre, veritabanındaki hashlenmiş şifre ile karşılaştırılır.
3. Eşleşme sağlanırsa kullanıcı giriş yapar, aksi takdirde hata döndürülür.

#### 3. Rate Limiting
**Rate limiting** sayesinde belirli bir zaman aralığında yapılan istekler sınırlandırılarak kötüye kullanımın (Brute-force saldırıları, DDoS vb.) önüne geçilir.

**30 sn de 10 istek sınırı** belirlenmiştir.
#### Örnek Log Kayıtları
```
16:17:58 | 200 |  103.272625ms | 127.0.0.1 | POST | /api/v1/auth/login | -
16:18:02 | 204 |      14.042µs | 127.0.0.1 | OPTIONS | /api/v1/users | -
16:18:02 | 200 |      3.2325ms | 127.0.0.1 | GET | /api/v1/users | Kullanıcılar veritabanından alındı
16:18:04 | 200 |    1.126375ms | 127.0.0.1 | GET | /api/v1/users | Kullanıcılar cache'den alındı (x9 tane)
16:18:06 | 429 |      12.583µs | 127.0.0.1 | GET | /api/v1/users | Rate limit aşıldı
16:18:19 | 429 |      34.792µs | 127.0.0.1 | GET | /api/v1/users | Rate limit aşıldı
```
- **HTTP 200**: Başarılı istekler
- **HTTP 429**: Rate limit aşıldığında dönen hata kodu

Rate limiting sayesinde belirlenen sınır aşıldığında kullanıcıya **HTTP 429 - Too Many Requests** hatası döndürülerek saldırılar önlenir.

#### 4. CORS (Cross-Origin Resource Sharing) Yapılandırması
**CORS** sayesinde yalnızca belirli domainlerden gelen isteklerin kabul edilmesi sağlanır.
- Yanlış yapılandırılmış CORS, güvenlik açıklarına yol açabilir.
- Güvenli olmayan kaynaklardan gelen istekler engellenir.

Örnek **CORS Konfigürasyonu:**
```go
app.Use(cors.New(cors.Config{
    AllowOrigins: "https://example.com, https://api.example.com",
    AllowMethods: "GET,POST,PUT,DELETE",
    AllowHeaders: "Content-Type, Authorization",
}))
```

---
Bu önlemler, uygulamanın güvenliğini artırmaya yardımcı olur ve yetkisiz erişimleri engeller.

## 9. Deployment

### 9.1. Production Ortamı

```bash
# Production build
docker-compose -f docker-compose.prod.yml up -d
```

### 9.2. Monitoring

- Prometheus metrics
- Grafana dashboard
- Log aggregation

Grafana Dashboard:  
<img width="1200" alt="metrics-1" src="https://github.com/user-attachments/assets/1a601a68-0eee-46e9-a1e9-93b7f152608c" />

Prometheus Metrics:
<img width="1426" alt="prometheus" src="https://github.com/user-attachments/assets/543ff99a-3a80-45ca-b640-7ca61bff073a" />

## 10. TODOS
1. Monitoring eklenecek &emsp; [✓]
2. Rate Limiting eklenecek &emsp; [✓]
3. Generate-structure.sh dosyası eklenecek (otomatik dosya oluşturucu) &emsp; [✓]
4. Frontend veya mobil app için mock datalar oluşturma (model ismi girilerek bir generate_mock.sh dosyası yardımıyla) eklenebilir. &emsp; [✓]
5. Migration ve db işlemleri düzenlenecek &emsp; [✓]

## 11. Katkıda Bulunma

1. Fork'layın
2. Feature branch oluşturun (`git checkout -b feature/amazing-feature`)
3. Değişikliklerinizi commit edin (`git commit -m 'feat: add amazing feature'`)
4. Branch'inizi push edin (`git push origin feature/amazing-feature`)
5. Pull Request oluşturun
