package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ydonggwui/blog-api/internal/database/sqlc"
	"github.com/ydonggwui/blog-api/internal/model"
	"github.com/ydonggwui/blog-api/internal/util"
)

var (
	ErrPostNotFound   = errors.New("post not found")
	ErrSlugExists     = errors.New("slug already exists")
)

type PostService struct {
	queries *sqlc.Queries
	db      *sql.DB
}

func NewPostService(queries *sqlc.Queries, db *sql.DB) *PostService {
	return &PostService{
		queries: queries,
		db:      db,
	}
}

// ListPublishedPosts returns paginated published posts
func (s *PostService) ListPublishedPosts(ctx context.Context, limit, offset int32) ([]model.PostListResponse, int64, error) {
	posts, err := s.queries.ListPublishedPosts(ctx, sqlc.ListPublishedPostsParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, 0, err
	}

	count, err := s.queries.CountPublishedPosts(ctx)
	if err != nil {
		return nil, 0, err
	}

	result := make([]model.PostListResponse, len(posts))
	for i, p := range posts {
		result[i] = s.toPostListResponse(ctx, p)
	}

	return result, count, nil
}

// ListPublishedPostsByCategory returns paginated published posts by category
func (s *PostService) ListPublishedPostsByCategory(ctx context.Context, categoryID int32, limit, offset int32) ([]model.PostListResponse, int64, error) {
	posts, err := s.queries.ListPublishedPostsByCategory(ctx, sqlc.ListPublishedPostsByCategoryParams{
		CategoryID: sql.NullInt32{Int32: categoryID, Valid: true},
		Limit:      limit,
		Offset:     offset,
	})
	if err != nil {
		return nil, 0, err
	}

	count, err := s.queries.CountPublishedPostsByCategory(ctx, sql.NullInt32{Int32: categoryID, Valid: true})
	if err != nil {
		return nil, 0, err
	}

	result := make([]model.PostListResponse, len(posts))
	for i, p := range posts {
		result[i] = s.toPostListResponseFromCategory(ctx, p)
	}

	return result, count, nil
}

// ListPublishedPostsByTag returns paginated published posts by tag
func (s *PostService) ListPublishedPostsByTag(ctx context.Context, tagID int32, limit, offset int32) ([]model.PostListResponse, int64, error) {
	posts, err := s.queries.ListPublishedPostsByTag(ctx, sqlc.ListPublishedPostsByTagParams{
		TagID:  tagID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, 0, err
	}

	count, err := s.queries.CountPublishedPostsByTag(ctx, tagID)
	if err != nil {
		return nil, 0, err
	}

	result := make([]model.PostListResponse, len(posts))
	for i, p := range posts {
		result[i] = s.toPostListResponseFromTag(ctx, p)
	}

	return result, count, nil
}

// GetPublishedPostBySlug returns a published post by slug
func (s *PostService) GetPublishedPostBySlug(ctx context.Context, slug string) (*model.PostResponse, error) {
	post, err := s.queries.GetPublishedPostBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrPostNotFound
		}
		return nil, err
	}

	return s.toPostResponse(ctx, post), nil
}

// SearchPublishedPosts searches published posts
func (s *PostService) SearchPublishedPosts(ctx context.Context, query string, limit, offset int32) ([]model.PostListResponse, int64, error) {
	queryParam := sql.NullString{String: query, Valid: query != ""}

	posts, err := s.queries.SearchPublishedPosts(ctx, sqlc.SearchPublishedPostsParams{
		Column1: queryParam,
		Limit:   limit,
		Offset:  offset,
	})
	if err != nil {
		return nil, 0, err
	}

	count, err := s.queries.CountSearchPublishedPosts(ctx, queryParam)
	if err != nil {
		return nil, 0, err
	}

	result := make([]model.PostListResponse, len(posts))
	for i, p := range posts {
		result[i] = s.toPostListResponseFromSearch(ctx, p)
	}

	return result, count, nil
}

// ListAllPosts returns all posts for admin (paginated)
func (s *PostService) ListAllPosts(ctx context.Context, limit, offset int32) ([]model.PostListResponse, int64, error) {
	posts, err := s.queries.ListAllPosts(ctx, sqlc.ListAllPostsParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, 0, err
	}

	count, err := s.queries.CountAllPosts(ctx)
	if err != nil {
		return nil, 0, err
	}

	result := make([]model.PostListResponse, len(posts))
	for i, p := range posts {
		result[i] = s.toPostListResponseFromAll(ctx, p)
	}

	return result, count, nil
}

// ListPostsByStatus returns posts by status for admin (paginated)
func (s *PostService) ListPostsByStatus(ctx context.Context, status string, limit, offset int32) ([]model.PostListResponse, int64, error) {
	posts, err := s.queries.ListPostsByStatus(ctx, sqlc.ListPostsByStatusParams{
		Status: sql.NullString{String: status, Valid: true},
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, 0, err
	}

	count, err := s.queries.CountPostsByStatus(ctx, sql.NullString{String: status, Valid: true})
	if err != nil {
		return nil, 0, err
	}

	result := make([]model.PostListResponse, len(posts))
	for i, p := range posts {
		result[i] = s.toPostListResponseFromStatus(ctx, p)
	}

	return result, count, nil
}

// GetPostByID returns a post by ID for admin
func (s *PostService) GetPostByID(ctx context.Context, id int32) (*model.PostResponse, error) {
	post, err := s.queries.GetPostByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrPostNotFound
		}
		return nil, err
	}

	return s.toPostResponseFromID(ctx, post), nil
}

// CreatePost creates a new post
func (s *PostService) CreatePost(ctx context.Context, req *model.CreatePostRequest) (*model.PostResponse, error) {
	// Generate slug if not provided
	slug := req.Slug
	if slug == "" {
		slug = util.GenerateSlug(req.Title)
	}

	// Check if slug exists
	exists, err := s.queries.CheckSlugExists(ctx, slug)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrSlugExists
	}

	// Calculate reading time
	readingTime := util.CalculateReadingTime(req.Content)

	// Determine status
	status := req.Status
	if status == "" {
		status = string(model.PostStatusDraft)
	}

	// Create post
	post, err := s.queries.CreatePost(ctx, sqlc.CreatePostParams{
		Title:       req.Title,
		Slug:        slug,
		Content:     req.Content,
		Excerpt:     sql.NullString{String: req.Excerpt, Valid: req.Excerpt != ""},
		CategoryID:  sql.NullInt32{Int32: ptrToInt32(req.CategoryID), Valid: req.CategoryID != nil},
		Status:      sql.NullString{String: status, Valid: true},
		ReadingTime: sql.NullInt32{Int32: int32(readingTime), Valid: true},
		Thumbnail:   sql.NullString{String: req.Thumbnail, Valid: req.Thumbnail != ""},
	})
	if err != nil {
		return nil, err
	}

	// Add tags
	if len(req.TagIDs) > 0 {
		for _, tagID := range req.TagIDs {
			_ = s.queries.AddPostTag(ctx, sqlc.AddPostTagParams{
				PostID: post.ID,
				TagID:  tagID,
			})
		}
	}

	return s.GetPostByID(ctx, post.ID)
}

// UpdatePost updates a post
func (s *PostService) UpdatePost(ctx context.Context, id int32, req *model.UpdatePostRequest) (*model.PostResponse, error) {
	// Check if post exists
	_, err := s.queries.GetPostByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrPostNotFound
		}
		return nil, err
	}

	// Generate slug if not provided
	slug := req.Slug
	if slug == "" {
		slug = util.GenerateSlug(req.Title)
	}

	// Check if slug exists (excluding current post)
	exists, err := s.queries.CheckSlugExistsExcept(ctx, sqlc.CheckSlugExistsExceptParams{
		Slug: slug,
		ID:   id,
	})
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrSlugExists
	}

	// Calculate reading time
	readingTime := util.CalculateReadingTime(req.Content)

	// Update post
	_, err = s.queries.UpdatePost(ctx, sqlc.UpdatePostParams{
		ID:          id,
		Title:       req.Title,
		Slug:        slug,
		Content:     req.Content,
		Excerpt:     sql.NullString{String: req.Excerpt, Valid: req.Excerpt != ""},
		CategoryID:  sql.NullInt32{Int32: ptrToInt32(req.CategoryID), Valid: req.CategoryID != nil},
		ReadingTime: sql.NullInt32{Int32: int32(readingTime), Valid: true},
		Thumbnail:   sql.NullString{String: req.Thumbnail, Valid: req.Thumbnail != ""},
	})
	if err != nil {
		return nil, err
	}

	// Update tags
	_ = s.queries.RemoveAllPostTags(ctx, id)
	if len(req.TagIDs) > 0 {
		for _, tagID := range req.TagIDs {
			_ = s.queries.AddPostTag(ctx, sqlc.AddPostTagParams{
				PostID: id,
				TagID:  tagID,
			})
		}
	}

	return s.GetPostByID(ctx, id)
}

// PublishPost publishes or unpublishes a post
func (s *PostService) PublishPost(ctx context.Context, id int32, publish bool) (*model.PostResponse, error) {
	var err error
	if publish {
		_, err = s.queries.PublishPost(ctx, id)
	} else {
		_, err = s.queries.UnpublishPost(ctx, id)
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrPostNotFound
		}
		return nil, err
	}

	return s.GetPostByID(ctx, id)
}

// DeletePost deletes a post
func (s *PostService) DeletePost(ctx context.Context, id int32) error {
	// Check if post exists
	_, err := s.queries.GetPostByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrPostNotFound
		}
		return err
	}

	return s.queries.DeletePost(ctx, id)
}

// IncrementViewCount increments the view count
func (s *PostService) IncrementViewCount(ctx context.Context, id int32) error {
	return s.queries.IncrementViewCount(ctx, id)
}

// GetPostIDBySlug returns post ID by slug
func (s *PostService) GetPostIDBySlug(ctx context.Context, slug string) (int32, error) {
	post, err := s.queries.GetPostBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrPostNotFound
		}
		return 0, err
	}
	return post.ID, nil
}

// Helper functions

func (s *PostService) getPostTags(ctx context.Context, postID int32) []model.TagBrief {
	tags, err := s.queries.GetPostTags(ctx, postID)
	if err != nil {
		return []model.TagBrief{}
	}

	result := make([]model.TagBrief, len(tags))
	for i, t := range tags {
		result[i] = model.TagBrief{
			ID:   t.ID,
			Name: t.Name,
			Slug: t.Slug,
		}
	}
	return result
}

func (s *PostService) toPostListResponse(ctx context.Context, p sqlc.ListPublishedPostsRow) model.PostListResponse {
	resp := model.PostListResponse{
		ID:          p.ID,
		Title:       p.Title,
		Slug:        p.Slug,
		Status:      p.Status.String,
		ViewCount:   p.ViewCount.Int32,
		CreatedAt:   p.CreatedAt.Time,
		Tags:        s.getPostTags(ctx, p.ID),
	}
	if p.Excerpt.Valid {
		resp.Excerpt = p.Excerpt.String
	}
	if p.CategoryID.Valid {
		resp.CategoryID = &p.CategoryID.Int32
	}
	if p.CategoryName.Valid {
		resp.CategoryName = &p.CategoryName.String
	}
	if p.CategorySlug.Valid {
		resp.CategorySlug = &p.CategorySlug.String
	}
	if p.ReadingTime.Valid {
		resp.ReadingTime = &p.ReadingTime.Int32
	}
	if p.Thumbnail.Valid {
		resp.Thumbnail = &p.Thumbnail.String
	}
	if p.PublishedAt.Valid {
		resp.PublishedAt = &p.PublishedAt.Time
	}
	return resp
}

func (s *PostService) toPostListResponseFromCategory(ctx context.Context, p sqlc.ListPublishedPostsByCategoryRow) model.PostListResponse {
	resp := model.PostListResponse{
		ID:          p.ID,
		Title:       p.Title,
		Slug:        p.Slug,
		Status:      p.Status.String,
		ViewCount:   p.ViewCount.Int32,
		CreatedAt:   p.CreatedAt.Time,
		Tags:        s.getPostTags(ctx, p.ID),
	}
	if p.Excerpt.Valid {
		resp.Excerpt = p.Excerpt.String
	}
	if p.CategoryID.Valid {
		resp.CategoryID = &p.CategoryID.Int32
	}
	if p.CategoryName.Valid {
		resp.CategoryName = &p.CategoryName.String
	}
	if p.CategorySlug.Valid {
		resp.CategorySlug = &p.CategorySlug.String
	}
	if p.ReadingTime.Valid {
		resp.ReadingTime = &p.ReadingTime.Int32
	}
	if p.Thumbnail.Valid {
		resp.Thumbnail = &p.Thumbnail.String
	}
	if p.PublishedAt.Valid {
		resp.PublishedAt = &p.PublishedAt.Time
	}
	return resp
}

func (s *PostService) toPostListResponseFromTag(ctx context.Context, p sqlc.ListPublishedPostsByTagRow) model.PostListResponse {
	resp := model.PostListResponse{
		ID:          p.ID,
		Title:       p.Title,
		Slug:        p.Slug,
		Status:      p.Status.String,
		ViewCount:   p.ViewCount.Int32,
		CreatedAt:   p.CreatedAt.Time,
		Tags:        s.getPostTags(ctx, p.ID),
	}
	if p.Excerpt.Valid {
		resp.Excerpt = p.Excerpt.String
	}
	if p.CategoryID.Valid {
		resp.CategoryID = &p.CategoryID.Int32
	}
	if p.CategoryName.Valid {
		resp.CategoryName = &p.CategoryName.String
	}
	if p.CategorySlug.Valid {
		resp.CategorySlug = &p.CategorySlug.String
	}
	if p.ReadingTime.Valid {
		resp.ReadingTime = &p.ReadingTime.Int32
	}
	if p.Thumbnail.Valid {
		resp.Thumbnail = &p.Thumbnail.String
	}
	if p.PublishedAt.Valid {
		resp.PublishedAt = &p.PublishedAt.Time
	}
	return resp
}

func (s *PostService) toPostListResponseFromSearch(ctx context.Context, p sqlc.SearchPublishedPostsRow) model.PostListResponse {
	resp := model.PostListResponse{
		ID:          p.ID,
		Title:       p.Title,
		Slug:        p.Slug,
		Status:      p.Status.String,
		ViewCount:   p.ViewCount.Int32,
		CreatedAt:   p.CreatedAt.Time,
		Tags:        s.getPostTags(ctx, p.ID),
	}
	if p.Excerpt.Valid {
		resp.Excerpt = p.Excerpt.String
	}
	if p.CategoryID.Valid {
		resp.CategoryID = &p.CategoryID.Int32
	}
	if p.CategoryName.Valid {
		resp.CategoryName = &p.CategoryName.String
	}
	if p.CategorySlug.Valid {
		resp.CategorySlug = &p.CategorySlug.String
	}
	if p.ReadingTime.Valid {
		resp.ReadingTime = &p.ReadingTime.Int32
	}
	if p.Thumbnail.Valid {
		resp.Thumbnail = &p.Thumbnail.String
	}
	if p.PublishedAt.Valid {
		resp.PublishedAt = &p.PublishedAt.Time
	}
	return resp
}

func (s *PostService) toPostListResponseFromAll(ctx context.Context, p sqlc.ListAllPostsRow) model.PostListResponse {
	resp := model.PostListResponse{
		ID:          p.ID,
		Title:       p.Title,
		Slug:        p.Slug,
		Status:      p.Status.String,
		ViewCount:   p.ViewCount.Int32,
		CreatedAt:   p.CreatedAt.Time,
		Tags:        s.getPostTags(ctx, p.ID),
	}
	if p.Excerpt.Valid {
		resp.Excerpt = p.Excerpt.String
	}
	if p.CategoryID.Valid {
		resp.CategoryID = &p.CategoryID.Int32
	}
	if p.CategoryName.Valid {
		resp.CategoryName = &p.CategoryName.String
	}
	if p.CategorySlug.Valid {
		resp.CategorySlug = &p.CategorySlug.String
	}
	if p.ReadingTime.Valid {
		resp.ReadingTime = &p.ReadingTime.Int32
	}
	if p.Thumbnail.Valid {
		resp.Thumbnail = &p.Thumbnail.String
	}
	if p.PublishedAt.Valid {
		resp.PublishedAt = &p.PublishedAt.Time
	}
	return resp
}

func (s *PostService) toPostListResponseFromStatus(ctx context.Context, p sqlc.ListPostsByStatusRow) model.PostListResponse {
	resp := model.PostListResponse{
		ID:          p.ID,
		Title:       p.Title,
		Slug:        p.Slug,
		Status:      p.Status.String,
		ViewCount:   p.ViewCount.Int32,
		CreatedAt:   p.CreatedAt.Time,
		Tags:        s.getPostTags(ctx, p.ID),
	}
	if p.Excerpt.Valid {
		resp.Excerpt = p.Excerpt.String
	}
	if p.CategoryID.Valid {
		resp.CategoryID = &p.CategoryID.Int32
	}
	if p.CategoryName.Valid {
		resp.CategoryName = &p.CategoryName.String
	}
	if p.CategorySlug.Valid {
		resp.CategorySlug = &p.CategorySlug.String
	}
	if p.ReadingTime.Valid {
		resp.ReadingTime = &p.ReadingTime.Int32
	}
	if p.Thumbnail.Valid {
		resp.Thumbnail = &p.Thumbnail.String
	}
	if p.PublishedAt.Valid {
		resp.PublishedAt = &p.PublishedAt.Time
	}
	return resp
}

func (s *PostService) toPostResponse(ctx context.Context, p sqlc.GetPublishedPostBySlugRow) *model.PostResponse {
	resp := &model.PostResponse{
		ID:          p.ID,
		Title:       p.Title,
		Slug:        p.Slug,
		Content:     p.Content,
		Status:      p.Status.String,
		ViewCount:   p.ViewCount.Int32,
		CreatedAt:   p.CreatedAt.Time,
		UpdatedAt:   p.UpdatedAt.Time,
		Tags:        s.getPostTags(ctx, p.ID),
	}
	if p.Excerpt.Valid {
		resp.Excerpt = p.Excerpt.String
	}
	if p.CategoryID.Valid {
		resp.CategoryID = &p.CategoryID.Int32
	}
	if p.CategoryName.Valid {
		resp.CategoryName = &p.CategoryName.String
	}
	if p.CategorySlug.Valid {
		resp.CategorySlug = &p.CategorySlug.String
	}
	if p.ReadingTime.Valid {
		resp.ReadingTime = &p.ReadingTime.Int32
	}
	if p.Thumbnail.Valid {
		resp.Thumbnail = &p.Thumbnail.String
	}
	if p.PublishedAt.Valid {
		resp.PublishedAt = &p.PublishedAt.Time
	}
	return resp
}

func (s *PostService) toPostResponseFromID(ctx context.Context, p sqlc.GetPostByIDRow) *model.PostResponse {
	resp := &model.PostResponse{
		ID:          p.ID,
		Title:       p.Title,
		Slug:        p.Slug,
		Content:     p.Content,
		Status:      p.Status.String,
		ViewCount:   p.ViewCount.Int32,
		CreatedAt:   p.CreatedAt.Time,
		UpdatedAt:   p.UpdatedAt.Time,
		Tags:        s.getPostTags(ctx, p.ID),
	}
	if p.Excerpt.Valid {
		resp.Excerpt = p.Excerpt.String
	}
	if p.CategoryID.Valid {
		resp.CategoryID = &p.CategoryID.Int32
	}
	if p.CategoryName.Valid {
		resp.CategoryName = &p.CategoryName.String
	}
	if p.CategorySlug.Valid {
		resp.CategorySlug = &p.CategorySlug.String
	}
	if p.ReadingTime.Valid {
		resp.ReadingTime = &p.ReadingTime.Int32
	}
	if p.Thumbnail.Valid {
		resp.Thumbnail = &p.Thumbnail.String
	}
	if p.PublishedAt.Valid {
		resp.PublishedAt = &p.PublishedAt.Time
	}
	return resp
}

func ptrToInt32(p *int32) int32 {
	if p == nil {
		return 0
	}
	return *p
}
