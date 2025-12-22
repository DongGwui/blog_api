# Blog API 사양서

## 개요

개인 블로그 서비스의 백엔드 API.
공개 API와 관리자 API를 단일 서버에서 라우팅으로 분리.

---

## 기술 스택

| 항목 | 기술 | 버전 |
|------|------|------|
| 언어 | Go | 1.22+ |
| 웹 프레임워크 | Gin | v1.9+ |
| SQL 생성 | sqlc | v1.25+ |
| 마이그레이션 | golang-migrate | v4 |
| 인증 | JWT (golang-jwt) | v5 |
| 검증 | go-playground/validator | v10 |
| API 문서 | swaggo/swag | |
| 환경 변수 | godotenv | |

---

## 프로젝트 구조

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
│   │   ├── post.go              # DB 접근 레이어
│   │   ├── category.go
│   │   ├── tag.go
│   │   ├── project.go
│   │   └── media.go
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
├── docs/
│   ├── SPEC.md                  # 이 문서
│   ├── SETUP.md                 # 개발 환경
│   ├── TASKS.md                 # 작업 목록
│   ├── DATABASE.md              # DB 스키마
│   └── ENDPOINTS.md             # API 명세
├── Dockerfile
├── docker-compose.dev.yml       # 로컬 개발용
├── .env.example
├── go.mod
├── go.sum
├── sqlc.yaml
└── Makefile
```

---

## API 구조

### 라우팅

```
/api
├── /public                      # 공개 API (인증 불필요)
│   ├── /posts
│   ├── /categories
│   ├── /tags
│   └── /projects
│
├── /admin                       # 관리자 API (JWT 필수)
│   ├── /auth
│   ├── /posts
│   ├── /categories
│   ├── /tags
│   ├── /projects
│   ├── /media
│   └── /dashboard
│
└── /health                      # 헬스 체크
```

### 네트워크 접근

| 경로 | 접근 | 인증 |
|------|------|------|
| /api/public/* | Cloudflare (공개) | 불필요 |
| /api/admin/* | Tailscale (비공개) | JWT 필수 |
| /api/health | 공개 | 불필요 |

---

## 핵심 기능

### 글 (Posts)
- 목록 조회 (페이지네이션, 카테고리/태그 필터)
- 상세 조회 (slug 기반)
- 검색 (pg_bigm 한국어 전문 검색)
- CRUD (관리자)
- 임시저장 / 발행 상태 관리
- 조회수 카운트 (Redis 중복 방지)

### 카테고리 & 태그
- 목록 조회 (글 개수 포함)
- CRUD (관리자)
- 글과 다대다 관계 (태그)

### 프로젝트
- 목록 조회 (정렬 순서)
- 상세 조회
- CRUD (관리자)
- 순서 변경

### 미디어
- 이미지 업로드 (MinIO)
- 이미지 목록 / 삭제
- 이미지 리사이징 (선택)

### 인증
- 관리자 로그인 (JWT 발급)
- 토큰 갱신
- 로그아웃

---

## 응답 형식

### 성공 응답

```json
{
  "data": { ... },
  "meta": {
    "page": 1,
    "per_page": 10,
    "total": 100,
    "total_pages": 10
  }
}
```

### 에러 응답

```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "Post not found"
  }
}
```

### HTTP 상태 코드

| 코드 | 의미 |
|------|------|
| 200 | 성공 |
| 201 | 생성됨 |
| 400 | 잘못된 요청 |
| 401 | 인증 필요 |
| 403 | 권한 없음 |
| 404 | 리소스 없음 |
| 500 | 서버 에러 |

---

## 외부 연동

### PostgreSQL
- 메인 데이터 저장
- pg_bigm 확장 (한국어 검색)

### Redis
- 조회수 중복 방지 (IP 기반, 24시간 TTL)
- 세션 캐시 (선택)

### MinIO
- 이미지 파일 저장
- 공개 읽기 정책

---

## 관련 문서

| 문서 | 내용 |
|------|------|
| SETUP.md | 개발 환경 세팅 |
| TASKS.md | 작업 체크리스트 |
| DATABASE.md | DB 스키마 상세 |
| ENDPOINTS.md | API 엔드포인트 명세 |
