# Blog API 개발 환경 세팅

## 사전 요구사항

- Go 1.22+
- Docker & Docker Compose
- Make (선택)
- IDE: GoLand 또는 VS Code + Go 확장

---

## 1. 프로젝트 클론 및 초기화

```bash
# 프로젝트 디렉토리로 이동
cd ~/blog

# 프로젝트 생성 (또는 git clone)
mkdir blog-api && cd blog-api

# Go 모듈 초기화
go mod init github.com/YOUR_USERNAME/blog-api
```

---

## 2. 의존성 설치

```bash
# 웹 프레임워크
go get -u github.com/gin-gonic/gin

# 데이터베이스
go get -u github.com/lib/pq
go get -u github.com/redis/go-redis/v9

# 인증
go get -u github.com/golang-jwt/jwt/v5

# 검증
go get -u github.com/go-playground/validator/v10

# 환경 변수
go get -u github.com/joho/godotenv

# MinIO
go get -u github.com/minio/minio-go/v7

# UUID
go get -u github.com/google/uuid
```

---

## 3. 개발 도구 설치

```bash
# sqlc (SQL → Go 코드 생성)
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# golang-migrate (DB 마이그레이션)
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# swag (Swagger 문서 생성)
go install github.com/swaggo/swag/cmd/swag@latest

# air (핫 리로드)
go install github.com/air-verse/air@latest
```

---

## 4. 환경 변수 설정

```bash
# .env.example 복사
cp .env.example .env
```

`.env` 파일:

```bash
# Server
PORT=8080
GIN_MODE=debug

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=blog
DB_PASSWORD=your_password
DB_NAME=blog
DB_SSLMODE=disable

# Redis
REDIS_HOST=localhost:6379
REDIS_PASSWORD=your_redis_password

# MinIO
MINIO_ENDPOINT=localhost:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=your_minio_password
MINIO_BUCKET=blog-images
MINIO_USE_SSL=false

# JWT
JWT_SECRET=your_jwt_secret_min_32_characters
JWT_EXPIRY=24h

# Admin (초기 계정)
ADMIN_USERNAME=admin
ADMIN_PASSWORD=your_admin_password
```

---

## 5. Docker Compose (로컬 개발용)

`docker-compose.dev.yml`:

```yaml
services:
  postgres:
    image: postgres:16-alpine
    container_name: blog-dev-postgres
    ports:
      - "5432:5432"
    environment:
      POSTGRES_USER: blog
      POSTGRES_PASSWORD: your_password
      POSTGRES_DB: blog
    volumes:
      - postgres_dev_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U blog"]
      interval: 5s
      timeout: 5s
      retries: 5

  redis:
    image: redis:7-alpine
    container_name: blog-dev-redis
    ports:
      - "6379:6379"
    command: redis-server --requirepass your_redis_password

  minio:
    image: minio/minio
    container_name: blog-dev-minio
    ports:
      - "9000:9000"
      - "9001:9001"
    environment:
      MINIO_ROOT_USER: minioadmin
      MINIO_ROOT_PASSWORD: your_minio_password
    command: server /data --console-address ":9001"
    volumes:
      - minio_dev_data:/data

volumes:
  postgres_dev_data:
  minio_dev_data:
```

---

## 6. 서비스 실행

```bash
# 개발용 DB, Redis, MinIO 실행
docker compose -f docker-compose.dev.yml up -d

# 상태 확인
docker compose -f docker-compose.dev.yml ps

# 로그 확인
docker compose -f docker-compose.dev.yml logs -f
```

---

## 7. 데이터베이스 마이그레이션

```bash
# 마이그레이션 파일 생성
migrate create -ext sql -dir migrations -seq init

# 마이그레이션 실행
migrate -path migrations -database "postgres://blog:your_password@localhost:5432/blog?sslmode=disable" up

# 롤백
migrate -path migrations -database "postgres://blog:your_password@localhost:5432/blog?sslmode=disable" down 1
```

---

## 8. sqlc 설정

`sqlc.yaml`:

```yaml
version: "2"
sql:
  - engine: "postgresql"
    queries: "internal/database/query.sql"
    schema: "migrations"
    gen:
      go:
        package: "sqlc"
        out: "internal/database/sqlc"
        sql_package: "pgx/v5"
        emit_json_tags: true
        emit_empty_slices: true
```

```bash
# SQL → Go 코드 생성
sqlc generate
```

---

## 9. MinIO 버킷 생성

```bash
# MinIO 클라이언트로 버킷 생성
docker exec blog-dev-minio mc alias set local http://localhost:9000 minioadmin your_minio_password
docker exec blog-dev-minio mc mb local/blog-images
docker exec blog-dev-minio mc anonymous set download local/blog-images
```

또는 http://localhost:9001 콘솔에서 수동 생성.

---

## 10. 개발 서버 실행

```bash
# 일반 실행
go run cmd/server/main.go

# 핫 리로드 (air)
air
```

`.air.toml` (air 설정):

```toml
root = "."
tmp_dir = "tmp"

[build]
cmd = "go build -o ./tmp/main ./cmd/server"
bin = "./tmp/main"
include_ext = ["go"]
exclude_dir = ["tmp", "vendor", "docs"]
delay = 1000

[log]
time = false

[misc]
clean_on_exit = true
```

---

## 11. Makefile (선택)

```makefile
.PHONY: dev run build migrate sqlc swagger test

# 개발 서버 (핫 리로드)
dev:
	air

# 일반 실행
run:
	go run cmd/server/main.go

# 빌드
build:
	go build -o bin/server cmd/server/main.go

# 마이그레이션
migrate-up:
	migrate -path migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path migrations -database "$(DB_URL)" down 1

# sqlc 생성
sqlc:
	sqlc generate

# Swagger 문서 생성
swagger:
	swag init -g cmd/server/main.go -o docs/swagger

# 테스트
test:
	go test -v ./...

# Docker 개발 환경
docker-up:
	docker compose -f docker-compose.dev.yml up -d

docker-down:
	docker compose -f docker-compose.dev.yml down
```

---

## 12. 개발 워크플로우

```
1. docker compose -f docker-compose.dev.yml up -d   # DB 실행
2. make migrate-up                                   # 마이그레이션
3. make sqlc                                         # SQL → Go
4. make dev                                          # 개발 서버 (핫 리로드)
5. http://localhost:8080/api/health                  # 확인
```

---

## 문제 해결

### DB 연결 실패
```bash
# PostgreSQL 상태 확인
docker compose -f docker-compose.dev.yml logs postgres

# 직접 연결 테스트
psql -h localhost -U blog -d blog
```

### Redis 연결 실패
```bash
docker exec blog-dev-redis redis-cli -a your_redis_password ping
```

### MinIO 연결 실패
```bash
# 콘솔 접속
open http://localhost:9001
```
