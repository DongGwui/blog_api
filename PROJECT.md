# Blog API 프로젝트

## 프로젝트 개요

개인 블로그 서비스의 백엔드 API. Go + Gin으로 구축하며, 공개 API와 관리자 API를 단일 서버에서 라우팅으로 분리합니다.

### 핵심 정보

| 항목 | 내용 |
|------|------|
| 언어 | Go 1.22+ |
| 프레임워크 | Gin |
| 데이터베이스 | PostgreSQL 16 (공유) |
| 캐시 | Redis 7 (블로그 전용) |
| 파일 저장소 | MinIO (공유) |
| SQL 생성 | sqlc |
| 인증 | JWT |

### 서비스 연결 정보 (로컬 개발)

```
PostgreSQL: localhost:5432, DB: blog, User: postgres
Redis: localhost:6379 (비밀번호 필요)
MinIO: localhost:9000, Bucket: blog-images
```

---

## 기술 스택 상세

### 필수 의존성

```go
// 웹 프레임워크
github.com/gin-gonic/gin

// 데이터베이스
github.com/lib/pq
github.com/redis/go-redis/v9

// 인증
github.com/golang-jwt/jwt/v5

// 검증
github.com/go-playground/validator/v10

// 환경 변수
github.com/joho/godotenv

// MinIO
github.com/minio/minio-go/v7

// UUID
github.com/google/uuid
```

### 개발 도구

```bash
sqlc          # SQL → Go 코드 생성
migrate       # DB 마이그레이션
swag          # Swagger 문서 생성
air           # 핫 리로드
```

---

## 디렉토리 구조

```
blog-api/
├── cmd/
│   └── server/
│       └── main.go              # 엔트리포인트
├── internal/
│   ├── config/
│   │   └── config.go            # 환경 변수 로드
│   ├── database/
│   │   ├── db.go                # DB 연결
│   │   ├── query.sql            # SQL 쿼리 (sqlc용)
│   │   └── sqlc/                # sqlc 생성 코드
│   ├── handler/
│   │   ├── public/              # 공개 API 핸들러
│   │   │   ├── post.go
│   │   │   ├── category.go
│   │   │   ├── tag.go
│   │   │   └── project.go
│   │   └── admin/               # 관리자 API 핸들러
│   │       ├── auth.go
│   │       ├── post.go
│   │       ├── category.go
│   │       ├── tag.go
│   │       ├── project.go
│   │       └── media.go
│   ├── middleware/
│   │   ├── auth.go              # JWT 인증
│   │   ├── cors.go              # CORS 설정
│   │   └── logger.go            # 요청 로깅
│   ├── model/
│   │   ├── post.go              # 요청/응답 구조체
│   │   ├── category.go
│   │   ├── tag.go
│   │   ├── project.go
│   │   └── media.go
│   ├── repository/
│   │   └── ...                  # DB 접근 레이어
│   ├── service/
│   │   ├── post.go              # 비즈니스 로직
│   │   ├── auth.go
│   │   ├── image.go             # MinIO 연동
│   │   └── view.go              # 조회수 (Redis)
│   └── router/
│       └── router.go            # 라우터 설정
├── migrations/
│   ├── 000001_init.up.sql
│   └── 000001_init.down.sql
├── docs/                        # 상세 문서
├── Dockerfile
├── docker-compose.dev.yml
├── .env.example
├── sqlc.yaml
├── Makefile
└── go.mod
```

---

## API 라우팅 구조

```
/api
├── /health                      # 헬스 체크
├── /public                      # 공개 API (인증 불필요)
│   ├── GET  /posts              # 글 목록 (페이지네이션)
│   ├── GET  /posts/:slug        # 글 상세
│   ├── GET  /posts/search       # 검색
│   ├── POST /posts/:slug/view   # 조회수 증가
│   ├── GET  /categories         # 카테고리 목록
│   ├── GET  /tags               # 태그 목록
│   ├── GET  /projects           # 프로젝트 목록
│   └── GET  /projects/:slug     # 프로젝트 상세
│
└── /admin                       # 관리자 API (JWT 필수)
    ├── POST /auth/login         # 로그인
    ├── GET  /auth/me            # 현재 사용자
    ├── CRUD /posts              # 글 관리
    ├── CRUD /categories         # 카테고리 관리
    ├── CRUD /tags               # 태그 관리
    ├── CRUD /projects           # 프로젝트 관리
    ├── /media                   # 미디어 관리
    └── GET  /dashboard/stats    # 대시보드 통계
```

---

## 데이터베이스 스키마

### 테이블 목록

| 테이블 | 용도 |
|--------|------|
| admins | 관리자 계정 (단일) |
| categories | 카테고리 |
| tags | 태그 |
| posts | 블로그 글 |
| post_tags | 글-태그 연결 (다대다) |
| projects | 포트폴리오 프로젝트 |
| media | 업로드된 미디어 |

### 주요 테이블 구조

**posts**
```sql
id, title, slug, content, excerpt, category_id, status,
view_count, reading_time, thumbnail,
created_at, updated_at, published_at
```

**projects**
```sql
id, title, slug, description, content, tech_stack (JSONB),
demo_url, github_url, thumbnail, images (JSONB),
is_featured, sort_order, created_at, updated_at
```

---

## 응답 형식

### 성공 (단일)
```json
{
  "data": { ... }
}
```

### 성공 (목록)
```json
{
  "data": [ ... ],
  "meta": {
    "page": 1,
    "per_page": 10,
    "total": 100,
    "total_pages": 10
  }
}
```

### 에러
```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "Post not found"
  }
}
```

---

## 환경 변수

```bash
# Server
PORT=8080
GIN_MODE=debug

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=비밀번호
DB_NAME=blog
DB_SSLMODE=disable

# Redis
REDIS_HOST=localhost:6379
REDIS_PASSWORD=비밀번호

# MinIO
MINIO_ENDPOINT=localhost:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=비밀번호
MINIO_BUCKET=blog-images
MINIO_USE_SSL=false

# JWT
JWT_SECRET=최소32자이상의시크릿키
JWT_EXPIRY=24h

# Admin
ADMIN_USERNAME=admin
ADMIN_PASSWORD=초기비밀번호
```

---

## 개발 환경 실행

### 사전 조건
- 공유 인프라 실행 중 (PostgreSQL, MinIO)
- 블로그 인프라 실행 중 (Redis)

### 실행 명령어

```bash
# 의존성 설치
go mod tidy

# DB 마이그레이션
make migrate-up

# sqlc 코드 생성
make sqlc

# 개발 서버 (핫 리로드)
make dev

# 또는 일반 실행
go run cmd/server/main.go
```

### 확인
```
http://localhost:8080/api/health
```

---

## 개발 규칙

### 코드 스타일
- 표준 Go 컨벤션 준수
- 에러는 항상 처리 (무시하지 않음)
- 핸들러 → 서비스 → 레포지토리 레이어 분리

### 커밋 메시지
```
feat: 글 목록 API 구현
fix: 조회수 중복 카운트 버그 수정
docs: API 문서 업데이트
refactor: 인증 미들웨어 개선
```

### API 설계
- RESTful 규칙 준수
- 복수형 명사 사용 (/posts, /categories)
- 적절한 HTTP 상태 코드 반환

---

## 참고 문서

| 파일 | 내용 |
|------|------|
| docs/SPEC.md | 상세 사양서 |
| docs/SETUP.md | 개발 환경 세팅 |
| docs/TASKS.md | 작업 체크리스트 |
| docs/DATABASE.md | DB 스키마 상세 |
| docs/ENDPOINTS.md | API 엔드포인트 명세 |
