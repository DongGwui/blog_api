-- Search Index Migration
-- pg_bigm: 2-gram 기반 전문 검색 (한국어/일본어에 최적화)

-- pg_bigm 확장 설치
CREATE EXTENSION IF NOT EXISTS pg_bigm;

-- posts 테이블에 검색 인덱스 생성
-- GIN 인덱스: 전문 검색에 최적화된 인덱스 타입

-- 제목 검색 인덱스
CREATE INDEX IF NOT EXISTS idx_posts_title_bigm
ON posts USING gin (title gin_bigm_ops);

-- 본문 검색 인덱스
CREATE INDEX IF NOT EXISTS idx_posts_content_bigm
ON posts USING gin (content gin_bigm_ops);
