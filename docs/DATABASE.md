# 데이터베이스 스키마

## 개요

PostgreSQL 16 사용. 한국어 검색을 위해 pg_bigm 확장 활용.

---

## ERD

```
┌─────────────┐       ┌─────────────┐       ┌─────────────┐
│   admins    │       │   posts     │──────<│ post_tags   │
└─────────────┘       └──────┬──────┘       └──────┬──────┘
                             │                     │
                             │ FK                  │ FK
                             ▼                     ▼
                      ┌─────────────┐       ┌─────────────┐
                      │ categories  │       │    tags     │
                      └─────────────┘       └─────────────┘

┌─────────────┐       ┌─────────────┐
│  projects   │       │   media     │
└─────────────┘       └─────────────┘
```

---

## 테이블 정의

### admins (관리자)

단일 사용자 블로그이므로 1개 레코드만 존재.

| 컬럼 | 타입 | 제약조건 | 설명 |
|------|------|----------|------|
| id | SERIAL | PK | |
| username | VARCHAR(50) | UNIQUE, NOT NULL | 로그인 ID |
| password | VARCHAR(255) | NOT NULL | bcrypt 해시 |
| created_at | TIMESTAMPTZ | DEFAULT NOW() | |
| updated_at | TIMESTAMPTZ | DEFAULT NOW() | |

---

### categories (카테고리)

| 컬럼 | 타입 | 제약조건 | 설명 |
|------|------|----------|------|
| id | SERIAL | PK | |
| name | VARCHAR(50) | NOT NULL | 표시 이름 |
| slug | VARCHAR(50) | UNIQUE, NOT NULL | URL용 |
| description | TEXT | | 설명 |
| sort_order | INT | DEFAULT 0 | 정렬 순서 |
| created_at | TIMESTAMPTZ | DEFAULT NOW() | |

---

### tags (태그)

| 컬럼 | 타입 | 제약조건 | 설명 |
|------|------|----------|------|
| id | SERIAL | PK | |
| name | VARCHAR(50) | NOT NULL | 표시 이름 |
| slug | VARCHAR(50) | UNIQUE, NOT NULL | URL용 |
| created_at | TIMESTAMPTZ | DEFAULT NOW() | |

---

### posts (글)

| 컬럼 | 타입 | 제약조건 | 설명 |
|------|------|----------|------|
| id | SERIAL | PK | |
| title | VARCHAR(200) | NOT NULL | 제목 |
| slug | VARCHAR(200) | UNIQUE, NOT NULL | URL용 |
| content | TEXT | NOT NULL | 마크다운 본문 |
| excerpt | VARCHAR(500) | | 요약 (목록용) |
| category_id | INT | FK → categories | 카테고리 |
| status | VARCHAR(20) | DEFAULT 'draft' | draft / published |
| view_count | INT | DEFAULT 0 | 조회수 |
| reading_time | INT | | 읽기 시간 (분) |
| thumbnail | VARCHAR(500) | | 썸네일 URL |
| created_at | TIMESTAMPTZ | DEFAULT NOW() | 작성일 |
| updated_at | TIMESTAMPTZ | DEFAULT NOW() | 수정일 |
| published_at | TIMESTAMPTZ | | 발행일 |

**status 값**
- `draft`: 임시저장
- `published`: 발행됨

---

### post_tags (글-태그 연결)

다대다 관계 테이블.

| 컬럼 | 타입 | 제약조건 | 설명 |
|------|------|----------|------|
| post_id | INT | FK → posts, ON DELETE CASCADE | |
| tag_id | INT | FK → tags, ON DELETE CASCADE | |

**PK**: (post_id, tag_id)

---

### projects (프로젝트)

| 컬럼 | 타입 | 제약조건 | 설명 |
|------|------|----------|------|
| id | SERIAL | PK | |
| title | VARCHAR(200) | NOT NULL | 프로젝트명 |
| slug | VARCHAR(200) | UNIQUE, NOT NULL | URL용 |
| description | TEXT | | 간단 설명 |
| content | TEXT | | 상세 설명 (마크다운) |
| tech_stack | JSONB | | ["Go", "Next.js"] |
| demo_url | VARCHAR(500) | | 데모 링크 |
| github_url | VARCHAR(500) | | GitHub 링크 |
| thumbnail | VARCHAR(500) | | 썸네일 URL |
| images | JSONB | | ["url1", "url2"] |
| is_featured | BOOLEAN | DEFAULT FALSE | 대표 프로젝트 |
| sort_order | INT | DEFAULT 0 | 정렬 순서 |
| created_at | TIMESTAMPTZ | DEFAULT NOW() | |
| updated_at | TIMESTAMPTZ | DEFAULT NOW() | |

---

### media (미디어)

| 컬럼 | 타입 | 제약조건 | 설명 |
|------|------|----------|------|
| id | SERIAL | PK | |
| filename | VARCHAR(255) | NOT NULL | 저장된 파일명 (UUID) |
| original_name | VARCHAR(255) | NOT NULL | 원본 파일명 |
| path | VARCHAR(500) | NOT NULL | MinIO 전체 경로 |
| url | VARCHAR(500) | NOT NULL | 공개 URL |
| mime_type | VARCHAR(100) | | image/jpeg 등 |
| size | BIGINT | | 바이트 크기 |
| width | INT | | 이미지 너비 |
| height | INT | | 이미지 높이 |
| created_at | TIMESTAMPTZ | DEFAULT NOW() | |

---

## 인덱스

```sql
-- 글 조회 최적화
CREATE INDEX idx_posts_status ON posts(status);
CREATE INDEX idx_posts_published_at ON posts(published_at DESC);
CREATE INDEX idx_posts_category ON posts(category_id);
CREATE INDEX idx_posts_slug ON posts(slug);

-- 카테고리/태그 조회
CREATE INDEX idx_categories_slug ON categories(slug);
CREATE INDEX idx_categories_sort ON categories(sort_order);
CREATE INDEX idx_tags_slug ON tags(slug);

-- 프로젝트 정렬
CREATE INDEX idx_projects_sort ON projects(sort_order);
CREATE INDEX idx_projects_featured ON projects(is_featured);

-- 한국어 전문 검색 (pg_bigm)
CREATE INDEX idx_posts_title_bigm ON posts USING gin (title gin_bigm_ops);
CREATE INDEX idx_posts_content_bigm ON posts USING gin (content gin_bigm_ops);
```

---

## 마이그레이션 파일

### 000001_init.up.sql

```sql
-- 확장 설치 (수동으로 먼저 필요할 수 있음)
-- CREATE EXTENSION IF NOT EXISTS pg_bigm;

-- 관리자
CREATE TABLE admins (
    id          SERIAL PRIMARY KEY,
    username    VARCHAR(50) UNIQUE NOT NULL,
    password    VARCHAR(255) NOT NULL,
    created_at  TIMESTAMPTZ DEFAULT NOW(),
    updated_at  TIMESTAMPTZ DEFAULT NOW()
);

-- 카테고리
CREATE TABLE categories (
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(50) NOT NULL,
    slug        VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    sort_order  INT DEFAULT 0,
    created_at  TIMESTAMPTZ DEFAULT NOW()
);

-- 태그
CREATE TABLE tags (
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(50) NOT NULL,
    slug        VARCHAR(50) UNIQUE NOT NULL,
    created_at  TIMESTAMPTZ DEFAULT NOW()
);

-- 글
CREATE TABLE posts (
    id           SERIAL PRIMARY KEY,
    title        VARCHAR(200) NOT NULL,
    slug         VARCHAR(200) UNIQUE NOT NULL,
    content      TEXT NOT NULL,
    excerpt      VARCHAR(500),
    category_id  INT REFERENCES categories(id) ON DELETE SET NULL,
    status       VARCHAR(20) DEFAULT 'draft',
    view_count   INT DEFAULT 0,
    reading_time INT,
    thumbnail    VARCHAR(500),
    created_at   TIMESTAMPTZ DEFAULT NOW(),
    updated_at   TIMESTAMPTZ DEFAULT NOW(),
    published_at TIMESTAMPTZ
);

-- 글-태그 연결
CREATE TABLE post_tags (
    post_id INT REFERENCES posts(id) ON DELETE CASCADE,
    tag_id  INT REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (post_id, tag_id)
);

-- 프로젝트
CREATE TABLE projects (
    id          SERIAL PRIMARY KEY,
    title       VARCHAR(200) NOT NULL,
    slug        VARCHAR(200) UNIQUE NOT NULL,
    description TEXT,
    content     TEXT,
    tech_stack  JSONB,
    demo_url    VARCHAR(500),
    github_url  VARCHAR(500),
    thumbnail   VARCHAR(500),
    images      JSONB,
    is_featured BOOLEAN DEFAULT FALSE,
    sort_order  INT DEFAULT 0,
    created_at  TIMESTAMPTZ DEFAULT NOW(),
    updated_at  TIMESTAMPTZ DEFAULT NOW()
);

-- 미디어
CREATE TABLE media (
    id            SERIAL PRIMARY KEY,
    filename      VARCHAR(255) NOT NULL,
    original_name VARCHAR(255) NOT NULL,
    path          VARCHAR(500) NOT NULL,
    url           VARCHAR(500) NOT NULL,
    mime_type     VARCHAR(100),
    size          BIGINT,
    width         INT,
    height        INT,
    created_at    TIMESTAMPTZ DEFAULT NOW()
);

-- 인덱스
CREATE INDEX idx_posts_status ON posts(status);
CREATE INDEX idx_posts_published_at ON posts(published_at DESC);
CREATE INDEX idx_posts_category ON posts(category_id);
CREATE INDEX idx_posts_slug ON posts(slug);
CREATE INDEX idx_categories_slug ON categories(slug);
CREATE INDEX idx_categories_sort ON categories(sort_order);
CREATE INDEX idx_tags_slug ON tags(slug);
CREATE INDEX idx_projects_sort ON projects(sort_order);
CREATE INDEX idx_projects_featured ON projects(is_featured);
```

### 000001_init.down.sql

```sql
DROP TABLE IF EXISTS post_tags;
DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS tags;
DROP TABLE IF EXISTS projects;
DROP TABLE IF EXISTS media;
DROP TABLE IF EXISTS admins;
```

---

## pg_bigm 설치 (한국어 검색)

PostgreSQL 기본 이미지에는 pg_bigm이 없음. 커스텀 이미지 또는 수동 설치 필요.

**옵션 1: 커스텀 Dockerfile**
```dockerfile
FROM postgres:16-alpine
RUN apk add --no-cache postgresql16-contrib
# pg_bigm 빌드 필요...
```

**옵션 2: LIKE 검색으로 대체 (소규모라면 충분)**
```sql
SELECT * FROM posts 
WHERE title LIKE '%검색어%' OR content LIKE '%검색어%';
```

**옵션 3: 나중에 pg_bigm 추가**
- 초기에는 LIKE 검색으로 시작
- 글이 많아지면 pg_bigm 또는 다른 검색 솔루션 도입
