package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	DefaultPage    = 1
	DefaultPerPage = 10
	MaxPerPage     = 100
)

type Pagination struct {
	Page    int
	PerPage int
	Offset  int
}

func GetPagination(c *gin.Context) Pagination {
	page := parseIntParam(c, "page", DefaultPage)
	perPage := parseIntParam(c, "per_page", DefaultPerPage)

	if page < 1 {
		page = DefaultPage
	}
	if perPage < 1 {
		perPage = DefaultPerPage
	}
	if perPage > MaxPerPage {
		perPage = MaxPerPage
	}

	return Pagination{
		Page:    page,
		PerPage: perPage,
		Offset:  (page - 1) * perPage,
	}
}

func (p Pagination) ToMeta(total int64) *Meta {
	totalPages := int(total) / p.PerPage
	if int(total)%p.PerPage > 0 {
		totalPages++
	}

	return &Meta{
		Page:       p.Page,
		PerPage:    p.PerPage,
		Total:      total,
		TotalPages: totalPages,
	}
}

func parseIntParam(c *gin.Context, key string, defaultValue int) int {
	value := c.Query(key)
	if value == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return intValue
}
