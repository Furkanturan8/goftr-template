# GOFTR-V1 Projesi

Bu proje, modern bir web uygulaması geliştirmek için Go ve Vue.js teknolojilerini kullanan full-stack bir projedir.

## Proje Yapısı

Proje iki ana bölümden oluşmaktadır:

### Backend (`/backend`)

Go dilinde Clean Architecture prensiplerine uygun olarak geliştirilmiş REST API. Backend [Furkan Turan](https://github.com/furkanturan8) tarafından geliştirilmiştir.

Özellikler:

- Go 1.21+
- Fiber web framework
- PostgreSQL veritabanı
- Redis önbellek sistemi
- JWT tabanlı kimlik doğrulama
- Clean Architecture
- Docker desteği
- Otomatik kod üretme araçları
- Detaylı loglama sistemi

Daha fazla bilgi için: [Backend README](/backend/README.md)

### Frontend (`/frontend`)

Vue.js tabanlı modern ve responsive admin paneli. Frontend tasarımı [ThemeSelection](https://themeselection.com/)'ın Sneat admin template'i kullanılarak geliştirilmiştir.

Özellikler:

- Vue 3
- Vuetify 3
- Vite build sistemi
- File-based routing
- Component auto-importing
- TypeScript desteği
- I18n desteği
- Modern ve kullanıcı dostu arayüz

Daha fazla bilgi için: [Frontend README](/frontend/README.md)

## Başlangıç

### Gereksinimler

- Go 1.21 veya üzeri
- Node.js 16 veya üzeri
- Docker ve Docker Compose
- PostgreSQL 15
- Redis 7

### Kurulum

1. Projeyi klonlayın:

```bash
git clone [repo-url]
cd goftr-v1
```

2. Backend servislerini başlatın:

```bash
cd backend
docker-compose up -d
go run cmd/server/main.go
```

3. Frontend uygulamasını başlatın:

```bash
cd frontend
npm install
npm run dev
```

## Geliştirme

Her iki proje de kendi içinde bağımsız olarak geliştirilebilir. Detaylı geliştirme kılavuzları için ilgili projelerin README dosyalarına başvurun:

- [Backend Geliştirme Kılavuzu](/backend/README.md)
- [Frontend Geliştirme Kılavuzu](/frontend/README.md)

## Lisans

### Backend

Backend kısmı [Furkan Turan](https://github.com/furkanturan8) tarafından geliştirilmiş olup, MIT lisansı altında lisanslanmıştır. Detaylar için [LICENSE](LICENSE) dosyasına bakın.

### Frontend

Frontend tasarımı [ThemeSelection](https://themeselection.com/)'ın Sneat admin template'i kullanılarak geliştirilmiştir.

- Copyright © [ThemeSelection](https://themeselection.com/)
- MIT Lisansı altında lisanslanmıştır
- ThemeSelection'ın ücretsiz ürünleri MIT lisansı altında açık kaynak kodludur. Bu ürünleri hem kişisel hem de ticari amaçlar için kullanabilirsiniz. Tek gereken, aşağıdaki bağlantıyı web uygulamanızın veya projenizin alt bilgisine eklemenizdir:

```html
<a href="https://themeselection.com/">ThemeSelection</a>
```
