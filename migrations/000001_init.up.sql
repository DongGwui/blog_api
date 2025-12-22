-- Blog API Initial Schema

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
