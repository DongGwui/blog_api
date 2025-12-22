# API 엔드포인트 명세

## 공통 사항

### Base URL
- 개발: `http://localhost:8080`
- 프로덕션 (Public): `https://api.dltmxm.link`
- 프로덕션 (Admin): Tailscale IP 사용

### 응답 형식

**성공 (단일)**
```json
{
  "data": { ... }
}
```

**성공 (목록)**
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

**에러**
```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "Post not found"
  }
}
```

---

## Health Check

### GET /api/health

서비스 상태 확인.

**Response 200**
```json
{
  "status": "ok",
  "services": {
    "database": "ok",
    "redis": "ok",
    "minio": "ok"
  }
}
```

---

## Public API

### Posts

#### GET /api/public/posts

글 목록 조회.

**Query Parameters**
| 파라미터 | 타입 | 기본값 | 설명 |
|----------|------|--------|------|
| page | int | 1 | 페이지 번호 |
| per_page | int | 10 | 페이지당 개수 (max: 50) |
| category | string | | 카테고리 slug |
| tag | string | | 태그 slug |

**Response 200**
```json
{
  "data": [
    {
      "id": 1,
      "title": "Go로 REST API 만들기",
      "slug": "go-rest-api",
      "excerpt": "Go와 Gin을 사용한 API 개발...",
      "category": {
        "id": 1,
        "name": "Backend",
        "slug": "backend"
      },
      "tags": [
        { "id": 1, "name": "Go", "slug": "go" }
      ],
      "thumbnail": "https://...",
      "view_count": 123,
      "reading_time": 5,
      "published_at": "2025-01-15T09:00:00Z"
    }
  ],
  "meta": {
    "page": 1,
    "per_page": 10,
    "total": 25,
    "total_pages": 3
  }
}
```

---

#### GET /api/public/posts/:slug

글 상세 조회.

**Response 200**
```json
{
  "data": {
    "id": 1,
    "title": "Go로 REST API 만들기",
    "slug": "go-rest-api",
    "content": "# 소개\n\nGo는 ...",
    "excerpt": "Go와 Gin을 사용한 API 개발...",
    "category": {
      "id": 1,
      "name": "Backend",
      "slug": "backend"
    },
    "tags": [
      { "id": 1, "name": "Go", "slug": "go" }
    ],
    "thumbnail": "https://...",
    "view_count": 123,
    "reading_time": 5,
    "created_at": "2025-01-15T08:00:00Z",
    "updated_at": "2025-01-15T09:00:00Z",
    "published_at": "2025-01-15T09:00:00Z",
    "prev_post": { "slug": "prev-post", "title": "이전 글" },
    "next_post": { "slug": "next-post", "title": "다음 글" }
  }
}
```

**Response 404**
```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "Post not found"
  }
}
```

---

#### GET /api/public/posts/search

글 검색.

**Query Parameters**
| 파라미터 | 타입 | 필수 | 설명 |
|----------|------|------|------|
| q | string | ✓ | 검색어 |
| page | int | | 페이지 번호 |
| per_page | int | | 페이지당 개수 |

**Response 200**
```json
{
  "data": [ ... ],
  "meta": { ... }
}
```

---

#### POST /api/public/posts/:slug/view

조회수 증가 (IP 기반 중복 방지).

**Response 200**
```json
{
  "data": {
    "view_count": 124
  }
}
```

---

### Categories

#### GET /api/public/categories

카테고리 목록.

**Response 200**
```json
{
  "data": [
    {
      "id": 1,
      "name": "Backend",
      "slug": "backend",
      "description": "백엔드 개발 관련 글",
      "post_count": 15
    }
  ]
}
```

---

### Tags

#### GET /api/public/tags

태그 목록.

**Response 200**
```json
{
  "data": [
    {
      "id": 1,
      "name": "Go",
      "slug": "go",
      "post_count": 8
    }
  ]
}
```

---

### Projects

#### GET /api/public/projects

프로젝트 목록.

**Response 200**
```json
{
  "data": [
    {
      "id": 1,
      "title": "개인 블로그",
      "slug": "personal-blog",
      "description": "Next.js와 Go로 만든 블로그",
      "tech_stack": ["Next.js", "Go", "PostgreSQL"],
      "thumbnail": "https://...",
      "is_featured": true
    }
  ]
}
```

#### GET /api/public/projects/:slug

프로젝트 상세.

**Response 200**
```json
{
  "data": {
    "id": 1,
    "title": "개인 블로그",
    "slug": "personal-blog",
    "description": "Next.js와 Go로 만든 블로그",
    "content": "## 프로젝트 소개\n\n...",
    "tech_stack": ["Next.js", "Go", "PostgreSQL"],
    "demo_url": "https://blog.dltmxm.link",
    "github_url": "https://github.com/...",
    "thumbnail": "https://...",
    "images": ["https://...", "https://..."]
  }
}
```

---

## Admin API

**인증**: 모든 Admin API는 JWT 토큰 필요.

```
Authorization: Bearer <token>
```

---

### Auth

#### POST /api/admin/auth/login

로그인.

**Request Body**
```json
{
  "username": "admin",
  "password": "password123"
}
```

**Response 200**
```json
{
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_at": "2025-01-16T09:00:00Z"
  }
}
```

**Response 401**
```json
{
  "error": {
    "code": "INVALID_CREDENTIALS",
    "message": "Invalid username or password"
  }
}
```

---

#### GET /api/admin/auth/me

현재 로그인 정보.

**Response 200**
```json
{
  "data": {
    "id": 1,
    "username": "admin"
  }
}
```

---

### Admin Posts

#### GET /api/admin/posts

글 목록 (전체, 상태 필터 가능).

**Query Parameters**
| 파라미터 | 타입 | 설명 |
|----------|------|------|
| status | string | draft / published / all |
| page | int | |
| per_page | int | |

---

#### POST /api/admin/posts

글 생성.

**Request Body**
```json
{
  "title": "새 글 제목",
  "content": "# 소개\n\n본문 내용...",
  "excerpt": "요약",
  "category_id": 1,
  "tag_ids": [1, 2, 3],
  "status": "draft",
  "thumbnail": "https://..."
}
```

**Response 201**
```json
{
  "data": {
    "id": 10,
    "slug": "새-글-제목",
    ...
  }
}
```

---

#### PUT /api/admin/posts/:id

글 수정.

**Request Body** (부분 업데이트 가능)
```json
{
  "title": "수정된 제목",
  "content": "수정된 내용"
}
```

---

#### DELETE /api/admin/posts/:id

글 삭제.

**Response 200**
```json
{
  "data": {
    "message": "Post deleted"
  }
}
```

---

#### PATCH /api/admin/posts/:id/publish

발행 상태 변경.

**Request Body**
```json
{
  "status": "published"
}
```

---

### Admin Categories

#### POST /api/admin/categories

**Request Body**
```json
{
  "name": "Frontend",
  "slug": "frontend",
  "description": "프론트엔드 개발"
}
```

#### PUT /api/admin/categories/:id

#### DELETE /api/admin/categories/:id

---

### Admin Tags

#### POST /api/admin/tags

**Request Body**
```json
{
  "name": "TypeScript",
  "slug": "typescript"
}
```

#### PUT /api/admin/tags/:id

#### DELETE /api/admin/tags/:id

---

### Admin Projects

#### POST /api/admin/projects

**Request Body**
```json
{
  "title": "새 프로젝트",
  "slug": "new-project",
  "description": "프로젝트 설명",
  "content": "상세 내용...",
  "tech_stack": ["React", "Node.js"],
  "demo_url": "https://...",
  "github_url": "https://...",
  "is_featured": false
}
```

#### PATCH /api/admin/projects/reorder

순서 변경.

**Request Body**
```json
{
  "orders": [
    { "id": 3, "sort_order": 0 },
    { "id": 1, "sort_order": 1 },
    { "id": 2, "sort_order": 2 }
  ]
}
```

---

### Admin Media

#### GET /api/admin/media

미디어 목록.

**Query Parameters**
| 파라미터 | 타입 | 설명 |
|----------|------|------|
| page | int | |
| per_page | int | |

**Response 200**
```json
{
  "data": [
    {
      "id": 1,
      "filename": "abc123.jpg",
      "original_name": "screenshot.jpg",
      "url": "https://minio.../blog-images/2025/01/abc123.jpg",
      "mime_type": "image/jpeg",
      "size": 102400,
      "width": 1920,
      "height": 1080,
      "created_at": "2025-01-15T09:00:00Z"
    }
  ]
}
```

---

#### POST /api/admin/media/upload

이미지 업로드.

**Request**: `multipart/form-data`

| 필드 | 타입 | 설명 |
|------|------|------|
| file | file | 이미지 파일 |

**Response 201**
```json
{
  "data": {
    "id": 5,
    "url": "https://minio.../blog-images/2025/01/xyz789.jpg",
    "filename": "xyz789.jpg"
  }
}
```

---

#### DELETE /api/admin/media/:id

미디어 삭제.

---

### Admin Dashboard

#### GET /api/admin/dashboard/stats

대시보드 통계.

**Response 200**
```json
{
  "data": {
    "posts": {
      "total": 50,
      "published": 45,
      "draft": 5
    },
    "categories": 8,
    "tags": 25,
    "projects": 6,
    "recent_posts": [
      { "id": 50, "title": "최근 글", "status": "published", "created_at": "..." }
    ]
  }
}
```
