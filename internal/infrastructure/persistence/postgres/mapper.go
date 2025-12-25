package postgres

import (
	"database/sql"

	"github.com/ydonggwui/blog-api/internal/database/sqlc"
	"github.com/ydonggwui/blog-api/internal/domain/entity"
)

// Category mappers

func toCategoryEntity(c sqlc.Category) *entity.Category {
	category := &entity.Category{
		ID:        c.ID,
		Name:      c.Name,
		Slug:      c.Slug,
		SortOrder: c.SortOrder.Int32,
	}
	if c.Description.Valid {
		category.Description = c.Description.String
	}
	if c.CreatedAt.Valid {
		category.CreatedAt = c.CreatedAt.Time
	}
	return category
}

func toCategoryEntities(categories []sqlc.Category) []entity.Category {
	result := make([]entity.Category, len(categories))
	for i, c := range categories {
		result[i] = *toCategoryEntity(c)
	}
	return result
}

func toCreateCategoryParams(c *entity.Category) sqlc.CreateCategoryParams {
	return sqlc.CreateCategoryParams{
		Name:        c.Name,
		Slug:        c.Slug,
		Description: sql.NullString{String: c.Description, Valid: c.Description != ""},
		SortOrder:   sql.NullInt32{Int32: c.SortOrder, Valid: true},
	}
}

func toUpdateCategoryParams(c *entity.Category) sqlc.UpdateCategoryParams {
	return sqlc.UpdateCategoryParams{
		ID:          c.ID,
		Name:        c.Name,
		Slug:        c.Slug,
		Description: sql.NullString{String: c.Description, Valid: c.Description != ""},
		SortOrder:   sql.NullInt32{Int32: c.SortOrder, Valid: true},
	}
}

// Tag mappers

func toTagEntity(t sqlc.Tag) *entity.Tag {
	tag := &entity.Tag{
		ID:   t.ID,
		Name: t.Name,
		Slug: t.Slug,
	}
	if t.CreatedAt.Valid {
		tag.CreatedAt = t.CreatedAt.Time
	}
	return tag
}

func toTagEntities(tags []sqlc.Tag) []entity.Tag {
	result := make([]entity.Tag, len(tags))
	for i, t := range tags {
		result[i] = *toTagEntity(t)
	}
	return result
}

func toTagEntityWithPostCount(t sqlc.ListTagsWithPostCountRow) *entity.Tag {
	tag := &entity.Tag{
		ID:        t.ID,
		Name:      t.Name,
		Slug:      t.Slug,
		PostCount: t.PostCount,
	}
	if t.CreatedAt.Valid {
		tag.CreatedAt = t.CreatedAt.Time
	}
	return tag
}

func toTagEntitiesWithPostCount(tags []sqlc.ListTagsWithPostCountRow) []entity.Tag {
	result := make([]entity.Tag, len(tags))
	for i, t := range tags {
		result[i] = *toTagEntityWithPostCount(t)
	}
	return result
}

func toCreateTagParams(t *entity.Tag) sqlc.CreateTagParams {
	return sqlc.CreateTagParams{
		Name: t.Name,
		Slug: t.Slug,
	}
}

func toUpdateTagParams(t *entity.Tag) sqlc.UpdateTagParams {
	return sqlc.UpdateTagParams{
		ID:   t.ID,
		Name: t.Name,
		Slug: t.Slug,
	}
}

// Post mappers

func toPostEntity(p sqlc.Post) *entity.Post {
	post := &entity.Post{
		ID:      p.ID,
		Title:   p.Title,
		Slug:    p.Slug,
		Content: p.Content,
		Status:  entity.PostStatus(p.Status.String),
	}
	if p.Excerpt.Valid {
		post.Excerpt = p.Excerpt.String
	}
	if p.CategoryID.Valid {
		post.CategoryID = &p.CategoryID.Int32
	}
	if p.ViewCount.Valid {
		post.ViewCount = p.ViewCount.Int32
	}
	if p.ReadingTime.Valid {
		post.ReadingTime = p.ReadingTime.Int32
	}
	if p.Thumbnail.Valid {
		post.Thumbnail = p.Thumbnail.String
	}
	if p.CreatedAt.Valid {
		post.CreatedAt = p.CreatedAt.Time
	}
	if p.UpdatedAt.Valid {
		post.UpdatedAt = p.UpdatedAt.Time
	}
	if p.PublishedAt.Valid {
		post.PublishedAt = &p.PublishedAt.Time
	}
	return post
}

func toPostWithDetails(p sqlc.GetPostByIDRow, tags []entity.TagBrief) *entity.PostWithDetails {
	post := &entity.PostWithDetails{
		Post: entity.Post{
			ID:      p.ID,
			Title:   p.Title,
			Slug:    p.Slug,
			Content: p.Content,
			Status:  entity.PostStatus(p.Status.String),
		},
		Tags: tags,
	}
	if p.Excerpt.Valid {
		post.Excerpt = p.Excerpt.String
	}
	if p.CategoryID.Valid {
		post.CategoryID = &p.CategoryID.Int32
	}
	if p.CategoryName.Valid {
		post.CategoryName = p.CategoryName.String
	}
	if p.CategorySlug.Valid {
		post.CategorySlug = p.CategorySlug.String
	}
	if p.ViewCount.Valid {
		post.ViewCount = p.ViewCount.Int32
	}
	if p.ReadingTime.Valid {
		post.ReadingTime = p.ReadingTime.Int32
	}
	if p.Thumbnail.Valid {
		post.Thumbnail = p.Thumbnail.String
	}
	if p.CreatedAt.Valid {
		post.CreatedAt = p.CreatedAt.Time
	}
	if p.UpdatedAt.Valid {
		post.UpdatedAt = p.UpdatedAt.Time
	}
	if p.PublishedAt.Valid {
		post.PublishedAt = &p.PublishedAt.Time
	}
	return post
}

func toPostWithDetailsFromSlug(p sqlc.GetPublishedPostBySlugRow, tags []entity.TagBrief) *entity.PostWithDetails {
	post := &entity.PostWithDetails{
		Post: entity.Post{
			ID:      p.ID,
			Title:   p.Title,
			Slug:    p.Slug,
			Content: p.Content,
			Status:  entity.PostStatus(p.Status.String),
		},
		Tags: tags,
	}
	if p.Excerpt.Valid {
		post.Excerpt = p.Excerpt.String
	}
	if p.CategoryID.Valid {
		post.CategoryID = &p.CategoryID.Int32
	}
	if p.CategoryName.Valid {
		post.CategoryName = p.CategoryName.String
	}
	if p.CategorySlug.Valid {
		post.CategorySlug = p.CategorySlug.String
	}
	if p.ViewCount.Valid {
		post.ViewCount = p.ViewCount.Int32
	}
	if p.ReadingTime.Valid {
		post.ReadingTime = p.ReadingTime.Int32
	}
	if p.Thumbnail.Valid {
		post.Thumbnail = p.Thumbnail.String
	}
	if p.CreatedAt.Valid {
		post.CreatedAt = p.CreatedAt.Time
	}
	if p.UpdatedAt.Valid {
		post.UpdatedAt = p.UpdatedAt.Time
	}
	if p.PublishedAt.Valid {
		post.PublishedAt = &p.PublishedAt.Time
	}
	return post
}

func toPostWithDetailsFromList(p sqlc.ListPublishedPostsRow, tags []entity.TagBrief) entity.PostWithDetails {
	post := entity.PostWithDetails{
		Post: entity.Post{
			ID:     p.ID,
			Title:  p.Title,
			Slug:   p.Slug,
			Status: entity.PostStatus(p.Status.String),
		},
		Tags: tags,
	}
	if p.Excerpt.Valid {
		post.Excerpt = p.Excerpt.String
	}
	if p.CategoryID.Valid {
		post.CategoryID = &p.CategoryID.Int32
	}
	if p.CategoryName.Valid {
		post.CategoryName = p.CategoryName.String
	}
	if p.CategorySlug.Valid {
		post.CategorySlug = p.CategorySlug.String
	}
	if p.ViewCount.Valid {
		post.ViewCount = p.ViewCount.Int32
	}
	if p.ReadingTime.Valid {
		post.ReadingTime = p.ReadingTime.Int32
	}
	if p.Thumbnail.Valid {
		post.Thumbnail = p.Thumbnail.String
	}
	if p.CreatedAt.Valid {
		post.CreatedAt = p.CreatedAt.Time
	}
	if p.PublishedAt.Valid {
		post.PublishedAt = &p.PublishedAt.Time
	}
	return post
}

func toPostWithDetailsFromCategory(p sqlc.ListPublishedPostsByCategoryRow, tags []entity.TagBrief) entity.PostWithDetails {
	post := entity.PostWithDetails{
		Post: entity.Post{
			ID:     p.ID,
			Title:  p.Title,
			Slug:   p.Slug,
			Status: entity.PostStatus(p.Status.String),
		},
		Tags: tags,
	}
	if p.Excerpt.Valid {
		post.Excerpt = p.Excerpt.String
	}
	if p.CategoryID.Valid {
		post.CategoryID = &p.CategoryID.Int32
	}
	if p.CategoryName.Valid {
		post.CategoryName = p.CategoryName.String
	}
	if p.CategorySlug.Valid {
		post.CategorySlug = p.CategorySlug.String
	}
	if p.ViewCount.Valid {
		post.ViewCount = p.ViewCount.Int32
	}
	if p.ReadingTime.Valid {
		post.ReadingTime = p.ReadingTime.Int32
	}
	if p.Thumbnail.Valid {
		post.Thumbnail = p.Thumbnail.String
	}
	if p.CreatedAt.Valid {
		post.CreatedAt = p.CreatedAt.Time
	}
	if p.PublishedAt.Valid {
		post.PublishedAt = &p.PublishedAt.Time
	}
	return post
}

func toPostWithDetailsFromTag(p sqlc.ListPublishedPostsByTagRow, tags []entity.TagBrief) entity.PostWithDetails {
	post := entity.PostWithDetails{
		Post: entity.Post{
			ID:     p.ID,
			Title:  p.Title,
			Slug:   p.Slug,
			Status: entity.PostStatus(p.Status.String),
		},
		Tags: tags,
	}
	if p.Excerpt.Valid {
		post.Excerpt = p.Excerpt.String
	}
	if p.CategoryID.Valid {
		post.CategoryID = &p.CategoryID.Int32
	}
	if p.CategoryName.Valid {
		post.CategoryName = p.CategoryName.String
	}
	if p.CategorySlug.Valid {
		post.CategorySlug = p.CategorySlug.String
	}
	if p.ViewCount.Valid {
		post.ViewCount = p.ViewCount.Int32
	}
	if p.ReadingTime.Valid {
		post.ReadingTime = p.ReadingTime.Int32
	}
	if p.Thumbnail.Valid {
		post.Thumbnail = p.Thumbnail.String
	}
	if p.CreatedAt.Valid {
		post.CreatedAt = p.CreatedAt.Time
	}
	if p.PublishedAt.Valid {
		post.PublishedAt = &p.PublishedAt.Time
	}
	return post
}

func toPostWithDetailsFromSearch(p sqlc.SearchPublishedPostsRow, tags []entity.TagBrief) entity.PostWithDetails {
	post := entity.PostWithDetails{
		Post: entity.Post{
			ID:     p.ID,
			Title:  p.Title,
			Slug:   p.Slug,
			Status: entity.PostStatus(p.Status.String),
		},
		Tags: tags,
	}
	if p.Excerpt.Valid {
		post.Excerpt = p.Excerpt.String
	}
	if p.CategoryID.Valid {
		post.CategoryID = &p.CategoryID.Int32
	}
	if p.CategoryName.Valid {
		post.CategoryName = p.CategoryName.String
	}
	if p.CategorySlug.Valid {
		post.CategorySlug = p.CategorySlug.String
	}
	if p.ViewCount.Valid {
		post.ViewCount = p.ViewCount.Int32
	}
	if p.ReadingTime.Valid {
		post.ReadingTime = p.ReadingTime.Int32
	}
	if p.Thumbnail.Valid {
		post.Thumbnail = p.Thumbnail.String
	}
	if p.CreatedAt.Valid {
		post.CreatedAt = p.CreatedAt.Time
	}
	if p.PublishedAt.Valid {
		post.PublishedAt = &p.PublishedAt.Time
	}
	return post
}

func toPostWithDetailsFromAll(p sqlc.ListAllPostsRow, tags []entity.TagBrief) entity.PostWithDetails {
	post := entity.PostWithDetails{
		Post: entity.Post{
			ID:     p.ID,
			Title:  p.Title,
			Slug:   p.Slug,
			Status: entity.PostStatus(p.Status.String),
		},
		Tags: tags,
	}
	if p.Excerpt.Valid {
		post.Excerpt = p.Excerpt.String
	}
	if p.CategoryID.Valid {
		post.CategoryID = &p.CategoryID.Int32
	}
	if p.CategoryName.Valid {
		post.CategoryName = p.CategoryName.String
	}
	if p.CategorySlug.Valid {
		post.CategorySlug = p.CategorySlug.String
	}
	if p.ViewCount.Valid {
		post.ViewCount = p.ViewCount.Int32
	}
	if p.ReadingTime.Valid {
		post.ReadingTime = p.ReadingTime.Int32
	}
	if p.Thumbnail.Valid {
		post.Thumbnail = p.Thumbnail.String
	}
	if p.CreatedAt.Valid {
		post.CreatedAt = p.CreatedAt.Time
	}
	if p.PublishedAt.Valid {
		post.PublishedAt = &p.PublishedAt.Time
	}
	return post
}

func toPostWithDetailsFromStatus(p sqlc.ListPostsByStatusRow, tags []entity.TagBrief) entity.PostWithDetails {
	post := entity.PostWithDetails{
		Post: entity.Post{
			ID:     p.ID,
			Title:  p.Title,
			Slug:   p.Slug,
			Status: entity.PostStatus(p.Status.String),
		},
		Tags: tags,
	}
	if p.Excerpt.Valid {
		post.Excerpt = p.Excerpt.String
	}
	if p.CategoryID.Valid {
		post.CategoryID = &p.CategoryID.Int32
	}
	if p.CategoryName.Valid {
		post.CategoryName = p.CategoryName.String
	}
	if p.CategorySlug.Valid {
		post.CategorySlug = p.CategorySlug.String
	}
	if p.ViewCount.Valid {
		post.ViewCount = p.ViewCount.Int32
	}
	if p.ReadingTime.Valid {
		post.ReadingTime = p.ReadingTime.Int32
	}
	if p.Thumbnail.Valid {
		post.Thumbnail = p.Thumbnail.String
	}
	if p.CreatedAt.Valid {
		post.CreatedAt = p.CreatedAt.Time
	}
	if p.PublishedAt.Valid {
		post.PublishedAt = &p.PublishedAt.Time
	}
	return post
}

func toTagBriefs(tags []sqlc.Tag) []entity.TagBrief {
	result := make([]entity.TagBrief, len(tags))
	for i, t := range tags {
		result[i] = entity.TagBrief{
			ID:   t.ID,
			Name: t.Name,
			Slug: t.Slug,
		}
	}
	return result
}

func toCreatePostParams(p *entity.Post) sqlc.CreatePostParams {
	return sqlc.CreatePostParams{
		Title:       p.Title,
		Slug:        p.Slug,
		Content:     p.Content,
		Excerpt:     sql.NullString{String: p.Excerpt, Valid: p.Excerpt != ""},
		CategoryID:  sql.NullInt32{Int32: ptrToInt32(p.CategoryID), Valid: p.CategoryID != nil},
		Status:      sql.NullString{String: string(p.Status), Valid: true},
		ReadingTime: sql.NullInt32{Int32: p.ReadingTime, Valid: true},
		Thumbnail:   sql.NullString{String: p.Thumbnail, Valid: p.Thumbnail != ""},
	}
}

func toUpdatePostParams(p *entity.Post) sqlc.UpdatePostParams {
	return sqlc.UpdatePostParams{
		ID:          p.ID,
		Title:       p.Title,
		Slug:        p.Slug,
		Content:     p.Content,
		Excerpt:     sql.NullString{String: p.Excerpt, Valid: p.Excerpt != ""},
		CategoryID:  sql.NullInt32{Int32: ptrToInt32(p.CategoryID), Valid: p.CategoryID != nil},
		ReadingTime: sql.NullInt32{Int32: p.ReadingTime, Valid: true},
		Thumbnail:   sql.NullString{String: p.Thumbnail, Valid: p.Thumbnail != ""},
	}
}

func ptrToInt32(p *int32) int32 {
	if p == nil {
		return 0
	}
	return *p
}

// Media mappers

func toMediaEntity(m sqlc.Medium) *entity.Media {
	media := &entity.Media{
		ID:           m.ID,
		Filename:     m.Filename,
		OriginalName: m.OriginalName,
		Path:         m.Path,
		URL:          m.Url,
	}
	if m.MimeType.Valid {
		media.MimeType = m.MimeType.String
	}
	if m.Size.Valid {
		media.Size = m.Size.Int64
	}
	if m.Width.Valid {
		media.Width = m.Width.Int32
	}
	if m.Height.Valid {
		media.Height = m.Height.Int32
	}
	if m.CreatedAt.Valid {
		media.CreatedAt = m.CreatedAt.Time
	}
	return media
}

func toMediaEntities(media []sqlc.Medium) []entity.Media {
	result := make([]entity.Media, len(media))
	for i, m := range media {
		result[i] = *toMediaEntity(m)
	}
	return result
}

func toCreateMediaParams(m *entity.Media) sqlc.CreateMediaParams {
	return sqlc.CreateMediaParams{
		Filename:     m.Filename,
		OriginalName: m.OriginalName,
		Path:         m.Path,
		Url:          m.URL,
		MimeType:     sql.NullString{String: m.MimeType, Valid: m.MimeType != ""},
		Size:         sql.NullInt64{Int64: m.Size, Valid: m.Size > 0},
		Width:        sql.NullInt32{Int32: m.Width, Valid: m.Width > 0},
		Height:       sql.NullInt32{Int32: m.Height, Valid: m.Height > 0},
	}
}

// Admin mappers

func toAdminEntity(a sqlc.Admin) *entity.Admin {
	admin := &entity.Admin{
		ID:       a.ID,
		Username: a.Username,
		Password: a.Password,
	}
	if a.CreatedAt.Valid {
		admin.CreatedAt = a.CreatedAt.Time
	}
	if a.UpdatedAt.Valid {
		admin.UpdatedAt = &a.UpdatedAt.Time
	}
	return admin
}
