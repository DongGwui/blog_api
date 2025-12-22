-- Search Index Migration
-- pg_bigm: 2-gram 기반 전문 검색 (한국어/일본어에 최적화)
-- pg_trgm: 3-gram 기반 전문 검색 (대안)

-- pg_bigm 확장 설치 (한국어 검색에 최적화)
-- 주의: pg_bigm이 설치되어 있지 않은 경우 pg_trgm을 사용하세요
-- CREATE EXTENSION IF NOT EXISTS pg_bigm;

-- pg_trgm 확장 설치 (PostgreSQL 기본 제공, 더 범용적)
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- posts 테이블에 검색 인덱스 생성
-- GIN 인덱스: 전문 검색에 최적화된 인덱스 타입

-- 제목 검색 인덱스
CREATE INDEX IF NOT EXISTS idx_posts_title_trgm
ON posts USING gin (title gin_trgm_ops);

-- 본문 검색 인덱스
CREATE INDEX IF NOT EXISTS idx_posts_content_trgm
ON posts USING gin (content gin_trgm_ops);

-- 복합 검색을 위한 tsvector 컬럼 추가 (선택적)
-- ALTER TABLE posts ADD COLUMN IF NOT EXISTS search_vector tsvector;

-- 검색 벡터 업데이트 트리거 (선택적)
-- CREATE OR REPLACE FUNCTION posts_search_vector_update() RETURNS trigger AS $$
-- BEGIN
--     NEW.search_vector := to_tsvector('simple', coalesce(NEW.title, '') || ' ' || coalesce(NEW.content, ''));
--     RETURN NEW;
-- END;
-- $$ LANGUAGE plpgsql;

-- CREATE TRIGGER posts_search_vector_trigger
-- BEFORE INSERT OR UPDATE ON posts
-- FOR EACH ROW EXECUTE FUNCTION posts_search_vector_update();
