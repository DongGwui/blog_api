package entity

import "time"

// Media represents a media file entity
type Media struct {
	ID           int32
	Filename     string
	OriginalName string
	Path         string
	URL          string
	MimeType     string
	Size         int64
	Width        int32
	Height       int32
	ThumbnailSM  string
	ThumbnailMD  string
	CreatedAt    time.Time
}

// UploadedFile represents the result of a file upload
type UploadedFile struct {
	ID           int32
	Filename     string
	OriginalName string
	URL          string
	MimeType     string
	Size         int64
	ThumbnailSM  string
	ThumbnailMD  string
}
