-- Add thumbnail columns to media table
ALTER TABLE media ADD COLUMN thumbnail_sm VARCHAR(500);
ALTER TABLE media ADD COLUMN thumbnail_md VARCHAR(500);
