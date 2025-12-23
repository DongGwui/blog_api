# Blog API 서버 설정 가이드

## 목차
1. [요구사항](#요구사항)
2. [개발 환경 설정](#개발-환경-설정)
3. [홈 서버 배포](#홈-서버-배포)
4. [환경 변수 설정](#환경-변수-설정)
5. [서비스 관리](#서비스-관리)
6. [문제 해결](#문제-해결)

---

## 요구사항

### 필수 소프트웨어
- **Go** 1.22 이상
- **Docker** & **Docker Compose**
- **Git**

### 시스템 요구사항
- CPU: 2코어 이상
- RAM: 2GB 이상
- 디스크: 10GB 이상

---

## 개발 환경 설정

### 1. 저장소 클론
```bash
git clone https://github.com/ydonggwui/blog-api.git
cd blog-api
```

### 2. 환경 변수 설정
```bash
cp .env.example .env
# .env 파일을 편집하여 필요한 값 설정
```

### 3. Docker 컨테이너 시작
```bash
# 모든 서비스 시작 (PostgreSQL + pg_bigm, Redis, MinIO)
docker compose -f docker-compose.dev.yml up -d

# 상태 확인
docker compose -f docker-compose.dev.yml ps
```

### 4. 데이터베이스 마이그레이션
```bash
# migrate 도구 설치 (최초 1회)
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# 마이그레이션 실행
migrate -path ./migrations -database "postgres://postgres:YOUR_PASSWORD@localhost:5432/blog?sslmode=disable" up
```

### 5. MinIO 버킷 설정
```bash
# 버킷 생성 및 공개 설정
docker exec blog_minio mc alias set local http://localhost:9000 minioadmin YOUR_MINIO_PASSWORD
docker exec blog_minio mc mb local/blog-images
docker exec blog_minio mc anonymous set download local/blog-images
```

### 6. 서버 실행
```bash
# 개발 모드 실행
go run ./cmd/server

# 또는 빌드 후 실행
go build -o blog-api ./cmd/server
./blog-api
```

### 7. 접속 확인
- **API 서버**: http://localhost:8080
- **Swagger UI**: http://localhost:8080/swagger/index.html
- **Health Check**: http://localhost:8080/api/health
- **MinIO Console**: http://localhost:9001

---

## 홈 서버 배포

### 방법 1: Docker Compose (권장)

#### 1. 프로젝트 준비
```bash
# 서버에 접속
ssh user@your-home-server

# 프로젝트 클론
git clone https://github.com/ydonggwui/blog-api.git
cd blog-api
```

#### 2. 프로덕션 환경 변수 설정
```bash
cp .env.example .env
nano .env
```

**.env 파일 설정 (프로덕션용)**:
```env
# Server
PORT=8080
GIN_MODE=release

# Database
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=강력한_비밀번호_사용
DB_NAME=blog
DB_SSLMODE=disable

# Redis
REDIS_HOST=redis:6379
REDIS_PASSWORD=강력한_비밀번호_사용

# MinIO
MINIO_ENDPOINT=minio:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=강력한_비밀번호_8자이상
MINIO_BUCKET=blog-images
MINIO_USE_SSL=false
MINIO_PUBLIC_URL=http://your-domain.com:9000

# JWT
JWT_SECRET=최소_32자_이상의_랜덤_문자열
JWT_EXPIRY=24h

# Admin
ADMIN_USERNAME=admin
ADMIN_PASSWORD=강력한_비밀번호_사용
```

#### 3. 프로덕션용 docker-compose.yml 생성
```bash
nano docker-compose.prod.yml
```

```yaml
services:
  app:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: blog_api
    restart: always
    ports:
      - "8080:8080"
    env_file:
      - .env
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
      minio:
        condition: service_healthy

  postgres:
    build:
      context: ./docker/postgres
      dockerfile: Dockerfile
    image: blog-postgres-bigm:16
    container_name: blog_postgres
    restart: always
    environment:
      POSTGRES_USER: ${DB_USER}
      POSTGRES_PASSWORD: ${DB_PASSWORD}
      POSTGRES_DB: ${DB_NAME}
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 5

  redis:
    image: redis:7-alpine
    container_name: blog_redis
    restart: always
    command: redis-server --requirepass ${REDIS_PASSWORD}
    volumes:
      - redis_data:/data
    healthcheck:
      test: ["CMD", "redis-cli", "-a", "${REDIS_PASSWORD}", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5

  minio:
    image: minio/minio:latest
    container_name: blog_minio
    restart: always
    command: server /data --console-address ":9001"
    environment:
      MINIO_ROOT_USER: ${MINIO_ACCESS_KEY}
      MINIO_ROOT_PASSWORD: ${MINIO_SECRET_KEY}
    ports:
      - "9000:9000"
      - "9001:9001"
    volumes:
      - minio_data:/data
    healthcheck:
      test: ["CMD", "mc", "ready", "local"]
      interval: 10s
      timeout: 5s
      retries: 5

volumes:
  postgres_data:
  redis_data:
  minio_data:
```

#### 4. Dockerfile 생성 (아직 없는 경우)
```bash
nano Dockerfile
```

```dockerfile
# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# 의존성 설치
COPY go.mod go.sum ./
RUN go mod download

# 소스 복사 및 빌드
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o blog-api ./cmd/server

# Run stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
ENV TZ=Asia/Seoul

WORKDIR /root/

COPY --from=builder /app/blog-api .

EXPOSE 8080

CMD ["./blog-api"]
```

#### 5. 서비스 시작
```bash
# 빌드 및 시작
docker compose -f docker-compose.prod.yml up -d --build

# 로그 확인
docker compose -f docker-compose.prod.yml logs -f app
```

#### 6. 마이그레이션 실행
```bash
# 컨테이너 내에서 마이그레이션 실행
docker exec blog_api migrate -path /migrations -database "postgres://postgres:YOUR_PASSWORD@postgres:5432/blog?sslmode=disable" up
```

#### 7. MinIO 버킷 설정
```bash
docker exec blog_minio mc alias set local http://localhost:9000 minioadmin YOUR_MINIO_PASSWORD
docker exec blog_minio mc mb local/blog-images
docker exec blog_minio mc anonymous set download local/blog-images
```

---

### 방법 2: Systemd 서비스 (직접 실행)

#### 1. 바이너리 빌드
```bash
go build -o /usr/local/bin/blog-api ./cmd/server
```

#### 2. Systemd 서비스 파일 생성
```bash
sudo nano /etc/systemd/system/blog-api.service
```

```ini
[Unit]
Description=Blog API Server
After=network.target postgresql.service redis.service

[Service]
Type=simple
User=blog
Group=blog
WorkingDirectory=/opt/blog-api
ExecStart=/usr/local/bin/blog-api
Restart=always
RestartSec=5
EnvironmentFile=/opt/blog-api/.env

[Install]
WantedBy=multi-user.target
```

#### 3. 서비스 활성화 및 시작
```bash
sudo systemctl daemon-reload
sudo systemctl enable blog-api
sudo systemctl start blog-api
sudo systemctl status blog-api
```

---

## 환경 변수 설정

| 변수명 | 설명 | 기본값 | 필수 |
|--------|------|--------|------|
| `PORT` | API 서버 포트 | 8080 | ✗ |
| `GIN_MODE` | Gin 모드 (debug/release) | debug | ✗ |
| `DB_HOST` | PostgreSQL 호스트 | localhost | ✓ |
| `DB_PORT` | PostgreSQL 포트 | 5432 | ✗ |
| `DB_USER` | PostgreSQL 사용자 | postgres | ✓ |
| `DB_PASSWORD` | PostgreSQL 비밀번호 | - | ✓ |
| `DB_NAME` | 데이터베이스 이름 | blog | ✓ |
| `DB_SSLMODE` | SSL 모드 | disable | ✗ |
| `REDIS_HOST` | Redis 호스트:포트 | localhost:6379 | ✓ |
| `REDIS_PASSWORD` | Redis 비밀번호 | - | ✓ |
| `MINIO_ENDPOINT` | MinIO 엔드포인트 | localhost:9000 | ✓ |
| `MINIO_ACCESS_KEY` | MinIO 접근 키 | minioadmin | ✓ |
| `MINIO_SECRET_KEY` | MinIO 비밀 키 (8자 이상) | - | ✓ |
| `MINIO_BUCKET` | 이미지 버킷 이름 | blog-images | ✓ |
| `MINIO_USE_SSL` | MinIO SSL 사용 | false | ✗ |
| `MINIO_PUBLIC_URL` | 이미지 공개 URL | - | ✓ |
| `JWT_SECRET` | JWT 서명 키 (32자 이상) | - | ✓ |
| `JWT_EXPIRY` | JWT 만료 시간 | 24h | ✗ |
| `ADMIN_USERNAME` | 초기 관리자 아이디 | admin | ✗ |
| `ADMIN_PASSWORD` | 초기 관리자 비밀번호 | - | ✓ |

---

## 서비스 관리

### Docker Compose 명령어

```bash
# 서비스 시작
docker compose -f docker-compose.dev.yml up -d

# 서비스 중지
docker compose -f docker-compose.dev.yml down

# 서비스 재시작
docker compose -f docker-compose.dev.yml restart

# 로그 확인
docker compose -f docker-compose.dev.yml logs -f

# 특정 서비스 로그
docker compose -f docker-compose.dev.yml logs -f postgres

# 컨테이너 상태 확인
docker compose -f docker-compose.dev.yml ps

# 볼륨 포함 전체 삭제 (주의: 데이터 삭제됨)
docker compose -f docker-compose.dev.yml down -v
```

### 데이터 백업

```bash
# PostgreSQL 백업
docker exec blog_postgres pg_dump -U postgres blog > backup_$(date +%Y%m%d).sql

# PostgreSQL 복원
docker exec -i blog_postgres psql -U postgres blog < backup_20241224.sql

# MinIO 데이터 백업 (볼륨 복사)
docker run --rm -v blog_api_minio_data:/data -v $(pwd):/backup alpine tar cvf /backup/minio_backup.tar /data
```

---

## 문제 해결

### 포트 충돌
```bash
# 포트 사용 확인
lsof -i :8080
netstat -tlnp | grep 8080

# 프로세스 종료
kill -9 $(lsof -ti:8080)
```

### 데이터베이스 연결 오류
```bash
# PostgreSQL 연결 테스트
docker exec blog_postgres psql -U postgres -c "SELECT 1"

# 비밀번호 확인 (8진수 해석 방지를 위해 따옴표 사용)
# docker-compose.yml에서 POSTGRES_PASSWORD: "0303" 형태로 작성
```

### MinIO 시작 실패
```bash
# 비밀번호 길이 확인 (최소 8자)
# MINIO_ROOT_PASSWORD는 8자 이상이어야 함

# 로그 확인
docker logs blog_minio
```

### 마이그레이션 오류
```bash
# dirty 상태 해제
migrate -path ./migrations -database "postgres://..." force VERSION_NUMBER

# 버전 확인
migrate -path ./migrations -database "postgres://..." version
```

### pg_bigm 확장 오류
```bash
# pg_bigm 설치 확인
docker exec blog_postgres psql -U postgres -d blog -c "SELECT * FROM pg_extension WHERE extname = 'pg_bigm';"

# 인덱스 확인
docker exec blog_postgres psql -U postgres -d blog -c "SELECT indexname FROM pg_indexes WHERE indexname LIKE '%bigm%';"
```

---

## 빠른 시작 체크리스트

- [ ] Docker Desktop 실행
- [ ] `.env` 파일 생성 및 설정
- [ ] `docker compose up -d` 실행
- [ ] `docker compose ps`로 모든 서비스 healthy 확인
- [ ] `migrate up` 실행
- [ ] MinIO 버킷 생성 및 공개 설정
- [ ] `go run ./cmd/server` 또는 `./blog-api` 실행
- [ ] http://localhost:8080/api/health 접속 확인
- [ ] http://localhost:8080/swagger/index.html 에서 API 문서 확인
