# Blog API 작업 체크리스트

## Phase 1: 프로젝트 초기화

### 기본 설정
- [ ] Go 모듈 초기화
- [ ] 의존성 설치
- [ ] 디렉토리 구조 생성
- [ ] .env.example 작성
- [ ] .gitignore 작성
- [ ] docker-compose.dev.yml 작성
- [ ] Makefile 작성

### 개발 도구
- [ ] sqlc 설치 및 설정
- [ ] golang-migrate 설치
- [ ] air (핫 리로드) 설치 및 설정
- [ ] swag (Swagger) 설치

### 인프라 연동
- [ ] Docker 개발 환경 실행 (postgres, redis, minio)
- [ ] DB 연결 테스트
- [ ] Redis 연결 테스트
- [ ] MinIO 버킷 생성

---

## Phase 2: 기반 코드

### 설정 및 연결
- [ ] config.go - 환경 변수 로드
- [ ] database/db.go - PostgreSQL 연결
- [ ] redis 클라이언트 설정
- [ ] minio 클라이언트 설정

### 라우터 및 미들웨어
- [ ] router.go - 기본 라우터 설정
- [ ] middleware/cors.go - CORS 설정
- [ ] middleware/logger.go - 요청 로깅
- [ ] middleware/auth.go - JWT 인증 미들웨어

### 공통 유틸
- [ ] 응답 헬퍼 (성공/에러 응답 포맷)
- [ ] 페이지네이션 헬퍼
- [ ] 슬러그 생성 유틸
- [ ] 읽기 시간 계산 유틸

### 헬스 체크
- [ ] GET /api/health 엔드포인트
- [ ] DB, Redis, MinIO 연결 상태 확인

---

## Phase 3: 데이터베이스

### 마이그레이션
- [ ] 000001_init.up.sql - 초기 스키마
- [ ] 000001_init.down.sql - 롤백
- [ ] 마이그레이션 실행 테스트

### sqlc 쿼리 작성
- [ ] posts 쿼리 (CRUD, 검색, 필터)
- [ ] categories 쿼리
- [ ] tags 쿼리
- [ ] post_tags 쿼리 (다대다)
- [ ] projects 쿼리
- [ ] media 쿼리
- [ ] admins 쿼리

### sqlc 생성
- [ ] sqlc generate 실행
- [ ] 생성된 코드 확인

---

## Phase 4: 인증 (Admin)

### 모델
- [ ] 로그인 요청/응답 구조체
- [ ] JWT 클레임 구조체

### 서비스
- [ ] 비밀번호 해싱 (bcrypt)
- [ ] JWT 토큰 생성
- [ ] JWT 토큰 검증

### 핸들러
- [ ] POST /api/admin/auth/login
- [ ] POST /api/admin/auth/logout
- [ ] GET /api/admin/auth/me

### 초기 관리자
- [ ] 환경 변수로 초기 계정 생성 (앱 시작 시)

---

## Phase 5: 글 (Posts)

### Public API
- [ ] GET /api/public/posts - 목록 (페이지네이션)
- [ ] GET /api/public/posts/:slug - 상세
- [ ] GET /api/public/posts/search - 검색
- [ ] POST /api/public/posts/:slug/view - 조회수 증가

### Admin API
- [ ] GET /api/admin/posts - 목록 (전체, 상태 필터)
- [ ] GET /api/admin/posts/:id - 상세
- [ ] POST /api/admin/posts - 생성
- [ ] PUT /api/admin/posts/:id - 수정
- [ ] DELETE /api/admin/posts/:id - 삭제
- [ ] PATCH /api/admin/posts/:id/publish - 발행 상태 변경

### 비즈니스 로직
- [ ] 슬러그 자동 생성
- [ ] 읽기 시간 계산
- [ ] 태그 연결 처리
- [ ] 발행일 자동 설정

### 조회수 (Redis)
- [ ] IP 해시 기반 중복 체크
- [ ] 24시간 TTL 설정
- [ ] DB 업데이트 (새 조회 시)

---

## Phase 6: 카테고리 & 태그

### Categories Public API
- [ ] GET /api/public/categories - 목록
- [ ] GET /api/public/categories/:slug/posts - 카테고리별 글

### Categories Admin API
- [ ] GET /api/admin/categories
- [ ] POST /api/admin/categories
- [ ] PUT /api/admin/categories/:id
- [ ] DELETE /api/admin/categories/:id

### Tags Public API
- [ ] GET /api/public/tags - 목록
- [ ] GET /api/public/tags/:slug/posts - 태그별 글

### Tags Admin API
- [ ] GET /api/admin/tags
- [ ] POST /api/admin/tags
- [ ] PUT /api/admin/tags/:id
- [ ] DELETE /api/admin/tags/:id

---

## Phase 7: 프로젝트 (Projects)

### Public API
- [ ] GET /api/public/projects - 목록
- [ ] GET /api/public/projects/:slug - 상세

### Admin API
- [ ] GET /api/admin/projects
- [ ] POST /api/admin/projects
- [ ] PUT /api/admin/projects/:id
- [ ] DELETE /api/admin/projects/:id
- [ ] PATCH /api/admin/projects/reorder - 순서 변경

---

## Phase 8: 미디어 (Media)

### MinIO 연동
- [ ] 이미지 업로드 서비스
- [ ] 파일명 UUID 생성
- [ ] 경로 생성 (년/월/파일명)

### Admin API
- [ ] GET /api/admin/media - 목록
- [ ] POST /api/admin/media/upload - 업로드
- [ ] DELETE /api/admin/media/:id - 삭제

### 이미지 처리 (선택)
- [ ] 리사이징 (썸네일 생성)
- [ ] WebP 변환

---

## Phase 9: 검색

### pg_bigm 설정
- [ ] 확장 설치 확인
- [ ] 인덱스 생성

### 검색 구현
- [ ] 제목 + 본문 검색
- [ ] 결과 하이라이팅 (선택)
- [ ] 페이지네이션

---

## Phase 10: 대시보드 & 기타

### 대시보드 API
- [ ] GET /api/admin/dashboard/stats
  - 전체 글 수
  - 발행/임시저장 글 수
  - 카테고리별 글 수
  - 최근 글 목록

### 기타
- [ ] GET /api/public/about - About 정보 (정적 또는 DB)
- [ ] RSS 피드 생성 엔드포인트 (선택)

---

## Phase 11: 문서화 & 테스트

### Swagger
- [ ] 핸들러 주석 작성
- [ ] swag init 실행
- [ ] Swagger UI 확인

### 테스트
- [ ] 핸들러 단위 테스트
- [ ] 서비스 단위 테스트
- [ ] 통합 테스트 (선택)

---

## Phase 12: 배포 준비

### Dockerfile
- [ ] 멀티 스테이지 빌드
- [ ] 경량 이미지 (scratch 또는 alpine)

### CI/CD
- [ ] GitHub Actions 워크플로우
- [ ] 빌드 및 푸시 설정

---

## 진행 상태

| Phase | 상태 | 예상 기간 |
|-------|------|-----------|
| Phase 1: 초기화 | ⬜ 대기 | 0.5일 |
| Phase 2: 기반 코드 | ⬜ 대기 | 1일 |
| Phase 3: 데이터베이스 | ⬜ 대기 | 1일 |
| Phase 4: 인증 | ⬜ 대기 | 0.5일 |
| Phase 5: 글 | ⬜ 대기 | 2일 |
| Phase 6: 카테고리 & 태그 | ⬜ 대기 | 1일 |
| Phase 7: 프로젝트 | ⬜ 대기 | 0.5일 |
| Phase 8: 미디어 | ⬜ 대기 | 1일 |
| Phase 9: 검색 | ⬜ 대기 | 0.5일 |
| Phase 10: 대시보드 | ⬜ 대기 | 0.5일 |
| Phase 11: 문서화 | ⬜ 대기 | 1일 |
| Phase 12: 배포 | ⬜ 대기 | 0.5일 |

**총 예상: 약 2주**

**상태**: ⬜ 대기 | 🔄 진행중 | ✅ 완료
