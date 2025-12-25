# Interfaces Layer

외부 세계(HTTP, CLI 등)와 애플리케이션 간의 데이터 변환을 담당하는 계층입니다.

## 원칙

- **DTO 정의**: HTTP 요청/응답에 특화된 데이터 구조
- **Mapper 제공**: Entity ↔ DTO 변환 로직
- **프레임워크 중립**: 가능한 Gin 의존성 최소화
- **검증 태그**: JSON 바인딩 및 유효성 검사 태그 사용

## 디렉토리 구조

```
interfaces/
└── http/
    ├── dto/              # Data Transfer Objects
    │   ├── category.go
    │   ├── tag.go
    │   ├── post.go
    │   ├── project.go
    │   ├── media.go
    │   ├── auth.go
    │   └── dashboard.go
    └── mapper/           # Entity ↔ DTO 변환
        ├── category_mapper.go
        ├── tag_mapper.go
        ├── post_mapper.go
        ├── project_mapper.go
        ├── media_mapper.go
        ├── auth_mapper.go
        └── dashboard_mapper.go
```

## DTO (Data Transfer Objects)

### 파일 구조

각 도메인별로 Request와 Response DTO를 정의합니다:

```go
// dto/category.go
package dto

// Request DTOs - 클라이언트 → 서버
type CreateCategoryRequest struct {
    Name        string  `json:"name" binding:"required"`
    Slug        string  `json:"slug" binding:"required"`
    Description *string `json:"description"`
}

type UpdateCategoryRequest struct {
    Name        *string `json:"name"`
    Slug        *string `json:"slug"`
    Description *string `json:"description"`
}

// Response DTOs - 서버 → 클라이언트
type CategoryResponse struct {
    ID          int32   `json:"id"`
    Name        string  `json:"name"`
    Slug        string  `json:"slug"`
    Description *string `json:"description,omitempty"`
    PostCount   int64   `json:"post_count"`
    CreatedAt   string  `json:"created_at"`
    UpdatedAt   *string `json:"updated_at,omitempty"`
}
```

### 파일별 내용

| 파일 | Request DTOs | Response DTOs |
|------|--------------|---------------|
| `category.go` | Create, Update | CategoryResponse |
| `tag.go` | Create, Update | TagResponse |
| `post.go` | Create, Update | PostResponse, PostListResponse |
| `project.go` | Create, Update, Reorder | ProjectResponse |
| `media.go` | (multipart) | MediaResponse |
| `auth.go` | LoginRequest | LoginResponse, AdminResponse |
| `dashboard.go` | - | DashboardStatsResponse |

### JSON 태그 컨벤션

```go
// 필수 필드
Name string `json:"name" binding:"required"`

// 선택 필드 (nullable)
Description *string `json:"description,omitempty"`

// 읽기 전용 (응답에만)
ID int32 `json:"id"`

// 시간 필드는 문자열로 포맷
CreatedAt string `json:"created_at"`
```

## Mapper

### 역할

Entity와 DTO 간의 변환을 담당합니다:

```go
// mapper/category_mapper.go
package mapper

import (
    "github.com/ydonggwui/blog-api/internal/domain/entity"
    domainService "github.com/ydonggwui/blog-api/internal/domain/service"
    "github.com/ydonggwui/blog-api/internal/interfaces/http/dto"
)

// Entity → Response DTO
func ToCategoryResponse(cat *entity.Category) *dto.CategoryResponse {
    if cat == nil {
        return nil
    }
    return &dto.CategoryResponse{
        ID:          cat.ID,
        Name:        cat.Name,
        Slug:        cat.Slug,
        Description: cat.Description,
        PostCount:   cat.PostCount,
        CreatedAt:   cat.CreatedAt.Format(time.RFC3339),
        UpdatedAt:   formatTimePtr(cat.UpdatedAt),
    }
}

// Entity 슬라이스 → Response DTO 슬라이스
func ToCategoryResponses(cats []entity.Category) []*dto.CategoryResponse {
    result := make([]*dto.CategoryResponse, len(cats))
    for i := range cats {
        result[i] = ToCategoryResponse(&cats[i])
    }
    return result
}

// Request DTO → Command (Service 입력)
func ToCreateCategoryCommand(req *dto.CreateCategoryRequest) *domainService.CreateCategoryCommand {
    return &domainService.CreateCategoryCommand{
        Name:        req.Name,
        Slug:        req.Slug,
        Description: req.Description,
    }
}
```

### 변환 방향

```
┌─────────────┐      ┌─────────────┐      ┌─────────────┐
│   HTTP      │      │   Mapper    │      │   Domain    │
│   Layer     │◄────►│             │◄────►│   Layer     │
└─────────────┘      └─────────────┘      └─────────────┘

Request DTO ──────► Command ──────► Service
    │                                  │
    │                                  ▼
    │                              Entity
    │                                  │
Response DTO ◄───── Mapper ◄──────────┘
```

### 시간 포맷팅 헬퍼

```go
func formatTimePtr(t *time.Time) *string {
    if t == nil {
        return nil
    }
    s := t.Format(time.RFC3339)
    return &s
}
```

## Handler에서의 사용

```go
// handler/admin/category.go
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
    // 1. Request DTO로 바인딩
    var req dto.CreateCategoryRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        handler.BadRequest(c, "Invalid request body")
        return
    }

    // 2. DTO → Command 변환
    cmd := mapper.ToCreateCategoryCommand(&req)

    // 3. Service 호출 (Entity 반환)
    category, err := h.categoryService.CreateCategory(c.Request.Context(), cmd)
    if err != nil {
        // 에러 처리
        return
    }

    // 4. Entity → Response DTO 변환
    handler.Created(c, mapper.ToCategoryResponse(category))
}
```

## 새 DTO/Mapper 추가 시

1. `dto/` 에 Request/Response 구조체 정의
2. `mapper/` 에 변환 함수 작성:
   - `ToXxxResponse`: Entity → Response
   - `ToXxxResponses`: []Entity → []Response
   - `ToXxxCommand`: Request → Command
3. Handler에서 import하여 사용

## 주의사항

- DTO는 JSON 직렬화에 최적화
- Entity의 모든 필드를 노출할 필요 없음 (예: 비밀번호 제외)
- 날짜/시간은 RFC3339 문자열로 변환
- nullable 필드는 `omitempty` 사용
- Mapper는 순수 함수로 작성 (부작용 없음)
