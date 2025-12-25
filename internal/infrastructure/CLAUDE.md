# Infrastructure Layer

외부 시스템(데이터베이스, 캐시, 스토리지)과의 통신을 담당하는 계층입니다.

## 원칙

- **Domain 인터페이스 구현**: `domain/repository/` 인터페이스를 구현
- **기술 세부사항 캡슐화**: sqlc, Redis, MinIO 등의 세부사항을 숨김
- **Entity 변환**: 외부 시스템의 데이터 구조를 Domain Entity로 변환

## 디렉토리 구조

```
infrastructure/
├── persistence/           # 데이터 영속성
│   ├── postgres/          # PostgreSQL Repository 구현
│   │   ├── *_repository.go
│   │   └── mapper.go      # sqlc 모델 → Entity 변환
│   └── redis/             # Redis Repository 구현
│       └── view_repository.go
└── storage/               # 파일 스토리지
    └── minio/             # MinIO (S3 호환) 구현
        └── storage_repository.go
```

## PostgreSQL Repository

### 파일별 역할

| 파일 | 구현 인터페이스 | 설명 |
|------|----------------|------|
| `category_repository.go` | `CategoryRepository` | 카테고리 CRUD |
| `tag_repository.go` | `TagRepository` | 태그 CRUD |
| `post_repository.go` | `PostRepository` | 게시물 CRUD + 검색 |
| `project_repository.go` | `ProjectRepository` | 프로젝트 CRUD + 정렬 |
| `media_repository.go` | `MediaRepository` | 미디어 메타데이터 |
| `admin_repository.go` | `AdminRepository` | 관리자 계정 |
| `dashboard_repository.go` | `DashboardRepository` | 통계 집계 쿼리 |
| `mapper.go` | - | sqlc 모델 ↔ Entity 변환 함수 |

### 구현 패턴

```go
package postgres

import (
    "context"
    "database/sql"

    "github.com/ydonggwui/blog-api/internal/database/sqlc"
    "github.com/ydonggwui/blog-api/internal/domain/entity"
    "github.com/ydonggwui/blog-api/internal/domain/repository"
)

type categoryRepository struct {
    queries *sqlc.Queries
}

// 생성자: 인터페이스 반환
func NewCategoryRepository(queries *sqlc.Queries) repository.CategoryRepository {
    return &categoryRepository{queries: queries}
}

// 메서드 구현: sqlc 호출 후 Entity로 변환
func (r *categoryRepository) FindByID(ctx context.Context, id int32) (*entity.Category, error) {
    row, err := r.queries.GetCategoryByID(ctx, id)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil // 또는 도메인 에러
        }
        return nil, err
    }
    return categoryRowToEntity(row), nil
}
```

### Mapper 패턴

```go
// mapper.go
package postgres

import (
    "github.com/ydonggwui/blog-api/internal/database/sqlc"
    "github.com/ydonggwui/blog-api/internal/domain/entity"
)

// sqlc Row → Entity
func categoryRowToEntity(row sqlc.Category) *entity.Category {
    return &entity.Category{
        ID:          row.ID,
        Name:        row.Name,
        Slug:        row.Slug,
        Description: nullStringToPtr(row.Description),
        CreatedAt:   row.CreatedAt,
        UpdatedAt:   nullTimeToPtr(row.UpdatedAt),
    }
}

// nullable 타입 변환 헬퍼
func nullStringToPtr(ns sql.NullString) *string {
    if ns.Valid {
        return &ns.String
    }
    return nil
}

func nullTimeToPtr(nt sql.NullTime) *time.Time {
    if nt.Valid {
        return &nt.Time
    }
    return nil
}
```

## Redis Repository

### view_repository.go

조회수 중복 방지를 위한 Redis 구현:

```go
package redis

import (
    "context"
    "time"

    "github.com/redis/go-redis/v9"
    "github.com/ydonggwui/blog-api/internal/domain/repository"
)

type viewRepository struct {
    client *redis.Client
}

func NewViewRepository(client *redis.Client) repository.ViewRepository {
    return &viewRepository{client: client}
}

// SetNX: 키가 없을 때만 설정 (중복 방지)
func (r *viewRepository) SetViewIfNotExists(ctx context.Context, key string, ttl time.Duration) (bool, error) {
    return r.client.SetNX(ctx, key, "1", ttl).Result()
}
```

## MinIO Storage

### storage_repository.go

S3 호환 파일 스토리지 구현:

```go
package minio

import (
    "context"
    "mime/multipart"

    "github.com/minio/minio-go/v7"
    "github.com/ydonggwui/blog-api/internal/config"
    "github.com/ydonggwui/blog-api/internal/domain/entity"
    "github.com/ydonggwui/blog-api/internal/domain/repository"
)

type storageRepository struct {
    client *minio.Client
    config *config.MinIOConfig
}

func NewStorageRepository(client *minio.Client, cfg *config.MinIOConfig) repository.StorageRepository {
    return &storageRepository{client: client, config: cfg}
}

func (r *storageRepository) Upload(ctx context.Context, file *entity.UploadedFile, header *multipart.FileHeader) (string, error) {
    // MinIO 업로드 로직
    // ...
    return url, nil
}
```

## 새 Repository 추가 시

1. `domain/repository/` 에 인터페이스 정의
2. 해당 기술 디렉토리에 구현체 작성:
   - PostgreSQL: `persistence/postgres/`
   - Redis: `persistence/redis/`
   - MinIO: `storage/minio/`
3. `mapper.go`에 변환 함수 추가 (필요시)
4. `router/router.go`에서 DI 설정

## sqlc 쿼리 추가 시

1. `database/query.sql`에 쿼리 추가
2. `sqlc generate` 실행
3. 생성된 메서드를 Repository에서 호출
4. `mapper.go`에서 Entity 변환

## 주의사항

- Domain 계층 import만 허용 (`domain/entity`, `domain/repository`)
- Application 계층 import 금지
- 에러는 그대로 반환하거나, 필요시 도메인 에러로 래핑
- sql.ErrNoRows는 nil 또는 도메인 에러로 변환
