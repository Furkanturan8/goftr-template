# Build aşaması
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Gerekli paketleri yükle
RUN apk add --no-cache gcc musl-dev

# Go modüllerini kopyala ve indir
COPY go.mod go.sum ./
RUN go mod download

# Kaynak kodları kopyala
COPY . .

# Uygulamayı derle
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/api/main.go

# Çalışma aşaması
FROM alpine:latest

WORKDIR /app

# SSL sertifikaları için gerekli
RUN apk --no-cache add ca-certificates

# Builder aşamasından derlenmiş uygulamayı kopyala
COPY --from=builder /app/main .
COPY --from=builder /app/config/config.yaml ./config/

# Uygulama için gerekli dizinleri oluştur
RUN mkdir -p /app/logs

# Uygulamayı çalıştır
CMD ["./main"] 