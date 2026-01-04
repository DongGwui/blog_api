package dto

import "time"

// MediaResponse represents the response for a media file
type MediaResponse struct {
	ID           int32     `json:"id"`
	Filename     string    `json:"filename"`
	OriginalName string    `json:"original_name"`
	Path         string    `json:"path"`
	URL          string    `json:"url"`
	MimeType     string    `json:"mime_type,omitempty"`
	Size         int64     `json:"size,omitempty"`
	Width        int32     `json:"width,omitempty"`
	Height       int32     `json:"height,omitempty"`
	ThumbnailSM  string    `json:"thumbnail_sm,omitempty"`
	ThumbnailMD  string    `json:"thumbnail_md,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

// MediaListResponse represents a paginated list of media files
type MediaListResponse struct {
	Items []MediaResponse `json:"items"`
	Total int64           `json:"total"`
}

// UploadMediaResponse represents the response after uploading a file
type UploadMediaResponse struct {
	ID           int32  `json:"id"`
	Filename     string `json:"filename"`
	OriginalName string `json:"original_name"`
	URL          string `json:"url"`
	MimeType     string `json:"mime_type,omitempty"`
	Size         int64  `json:"size,omitempty"`
	ThumbnailSM  string `json:"thumbnail_sm,omitempty"`
	ThumbnailMD  string `json:"thumbnail_md,omitempty"`
}
