package domain

import "errors"

// Category errors
var (
	ErrCategoryNotFound   = errors.New("category not found")
	ErrCategorySlugExists = errors.New("category slug already exists")
	ErrCategoryHasPosts   = errors.New("category has posts")
)

// Tag errors
var (
	ErrTagNotFound   = errors.New("tag not found")
	ErrTagSlugExists = errors.New("tag slug already exists")
	ErrTagHasPosts   = errors.New("tag has posts")
)

// Post errors
var (
	ErrPostNotFound = errors.New("post not found")
	ErrSlugExists   = errors.New("slug already exists")
)

// Project errors
var (
	ErrProjectNotFound   = errors.New("project not found")
	ErrProjectSlugExists = errors.New("project slug already exists")
)

// Media errors
var (
	ErrMediaNotFound   = errors.New("media not found")
	ErrInvalidFileType = errors.New("invalid file type")
	ErrFileTooLarge    = errors.New("file too large")
	ErrUploadFailed    = errors.New("upload failed")
)

// Auth errors
var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrAdminNotFound      = errors.New("admin not found")
)
