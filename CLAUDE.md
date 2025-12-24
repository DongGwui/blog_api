# Blog API 프로젝트

개인 블로그 백엔드 API 서버입니다.

## 기술 스택

- **언어**: Go 1.24
- **프레임워크**: Gin
- **데이터베이스**: PostgreSQL 16 + pg_bigm (한국어 검색)
- **캐시**: Redis 7
- **스토리지**: MinIO (S3 호환)
- **ORM**: sqlc (타입 안전 SQL)
- **인증**: JWT
- **문서화**: Swagger (swaggo)

## 프로젝트 구조

```
blog_api/
├── cmd/server/          # 애플리케이션 진입점
├── internal/
│   ├── config/          # 설정 관리
│   ├── database/        # DB 연결 및 sqlc 생성 코드
│   ├── handler/         # HTTP 핸들러 (admin/, public/)
│   ├── middleware/      # 미들웨어 (auth, cors, logger)
│   ├── model/           # 요청/응답 모델
│   ├── router/          # 라우터 설정
│   ├── service/         # 비즈니스 로직
│   └── util/            # 유틸리티 함수
├── migrations/          # DB 마이그레이션
├── docker/postgres/     # PostgreSQL + pg_bigm Dockerfile
└── docs/                # 문서
```

## 주요 API 엔드포인트

- `GET /api/health` - 헬스 체크
- `/api/public/*` - 공개 API (인증 불필요)
- `/api/admin/*` - 관리자 API (JWT 인증 필요)

## 개발 명령어

```bash
# Docker 컨테이너 시작
docker compose -f docker-compose.dev.yml up -d

# 마이그레이션 실행
migrate -path ./migrations -database "postgres://postgres:0303@localhost:5432/blog?sslmode=disable" up

# 서버 실행
go run ./cmd/server

# 테스트 실행
go test ./...

# sqlc 코드 생성
sqlc generate

# Swagger 문서 생성
swag init -g cmd/server/main.go -o docs/swagger
```

## 환경 변수

`.env` 파일 참조. 필수 항목:
- `DB_PASSWORD` - PostgreSQL 비밀번호
- `REDIS_PASSWORD` - Redis 비밀번호
- `MINIO_SECRET_KEY` - MinIO 비밀번호 (8자 이상)
- `JWT_SECRET` - JWT 서명 키 (32자 이상)
- `ADMIN_PASSWORD` - 초기 관리자 비밀번호

## 참고 문서

- `docs/SERVER_GUIDE.md` - 서버 설정 및 배포 가이드
- `docs/swagger/` - API 문서 (Swagger UI: http://localhost:8080/swagger/index.html)
