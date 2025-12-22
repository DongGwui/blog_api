-- Rollback Search Index Migration

-- 인덱스 삭제
DROP INDEX IF EXISTS idx_posts_title_bigm;
DROP INDEX IF EXISTS idx_posts_content_bigm;

-- 확장 삭제 (다른 곳에서 사용 중일 수 있으므로 주의)
-- DROP EXTENSION IF EXISTS pg_bigm;
