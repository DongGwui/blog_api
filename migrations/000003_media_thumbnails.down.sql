-- Remove thumbnail columns from media table
ALTER TABLE media DROP COLUMN thumbnail_sm;
ALTER TABLE media DROP COLUMN thumbnail_md;
