-- Blog API SQL Queries
-- This file contains all SQL queries for sqlc code generation

-- ============================================================================
-- ADMINS
-- ============================================================================

-- name: GetAdminByUsername :one
SELECT * FROM admins WHERE username = $1;

-- name: GetAdminByID :one
SELECT * FROM admins WHERE id = $1;

-- name: CreateAdmin :one
INSERT INTO admins (username, password)
VALUES ($1, $2)
RETURNING *;

-- name: UpdateAdminPassword :exec
UPDATE admins SET password = $2, updated_at = NOW() WHERE id = $1;

-- ============================================================================
-- CATEGORIES
-- ============================================================================

-- name: ListCategories :many
SELECT * FROM categories ORDER BY sort_order ASC, id ASC;

-- name: GetCategoryByID :one
SELECT * FROM categories WHERE id = $1;

-- name: GetCategoryBySlug :one
SELECT * FROM categories WHERE slug = $1;

-- name: CreateCategory :one
INSERT INTO categories (name, slug, description, sort_order)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateCategory :one
UPDATE categories
SET name = $2, slug = $3, description = $4, sort_order = $5
WHERE id = $1
RETURNING *;

-- name: DeleteCategory :exec
DELETE FROM categories WHERE id = $1;

-- name: GetCategoryPostCount :one
SELECT COUNT(*) FROM posts WHERE category_id = $1 AND status = 'published';

-- ============================================================================
-- TAGS
-- ============================================================================

-- name: ListTags :many
SELECT * FROM tags ORDER BY name ASC;

-- name: GetTagByID :one
SELECT * FROM tags WHERE id = $1;

-- name: GetTagBySlug :one
SELECT * FROM tags WHERE slug = $1;

-- name: CreateTag :one
INSERT INTO tags (name, slug)
VALUES ($1, $2)
RETURNING *;

-- name: UpdateTag :one
UPDATE tags
SET name = $2, slug = $3
WHERE id = $1
RETURNING *;

-- name: DeleteTag :exec
DELETE FROM tags WHERE id = $1;

-- name: GetTagPostCount :one
SELECT COUNT(*) FROM post_tags pt
JOIN posts p ON pt.post_id = p.id
WHERE pt.tag_id = $1 AND p.status = 'published';

-- name: ListTagsWithPostCount :many
SELECT t.*, COUNT(pt.post_id) as post_count
FROM tags t
LEFT JOIN post_tags pt ON t.id = pt.tag_id
LEFT JOIN posts p ON pt.post_id = p.id AND p.status = 'published'
GROUP BY t.id
ORDER BY t.name ASC;

-- ============================================================================
-- POSTS
-- ============================================================================

-- name: ListPublishedPosts :many
SELECT p.*, c.name as category_name, c.slug as category_slug
FROM posts p
LEFT JOIN categories c ON p.category_id = c.id
WHERE p.status = 'published'
ORDER BY p.published_at DESC
LIMIT $1 OFFSET $2;

-- name: CountPublishedPosts :one
SELECT COUNT(*) FROM posts WHERE status = 'published';

-- name: ListPublishedPostsByCategory :many
SELECT p.*, c.name as category_name, c.slug as category_slug
FROM posts p
LEFT JOIN categories c ON p.category_id = c.id
WHERE p.status = 'published' AND p.category_id = $1
ORDER BY p.published_at DESC
LIMIT $2 OFFSET $3;

-- name: CountPublishedPostsByCategory :one
SELECT COUNT(*) FROM posts WHERE status = 'published' AND category_id = $1;

-- name: ListPublishedPostsByTag :many
SELECT p.*, c.name as category_name, c.slug as category_slug
FROM posts p
LEFT JOIN categories c ON p.category_id = c.id
JOIN post_tags pt ON p.id = pt.post_id
WHERE p.status = 'published' AND pt.tag_id = $1
ORDER BY p.published_at DESC
LIMIT $2 OFFSET $3;

-- name: CountPublishedPostsByTag :one
SELECT COUNT(*) FROM posts p
JOIN post_tags pt ON p.id = pt.post_id
WHERE p.status = 'published' AND pt.tag_id = $1;

-- name: GetPublishedPostBySlug :one
SELECT p.*, c.name as category_name, c.slug as category_slug
FROM posts p
LEFT JOIN categories c ON p.category_id = c.id
WHERE p.slug = $1 AND p.status = 'published';

-- name: SearchPublishedPosts :many
SELECT p.*, c.name as category_name, c.slug as category_slug
FROM posts p
LEFT JOIN categories c ON p.category_id = c.id
WHERE p.status = 'published'
  AND (p.title ILIKE '%' || $1 || '%' OR p.content ILIKE '%' || $1 || '%')
ORDER BY p.published_at DESC
LIMIT $2 OFFSET $3;

-- name: CountSearchPublishedPosts :one
SELECT COUNT(*) FROM posts
WHERE status = 'published'
  AND (title ILIKE '%' || $1 || '%' OR content ILIKE '%' || $1 || '%');

-- name: ListAllPosts :many
SELECT p.*, c.name as category_name, c.slug as category_slug
FROM posts p
LEFT JOIN categories c ON p.category_id = c.id
ORDER BY p.created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountAllPosts :one
SELECT COUNT(*) FROM posts;

-- name: ListPostsByStatus :many
SELECT p.*, c.name as category_name, c.slug as category_slug
FROM posts p
LEFT JOIN categories c ON p.category_id = c.id
WHERE p.status = $1
ORDER BY p.created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountPostsByStatus :one
SELECT COUNT(*) FROM posts WHERE status = $1;

-- name: GetPostByID :one
SELECT p.*, c.name as category_name, c.slug as category_slug
FROM posts p
LEFT JOIN categories c ON p.category_id = c.id
WHERE p.id = $1;

-- name: GetPostBySlug :one
SELECT p.*, c.name as category_name, c.slug as category_slug
FROM posts p
LEFT JOIN categories c ON p.category_id = c.id
WHERE p.slug = $1;

-- name: CreatePost :one
INSERT INTO posts (title, slug, content, excerpt, category_id, status, reading_time, thumbnail)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: UpdatePost :one
UPDATE posts
SET title = $2, slug = $3, content = $4, excerpt = $5, category_id = $6,
    reading_time = $7, thumbnail = $8, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: PublishPost :one
UPDATE posts
SET status = 'published', published_at = NOW(), updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UnpublishPost :one
UPDATE posts
SET status = 'draft', updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeletePost :exec
DELETE FROM posts WHERE id = $1;

-- name: IncrementViewCount :exec
UPDATE posts SET view_count = view_count + 1 WHERE id = $1;

-- name: CheckSlugExists :one
SELECT EXISTS(SELECT 1 FROM posts WHERE slug = $1);

-- name: CheckSlugExistsExcept :one
SELECT EXISTS(SELECT 1 FROM posts WHERE slug = $1 AND id != $2);

-- ============================================================================
-- POST_TAGS
-- ============================================================================

-- name: GetPostTags :many
SELECT t.* FROM tags t
JOIN post_tags pt ON t.id = pt.tag_id
WHERE pt.post_id = $1
ORDER BY t.name ASC;

-- name: AddPostTag :exec
INSERT INTO post_tags (post_id, tag_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- name: RemovePostTag :exec
DELETE FROM post_tags WHERE post_id = $1 AND tag_id = $2;

-- name: RemoveAllPostTags :exec
DELETE FROM post_tags WHERE post_id = $1;

-- name: SetPostTags :exec
DELETE FROM post_tags WHERE post_id = $1;

-- ============================================================================
-- PROJECTS
-- ============================================================================

-- name: ListProjects :many
SELECT * FROM projects ORDER BY sort_order ASC, id ASC;

-- name: ListFeaturedProjects :many
SELECT * FROM projects WHERE is_featured = true ORDER BY sort_order ASC, id ASC;

-- name: GetProjectByID :one
SELECT * FROM projects WHERE id = $1;

-- name: GetProjectBySlug :one
SELECT * FROM projects WHERE slug = $1;

-- name: CreateProject :one
INSERT INTO projects (title, slug, description, content, tech_stack, demo_url, github_url, thumbnail, images, is_featured, sort_order)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING *;

-- name: UpdateProject :one
UPDATE projects
SET title = $2, slug = $3, description = $4, content = $5, tech_stack = $6,
    demo_url = $7, github_url = $8, thumbnail = $9, images = $10,
    is_featured = $11, sort_order = $12, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteProject :exec
DELETE FROM projects WHERE id = $1;

-- name: UpdateProjectOrder :exec
UPDATE projects SET sort_order = $2 WHERE id = $1;

-- ============================================================================
-- MEDIA
-- ============================================================================

-- name: ListMedia :many
SELECT * FROM media ORDER BY created_at DESC LIMIT $1 OFFSET $2;

-- name: CountMedia :one
SELECT COUNT(*) FROM media;

-- name: GetMediaByID :one
SELECT * FROM media WHERE id = $1;

-- name: CreateMedia :one
INSERT INTO media (filename, original_name, path, url, mime_type, size, width, height)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: DeleteMedia :exec
DELETE FROM media WHERE id = $1;

-- ============================================================================
-- DASHBOARD STATS
-- ============================================================================

-- name: GetPostStats :one
SELECT
    COUNT(*) as total,
    COUNT(*) FILTER (WHERE status = 'published') as published,
    COUNT(*) FILTER (WHERE status = 'draft') as draft
FROM posts;

-- name: GetCategoryStats :many
SELECT c.id, c.name, c.slug, COUNT(p.id) as post_count
FROM categories c
LEFT JOIN posts p ON c.id = p.category_id AND p.status = 'published'
GROUP BY c.id
ORDER BY post_count DESC;

-- name: GetRecentPosts :many
SELECT p.*, c.name as category_name, c.slug as category_slug
FROM posts p
LEFT JOIN categories c ON p.category_id = c.id
ORDER BY p.created_at DESC
LIMIT $1;

-- name: GetTotalViews :one
SELECT COALESCE(SUM(view_count), 0) as total_views FROM posts;
